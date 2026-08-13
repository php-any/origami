
    /**
     * 预处理单词, 给一些词打 tag, 补充省却结构+预测部分关键字类型
     * 删除空白符号
     */
    public func readAllTokens():Unit {
        var list = ArrayList<Token>()
        var got = ArrayList<Token>()

        var lastNewline = false
        if (this.tokens[0].Type == TokenType.DOCTYPE) {
            this.isHtml = true
        } else {
            // 读取完毕, 开始处理
            for (i in 0..this.tokens.size) {
                var t = this.tokens[i]
                if (lastNewline) {
                    if (t.Type == TokenType.NEWLINE) {
                        // 如果多个换行保留一个
                        continue
                    } else {
                        lastNewline = false
                        list.add(Token(TokenType.NEWLINE, "\n", t.Line, t.Pos))
                    }
                }

                match (t.Type) {
                    // 字符串插值检查; 进行转义
                    case TokenType.STRINGS =>
                        t.Literal = t.Literal.trimStart("\"")
                        t.Literal = t.Literal.trimEnd("\"")

                        t.Type = TokenType.T_STRING

                        if(t.Literal.size <= 0) {
                            list.add(t) // 空的，无需处理
                            continue
                        }

                        var newStr = ArrayList<Rune>()

                        var i = 0
                        var newLiteral = t.Literal.toRuneArray()
                        while (i < newLiteral.size) {
                            let rs = newLiteral[i]
                            match (rs) {
                                case '{' =>
                                    // {、{} 不会是插值
                                    if (newLiteral.size <= 2) {
                                        newStr.add(all: newLiteral)
                                        break
                                    }
                                    if (i == 0 && newLiteral.size > (i + 3) && newLiteral[i + 1] == r'$' &&
                                        !isNumber(newLiteral[i + 2]) && !isSpecialSymbols(newLiteral[i + 2])) {
                                        // 必须是 {$ 作为开始才算插值
                                        (i, newStr) = this.checkStringInsertValue(i, newLiteral, t, list, newStr)
                                    } else if (i >= 1 && newLiteral.size > (i + 3) && newLiteral[i + 1] == r'$' &&
                                        !isNumber(newLiteral[i + 2]) && !isSpecialSymbols(newLiteral[i + 2])) {
                                        // 必须是 {$ 作为开始才算插值
                                        (i, newStr) = this.checkStringInsertValue(i, newLiteral, t, list, newStr)
                                    } else {
                                        var oldIndex = i
                                        if (newLiteral.size > (i + 3) && newLiteral[i + 1] == r'@' &&
                                            !isNumber(newLiteral[i + 2]) && !isSpecialSymbols(newLiteral[i + 2])) {
                                            // 必须是 {@标识符( 作为开始才算插值
                                            (i, newStr) = this.checkStringInsertFunc(i, newLiteral, t, list, newStr)
                                        }
                                        if (oldIndex == i) {
                                            newStr.add(rs)
                                        }
                                    }

                                case _ => newStr.add(rs)
                            }

                            i++
                        }
                        if (newStr.size > 0) {
                            var newT = t
                            newT.Literal = String(newStr.toArray())
                            list.add(newT)
                        } else if(list[list.size - 1].Type == TokenType.PLUS) {
                            var newT = t
                            newT.Literal = ""
                            list.add(newT)
                        }
                        continue
                    // TODO 待实现
                    case TokenType.T_STRING =>
                        t.Literal = t.Literal.trimStart("'")
                        t.Literal = t.Literal.trimEnd("'")
                    // let、var 语句, 后一个 token 改成变量(非$开始也支持)
                    // 忽略空白符号
                    case TokenType.SPACE | TokenType.TAB => continue
                    // 注释, 当前不支持反射，无需注释保留
                    case TokenType.COMMENT | TokenType.MULTILINE_COMMENT => continue
                    // 换行, 如果多个换行保留一个
                    case TokenType.NEWLINE =>
                        lastNewline = true
                        continue

                    // TODO 待实现
                    case _ => ()
                }

                list.add(t)
            }

            // 检查换行符号, 判断是否需要补充;符号
            var checkLast = Token(TokenType.SPACE, "", 0, 0)
            for (i in 0..list.size) {
                var t = list[i]
                if (t.Type == TokenType.NEWLINE) {
                    // 不用补充;的前一个符号
                    match (checkLast.Type) {
                        // 前一个已经是;结束了, 不用补充
                        case TokenType.SENTINEL => continue
                        // , . | ->
                        case TokenType.COMMA | TokenType.DOT | TokenType.RBRACE | ARROW => continue
                        // + - * / %
                        case TokenType.PLUS | TokenType.MINUS | TokenType.MULTIPLY | TokenType.SLASH | TokenType
                                .REMAINDER => continue
                        // & |
                        case TokenType.BITAND | TokenType.BITOR => continue
                        // && ||
                        case TokenType.LOGICAL_AND | TokenType.LOGICAL_OR => continue
                        // [ { (
                        case TokenType.LBRACKET | TokenType.LBRACE | TokenType.LPAREN => continue
                        case _ =>
                            // 判断下一个单词
                            if (i + 1 <= list.size) {
                                let peek = list[i + 1]
                                match (peek.Type) {
                                    // ] } )   还没有闭合的语句, 也跳过补充;
                                    case TokenType.LBRACKET | TokenType.RBRACKET | TokenType.LBRACE | TokenType.RBRACE | TokenType
                                            .LPAREN | TokenType.RPAREN => continue
                                    case _ =>
                                        t.Type = TokenType.SENTINEL
                                        t.Literal = ";"
                                }
                            }
                    }
                }
                got.add(t)
                checkLast = t
            }
            this.tokens = got.toArray()
        }

        // 重置循环状态
        this.curIndex = 0
        this.curToken = this.tokens[0]
        if (tokens.size > 1) {
            this.peekToken = this.tokens[1]
        } else {
            this.peekToken = eofToken(this)
        }
        // debugCodes(this.tokens, 0)
    }