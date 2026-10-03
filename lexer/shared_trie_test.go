package lexer

import (
	"fmt"
	"reflect"
	"sync"
	"testing"
)

func TestIndependentLexersShareReadOnlyTrie(t *testing.T) {
	first, second := NewLexer(), NewLexer()
	if first.root != second.root {
		t.Fatal("lexers retained independent token tries")
	}
	const source = `<?php $a = fn($x) => $x + 2; echo "hello {$a(3)}"; ?>tail`
	want := first.TokenizeTemplate(source)
	var workers sync.WaitGroup
	for i := 0; i < 8; i++ {
		workers.Add(1)
		go func(i int) {
			defer workers.Done()
			lexer := NewLexer()
			for j := 0; j < 10; j++ {
				lexer.Tokenize(fmt.Sprintf("$other%d = %d;", i, j))
				if got := lexer.TokenizeTemplate(source); !reflect.DeepEqual(got, want) {
					t.Errorf("worker %d changed token output", i)
					return
				}
			}
		}(i)
	}
	workers.Wait()
}
