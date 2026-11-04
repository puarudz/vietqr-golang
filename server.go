package main

import (
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"io/fs"
	"log/slog"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"time"
)

type server struct {
	mux       *http.ServeMux
	index     []byte
	banksJSON []byte
}

func newServer(static fs.FS) (http.Handler, error) {
	index, err := fs.ReadFile(static, "index.html")
	if err != nil {
		return nil, fmt.Errorf("đọc index.html: %w", err)
	}

	banksJSON, err := json.Marshal(struct {
		Count int    `json:"count"`
		Banks []Bank `json:"banks"`
	}{Count: len(banks), Banks: banks})
	if err != nil {
		return nil, fmt.Errorf("mã hoá danh sách ngân hàng: %w", err)
	}

	s := &server{
		mux:       http.NewServeMux(),
		index:     index,
		banksJSON: banksJSON,
	}

	s.mux.HandleFunc("GET /{$}", s.handleIndex)
	s.mux.HandleFunc("GET /api/banks", s.handleBanks)
	s.mux.HandleFunc("POST /api/qr", s.handleGenerate)
	s.mux.HandleFunc("GET /api/qr.png", s.handleImage("png"))
	s.mux.HandleFunc("GET /api/qr.svg", s.handleImage("svg"))
	s.mux.HandleFunc("GET /api/decode", s.handleDecode)
	s.mux.HandleFunc("GET /healthz", s.handleHealth)
	s.mux.HandleFunc("GET /robots.txt", s.handleRobots)
	s.mux.HandleFunc("GET /favicon.ico", func(w http.ResponseWriter, r *http.Request) {
		http.Redirect(w, r, "/assets/favicon.svg", http.StatusMovedPermanently)
	})
	s.mux.Handle("GET /assets/", cacheStatic(http.StripPrefix("/assets/", http.FileServer(http.FS(static)))))
	s.mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		if strings.HasPrefix(r.URL.Path, "/api/") {
			writeError(w, http.StatusNotFound, "", "Không tìm thấy endpoint.")
			return
		}
		http.NotFound(w, r)
	})

	return s, nil
}

func (s *server) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	start := time.Now()
	sw := &statusWriter{ResponseWriter: w}
	defer func() {
		if rec := recover(); rec != nil {
			slog.Error("panic khi xử lý request", "path", r.URL.Path, "panic", rec)
			if !sw.wrote {
				writeError(sw, http.StatusInternalServerError, "", "Có lỗi xảy ra, vui lòng thử lại.")
			}
		}
		if strings.HasPrefix(r.URL.Path, "/api/") {
			slog.Info("api", "method", r.Method, "path", r.URL.Path, "status", sw.status, "duration", time.Since(start).Round(time.Millisecond))
		}
	}()

	sw.Header().Set("X-Content-Type-Options", "nosniff")
	sw.Header().Set("Referrer-Policy", "no-referrer")
	sw.Header().Set("X-Robots-Tag", "noindex, nofollow")
	sw.Header().Set("Content-Security-Policy", "default-src 'self'; img-src 'self' data:; style-src 'self'; font-src 'self'; script-src 'self'; connect-src 'self'; base-uri 'none'; form-action 'self'; frame-ancestors 'none'")

	if strings.HasPrefix(r.URL.Path, "/api/") {
		sw.Header().Set("Access-Control-Allow-Origin", "*")
		sw.Header().Set("Access-Control-Allow-Methods", "GET, POST, OPTIONS")
		sw.Header().Set("Access-Control-Allow-Headers", "Content-Type")
		if r.Method == http.MethodOptions {
			sw.WriteHeader(http.StatusNoContent)
			return
		}
	}

	s.mux.ServeHTTP(sw, r)
}

func (s *server) handleIndex(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	w.Header().Set("Cache-Control", "no-cache")
	w.Write(s.index)
}

func (s *server) handleBanks(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.Header().Set("Cache-Control", "public, max-age=3600")
	w.Write(s.banksJSON)
}

