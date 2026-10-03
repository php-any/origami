package runtime_test

import (
	"sync"
	"testing"

	"github.com/php-any/origami/data"
	"github.com/php-any/origami/node"
	"github.com/php-any/origami/parser"
	"github.com/php-any/origami/runtime"
)

func TestStringCoercionUsesCallerRequest(t *testing.T) {
	p := parser.NewParser()
	base := runtime.NewVM(p)
	program, ctl := p.ParseString("function conversionLabel() { return 'worker'; } class ConversionParent { private const LABEL = 'parent'; public int $count = 0; public function __toString(): string { $this->count++; return self::LABEL . ':' . conversionLabel() . ':' . $this->count; } } class ConversionChild extends ConversionParent {} return new ConversionChild();", "string_conversion.php")
	if ctl != nil {
		t.Fatal(ctl.AsString())
	}
	value, ctl := program.GetValue(base.CreateContext(p.GetVariables()))
	if ctl != nil {
		t.Fatal(ctl.AsString())
	}
	object := value.(*data.ClassValue)
	var group sync.WaitGroup
	for i := 0; i < 8; i++ {
		group.Add(1)
		go func() {
			defer group.Done()
			request := runtime.NewRequestVM(base)
			function := node.NewFunctionStatement(nil, "conversionLabel", nil, []data.GetValue{node.NewReturnStatement(nil, data.NewStringValue("request"))}, nil, nil, false)
			if ctl := request.AddFunc(function); ctl != nil {
				t.Error(ctl.AsString())
				return
			}
			ctx := request.CreateContext(nil)
			for i, ref := range []data.TypeRef{data.TypeString, data.DeclaredTypeRef(data.NewDeclaredUnionType([]data.Types{data.TypeString, data.TypeInt}))} {
				converted, accepted, ctl := data.PrepareDeclaredValueInContext(ref, object, ctx)
				want := "parent:request:" + data.NewIntValue(i+1).AsString()
				if ctl != nil || !accepted || converted == nil || converted.AsString() != want {
					t.Errorf("request conversion = %v, %v, %v; want %s", converted, accepted, ctl, want)
				}
			}
			fresh, ctl := node.NewNewExpression(nil, "ConversionChild", nil).GetValue(ctx)
			if ctl != nil {
				t.Error(ctl.AsString())
				return
			}
			if _, accepted, ctl := data.PrepareDeclaredValueInContext(data.TypeString, fresh.(data.Value), ctx); ctl != nil || !accepted {
				t.Errorf("fresh object conversion failed: %v", ctl)
			}
			count, ctl := fresh.(*data.ClassValue).GetProperty("count")
			if ctl != nil || count.AsString() != "1" {
				t.Error("conversion lost the request object's identity")
			}
		}()
	}
	group.Wait()
	count, ctl := object.GetProperty("count")
	if ctl != nil || count.AsString() != "0" {
		t.Fatal("string conversion mutated the worker object")
	}
}
