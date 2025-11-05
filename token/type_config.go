package token

// TokenType 表示 token 的类型
type TokenType int

// Token 类型常量（使用更常见的简写命名）
const (
	// 关键字
	KEYWORD_START  TokenType = iota + 1 // if 条件语句
	IF                                  // if 条件语句
	ELSE                                // else 条件语句
	ELSE_IF                             // else if 条件语句
	WHILE                               // while 循环语句
	FOR                                 // for 循环语句
	FOREACH                             // foreach 循环语句
	DO                                  // do-while 循环语句
	SWITCH                              // switch 选择语句
	CASE                                // case 选择分支
	BREAK                               // break 跳出循环
	CONTINUE                            // continue 继续循环
	RETURN                              // return 返回语句
	FUNC                                // function 函数定义
	CLASS                               // class 类定义
	PUBLIC                              // public 访问修饰符
	PRIVATE                             // private 访问修饰符
	PROTECTED                           // protected 访问修饰符
	STATIC                              // static 静态成员
	FINAL                               // final 最终类/方法
	ABSTRACT                            // abstract 抽象类/方法
	INTERFACE                           // interface 接口定义
	TRAIT                               // trait 特性定义
	NAMESPACE                           // namespace 命名空间
	USE                                 // use 使用命名空间/特性
	AS                                  // as 别名定义
	NEW                                 // new 实例化对象
	INSTANCEOF                          // instanceof 类型检查
	LIKE                                // like 类型检查
	CONST                               // const 常量定义
	VAR                                 // var 变量定义
	ECHO                                // echo 输出语句
	THROW                               // throw 抛出异常
	TRY                                 // try 异常处理
	CATCH                               // catch 捕获异常
	FINALLY                             // finally 最终处理
	CLONE                               // clone 克隆对象
	YIELD                               // yield 生成器
	FROM                                // from 生成器
	INSTEAD_OF                          // insteadof 特性冲突解决
	EXTENDS                             // extends 继承
	IMPLEMENTS                          // implements 实现接口
	LIST                                // list 列表赋值
	ARRAY                               // array 数组定义
	ISSET                               // isset 变量存在检查
	UNSET                               // unset 销毁变量
	EMPTY                               // empty 空值检查
	DIE                                 // die 终止脚本
	EXIT                                // exit 终止脚本
	EVAL                                // eval 执行字符串
	HALT_COMPILER                       // __halt_compiler 停止编译
	MATCH                               // match 匹配表达式 (PHP 8.0+)
	ENUM                                // enum 枚举定义 (PHP 8.1+)
	READONLY                            // readonly 只读属性 (PHP 8.1+)
	FN                                  // fn 箭头函数 (PHP 7.4+)
	SPAWN                               // spawn 异步
	THIS                                // this
	PARENT                              // parent
	IN                                  // in 运算符
	DEFAULT                             // default 默认值
	UNUSED                              // _ 符号, 跳过变量
	DIR                                 // __DIR__ 魔术常量
	FILE                                // __FILE__
	LINE                                // __LINE__
	GENERIC_TYPE                        // 泛型类型参数
	CLASS_CONSTANT                      // ::class 类常量关键字

	KEYWORD_END TokenType = iota + 100

	// 运算符
	ADD                 // + 加法运算符
	SUB                 // - 减法运算符
	MUL                 // * 乘法运算符
	QUO                 // / 除法运算符（quotient 缩写）
	REM                 // % 取模运算符（remainder 缩写）
	ASSIGN              // = 赋值运算符
	EQ                  // == 等于运算符
	NE                  // != 不等于运算符
	EQ_STRICT           // === 全等运算符
	NE_STRICT           // !== 不全等运算符
	LT                  // < 小于运算符
	GT                  // > 大于运算符
	LE                  // <= 小于等于运算符
	GE                  // >= 大于等于运算符
	LAND                // && 逻辑与运算符（logical and）
	LOR                 // || 逻辑或运算符（logical or）
	NOT                 // ! 逻辑非运算符
	BIT_AND             // & 按位与运算符
	BIT_OR              // | 按位或运算符
	BIT_XOR             // ^ 按位异或运算符
	BIT_NOT             // ~ 按位取反运算符
	SHL                 // << 左移运算符（shift left）
	SHR                 // >> 右移运算符（shift right）
	INCR                // ++ 自增运算符
	DECR                // -- 自减运算符
	OBJECT_OPERATOR     // -> 对象成员访问运算符
	ARRAY_KEY_VALUE     // => 数组键值对分隔符
	TERNARY             // ? : 三元运算符
	COLON               // : 标签或三元运算符部分
	SCOPE_RESOLUTION    // :: 静态成员访问运算符
	AT                  // @ 错误抑制运算符
	DOLLAR              // $ 变量标识符
	COMMA               // , 分隔符
	SEMICOLON           // ; 语句结束符
	LPAREN              // ( 左括号
	RPAREN              // ) 右括号
	LBRACE              // { 左花括号
	RBRACE              // } 右花括号
	LBRACKET            // [ 左方括号
	RBRACKET            // ] 右方括号
	SPACESHIP           // <=> 太空船运算符 (PHP 7+)
	NULLSAFE_CALL       // ??-> 空安全对象运算符 (PHP 8.0+)
	NULL_COALESCE       // ?? 空合并运算符 (PHP 7+)
	POWER               // ** 幂运算符 (PHP 5.6+)
	POWER_EQ            // **= 幂赋值运算符 (PHP 5.6+)
	ADD_EQ              // += 加法赋值运算符
	SUB_EQ              // -= 减法赋值运算符
	MUL_EQ              // *= 乘法赋值运算符
	QUO_EQ              // /= 除法赋值运算符
	REM_EQ              // %= 取模赋值运算符
	CONCAT_EQ           // .= 字符串连接赋值运算符
	BIT_AND_EQ          // &= 按位与赋值运算符
	BIT_OR_EQ           // |= 按位或赋值运算符
	BIT_XOR_EQ          // ^= 按位异或赋值运算符
	SHL_EQ              // <<= 左移赋值运算符
	SHR_EQ              // >>= 右移赋值运算符
	NAMESPACE_SEPARATOR // \ 命名空间分隔符
	DOT                 // . 点号（注意与 CONCAT 冲突）
	ELLIPSIS            // ... 省略号
	DOUBLE_DOT          // .. 双点号

	// 专用：字符串插值连接符（仅预处理阶段注入，不由词法直接产生）
	INTERPOLATION_LINK

	VALUE_START
	// 字面量
	NUMBER  // 复杂数字面量
	INT     // 整数字面量
	FLOAT   // 浮点数字面量
	STRING  // 字符串字面量
	BOOL    // bool 类型
	HEREDOC // 定界符字符串
	NOWDOC  // 定界符字符串(不解析变量)
	NULL    // null 字面量
	TRUE    // true 字面量
	FALSE   // false 字面量
	BYTE    // byte 字面量

	VALUE_END TokenType = iota + 100

	// 其他
	IDENTIFIER        // 标识符
	VARIABLE          // 变量
	COMMENT           // 单行注释
	MULTILINE_COMMENT // 多行注释
	WHITESPACE        // 空白字符
	EOF               // 文件结束
	NEWLINE           // 换行符

	START_TAG // <?php 开始标签
	END_TAG   // ?> 结束标签
	HTML_TAG

	UNKNOWN // 未知标记
)
