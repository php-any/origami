package runtime

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/php-any/origami/data"
	"github.com/php-any/origami/parser"
)

func TestRequestOverlayLSBContainerInstance(t *testing.T) {
	dir := t.TempDir()
	boot := filepath.Join(dir, "OverlayLsb_boot.php")
	req := filepath.Join(dir, "OverlayLsb_req.php")
	if err := os.WriteFile(boot, []byte(`<?php
class OverlayLsb_Box {
    protected static $instance;
    public static $probe;
    public $tok = '';
    public static function getInstance() {
        return static::$instance ??= new static;
    }
    public static function setInstance($c = null) {
        return static::$instance = $c;
    }
}
class OverlayLsb_App extends OverlayLsb_Box {}
$boot = new OverlayLsb_App();
$boot->tok = 'BOOT';
OverlayLsb_Box::setInstance($boot);
`), 0o644); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(req, []byte(`<?php
$req = new OverlayLsb_App();
$req->tok = 'REQUEST';
$req->setInstance($req);
$got = OverlayLsb_Box::getInstance();
OverlayLsb_Box::$probe = $got->tok;
`), 0o644); err != nil {
		t.Fatal(err)
	}

	vm := NewVM(parser.NewParser()).(*VM)
	if _, ctl := vm.LoadAndRun(boot); ctl != nil {
		t.Fatalf("boot: %v", ctl)
	}

	restore := BeginRequestOutput()
	defer restore()
	if _, ctl := vm.LoadAndRun(req); ctl != nil {
		t.Fatalf("request: %v", ctl)
	}

	cls, ok := vm.GetClass("OverlayLsb_Box")
	if !ok {
		t.Fatal("OverlayLsb_Box missing")
	}
	getter, ok := cls.(data.GetStaticProperty)
	if !ok {
		t.Fatal("no GetStaticProperty")
	}
	v, ok := getter.GetStaticProperty("probe")
	if !ok || v == nil {
		t.Fatal("probe missing")
	}
	if v.AsString() != "REQUEST" {
		t.Fatalf("HTTP overlay 下 Container::getInstance 应是请求 Application，实际 %q", v.AsString())
	}
}
