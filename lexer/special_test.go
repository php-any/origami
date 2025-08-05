package lexer

import (
	"reflect"
	"testing"

	"github.com/php-any/origami/token"
)

func Test_handleNumber(t *testing.T) {
	type args struct {
		input string
		start int
	}
	tests := []struct {
		name  string
		args  args
		want  SpecialToken
		want1 int
		want2 bool
	}{
		{
			name: "整数",
			args: args{
				input: "123",
				start: 0,
			},
			want: SpecialToken{
				Type:    token.INT,
				Literal: "123",
				Length:  3,
			},
			want1: 3,
			want2: true,
		},
		{
			name: "浮点数",
			args: args{
				input: "123.45",
				start: 0,
			},
			want: SpecialToken{
				Type:    token.FLOAT,
				Literal: "123.45",
				Length:  6,
			},
			want1: 6,
			want2: true,
		},
		{
			name: "科学计数法",
			args: args{
				input: "1.23e+10",
				start: 0,
			},
			want: SpecialToken{
				Type:    token.NUMBER,
				Literal: "1.23e+10",
				Length:  8,
			},
			want1: 8,
			want2: true,
		},
		{
			name: "十六进制",
			args: args{
				input: "0xFF",
				start: 0,
			},
			want: SpecialToken{
				Type:    token.NUMBER,
				Literal: "0xFF",
				Length:  4,
			},
			want1: 4,
			want2: true,
		},
		{
			name: "二进制",
			args: args{
				input: "0b1010",
				start: 0,
			},
			want: SpecialToken{
				Type:    token.NUMBER,
				Literal: "0b1010",
				Length:  6,
			},
			want1: 6,
			want2: true,
		},
		{
			name: "八进制",
			args: args{
				input: "0777",
				start: 0,
			},
			want: SpecialToken{
				Type:    token.NUMBER,
				Literal: "0777",
				Length:  4,
			},
			want1: 4,
			want2: true,
		},
		{
			name: "负数",
			args: args{
				input: "-123",
				start: 0,
			},
			want: SpecialToken{
				Type:    token.INT,
				Literal: "-123",
				Length:  4,
			},
			want1: 4,
			want2: true,
		},
		{
			name: "负浮点数",
			args: args{
				input: "-123.45",
				start: 0,
			},
			want: SpecialToken{
				Type:    token.FLOAT,
				Literal: "-123.45",
				Length:  7,
			},
			want1: 7,
			want2: true,
		},
		{
			name: "带分隔符的数字",
			args: args{
				input: "123;",
				start: 0,
			},
			want: SpecialToken{
				Type:    token.INT,
				Literal: "123",
				Length:  3,
			},
			want1: 3,
			want2: true,
		},
		{
			name: "无效数字",
			args: args{
				input: "123abc",
				start: 0,
			},
			want: SpecialToken{
				Type:    token.NUMBER,
				Literal: "123abc",
				Length:  6,
			},
			want1: 6,
			want2: true,
		},
		{
			name: "空字符串",
			args: args{
				input: "",
				start: 0,
			},
			want:  SpecialToken{},
			want1: 0,
			want2: false,
		},
		{
			name: "非数字开始",
			args: args{
				input: "abc123",
				start: 0,
			},
			want:  SpecialToken{},
			want1: 0,
			want2: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, got1, got2 := handleNumber(tt.args.input, tt.args.start)
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("handleNumber() got = %v, want %v", got, tt.want)
			}
			if got1 != tt.want1 {
				t.Errorf("handleNumber() got1 = %v, want %v", got1, tt.want1)
			}
			if got2 != tt.want2 {
				t.Errorf("handleNumber() got2 = %v, want %v", got2, tt.want2)
			}
		})
	}
}
