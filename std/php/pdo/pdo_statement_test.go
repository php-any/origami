package pdo

import "testing"

func TestDriverValueString(t *testing.T) {
	tests := []struct {
		name  string
		value any
		want  string
	}{
		{name: "nil", value: nil, want: ""},
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

func TestDriverValue(t *testing.T) {
	tests := []struct {
		name  string
		value any
		isNil bool
		want  string
	}{
		{name: "nil preserves null", value: nil, isNil: true},
		{name: "bytes", value: []byte("value"), want: "value"},
		{name: "string", value: "value", want: "value"},
		{name: "integer", value: int64(42), want: "42"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := driverValue(tt.value)
			if tt.isNil {
				if got != nil {
					t.Fatalf("driverValue(%v) = %v, want nil", tt.value, got)
				}
				return
			}
			if s, ok := got.(string); !ok || s != tt.want {
				t.Fatalf("driverValue(%v) = %v, want %q", tt.value, got, tt.want)
			}
		})
	}
}
