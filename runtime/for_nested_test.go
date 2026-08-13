package runtime_test

import (
	"bytes"
	"os"
	"strings"
	"testing"

	"github.com/php-any/origami/data"
	"github.com/php-any/origami/parser"
	"github.com/php-any/origami/runtime"
)

func TestNestedForSecondIteration(t *testing.T) {
	src := `<?php
for ($i = 0; $i < 2; $i++) {
    for ($k = 0; $k < 2; $k++) {
        echo "i=$i k=$k\n";
    }
}
`
	p := parser.NewParser()
	program, acl := p.ParseString(src, "test.zy")
	if acl != nil {
		t.Fatalf("parse error: %v", acl)
	}

	vm := runtime.NewVM(p)
	ctx := vm.CreateContext(p.GetVariables())

	oldStdout := os.Stdout
	r, w, err := os.Pipe()
	if err != nil {
		t.Fatal(err)
	}
	os.Stdout = w

	data.ResetUserOutput()
	_, acl = program.GetValue(ctx)
	w.Close()
	os.Stdout = oldStdout

	if acl != nil {
		t.Fatalf("unexpected control: %v", acl)
	}

	var buf bytes.Buffer
	_, _ = buf.ReadFrom(r)
	out := buf.String()

	want := []string{
		"i=0 k=0",
		"i=0 k=1",
		"i=1 k=0",
		"i=1 k=1",
	}
	for _, line := range want {
		if !strings.Contains(out, line) {
			t.Fatalf("missing %q in output:\n%s", line, out)
		}
	}
}