type generateRequest struct {
	Bank        string `json:"bank"`
	Account     string `json:"account"`
	AccountName string `json:"accountName"`
	Amount      int64  `json:"amount"`
	Message     string `json:"message"`
}

type generateResponse struct {
	Payload     string `json:"payload"`
	Bank        Bank   `json:"bank"`
	Account     string `json:"account"`
	AccountName string `json:"accountName,omitempty"`
	Amount      int64  `json:"amount"`
	Message     string `json:"message"`
	QR          struct {
		PNG string `json:"png"`
		SVG string `json:"svg"`
	} `json:"qr"`
}

func (s *server) handleGenerate(w http.ResponseWriter, r *http.Request) {
	r.Body = http.MaxBytesReader(w, r.Body, 64<<10)
	dec := json.NewDecoder(r.Body)
	dec.DisallowUnknownFields()

	var req generateRequest
	if err := dec.Decode(&req); err != nil {
		message := "Dữ liệu JSON không hợp lệ."
		var tooLarge *http.MaxBytesError
		if errors.As(err, &tooLarge) {
			message = "Dữ liệu gửi lên quá lớn."
		}
		writeError(w, http.StatusBadRequest, "", message)
		return
	}
	if err := dec.Decode(&struct{}{}); !errors.Is(err, io.EOF) {
		writeError(w, http.StatusBadRequest, "", "Chỉ chấp nhận một object JSON.")
		return
	}

	qr, ferr := normalizeQRRequest(req.Bank, req.Account, req.AccountName, req.Amount, req.Message)
	if ferr != nil {
		writeError(w, http.StatusBadRequest, ferr.Field, ferr.Message)
		return
	}

	var resp generateResponse
	resp.Payload = qr.payload()
	resp.Bank = qr.Bank
	resp.Account = qr.Account
	resp.AccountName = qr.AccountName
	resp.Amount = qr.Amount
	resp.Message = qr.Message
	resp.QR.PNG, resp.QR.SVG = imageURLs(qr)

	w.Header().Set("Cache-Control", "no-store")
	writeJSON(w, http.StatusOK, resp)
}

func (s *server) handleImage(kind string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		values := r.URL.Query()
		qr, ferr := normalizeQRValues(values)
		if ferr != nil {
			writeError(w, http.StatusBadRequest, ferr.Field, ferr.Message)
			return
		}

		size := 720
		if raw := strings.TrimSpace(values.Get("size")); raw != "" {
			parsed, err := strconv.Atoi(raw)
			if err != nil || parsed < 256 || parsed > 2048 {
				writeError(w, http.StatusBadRequest, "size", "Kích thước ảnh phải từ 256 đến 2048 px.")
				return
			}
			size = parsed
		}

		payload := qr.payload()
		sum := sha256.Sum256([]byte(kind + ":" + strconv.Itoa(size) + ":" + payload))
		etag := `"` + hex.EncodeToString(sum[:12]) + `"`

		w.Header().Set("ETag", etag)
		w.Header().Set("Cache-Control", "public, max-age=86400")
		if r.Header.Get("If-None-Match") == etag {
			w.WriteHeader(http.StatusNotModified)
			return
		}

		var (
			data        []byte
			contentType string
			ext         string
		)
		switch kind {
		case "png":
			var err error
			data, err = renderPNG(payload, size)
			if err != nil {
				writeError(w, http.StatusInternalServerError, "", "Không tạo được ảnh QR.")
				return
			}
			contentType, ext = "image/png", ".png"
		case "svg":
			var err error
			data, err = renderSVG(payload, size)
			if err != nil {
				writeError(w, http.StatusInternalServerError, "", "Không tạo được ảnh QR.")
				return
			}
			contentType, ext = "image/svg+xml; charset=utf-8", ".svg"
		default:
			writeError(w, http.StatusNotFound, "", "Định dạng ảnh không hỗ trợ.")
			return
		}

		w.Header().Set("Content-Type", contentType)
		w.Header().Set("Content-Disposition", `inline; filename="`+imageFilename(qr)+ext+`"`)
		w.Write(data)
	}
}

