package parser_json

import "gt-bots/libs/i18n"

const PARSER = "json"

// 注册解析器, 使用的时候需要引入
func init() {
    i18n.NewParser().Register(PARSER, NewParserJson())
}
