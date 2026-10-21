// Package pharfs reads native PHP archives without extracting or executing them.
package pharfs

import (
	"bytes"
	"compress/bzip2"
	"compress/flate"
	"crypto/md5"
	"crypto/sha1"
	"crypto/sha256"
	"crypto/sha512"
	"encoding/binary"
	"fmt"
	"hash"
	"hash/crc32"
	"io"
	"os"
	"path"
	"path/filepath"
	"strings"
	"sync"
	"time"

	"github.com/php-any/origami/lexer"
	"github.com/php-any/origami/token"
)

type entry struct {
	name                    string
	size, stamp, flags, crc uint32
	body                    []byte
}

func (e *entry) Name() string { return path.Base(e.name) }
func (e *entry) Size() int64  { return int64(e.size) }
func (e *entry) Mode() os.FileMode {
	mode := os.FileMode(e.flags & 0777)
	if e.IsDir() {
		mode |= os.ModeDir
	}
	return mode
}
func (e *entry) ModTime() time.Time { return time.Unix(int64(e.stamp), 0) }
func (e *entry) IsDir() bool        { return strings.HasSuffix(e.name, "/") }
func (e *entry) Sys() any           { return nil }

type Archive struct {
	Filename, Alias string
	entries         map[string]*entry
}

// cursor keeps malformed lengths from slicing outside the manifest.
type cursor struct {
	b   []byte
	err error
}

func (c *cursor) take(n uint32) []byte {
	if c.err != nil || uint64(n) > uint64(len(c.b)) {
		c.err = io.ErrUnexpectedEOF
		return nil
	}
	v := c.b[:int(n)]
	c.b = c.b[int(n):]
	return v
}
func (c *cursor) u32() uint32 {
	b := c.take(4)
	if len(b) != 4 {
		return 0
	}
	return binary.LittleEndian.Uint32(b)
}
func Open(filename string) (*Archive, error) {
	b, err := os.ReadFile(filename)
	if err != nil {
		return nil, err
	}
	offset := 0
	tokens := lexer.NewLexer().TokenizeTemplate(string(b))
	for i, tok := range tokens {
		if tok.Type() == token.HALT_COMPILER && i+3 < len(tokens) &&
			tokens[i+1].Type() == token.LPAREN && tokens[i+2].Type() == token.RPAREN && tokens[i+3].Type() == token.SEMICOLON {
			offset = tokens[i+3].End()
			break
		}
	}
	if offset == 0 {
		return nil, fmt.Errorf("phar: missing compiler halt")
	}
	// PHP permits a closing tag and one line ending between stub and manifest.
	tail := b[offset:]
	trimmed := bytes.TrimLeft(tail, " \t")
	if bytes.HasPrefix(trimmed, []byte("?>")) {
		tail = trimmed[2:]
		if bytes.HasPrefix(tail, []byte("\r\n")) {
			tail = tail[2:]
		} else if bytes.HasPrefix(tail, []byte("\n")) {
			tail = tail[1:]
		}
	}
	c := &cursor{b: tail}
	length := c.u32()
	if length > 1024*1024 {
		return nil, fmt.Errorf("phar: manifest too large")
	}
	m := &cursor{b: c.take(length)}
	count := m.u32()
	m.take(2)
	flags := m.u32()
	alias := string(m.take(m.u32()))
	m.take(m.u32())
	if c.err != nil || m.err != nil || count > length/24 {
		return nil, fmt.Errorf("phar: invalid manifest")
	}
	if flags&0x10000 != 0 {
		if err := verifySignature(b); err != nil {
			return nil, err
		}
	} else {
		return nil, fmt.Errorf("phar: unsigned archive")
	}
	abs, err := filepath.Abs(filename)
	if err != nil {
		return nil, err
	}
	a := &Archive{Filename: filepath.ToSlash(abs), Alias: alias, entries: make(map[string]*entry)}
	for i := uint32(0); i < count; i++ {
		name := string(m.take(m.u32()))
		e := &entry{name: name, size: m.u32(), stamp: m.u32()}
		compressed := m.u32()
		e.crc = m.u32()
		e.flags = m.u32()
		m.take(m.u32())
		if m.err != nil || name == "" || strings.HasPrefix(name, "/") || strings.Contains(name, "\\") || path.Clean(name) != strings.TrimSuffix(name, "/") {
			return nil, fmt.Errorf("phar: invalid entry %q", name)
		}
		if _, exists := a.entries[name]; exists {
			return nil, fmt.Errorf("phar: duplicate entry %q", name)
		}
		e.body = c.take(compressed)
		if c.err != nil {
			return nil, c.err
		}
		a.entries[name] = e
	}
	if m.err != nil || len(m.b) != 0 {
		return nil, fmt.Errorf("phar: invalid manifest size")
	}
	return a, nil
}
func verifySignature(b []byte) error {
	if len(b) < 8 || string(b[len(b)-4:]) != "GBMB" {
		return fmt.Errorf("phar: missing signature")
	}
	var h hash.Hash
	switch binary.LittleEndian.Uint32(b[len(b)-8:]) {
	case 1:
		h = md5.New()
	case 2:
		h = sha1.New()
	case 3:
		h = sha256.New()
	case 4:
		h = sha512.New()
	default:
		return fmt.Errorf("phar: unsupported signature algorithm")
	}
	end := len(b) - 8 - h.Size()
	if end < 0 {
		return io.ErrUnexpectedEOF
	}
	h.Write(b[:end])
	if !bytes.Equal(h.Sum(nil), b[end:len(b)-8]) {
		return fmt.Errorf("phar: invalid signature")
	}
	return nil
}
func (a *Archive) read(name string) ([]byte, error) {
	e, ok := a.entries[name]
	if !ok || e.IsDir() {
		return nil, os.ErrNotExist
	}
	var r io.Reader = bytes.NewReader(e.body)
	switch e.flags & 0x3000 {
	case 0:
	case 0x1000:
		z := flate.NewReader(r)
		defer z.Close()
		r = z
	case 0x2000:
		r = bzip2.NewReader(r)
	default:
		return nil, fmt.Errorf("phar: invalid compression")
	}
	b, err := io.ReadAll(io.LimitReader(r, int64(e.size)+1))
	if err != nil {
		return nil, err
	}
	if len(b) != int(e.size) || crc32.ChecksumIEEE(b) != e.crc {
		return nil, fmt.Errorf("phar: corrupt entry %q", name)
	}
	return b, nil
}

