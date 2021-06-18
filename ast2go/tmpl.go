package ast2go

import (
	"log"

	"github.com/wrype/asn1go/common"

	"github.com/wrype/golines"
)

var gShortener = golines.NewShortener(golines.ShortenerConfig{
	KeepAnnotations:  false,
	MaxLen:           100,
	ReformatTags:     true,
	ShortenComments:  false,
	TabLen:           4,
	BaseFormatterCmd: "gofmt",
})

var initTmpl string = `package {{.ModuleName}}

import (
	"github.com/wrype/asn1go/common"
)

var(
{{range .Members}}
gAst{{.Identifier}}={{ ToVar . }}
{{end}}
)

func init(){
}
`

var coderHerder string = `package %s

import (
	"encoding/json"
	"github.com/wrype/asn1go/common"
	"github.com/wrype/asn1go/uper"
	"fmt"
)
`

func getUperEncoder(exp common.ExprType) string {
	switch exp {
	case common.AMT_TYPE_INTEGER:
		return "UperEncodeInteger"
	case common.AMT_TYPE_BOOLEAN:
		return "UperEncodeBoolean"
	case common.AMT_TYPE_CHOICE:
		return "UperEncodeChoice"
	case common.AMT_TYPE_ENUMERATED:
		return "UperEncodeEnumerated"
	case common.AMT_TYPE_SEQUENCE:
		return "UperEncodeSequence"
	case common.AMT_TYPE_IA5String:
		return "UperEncodeIA5String"
	case common.AMT_TYPE_OCTET_STRING:
		return "UperEncodeOctetString"
	case common.AMT_TYPE_BIT_STRING:
		return "UperEncodeBitString"
	case common.AMT_TYPE_REAL:
		return "UperEncodeReal"
	case common.AMT_TYPE_SEQUENCE_OF:
		return "UperEncodeSequenceOf"
	default:
		log.Panicf("unsupported AMT_TYPE: %v for gen coder", exp)
	}
	return ""
}

func getUperDecoder(exp common.ExprType) string {
	switch exp {
	case common.AMT_TYPE_INTEGER:
		return "UperDecodeInteger"
	case common.AMT_TYPE_BOOLEAN:
		return "UperDecodeBoolean"
	case common.AMT_TYPE_CHOICE:
		return "UperDecodeChoice"
	case common.AMT_TYPE_ENUMERATED:
		return "UperDecodeEnumerated"
	case common.AMT_TYPE_SEQUENCE:
		return "UperDecodeSequence"
	case common.AMT_TYPE_IA5String:
		return "UperDecodeIA5String"
	case common.AMT_TYPE_OCTET_STRING:
		return "UperDecodeOctetString"
	case common.AMT_TYPE_BIT_STRING:
		return "UperDecodeBitString"
	case common.AMT_TYPE_REAL:
		return "UperDecodeReal"
	case common.AMT_TYPE_SEQUENCE_OF:
		return "UperDecodeSequenceOf"
	default:
		log.Panicf("unsupported AMT_TYPE: %v for gen coder", exp)
	}
	return ""
}

