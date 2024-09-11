package client

import (
	"io"
	"net/http"
)

func ForwardBackend(request string) (interface{}, error) {
	resp, err := http.Get(request)
	if err != nil {
		return "Backend serveriga kirish imkonsiz", err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "Backend javobini o'qib bo'lmadi", err
	}

	return body, nil
}
