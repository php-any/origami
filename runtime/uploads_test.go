package runtime_test

import (
	"bytes"
	"context"
	"github.com/php-any/origami/data"
	"github.com/php-any/origami/node"
	"github.com/php-any/origami/parser"
	"github.com/php-any/origami/runtime"
	"github.com/php-any/origami/std/php"
	"mime/multipart"
	"net/http/httptest"
	"os"
	"testing"
	"time"
)

func TestRequestMultipartUploadsAndCleanup(t *testing.T) {
	base := runtime.NewVM(parser.NewParser())
	php.Load(base)
	vm := runtime.NewRequestVM(base).(*runtime.RequestVM)
	var body bytes.Buffer
	form := multipart.NewWriter(&body)
	for _, text := range []string{"first", "last"} {
		_ = form.WriteField("field", text)
	}
	for _, filename := range []string{"one.txt", "two.txt"} {
		part, err := form.CreateFormFile("uploads[]", filename)
		if err != nil {
			t.Fatal(err)
		}
		_, _ = part.Write([]byte(filename))
	}
	_ = form.Close()
	request := httptest.NewRequest("POST", "http://example.test/upload", &body)
	request.Header.Set("Content-Type", form.FormDataContentType())
	client, cancel := context.WithCancel(request.Context())
	defer cancel()
	vm.BindHTTP(request.WithContext(client), httptest.NewRecorder())
	ctx := vm.CreateContext(nil)
	files, ctl := node.NewFilesVariable(nil).GetValue(ctx)
	if ctl != nil {
		t.Fatal(ctl.AsString())
	}
	upload, exists := files.(*data.ArrayValue).LookupZValByStringKey("uploads")
	if !exists {
		t.Fatal("files not initialized before POST")
	}
	info := upload.ReadValue().(*data.ArrayValue)
	names, _ := info.LookupZValByStringKey("name")
	paths, _ := info.LookupZValByStringKey("tmp_name")
	if names.ReadValue().(*data.ArrayValue).Len() != 2 || paths.ReadValue().(*data.ArrayValue).Len() != 2 {
		t.Fatal("incorrect PHP upload column layout")
	}
	temporary := []string{}
	for _, slot := range paths.ReadValue().(*data.ArrayValue).Range() {
		path := slot.ReadValue().AsString()
		temporary = append(temporary, path)
		if !vm.IsUploadedFile(path) {
			t.Fatal("upload not registered")
		}
		content, err := os.ReadFile(path)
		if err != nil || len(content) == 0 {
			t.Fatalf("temporary upload: %q %v", content, err)
		}
	}
	post, ctl := node.NewPostVariable(nil).GetValue(ctx)
	if ctl != nil {
		t.Fatal(ctl.AsString())
	}
	field, _ := post.(*data.ArrayValue).LookupZValByStringKey("field")
	if field.ReadValue().AsString() != "last" {
		t.Fatal("multipart scalar order lost")
	}
	other := runtime.NewRequestVM(base).(*runtime.RequestVM)
	if other.IsUploadedFile(temporary[0]) {
		t.Fatal("upload identity leaked to other request")
	}
	vm.StorePHPIni("ignore_user_abort", "1")
	cancel()
	time.Sleep(10 * time.Millisecond)
	if !vm.IsUploadedFile(temporary[0]) {
		t.Fatal("ignored client abort removed upload")
	}
	vm.RunShutdownCallbacks()
	for _, path := range temporary {
		if _, err := os.Stat(path); !os.IsNotExist(err) {
			t.Fatalf("upload retained after request shutdown: %v", err)
		}
	}
	vm.RunShutdownCallbacks()
}
