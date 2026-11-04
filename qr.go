package main

import (
	"errors"
	"fmt"
	"strings"

	qrcode "github.com/skip2/go-qrcode"
	"github.com/sunary/vietqr"
)

const (
	maxAmount       = int64(999_999_999_999)
	maxMessageRunes = 50
	maxAccountName  = 50
)

type qrRequest struct {
	Bank        Bank
	Account     string
	AccountName string
	Amount      int64
	Message     string
}

func (q qrRequest) payload() string {
	return vietqr.Encode(vietqr.TransferInfo{
		BankCode: q.Bank.Code,
		BankNo:   q.Account,
		Amount:   q.Amount,
		Message:  q.Message,
	})
}

func renderPNG(payload string, size int) ([]byte, error) {
	code, err := qrcode.New(payload, qrcode.High)
	if err != nil {
		return nil, err
	}
	return code.PNG(size)
}

func renderSVG(payload string, size int) ([]byte, error) {
	code, err := qrcode.New(payload, qrcode.High)
	if err != nil {
		return nil, err
	}

	bitmap := code.Bitmap()
	modules := len(bitmap)
	if modules == 0 {
		return nil, errors.New("mã QR rỗng")
	}

	var path strings.Builder
	for y, row := range bitmap {
		for x := 0; x < len(row); {
			if !row[x] {
				x++
				continue
			}
			start := x
			for x < len(row) && row[x] {
				x++
			}
			fmt.Fprintf(&path, "M%d %dh%dv1h-%dz", start, y, x-start, x-start)
		}
	}

	var out strings.Builder
	fmt.Fprintf(&out, `<svg xmlns="http://www.w3.org/2000/svg" width="%d" height="%d" viewBox="0 0 %d %d" shape-rendering="crispEdges" role="img" aria-label="VietQR">`, size, size, modules, modules)
	fmt.Fprintf(&out, `<rect width="%d" height="%d" fill="#fff"/>`, modules, modules)
	fmt.Fprintf(&out, `<path fill="#000" d="%s"/></svg>`, path.String())
	return []byte(out.String()), nil
}

// decodePayload wraps the upstream decoder, which does not guard against
// malformed input and can panic on truncated payloads.
func decodePayload(payload string) (info *vietqr.TransferInfo, err error) {
	defer func() {
		if recover() != nil {
			info, err = nil, errors.New("payload QR không hợp lệ")
		}
	}()

	payload = strings.TrimSpace(payload)
	if len(payload) < 24 || len(payload) > 1024 {
		return nil, errors.New("payload QR không hợp lệ")
	}
	return vietqr.Decode(payload)
}
