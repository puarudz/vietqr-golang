package main

import (
	"fmt"
	"net/url"
	"strconv"
	"strings"
	"unicode"
	"unicode/utf8"
)

type fieldError struct {
	Field   string
	Message string
}

func (e *fieldError) Error() string { return e.Message }

func normalizeQRRequest(bankInput, account, accountName string, amount int64, message string) (*qrRequest, *fieldError) {
	bankInput = strings.TrimSpace(bankInput)
	if bankInput == "" {
		return nil, &fieldError{"bank", "Vui lòng chọn ngân hàng thụ hưởng."}
	}
	bank, ok := findBank(bankInput)
	if !ok {
		return nil, &fieldError{"bank", "Ngân hàng không hợp lệ hoặc chưa được hỗ trợ."}
	}

	account = stripSeparators(account)
	if account == "" {
		return nil, &fieldError{"account", "Vui lòng nhập số tài khoản."}
	}
	if len(account) < 4 || len(account) > 19 {
		return nil, &fieldError{"account", "Số tài khoản phải từ 4 đến 19 ký tự."}
	}
	for _, r := range account {
		if !isASCIIAlnum(r) {
			return nil, &fieldError{"account", "Số tài khoản chỉ gồm chữ và số, không dấu cách."}
		}
	}

	if amount < 0 {
		return nil, &fieldError{"amount", "Số tiền không được âm."}
	}
	if amount > maxAmount {
		return nil, &fieldError{"amount", "Số tiền vượt quá giới hạn 999.999.999.999 đ."}
	}

	accountName = collapseSpaces(stripControl(accountName))
	if utf8.RuneCountInString(accountName) > maxAccountName {
		return nil, &fieldError{"accountName", fmt.Sprintf("Chủ tài khoản tối đa %d ký tự.", maxAccountName)}
	}

	message = collapseSpaces(stripControl(message))
	if n := utf8.RuneCountInString(message); n > maxMessageRunes {
		return nil, &fieldError{"message", fmt.Sprintf("Nội dung chuyển khoản tối đa %d ký tự.", maxMessageRunes)}
	}

	return &qrRequest{
		Bank:        *bank,
		Account:     account,
		AccountName: accountName,
		Amount:      amount,
		Message:     message,
	}, nil
}

func normalizeQRValues(values url.Values) (*qrRequest, *fieldError) {
	var amount int64
	if raw := strings.TrimSpace(values.Get("amount")); raw != "" {
		parsed, err := strconv.ParseInt(raw, 10, 64)
		if err != nil {
			return nil, &fieldError{"amount", "Số tiền không hợp lệ."}
		}
		amount = parsed
	}
	return normalizeQRRequest(
		values.Get("bank"),
		values.Get("account"),
		values.Get("accountName"),
		amount,
		values.Get("message"),
	)
}

func stripSeparators(s string) string {
	return strings.Map(func(r rune) rune {
		switch r {
		case ' ', '.', '-', '_':
			return -1
		}
		return r
	}, strings.TrimSpace(s))
}

func stripControl(s string) string {
	return strings.Map(func(r rune) rune {
		if unicode.IsControl(r) {
			return ' '
		}
		return r
	}, s)
}

func collapseSpaces(s string) string {
	return strings.Join(strings.Fields(s), " ")
}

func isASCIIAlnum(r rune) bool {
	return (r >= '0' && r <= '9') || (r >= 'a' && r <= 'z') || (r >= 'A' && r <= 'Z')
}
