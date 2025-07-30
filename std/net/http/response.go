package http

import (
	"github.com/php-any/origami/data"
	"net/http"
)

type Response struct {
	w http.ResponseWriter
	r *http.Request
}

func (req *Response) Write(text string) (data.GetValue, data.Control) {
	n, err := req.w.Write([]byte(text))
	if err != nil {
		return nil, data.NewErrorThrow(nil, err)
	}
	return data.NewIntValue(n), nil
}
