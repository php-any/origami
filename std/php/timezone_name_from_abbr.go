package php

import (
	"time"

	"github.com/php-any/origami/data"
	"github.com/php-any/origami/node"
)

// timezoneNames 为 timezone_name_from_abbr 使用的常见 IANA 时区名列表。
// 与 PHP 行为一致：根据缩写或 UTC 偏移匹配并返回第一个匹配的时区名。
var timezoneNames = []string{
	"Africa/Cairo", "Africa/Johannesburg", "Africa/Lagos", "America/Chicago",
	"America/Denver", "America/Los_Angeles", "America/Mexico_City", "America/New_York",
	"America/Sao_Paulo", "America/Toronto", "Asia/Dubai", "Asia/Hong_Kong",
	"Asia/Jerusalem", "Asia/Kolkata", "Asia/Shanghai", "Asia/Singapore",
	"Asia/Tokyo", "Australia/Sydney", "Europe/Berlin", "Europe/London",
	"Europe/Paris", "Europe/Moscow", "Pacific/Auckland", "UTC",
}

func NewTimezoneNameFromAbbrFunction() data.FuncStmt {
	return &TimezoneNameFromAbbrFunction{}
}

type TimezoneNameFromAbbrFunction struct {
	data.Function
}

func (f *TimezoneNameFromAbbrFunction) Call(ctx data.Context) (data.GetValue, data.Control) {
	abbrVal, _ := ctx.GetIndexValue(0)
	utcOffsetVal, _ := ctx.GetIndexValue(1)
	isDSTVal, _ := ctx.GetIndexValue(2)

	abbr := ""
	if abbrVal != nil {
		abbr = abbrVal.AsString()
	}

	utcOffset := -1
	if utcOffsetVal != nil {
		if ai, ok := utcOffsetVal.(data.AsInt); ok {
			if o, err := ai.AsInt(); err == nil {
				utcOffset = o
			}
		}
	}

	isDST := -1
	if isDSTVal != nil {
		if di, ok := isDSTVal.(data.AsInt); ok {
			if d, err := di.AsInt(); err == nil {
				isDST = d
			}
		}
	}

	name := timezoneNameFromAbbr(abbr, utcOffset, isDST)
	if name == "" {
		return data.NewBoolValue(false), nil
	}
	return data.NewStringValue(name), nil
}

func (f *TimezoneNameFromAbbrFunction) GetName() string {
	return "timezone_name_from_abbr"
}

func (f *TimezoneNameFromAbbrFunction) GetParams() []data.GetValue {
	return []data.GetValue{
		node.NewParameter(nil, "abbr", 0, nil, nil),
		node.NewParameter(nil, "utc_offset", 1, node.NewIntLiteral(nil, "-1"), nil),
		node.NewParameter(nil, "is_dst", 2, node.NewIntLiteral(nil, "-1"), nil),
	}
}

func (f *TimezoneNameFromAbbrFunction) GetVariables() []data.Variable {
	return []data.Variable{
		node.NewVariable(nil, "abbr", 0, nil),
		node.NewVariable(nil, "utc_offset", 1, nil),
		node.NewVariable(nil, "is_dst", 2, nil),
	}
}

// timezoneNameFromAbbr 根据缩写、UTC 偏移（秒）和是否夏令时在预定义时区列表中查找并返回第一个匹配的 IANA 时区名。
// 与 PHP timezone_name_from_abbr 语义一致：utcOffset/isDST 为 -1 表示不参与匹配。
func timezoneNameFromAbbr(abbr string, utcOffset, isDST int) string {
	// 用于检测 DST 的参考时间：冬季（非 DST）与夏季（DST）
	jan := time.Date(2024, 1, 15, 12, 0, 0, 0, time.UTC)
	jul := time.Date(2024, 7, 15, 12, 0, 0, 0, time.UTC)

	for _, zoneName := range timezoneNames {
		loc, err := time.LoadLocation(zoneName)
		if err != nil {
			continue
		}

		if abbr != "" {
			// 先按缩写匹配：在冬季和夏季各取一次 Zone，看是否与 abbr 一致
			janInZone := jan.In(loc)
			julInZone := jul.In(loc)
			janAbbr, janOff := janInZone.Zone()
			julAbbr, julOff := julInZone.Zone()

			abbrMatch := (janAbbr == abbr || julAbbr == abbr)
			if !abbrMatch {
				continue
			}

			if utcOffset >= 0 {
				// 要求偏移也匹配：优先精确偏移
				janMatch := (janAbbr == abbr && janOff == utcOffset)
				julMatch := (julAbbr == abbr && julOff == utcOffset)
				if isDST == 0 {
					if !janMatch {
						continue
					}
				} else if isDST == 1 {
					if !julMatch {
						continue
					}
				} else {
					if !janMatch && !julMatch {
						continue
					}
				}
			}
			return zoneName
		}

		// 无缩写：仅按 utcOffset 和 isDST 匹配
		if utcOffset < 0 {
			continue
		}
		janInZone := jan.In(loc)
		julInZone := jul.In(loc)
		_, janOff := janInZone.Zone()
		_, julOff := julInZone.Zone()

		if isDST == 0 {
			if janOff == utcOffset {
				return zoneName
			}
		} else if isDST == 1 {
			if julOff == utcOffset {
				return zoneName
			}
		} else {
			if janOff == utcOffset || julOff == utcOffset {
				return zoneName
			}
		}
	}
	return ""
}
