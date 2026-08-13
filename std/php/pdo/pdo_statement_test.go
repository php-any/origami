package pdo

import "testing"

func TestDriverValueString(t *testing.T) {
	tests := []struct {
		name  string
		value any
		want  string
	}{
		{name: "nil", value: nil, want: ""},
		{name: "bytes", value: []byte("value"), want: "value"},
		{name: "string", value: "value", want: "value"},
		{name: "integer", value: int64(42), want: "42"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := driverValueString(tt.value); got != tt.want {
				t.Fatalf("driverValueString(%v) = %q, want %q", tt.value, got, tt.want)
			}
		})
	}
}
