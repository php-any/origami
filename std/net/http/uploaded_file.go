package http

import (
	"errors"
	"io"
	"mime/multipart"
	"os"
	"path/filepath"
	"strings"

	httpsrc "net/http"

	"github.com/php-any/origami/data"
)

const uploadedFileTypeFQN = "Net\\Http\\UploadedFile"

func requestFromProxy(reqProxy data.Value) *httpsrc.Request {
	reqCV, ok := reqProxy.(*data.ProxyValue)
	if !ok {
		return nil
	}
	if rc, ok := reqCV.Class.(*RequestClass); ok {
		return rc.source
	}
	if getter, ok := reqCV.Class.(interface{ GetSource() any }); ok {
		if r, ok := getter.GetSource().(*httpsrc.Request); ok {
			return r
		}
	}
	return nil
}

func ensureMultipartParsed(r *httpsrc.Request) error {
	if r == nil {
		return errors.New("request is nil")
	}
	if r.MultipartForm != nil {
		return nil
	}
	return r.ParseMultipartForm(32 << 20)
}

func newUploadedFileValue(ctx data.Context, header *multipart.FileHeader) data.Value {
	if header == nil {
		return data.NewNullValue()
	}
	return data.NewProxyValue(NewUploadedFileClassFrom(header), ctx.CreateBaseContext())
}

func uploadedFileFromForm(r *httpsrc.Request, field string) (*multipart.FileHeader, error) {
	if err := ensureMultipartParsed(r); err != nil {
		return nil, err
	}
	_, header, err := r.FormFile(field)
	if err != nil {
		return nil, err
	}
	return header, nil
}

func isUploadedFileTypeFQN(typeFQN string) bool {
	return strings.TrimPrefix(typeFQN, "\\") == uploadedFileTypeFQN
}

func typeFQNFromTypes(ty data.Types) (string, bool) {
	if ty == nil {
		return "", false
	}
	nullable := false
	if nt, ok := ty.(data.NullableType); ok {
		nullable = true
		ty = nt.BaseType
	}
	return strings.TrimPrefix(ty.String(), "\\"), nullable
}

func readUploadedFileContent(header *multipart.FileHeader) ([]byte, error) {
	if header == nil {
		return nil, errors.New("uploaded file is missing")
	}
	file, err := header.Open()
	if err != nil {
		return nil, err
	}
	defer file.Close()
	return io.ReadAll(file)
}

func storeUploadedFile(header *multipart.FileHeader, directory, name string) (string, error) {
	if header == nil {
		return "", errors.New("uploaded file is missing")
	}
	if strings.TrimSpace(directory) == "" {
		return "", errors.New("store directory is required")
	}
	if name == "" {
		name = filepath.Base(header.Filename)
	}
	target := filepath.Join(directory, name)
	file, err := header.Open()
	if err != nil {
		return "", err
	}
	defer file.Close()
	out, err := createFileForWrite(target)
	if err != nil {
		return "", err
	}
	defer out.Close()
	if _, err := io.Copy(out, file); err != nil {
		return "", err
	}
	return target, nil
}

func createFileForWrite(path string) (*os.File, error) {
	if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
		return nil, err
	}
	return os.Create(path)
}
