package pharfs

import (
	"bytes"
	"compress/flate"
	"crypto/sha256"
	"encoding/binary"
	"hash/crc32"
	"os"
	"path/filepath"
	"testing"
)

func archiveFixture(t *testing.T, compressed bool) string {
	t.Helper()
	content := []byte("<?php return 'archive';")
	body := content
	flags := uint32(0644)
	if compressed {
		var buf bytes.Buffer
		z, _ := flate.NewWriter(&buf, flate.DefaultCompression)
		z.Write(content)
		z.Close()
		body = buf.Bytes()
		flags |= 0x1000
	}
	var m bytes.Buffer
	put := func(n uint32) { binary.Write(&m, binary.LittleEndian, n) }
	put(1)
	m.Write([]byte{0x11, 0x00})
	put(0x10000)
	put(8)
	m.WriteString("app.phar")
	put(0)
	put(13)
	m.WriteString("src/hello.php")
	put(uint32(len(content)))
	put(123)
	put(uint32(len(body)))
	put(crc32.ChecksumIEEE(content))
	put(flags)
	put(0)
	var out bytes.Buffer
	out.WriteString("<?php __HALT_COMPILER(); ?>\r\n")
	binary.Write(&out, binary.LittleEndian, uint32(m.Len()))
	out.Write(m.Bytes())
	out.Write(body)
	sum := sha256.Sum256(out.Bytes())
	out.Write(sum[:])
	binary.Write(&out, binary.LittleEndian, uint32(3))
	out.WriteString("GBMB")
	filename := filepath.Join(t.TempDir(), "test.phar")
	if err := os.WriteFile(filename, out.Bytes(), 0600); err != nil {
		t.Fatal(err)
	}
	return filename
}

func TestArchiveReadAndIsolation(t *testing.T) {
	for _, compressed := range []bool{false, true} {
		filename := archiveFixture(t, compressed)
		var fs, other FS
		if err := fs.Map(filename, ""); err != nil {
			t.Fatal(err)
		}
		b, err := fs.ReadFile("phar://app.phar/src/hello.php")
		if err != nil || string(b) != "<?php return 'archive';" {
			t.Fatalf("%q %v", b, err)
		}
		if _, err := other.ReadFile("phar://app.phar/src/hello.php"); err == nil {
			t.Fatal("alias leaked")
		}
		if info, err := fs.Stat("phar://app.phar/src"); err != nil || !info.IsDir() {
			t.Fatalf("directory: %v %v", info, err)
		}
	}
}
func TestRejectCorruptArchive(t *testing.T) {
	filename := archiveFixture(t, false)
	b, _ := os.ReadFile(filename)
	b[len(b)-45] ^= 1
	os.WriteFile(filename, b, 0600)
	if _, err := Open(filename); err == nil {
		t.Fatal("accepted corrupt signature")
	}
	for _, n := range []int{0, 4, 32, len(b) - 1} {
		os.WriteFile(filename, b[:n], 0600)
		if _, err := Open(filename); err == nil {
			t.Fatalf("accepted truncation %d", n)
		}
	}
}
