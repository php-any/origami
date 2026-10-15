package parser

import (
	"strings"
	"testing"

	"github.com/php-any/origami/data"
)

func TestLoopControlRejectsInvalidLevels(t *testing.T) {
	for _, keyword := range []string{"break", "continue"} {
		for _, level := range []string{"0", "-1", "1.5", "1e2", "$level", "999999999999999999999999999"} {
			t.Run(keyword+level, func(t *testing.T) {
				_, ctl := NewParser().ParseString("while (true) { "+keyword+" "+level+"; }", "levels.php")
				failure, ok := ctl.(*data.ThrowValue)
				if !ok || !failure.PHPCompileFatal || !strings.Contains(failure.Error.Error(), "positive integers") {
					t.Fatalf("expected positive-integer compile fatal, got %v", ctl)
				}
			})
		}
	}
}

func TestLoopControlAcceptsIntegerBases(t *testing.T) {
	for _, keyword := range []string{"break", "continue"} {
		for _, level := range []string{"1", "0x1", "0b1", "01", "0o1"} {
			if _, ctl := NewParser().ParseString("while (true) { "+keyword+" "+level+"; }", "levels.php"); ctl != nil {
				t.Fatalf("%s %s: %v", keyword, level, ctl)
			}
		}
	}
}