var tmplCoderNormal = `{{$Ident:=.Identifier}} {{$ExprType:=.ExprType}}
func (t *{{$Ident}}) UperEncode() *common.BitBuffer {
	return uper.{{GetUperEncoder $ExprType}}(gAst{{$Ident}}, {{.PackType}}(*t))
}

func (t *{{$Ident}}) UperDecode(buf *common.BitBuffer) {
	if buf.Error() != nil {
		return
	}
	*t = {{$Ident}}(uper.{{GetUperDecoder $ExprType}}(gAst{{$Ident}}, buf))
}

{{if eq $ExprType "ENUMERATED"}} 
func (t *{{$Ident}}) MarshalJSON() ([]byte, error) {
	switch *t { 
	{{- range $i,$v := .Members -}}
		{{if not $v.IsExtLabel}} 
	case {{$v.Identifier}}:
		return []byte("\"{{index $.OriginEnumValMapping $v.Value}}\""), nil
		{{- end -}}	{{end}}
	default:
		return []byte("null"), nil
	}
}

func (t *{{$Ident}}) UnmarshalJSON(b []byte) error {
	switch string(b) {
	{{- range $i,$v := .Members -}}
		{{if not $v.IsExtLabel}}
	case "\"{{index $.OriginEnumValMapping $v.Value}}\"":
		*t = {{$v.Identifier}}
		{{- end -}}	{{end}}
	default:
		return fmt.Errorf("{{$Ident}}: unknown value set")
	}
	return nil
}
	{{range $i,$v := .Members}}
		{{if not $v.IsExtLabel}}
			{{$ElemIdent:=$v.Identifier}}
func (t *{{$Ident}}) Is{{Case2Camel $ElemIdent $Ident}}() bool {
	return *t == {{$ElemIdent}}
}
		{{end}}
	{{end}}
{{end}}

{{if eq $ExprType "BIT_STRING"}}
	{{range $i,$v := .Members}}
		{{if not $v.IsExtLabel}}
			{{$ElemIdent:=$v.Identifier}}
func (t *{{$Ident}}) Has{{Case2Camel $ElemIdent $Ident}}() bool {
	return *t&{{$ElemIdent}} == {{$ElemIdent}}
}
		{{end}}
	{{end}}
{{end}}

func (t *{{$Ident}}) ToString() string {
	data, _:=json.Marshal(t)
	return string(data)
}
`

var tmplCoderReflect = `{{$Ident := .Identifier}} {{$ExprType:=.ExprType}}
func (t *{{$Ident}}) UperEncode() *common.BitBuffer {
	return uper.{{GetUperEncoder $ExprType}}(gAst{{$Ident}}, t)
}

func (t *{{$Ident}}) UperDecode(buf *common.BitBuffer) {
	if buf.Error() != nil {
		return
	}
	uper.{{GetUperDecoder $ExprType}}(gAst{{$Ident}}, buf, t)
}

{{if eq $ExprType "CHOICE"}}
	{{range $i,$v := .Members}}
		{{if not $v.IsExtLabel}}
			{{$ElemIdent:=$v.Identifier}}
func (t *{{$Ident}}) Is{{Case2Camel $ElemIdent $Ident}}() bool {
	return t.{{$ElemIdent}} != nil
}
		{{end}}
	{{end}}
{{end}}

func (t *{{$Ident}}) ToString() string {
	data, _:=json.Marshal(t)
	return string(data)
}
`

var tmplConst = `{{if eq .ExprType "ENUMERATED"}}const (
	{{range .Members}} {{if eq .MetaType "AMT_VALUE"}} {{.Identifier}} {{$.Identifier}} = {{.Value}}
	{{end}}{{end}})
{{else if eq .ExprType "BIT_STRING"}}const (
	{{range .Members}} {{if eq .MetaType "AMT_VALUE"}} {{.Identifier}} {{$.Identifier}} = 1 << {{.Value}}
	{{end}}{{end}})
{{end}}
`

var tmplGlobalConst string = `package {{.ModuleName}}

const (
	{{range .Members}} {{.Identifier}} {{.ExprType}} = {{.Constraints.Value}}
	{{end}}
)
`

var tmplMaskMode string = `
{{range $Ident,$ChoiceElem:=.}}
/**
 * set numeric value in choice by elem index
 */
func (t *{{$Ident}}) SetMask(elemIdx int, val uint64) error {
	switch elemIdx{ {{range $idx,$Elem:=$ChoiceElem}}
	case {{$idx}}:
		*t.{{$Elem.ElemIdent}} = {{$Elem.Type}}(val) {{end}}
	default:
		return fmt.Errorf("{{$Ident}}: elem index out of range")
	}
	return nil
} {{end}}
`
