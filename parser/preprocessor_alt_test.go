package parser

import (
	"fmt"
	"os"
	"strings"
	"testing"
)

func TestAltConvertRealFile(t *testing.T) {
	content, err := os.ReadFile("/workspace/examples/laravel13/storage/framework/views/748237abb381b795bb452776ecf33d8fd51cc21aebbc3a3766f655fc6a988263.php")
	if err != nil {
		t.Fatal(err)
	}
	code := string(content)
	fmt.Println("hasControlColon:", hasControlColon(code))
	out := convertAltPHPSyntax("test.php", code)
	if strings.Contains(out, "if(session('error')):") {
		t.Fatal("conversion of if(...): failed on real file")
	}
	if strings.Contains(out, "if ($__bag->has($__errorArgs[0])) :") {
		t.Fatal("conversion of if (cond) : failed on real file")
	}
	t.Log("OK")
}