type decodeResponse struct {
	Payload string `json:"payload"`
	Bank    *Bank  `json:"bank,omitempty"`
	Account string `json:"account"`
	Amount  int64  `json:"amount"`
	Message string `json:"message"`
}

func (s *server) handleDecode(w http.ResponseWriter, r *http.Request) {
	payload := strings.TrimSpace(r.URL.Query().Get("payload"))
	if payload == "" {
		writeError(w, http.StatusBadRequest, "payload", "Thiếu tham số payload.")
		return
	}

	info, err := decodePayload(payload)
	if err != nil {
		writeError(w, http.StatusBadRequest, "payload", "Mã QR không hợp lệ hoặc sai checksum.")
		return
	}
	bank, ok := findBank(info.BankCode)
	if !ok || info.BankNo == "" {
		writeError(w, http.StatusBadRequest, "payload", "Chỉ hỗ trợ mã VietQR chuyển khoản đến tài khoản.")
		return
	}

	writeJSON(w, http.StatusOK, decodeResponse{
		Payload: payload,
		Bank:    bank,
		Account: info.BankNo,
		Amount:  info.Amount,
		Message: info.Message,
	})
}

func (s *server) handleHealth(w http.ResponseWriter, r *http.Request) {
	writeJSON(w, http.StatusOK, map[string]string{"status": "ok", "version": version})
}

func (s *server) handleRobots(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/plain; charset=utf-8")
	w.Write([]byte("User-agent: *\nDisallow: /\n"))
}

func imageURLs(qr *qrRequest) (pngURL, svgURL string) {
	values := url.Values{}
	values.Set("bank", qr.Bank.Code)
	values.Set("account", qr.Account)
	if qr.Amount > 0 {
		values.Set("amount", strconv.FormatInt(qr.Amount, 10))
	}
	if qr.Message != "" {
		values.Set("message", qr.Message)
	}

	svgURL = "/api/qr.svg?" + values.Encode()
	values.Set("size", "1024")
	pngURL = "/api/qr.png?" + values.Encode()
	return pngURL, svgURL
}

func imageFilename(qr *qrRequest) string {
	bankCode := qr.Bank.VietQRCode
	if bankCode == "" {
		bankCode = qr.Bank.Code
	}
	return "vietqr-" + asciiPart(bankCode) + "-" + asciiPart(qr.Account)
}

func asciiPart(s string) string {
	var b strings.Builder
	for _, r := range s {
		if (r >= 'a' && r <= 'z') || (r >= 'A' && r <= 'Z') || (r >= '0' && r <= '9') || r == '-' || r == '_' {
			b.WriteRune(r)
		}
	}
	return b.String()
}

type errorResponse struct {
	Error string `json:"error"`
	Field string `json:"field,omitempty"`
}

func writeError(w http.ResponseWriter, status int, field, message string) {
	w.Header().Set("Cache-Control", "no-store")
	writeJSON(w, status, errorResponse{Error: message, Field: field})
}

func writeJSON(w http.ResponseWriter, status int, value any) {
	data, err := json.Marshal(value)
	if err != nil {
		http.Error(w, `{"error":"internal server error"}`, http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(status)
	w.Write(data)
}

func cacheStatic(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Cache-Control", "public, max-age=3600")
		next.ServeHTTP(w, r)
	})
}

type statusWriter struct {
	http.ResponseWriter
	status int
	wrote  bool
}

func (w *statusWriter) WriteHeader(status int) {
	if w.wrote {
		return
	}
	w.status = status
	w.wrote = true
	w.ResponseWriter.WriteHeader(status)
}

func (w *statusWriter) Write(p []byte) (int, error) {
	if !w.wrote {
		w.status = http.StatusOK
		w.wrote = true
	}
	return w.ResponseWriter.Write(p)
}
