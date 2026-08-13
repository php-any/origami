package http

import (
	stdjson "encoding/json"
	httpsrc "net/http"
	"net/http/httptest"
	"testing"

	"github.com/php-any/origami/data"
	"github.com/php-any/origami/std/serializer/json"
)

func TestDefaultFormattedPayload(t *testing.T) {
	payload := defaultFormattedPayload(200, "success", data.NewStringValue("ok"))
	obj, ok := payload.(*data.ObjectValue)
	if !ok {
		t.Fatalf("payload type = %T, want *data.ObjectValue", payload)
	}

	codeVal, _ := obj.GetProperty("code")
	if got := codeVal.AsString(); got != "200" {
		t.Fatalf("code = %q, want 200", got)
	}
	msgVal, _ := obj.GetProperty("message")
	if got := msgVal.AsString(); got != "success" {
		t.Fatalf("message = %q, want success", got)
	}
	dataVal, _ := obj.GetProperty("data")
	if got := dataVal.AsString(); got != "ok" {
		t.Fatalf("data = %q, want ok", got)
	}
	tsVal, _ := obj.GetProperty("timestamp")
	if tsVal == nil {
		t.Fatal("timestamp missing")
	}
}

func TestWriteFormattedResponse_DefaultEnvelope(t *testing.T) {
	rec := httptest.NewRecorder()
	bw := newBufferedWriter(rec)

	if err := writeFormattedResponse(bw, nil, 201, "created", data.NewStringValue("item")); err != nil {
		t.Fatalf("writeFormattedResponse: %v", err)
	}

	if rec.Code != 201 {
		t.Fatalf("status = %d, want 201", rec.Code)
	}
	if got := rec.Header().Get("Content-Type"); got != "application/json; charset=utf-8" {
		t.Fatalf("Content-Type = %q", got)
	}

	var body map[string]any
	if err := stdjson.Unmarshal(rec.Body.Bytes(), &body); err != nil {
		t.Fatalf("json.Unmarshal: %v", err)
	}
	if body["code"] != float64(201) {
		t.Fatalf("body.code = %v, want 201", body["code"])
	}
	if body["message"] != "created" {
		t.Fatalf("body.message = %v, want created", body["message"])
	}
	if body["data"] != "item" {
		t.Fatalf("body.data = %v, want item", body["data"])
	}
	if _, ok := body["timestamp"]; !ok {
		t.Fatal("body.timestamp missing")
	}
}

func TestWriteFormattedResponse_ErrorStatus(t *testing.T) {
	rec := httptest.NewRecorder()
	bw := newBufferedWriter(rec)

	if err := writeFormattedResponse(bw, nil, 404, "not found", data.NewNullValue()); err != nil {
		t.Fatalf("writeFormattedResponse: %v", err)
	}
	if rec.Code != 404 {
		t.Fatalf("status = %d, want 404", rec.Code)
	}
}

func TestBufferedWriter_JsonBypassesFormatter(t *testing.T) {
	rec := httptest.NewRecorder()
	bw := newBufferedWriter(rec)
	bw.formatter = &formatHandlerSlot{}

	raw := data.NewObjectValue()
	raw.SetProperty("plain", data.NewStringValue("yes"))
	bytes, err := raw.Marshal(json.NewJsonSerializer())
	if err != nil {
		t.Fatalf("Marshal: %v", err)
	}
	if err := bw.WriteJSON(bytes); err != nil {
		t.Fatalf("WriteJSON: %v", err)
	}

	var body map[string]any
	if err := stdjson.Unmarshal(rec.Body.Bytes(), &body); err != nil {
		t.Fatalf("json.Unmarshal: %v", err)
	}
	if _, ok := body["code"]; ok {
		t.Fatalf("json() should not add envelope, got %v", body)
	}
	if body["plain"] != "yes" {
		t.Fatalf("body.plain = %v, want yes", body["plain"])
	}
}

func TestRequestFormatterAttach(t *testing.T) {
	req := httptest.NewRequest("GET", "/", nil)
	slot := &formatHandlerSlot{}
	attachRequestFormatter(req, slot)
	defer detachRequestAttrs(req)

	if got := requestFormatterFor(req); got != slot {
		t.Fatal("requestFormatterFor did not return attached slot")
	}

	detachRequestAttrs(req)
	if got := requestFormatterFor(req); got != nil {
		t.Fatal("formatter slot not cleaned up")
	}
}

func TestWithResponseFormatter_AttachesToRequest(t *testing.T) {
	server := &ServerClass{
		formatHandler: &formatHandlerSlot{},
	}
	var captured *formatHandlerSlot
	h := withResponseFormatter(server, httpsrc.HandlerFunc(func(w httpsrc.ResponseWriter, r *httpsrc.Request) {
		captured = requestFormatterFor(r)
	}))

	rec := httptest.NewRecorder()
	req := httptest.NewRequest("GET", "/api", nil)
	attachRequestAttrs(req)
	defer detachRequestAttrs(req)
	h.ServeHTTP(rec, req)

	if captured != server.formatHandler {
		t.Fatal("formatter not attached during request")
	}
}
