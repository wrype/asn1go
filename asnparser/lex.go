package asnparser

import (
	"errors"
	"fmt"
	"log"
	"regexp"
	"strconv"
)

//go:generate goyacc -o asn.go -p "ParsedGrammar" asn.y

type AsnLexer struct {
	text   []byte      //输入的内容
	err    error       //错误信息
	result *ModuleList //最后的语法树
}

func NewAsnLexer(buf []byte) *AsnLexer {
	return &AsnLexer{text: buf}
}

func (lex *AsnLexer) Lex(lval *ParsedGrammarSymType) int {
	//返回词法分析结果
	for {
		spaceReg, _ := regexp.Compile(`^\s+`)
		commentReg, _ := regexp.Compile(`^[-][-][^\n]*`)
		numberReg, _ := regexp.Compile(`^[-]?[0-9]+`)
		identifierReg, _ := regexp.Compile(`^[a-zA-Z0-9][-a-zA-Z0-9]*`)
		if lex.regexp(spaceReg, lval) {
			continue //忽略空白字符
		} else if lex.regexp(commentReg, lval) {
			continue //忽略注释
		} else if lex.keyword("ABSENT", lval) {
			return ABSENT
		} else if lex.keyword("ABSTRACT_SYNTAX", lval) {
			return ABSTRACT_SYNTAX
		} else if lex.keyword("ALL", lval) {
			return ALL
		} else if lex.keyword("ANY", lval) {
			return ANY
		} else if lex.keyword("APPLICATION", lval) {
			return APPLICATION
		} else if lex.keyword("AUTOMATIC", lval) {
			return AUTOMATIC
		} else if lex.keyword("BEGIN", lval) {
			return BEGIN
		} else if lex.keyword("BIT", lval) {
			return BIT
		} else if lex.keyword("BMPString", lval) {
			return BMPString
		} else if lex.keyword("BOOLEAN", lval) {
			return BOOLEAN
		} else if lex.keyword("BY", lval) {
			return BY
		} else if lex.keyword("CHARACTER", lval) {
			return CHARACTER
		} else if lex.keyword("CHOICE", lval) {
			return CHOICE
		} else if lex.keyword("CLASS", lval) {
			return CLASS
		} else if lex.keyword("COMPONENTS", lval) {
			return COMPONENTS
		} else if lex.keyword("COMPONENT", lval) {
			return COMPONENT
		} else if lex.keyword("CONSTRAINED", lval) {
			return CONSTRAINED
		} else if lex.keyword("CONTAINING", lval) {
			return CONTAINING
		} else if lex.keyword("DEFAULT", lval) {
			return DEFAULT
		} else if lex.keyword("DEFINITIONS", lval) {
			return DEFINITIONS
		} else if lex.keyword("DEFINED", lval) {
			return DEFINED
		} else if lex.keyword("EMBEDDED", lval) {
			return EMBEDDED
		} else if lex.keyword("ENCODED", lval) {
			return ENCODED
		} else if lex.keyword("ENCODING_CONTROL", lval) {
			return ENCODING_CONTROL
		} else if lex.keyword("END", lval) {
			return END
		} else if lex.keyword("ENUMERATED", lval) {
			return ENUMERATED
		} else if lex.keyword("EXPLICIT", lval) {
			return EXPLICIT
		} else if lex.keyword("EXPORTS", lval) {
			return EXPORTS
		} else if lex.keyword("EXTENSIBILITY", lval) {
			return EXTENSIBILITY
		} else if lex.keyword("EXTERNAL", lval) {
			return EXTERNAL
		} else if lex.keyword("FALSE", lval) {
			return FALSE
		} else if lex.keyword("FROM", lval) {
			return FROM
		} else if lex.keyword("GeneralizedTime", lval) {
			return GeneralizedTime
		} else if lex.keyword("GeneralString", lval) {
			return GeneralString
		} else if lex.keyword("GraphicString", lval) {
			return GraphicString
		} else if lex.keyword("IA5String", lval) {
			return IA5String
		} else if lex.keyword("IDENTIFIER", lval) {
			return IDENTIFIER
		} else if lex.keyword("IMPLICIT", lval) {
			return IMPLICIT
		} else if lex.keyword("IMPLIED", lval) {
			return IMPLIED
		} else if lex.keyword("IMPORTS", lval) {
			return IMPORTS
		} else if lex.keyword("INCLUDES", lval) {
			return INCLUDES
		} else if lex.keyword("INSTANCE", lval) {
			return INSTANCE
		} else if lex.keyword("INSTRUCTIONS", lval) {
			return INSTRUCTIONS
		} else if lex.keyword("INTEGER", lval) {
			return INTEGER
		} else if lex.keyword("ISO646String", lval) {
			return ISO646String
		} else if lex.keyword("MAX", lval) {
			return MAX
		} else if lex.keyword("MIN", lval) {
			return MIN
		} else if lex.keyword("MINUS_INFINITY", lval) {
			return MINUS_INFINITY
		} else if lex.keyword("NULL", lval) {
			return NULL
		} else if lex.keyword("NumericString", lval) {
			return NumericString
		} else if lex.keyword("OBJECT", lval) {
			return OBJECT
		} else if lex.keyword("ObjectDescriptor", lval) {
			return ObjectDescriptor
		} else if lex.keyword("OCTET", lval) {
			return OCTET
		} else if lex.keyword("OF", lval) {
			return OF
		} else if lex.keyword("OPTIONAL", lval) {
			return OPTIONAL
		} else if lex.keyword("PATTERN", lval) {
			return PATTERN
		} else if lex.keyword("PDV", lval) {
			return PDV
		} else if lex.keyword("PLUS_INFINITY", lval) {
			return PLUS_INFINITY
		} else if lex.keyword("PRESENT", lval) {
			return PRESENT
		} else if lex.keyword("PrintableString", lval) {
			return PrintableString
		} else if lex.keyword("PRIVATE", lval) {
			return PRIVATE
		} else if lex.keyword("REAL", lval) {
			return REAL
		} else if lex.keyword("RELATIVE OID", lval) {
			return RELATIVE_OID
		} else if lex.keyword("SEQUENCE", lval) {
			return SEQUENCE
		} else if lex.keyword("SET", lval) {
			return SET
		} else if lex.keyword("SIZE", lval) {
			return SIZE
		} else if lex.keyword("STRING", lval) {
			return STRING
		} else if lex.keyword("SYNTAX", lval) {
			return SYNTAX
		} else if lex.keyword("T61String", lval) {
			return T61String
		} else if lex.keyword("TAGS", lval) {
			return TAGS
		} else if lex.keyword("TeletexString", lval) {
			return TeletexString
		} else if lex.keyword("TRUE", lval) {
			return TRUE
		} else if lex.keyword("TYPE_IDENTIFIER", lval) {
			return TYPE_IDENTIFIER
		} else if lex.keyword("UNIQUE", lval) {
			return UNIQUE
		} else if lex.keyword("UNIVERSAL", lval) {
			return UNIVERSAL
		} else if lex.keyword("UniversalString", lval) {
			return UniversalString
		} else if lex.keyword("UTCTime", lval) {
			return UTCTime
		} else if lex.keyword("UTF8String", lval) {
			return UTF8String
		} else if lex.keyword("VideotexString", lval) {
			return VideotexString
		} else if lex.keyword("VisibleString", lval) {
			return VisibleString
		} else if lex.keyword("WITH", lval) {
			return WITH
		} else if lex.keyword("UTF8_BOM", lval) {
			return UTF8_BOM
		} else if lex.keyword("::=", lval) {
			return ASSIGNMENT
		} else if lex.keyword("{", lval) {
			return '{'
		} else if lex.keyword("}", lval) {
			return '}'
		} else if lex.keyword("[", lval) {
			return '['
		} else if lex.keyword("]", lval) {
			return ']'
		} else if lex.keyword("(", lval) {
			return '('
		} else if lex.keyword(")", lval) {
			return ')'
		} else if lex.keyword(",", lval) {
			return ','
		} else if lex.keyword("...", lval) {
			return ELLIPSISS
		} else if lex.keyword("..", lval) {
			return ELLIPSIS
		} else if lex.keyword(";", lval) {
			return ';'
		} else if lex.regexp(numberReg, lval) {
			lval.Int64 = ParseInt(lval.str)
			return Number
		} else if lex.regexp(identifierReg, lval) {
			lval.String = String(lval.str)
			return Identifier
		} else if len(lex.text) > 0 {
			fmt.Println("还有内容未匹配成功", string(lex.text))
			return -1 //还有内容未匹配成功
		} else {
			return 0
		}
	}
}

func (lex *AsnLexer) Error(s string) {
	lex.err = errors.New(s)
	log.Println("Lex error: ", s)
}

func (lex *AsnLexer) Err() error {
	return lex.err
}

func (lex *AsnLexer) Result() *ModuleList {
	return lex.result
}

func (lex *AsnLexer) Text() string {
	return string(lex.text)
}

//对 text 进行正则匹配
func (lex *AsnLexer) regexp(reg *regexp.Regexp, lval *ParsedGrammarSymType) bool {
	if reg.Match(lex.text) {
		findText := reg.Find(lex.text)
		lex.text = lex.text[len(findText):]
		lval.str = string(findText)
		return true
	}
	return false
}

//对 text 进行关键字匹配
func (lex *AsnLexer) keyword(str string, lval *ParsedGrammarSymType) bool {
	kw := []byte(str)
	if len(lex.text) < len(kw) {
		return false
	}
	for i, b := range kw {
		if lex.text[i] != b {
			return false
		}
	}
	lex.text = lex.text[len(kw):]
	lval.str = str
	lval.String = String(str)
	return true
}

func ParseInt(str string) Int64 {
	if i, e := strconv.ParseInt(str, 10, 64); e != nil {
		return -999
	} else {
		return Int64(i)
	}
}
