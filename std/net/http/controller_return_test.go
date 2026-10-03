package http

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/php-any/origami/data"
)

func TestWriteHandlerReturnValueArray(t *testing.T) {
	rr := httptest.NewRecorder()
	bw, responseClass := beginResponse(rr, nil)
	resProxy := data.NewProxyValue(responseClass, nil)

	payload := data.NewArrayValue(nil).(*data.ArrayValue)
	payload.SetStringKey("hello", data.NewStringValue("Origami"))
	body := data.NewArrayValue(nil).(*data.ArrayValue)
	body.SetStringKey("code", data.NewIntValue(200))
	body.SetStringKey("message", data.NewStringValue("success"))
	body.SetStringKey("data", payload)

	if acl := writeHandlerReturnValue(resProxy, body); acl != nil {
		t.Fatalf("write return: %v", acl)
	}
	bw.commitPending()

	if rr.Code != http.StatusOK {
		t.Fatalf("status = %d", rr.Code)
	}
	var decoded map[string]any
	if err := json.Unmarshal(rr.Body.Bytes(), &decoded); err != nil {
		t.Fatalf("json: %v", err)
	}
	if decoded["message"] != "success" {
		t.Fatalf("body = %#v", decoded)
	}
}

func TestWriteHandlerReturnValueSkipsWhenAlreadyWritten(t *testing.T) {
	rr := httptest.NewRecorder()
	bw, responseClass := beginResponse(rr, nil)
	resProxy := data.NewProxyValue(responseClass, nil)

	_ = bw.WriteJSON([]byte(`{"ok":true}`))
	ret := data.NewArrayValue(nil).(*data.ArrayValue)
	ret.SetStringKey("ignored", data.NewBoolValue(true))
	if acl := writeHandlerReturnValue(resProxy, ret); acl != nil {
		t.Fatalf("unexpected acl: %v", acl)
	}
	if rr.Body.String() != `{"ok":true}` {
		t.Fatalf("body mutated: %s", rr.Body.String())
	}
}
