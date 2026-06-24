package cli

import (
	"strings"
)

// Option 命令行选项
type Option struct {
	Name         string // 选项名称
	ShortName    string // 短名称（单字符）
	Description  string // 描述
	Required     bool   // 是否必需
	HasValue     bool   // 是否需要值
	DefaultValue string // 默认值
}

// OptionParser 选项解析器
type OptionParser struct {
	options []Option
}

// NewOptionParser 创建选项解析器
func NewOptionParser() *OptionParser {
	return &OptionParser{}
}

// AddOption 添加选项
func (p *OptionParser) AddOption(opt Option) {
	p.options = append(p.options, opt)
}

// Parse 解析命令行参数
func (p *OptionParser) Parse(args []string) (map[string]string, []string, error) {
	result := make(map[string]string)
	var remaining []string

	// 设置默认值
	for _, opt := range p.options {
		if opt.DefaultValue != "" {
			result[opt.Name] = opt.DefaultValue
		}
	}

	i := 0
	for i < len(args) {
		arg := args[i]

		if strings.HasPrefix(arg, "--") {
			// 长选项
			name := strings.TrimPrefix(arg, "--")
			opt := p.findOption(name)
			if opt == nil {
				remaining = append(remaining, arg)
				i++
				continue
			}

			if opt.HasValue {
				if i+1 >= len(args) {
					return nil, nil, &OptionError{
						Option:  name,
						Message: "选项 --" + name + " 需要一个值",
					}
				}
				result[opt.Name] = args[i+1]
				i += 2
			} else {
				result[opt.Name] = "true"
				i++
			}
		} else if strings.HasPrefix(arg, "-") && len(arg) > 1 {
			// 短选项
			name := arg[1:]
			opt := p.findOptionByShortName(name)
			if opt == nil {
				remaining = append(remaining, arg)
				i++
				continue
			}

			if opt.HasValue {
				if i+1 >= len(args) {
					return nil, nil, &OptionError{
						Option:  name,
						Message: "选项 -" + name + " 需要一个值",
					}
				}
				result[opt.Name] = args[i+1]
				i += 2
			} else {
				result[opt.Name] = "true"
				i++
			}
		} else {
			// 非选项参数
			remaining = append(remaining, arg)
			i++
		}
	}

	// 检查必需选项
	for _, opt := range p.options {
		if opt.Required {
			if _, ok := result[opt.Name]; !ok {
				return nil, nil, &OptionError{
					Option:  opt.Name,
					Message: "选项 --" + opt.Name + " 是必需的",
				}
			}
		}
	}

	return result, remaining, nil
}

// findOption 查找长选项
func (p *OptionParser) findOption(name string) *Option {
	for _, opt := range p.options {
		if opt.Name == name {
			return &opt
		}
	}
	return nil
}

// findOptionByShortName 查找短选项
func (p *OptionParser) findOptionByShortName(name string) *Option {
	for _, opt := range p.options {
		if opt.ShortName == name {
			return &opt
		}
	}
	return nil
}

// OptionError 选项错误
type OptionError struct {
	Option  string
	Message string
}

func (e *OptionError) Error() string {
	return e.Message
}