// FS belongs to one VM/request; aliases never escape that execution scope.
type FS struct {
	mu       sync.Mutex
	archives map[string]*Archive
}

func (fs *FS) Map(filename, alias string) error {
	a, err := Open(filename)
	if err != nil {
		return err
	}
	fs.mu.Lock()
	defer fs.mu.Unlock()
	if fs.archives == nil {
		fs.archives = make(map[string]*Archive)
	}
	if alias == "" {
		alias = a.Alias
	}
	if previous := fs.archives[alias]; alias != "" && previous != nil && previous.Filename != a.Filename {
		return fmt.Errorf("phar: alias %q already mapped", alias)
	}
	fs.archives[a.Filename] = a
	if alias != "" {
		fs.archives[alias] = a
	}
	return nil
}
func (fs *FS) resolve(uri string) (*Archive, string, error) {
	raw := strings.ReplaceAll(strings.TrimPrefix(uri, "phar://"), "\\", "/")
	fs.mu.Lock()
	var found *Archive
	var name string
	best := 0
	for key, a := range fs.archives {
		if len(key) > best && (raw == key || strings.HasPrefix(raw, key+"/")) {
			found, name, best = a, strings.TrimPrefix(raw[len(key):], "/"), len(key)
		}
	}
	fs.mu.Unlock()
	if found == nil {
		// Absolute archive paths can be used without mapPhar().
		for at := strings.Index(raw, "/"); at >= 0 && at < len(raw); {
			if at > 0 {
				if info, err := os.Stat(raw[:at]); err == nil && !info.IsDir() {
					if err := fs.Map(raw[:at], ""); err != nil {
						return nil, "", err
					}
					return fs.resolve(uri)
				}
			}
			next := strings.Index(raw[at+1:], "/")
			if next < 0 {
				break
			}
			at += next + 1
		}
		return nil, "", os.ErrNotExist
	}
	name = strings.TrimPrefix(path.Clean("/"+name), "/")
	return found, name, nil
}
func (fs *FS) Canonical(uri string) (string, error) {
	a, name, err := fs.resolve(uri)
	if err != nil {
		return "", err
	}
	return "phar://" + a.Filename + "/" + name, nil
}
func (fs *FS) ReadFile(uri string) ([]byte, error) {
	a, name, err := fs.resolve(uri)
	if err != nil {
		return nil, err
	}
	return a.read(name)
}
func (fs *FS) Stat(uri string) (os.FileInfo, error) {
	a, name, err := fs.resolve(uri)
	if err != nil {
		return nil, err
	}
	if e := a.entries[name]; e != nil {
		return e, nil
	}
	for key := range a.entries {
		if name == "" || strings.HasPrefix(key, name+"/") {
			return &entry{name: name + "/", flags: 0555}, nil
		}
	}
	return nil, os.ErrNotExist
}
