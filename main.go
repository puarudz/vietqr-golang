package main

import (
	"context"
	"embed"
	"errors"
	"flag"
	"io/fs"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"strings"
	"syscall"
	"time"
)

//go:embed web
var webFiles embed.FS

var version = "dev"

func main() {
	addr := flag.String("addr", defaultAddr(), "địa chỉ HTTP, ví dụ :8080")
	flag.Parse()

	static, err := fs.Sub(webFiles, "web")
	if err != nil {
		slog.Error("không đọc được tài nguyên nhúng", "error", err)
		os.Exit(1)
	}

	handler, err := newServer(static)
	if err != nil {
		slog.Error("khởi tạo server thất bại", "error", err)
		os.Exit(1)
	}

	srv := &http.Server{
		Addr:              *addr,
		Handler:           handler,
		ReadHeaderTimeout: 5 * time.Second,
		ReadTimeout:       15 * time.Second,
		WriteTimeout:      30 * time.Second,
		IdleTimeout:       60 * time.Second,
		MaxHeaderBytes:    1 << 20,
	}

	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	go func() {
		<-ctx.Done()
		shutdownCtx, cancel := context.WithTimeout(context.Background(), 8*time.Second)
		defer cancel()
		if err := srv.Shutdown(shutdownCtx); err != nil {
			slog.Error("tắt server không sạch", "error", err)
		}
	}()

	slog.Info("VietQR Core đang chạy", "addr", *addr, "version", version)
	if err := srv.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
		slog.Error("server dừng bất thường", "error", err)
		os.Exit(1)
	}
}

func defaultAddr() string {
	if addr := os.Getenv("ADDR"); addr != "" {
		return addr
	}
	if port := os.Getenv("PORT"); port != "" {
		if !strings.HasPrefix(port, ":") {
			port = ":" + port
		}
		return port
	}
	return ":8080"
}
