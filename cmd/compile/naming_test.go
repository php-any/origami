package compile

import "testing"

func TestPathToFuncSuffix(t *testing.T) {
	tests := []struct {
		path string
		want string
	}{
		{"user/login.php", "User_Login"},
		{"user_login.php", "UserLogin"},
		{"examples/spring/src/config/DatabaseBootstrap.php", "Examples_Spring_Src_Config_DatabaseBootstrap"},
		{"examples/spring/index.php", "Examples_Spring_Index"},
	}
	for _, tt := range tests {
		got := pathToFuncSuffix(tt.path)
		if got != tt.want {
			t.Errorf("pathToFuncSuffix(%q) = %q, want %q", tt.path, got, tt.want)
		}
	}
}

func TestFuncNamesDistinctForCollision(t *testing.T) {
	a := funcNameFromPath("user/login.php")
	b := funcNameFromPath("user_login.php")
	if a == b {
		t.Fatalf("expected distinct names, both %q", a)
	}
	if a != "AST_User_Login" {
		t.Fatalf("user/login: got %q", a)
	}
	if b != "AST_UserLogin" {
		t.Fatalf("user_login: got %q", b)
	}
}

func TestGoFileNamesDistinctForCollision(t *testing.T) {
	a := goFileNameFromPath("user/login.php")
	b := goFileNameFromPath("user_login.php")
	if a == b {
		t.Fatalf("expected distinct files, both %q", a)
	}
}
