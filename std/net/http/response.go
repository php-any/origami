package http

import (
	"bufio"
	"fmt"
	"io"
	"mime"
	"net"
	"os"
	"path/filepath"
	"strings"

	httpsrc "net/http"

	"github.com/php-any/origami/data"
)

// bufferedWriter 在 handler 生命周期内包装 net/http.ResponseWriter。
// 参考 Echo/Gin：状态码在首次写出响应时提交，避免 header/status 调用顺序耦合。
// 这不是为链式 API 服务，而是保证命令式写法 $res->header(); $res->status(); $res->write(); 语义正确。
type bufferedWriter struct {
	httpsrc.ResponseWriter
	status     int
	statusSet  bool
	headerSent bool
	formatter  *formatHandlerSlot
}

func newBufferedWriter(w httpsrc.ResponseWriter) *bufferedWriter {
	if w == nil {
		return &bufferedWriter{status: httpsrc.StatusOK}
	}
	if bw, ok := w.(*bufferedWriter); ok {
		return bw
	}
	return &bufferedWriter{
		ResponseWriter: w,
		status:         httpsrc.StatusOK,
	}
}

func (b *bufferedWriter) WriteHeader(code int) {
	if b.headerSent {
		return
	}
	b.status = code
	b.headerSent = true
	b.ResponseWriter.WriteHeader(code)
}

func (b *bufferedWriter) Write(p []byte) (int, error) {
	b.sendHeader()
	return b.ResponseWriter.Write(p)
}

// Hijack 委托底层 ResponseWriter，供 WebSocket 等协议升级使用。
func (b *bufferedWriter) Hijack() (net.Conn, *bufio.ReadWriter, error) {
	if h, ok := b.ResponseWriter.(httpsrc.Hijacker); ok {
		return h.Hijack()
	}
	return nil, nil, fmt.Errorf("response does not implement http.Hijacker")
}

// Flush 委托底层 ResponseWriter（若实现 http.Flusher）。
func (b *bufferedWriter) Flush() {
	if f, ok := b.ResponseWriter.(httpsrc.Flusher); ok {
		f.Flush()
	}
}

func (b *bufferedWriter) SetStatus(code int) {
	if b.headerSent {
		return
	}
	b.status = code
	b.statusSet = true
}

func (b *bufferedWriter) SetHeader(key, value string) {
	b.Header().Set(key, value)
}

func (b *bufferedWriter) sendHeader() {
	if !b.headerSent {
		b.WriteHeader(b.status)
	}
}

// commitPending 在 handler 返回时补发「只设置了 status、未写 body」的响应（如 204）。
func (b *bufferedWriter) commitPending() {
	if !b.headerSent && b.statusSet {
		b.WriteHeader(b.status)
	}
}

// Redirect 发送重定向（Laravel/Symfony 风格的完整响应操作）。
func (b *bufferedWriter) Redirect(url string, code int) {
	b.SetHeader("Location", url)
	if !b.headerSent {
		b.status = code
		b.statusSet = true
	}
	b.sendHeader()
	_, _ = b.ResponseWriter.Write(nil)
}

// NoContent 发送无 body 响应。
func (b *bufferedWriter) NoContent(code int) {
	if !b.headerSent {
		b.status = code
		b.statusSet = true
	}
	b.sendHeader()
}

// WriteHTML 写出 HTML body。
func (b *bufferedWriter) WriteHTML(body []byte) error {
	b.SetHeader("Content-Type", "text/html; charset=utf-8")
	_, err := b.Write(body)
	return err
}

// WriteJSON 写出 JSON body（设置 Content-Type 并写入）。
func (b *bufferedWriter) WriteJSON(body []byte) error {
	b.SetHeader("Content-Type", "application/json; charset=utf-8")
	_, err := b.Write(body)
	return err
}

// SendFile 发送文件下载响应。downloadName 为空时使用路径中的文件名。
func (b *bufferedWriter) SendFile(path string, downloadName string) error {
	resolved, err := resolveReadableFile(path)
	if err != nil {
		return err
	}

	name := downloadName
	if name == "" {
		name = filepath.Base(resolved)
	}

	ctype := mime.TypeByExtension(filepath.Ext(resolved))
	if ctype == "" {
		ctype = "application/octet-stream"
	}
	b.SetHeader("Content-Type", ctype)
	b.SetHeader("Content-Disposition", attachmentDisposition(name))

	f, err := os.Open(resolved)
	if err != nil {
		return err
	}
	defer f.Close()

	b.sendHeader()
	_, err = io.Copy(b.ResponseWriter, f)
	return err
}

func resolveReadableFile(path string) (string, error) {
	if path == "" {
		return "", fmt.Errorf("文件路径为空")
	}
	if !filepath.IsAbs(path) {
		cwd, err := os.Getwd()
		if err != nil {
			return "", err
		}
		path = filepath.Join(cwd, path)
	}
	path = filepath.Clean(path)
	info, err := os.Stat(path)
	if err != nil {
		return "", err
	}
	if info.IsDir() {
		return "", fmt.Errorf("路径是目录: %s", path)
	}
	return path, nil
}

func attachmentDisposition(filename string) string {
	safe := strings.ReplaceAll(filename, `"`, `_`)
	return `attachment; filename="` + safe + `"`
}

// SetCookie 通过 Set-Cookie 头发送 cookie。
func (b *bufferedWriter) SetCookie(c *httpsrc.Cookie) {
	httpsrc.SetCookie(b.ResponseWriter, c)
}

func beginResponse(w httpsrc.ResponseWriter, r *httpsrc.Request) (*bufferedWriter, data.ClassStmt) {
	bw := newBufferedWriter(w)
	if r != nil {
		bw.formatter = requestFormatterFor(r)
	}
	return bw, &ResponseWriterClass{w: bw}
}

func responseSelf(w *bufferedWriter, ctx data.Context) (data.GetValue, data.Control) {
	return data.NewProxyValue(&ResponseWriterClass{w: w}, ctx.CreateBaseContext()), nil
}
