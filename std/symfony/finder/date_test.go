package finder

import (
	"math"
	"testing"
	"time"
)

func TestParseDateExprRelative(t *testing.T) {
	// FileSessionHandler::gc 的写法：'<= now - '.$lifetime.' seconds'
	op, ts, ok := parseDateExpr("<= now - 7200 seconds")
	if !ok || op != "<=" {
		t.Fatalf("gc 表达式解析失败: op=%q ok=%v", op, ok)
	}
	want := time.Now().Add(-7200 * time.Second).Unix()
	if d := ts - want; d < -2 || d > 2 {
		t.Fatalf("阈值偏移 %ds，期望 ≈%d 得到 %d", d, want, ts)
	}

	cases := []struct {
		expr string
		op   string
		// 期望的偏移秒数（相对 now），用于断言方向
		shift int64
	}{
		{">= now", ">=", 0},
		{"> now - 1 day", ">", -86400},
		{"< 2 hours ago", "<", -7200},
		{">= now -  3   minutes", ">=", -180},
		{"<= yesterday", "<=", -86400}, // 粗略：昨天的零点在 24h 内
	}
	for _, c := range cases {
		op, ts, ok := parseDateExpr(c.expr)
		if !ok {
			t.Errorf("%q 解析失败", c.expr)
			continue
		}
		if op != c.op {
			t.Errorf("%q op=%q 期望 %q", c.expr, op, c.op)
		}
		if c.expr == "<= yesterday" {
			continue // 零点对齐，不按固定偏移断言
		}
		if d := ts - time.Now().Unix() - c.shift; d < -2 || d > 2 {
			t.Errorf("%q 偏移 %ds 期望 %d", c.expr, d, c.shift)
		}
	}
}

func TestParseDateExprAbsolute(t *testing.T) {
	op, ts, ok := parseDateExpr("> 2024-01-02 03:04:05")
	if !ok || op != ">" {
		t.Fatalf("绝对日期解析失败: op=%q ok=%v", op, ok)
	}
	if want := time.Date(2024, 1, 2, 3, 4, 5, 0, time.Local).Unix(); ts != want {
		t.Fatalf("绝对日期 = %d 期望 %d", ts, want)
	}
	if _, ts, ok := parseDateExpr("== 1700000000"); !ok || ts != 1700000000 {
		t.Fatalf("时间戳解析失败: ts=%d ok=%v", ts, ok)
	}
}

func TestParseDateExprInvalidNeverMatches(t *testing.T) {
	// 解析不了的条件必须「永不匹配」，否则 gc 会删掉全部文件。
	for _, expr := range []string{"", "   ", "!! oops", "<= ", "> soon"} {
		op, ts, ok := parseDateExpr(expr)
		if ok {
			t.Errorf("%q 不应解析成功: op=%q ts=%d", expr, op, ts)
			continue
		}
		specs := finderDateSpecs([]string{expr})
		if len(specs) != 1 || !specs[0].never {
			t.Errorf("%q 应得到 never 条件: %+v", expr, specs)
		}
		if specs[0].match(math.MaxInt64) || specs[0].match(0) {
			t.Errorf("%q 的 never 条件不匹配任何 mtime", expr)
		}
	}
	// 空参数列表才是「不过滤」。
	if specs := finderDateSpecs(nil); specs != nil {
		t.Errorf("空参数应为 nil（不过滤），得到 %+v", specs)
	}
}
