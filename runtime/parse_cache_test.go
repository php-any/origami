package runtime

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/php-any/origami/data"
	"github.com/php-any/origami/parser"
)

func TestParseFileCachedSkipsRepeatedDiskRead(t *testing.T) {
	dir := t.TempDir()
	file := filepath.Join(dir, "cached.php")
	if err := os.WriteFile(file, []byte(`<?php
$hits = 0;
function bump() { global $hits; $hits++; return $hits; }
echo bump();
`), 0o644); err != nil {
		t.Fatal(err)
	}

	p := parser.NewParser()
	vm := NewVM(p).(*VM)
	vm.SetThrowControl(func(data.Control) {})

	program1, vars1, acl := vm.ParseFileCached(file)
	if acl != nil {
		t.Fatalf("first parse: %v", acl)
	}
	if program1 == nil || len(vars1) == 0 {
		t.Fatal("expected program and vars")
	}

	if err := os.Remove(file); err != nil {
		t.Fatal(err)
	}

	program2, _, acl := vm.ParseFileCached(file)
	if acl != nil {
		t.Fatalf("cached parse should not touch disk: %v", acl)
	}
	if program1 != program2 {
		t.Fatal("expected same cached program instance")
	}
}

func TestTempVMLoadInCallerContextUsesParseCache(t *testing.T) {
	dir := t.TempDir()
	file := filepath.Join(dir, "view.php")
	if err := os.WriteFile(file, []byte(`<?php echo "ok";`), 0o644); err != nil {
		t.Fatal(err)
	}

	p := parser.NewParser()
	base := NewVM(p).(*VM)
	base.SetThrowControl(func(data.Control) {})
	temp := NewTempVM(base).(*TempVM)
	parent := temp.CreateContext(nil)

	if _, acl := temp.LoadInCallerContext(parent, file); acl != nil {
		t.Fatalf("first include: %v", acl)
	}
	if err := os.Remove(file); err != nil {
		t.Fatal(err)
	}
	if _, acl := temp.LoadInCallerContext(parent, file); acl != nil {
		t.Fatalf("cached include should not touch disk: %v", acl)
	}
}
