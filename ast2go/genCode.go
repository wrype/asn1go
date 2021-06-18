package ast2go

import (
	"bytes"
	"fmt"
	"go/format"
	"log"
	"os"
	"text/template"

	"github.com/wrype/asn1go/common"
)

func genCode(ctx *xctx, member *common.Member, buffer *bytes.Buffer) {
	if member.MetaType == "AMT_VALUE_ASSIGMENT" {
		ident, xtype := parseType(ctx, member)
		member.Identifier = ident
		member.ExprType = common.ExprType(xtype)
		return
	}

	ident, xtype := parseType(ctx, member)
	member.Identifier = ident
	fmt.Fprintf(buffer, "type %s %s\n", ident, xtype)
	if ctx.detectConst { //处理常量
		constTmpl, _ := template.New("constCode").Parse(tmplConst)
		constTmpl.Execute(buffer, member)
	}

	coder := template.FuncMap{
		"GetUperEncoder": getUperEncoder,
		"GetUperDecoder": getUperDecoder,
		"Case2Camel":     Case2Camel,
	}
	coderNormal := template.New("coderNormal")
	coderNormal.Funcs(coder)
	coderNormal, _ = coderNormal.Parse(tmplCoderNormal)

	coderReflect := template.New("coderReflect")
	coderReflect.Funcs(coder)
	coderReflect, _ = coderReflect.Parse(tmplCoderReflect)

	switch member.MetaType {
	case "AMT_TYPE":
		switch member.ExprType {
		case "INTEGER", "BOOLEAN", "REAL", "IA5String", "OCTET_STRING", "BIT_STRING":
			if err := coderNormal.Execute(buffer, struct {
				*common.Member
				PackType string
			}{member, xtype}); err != nil {
				log.Fatal(err)
			}
		case "ENUMERATED":
			if err := coderNormal.Execute(buffer, struct {
				*common.Member
				PackType             string
				OriginEnumValMapping map[int]string
			}{member, xtype, ctx.originEnumValMapping}); err != nil {
				log.Fatal(err)
			}
		case "CHOICE", "SEQUENCE", "SEQUENCE_OF":
			if err := coderReflect.Execute(buffer, member); err != nil {
				log.Fatal(err)
			}
		default:
			log.Printf("unsupported AMT_TYPE: %v for gen coder", member.ExprType)
		}
	case "ASN_TYPE":
	}
}

func GenPkg(mod *common.Module) {
	os.RemoveAll(mod.ModuleName)
	os.Mkdir(mod.ModuleName, 0666)
	oldWd, _ := os.Getwd()
	os.Chdir(mod.ModuleName)

	normMember := make([]common.Member, 0)
	gValMember := make([]common.Member, 0)
	func() {
		coderFile, err := os.OpenFile(mod.ModuleName+".go", os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0666)
		if err != nil {
			log.Fatal(err)
		}
		defer coderFile.Close()

		fmt.Fprintf(coderFile, coderHerder, mod.ModuleName)
		var buffer bytes.Buffer
		recorder := newMaskTypeRecorder()
		for idx := range mod.Members {
			ctx := xctx{
				recorder: recorder,
			}
			genCode(&ctx, &mod.Members[idx], &buffer)
			if ctx.detectGlobalConst {
				gValMember = append(gValMember, mod.Members[idx])
			} else {
				normMember = append(normMember, mod.Members[idx])
			}
		}

		genMaskModeCode(&buffer, *recorder)

		data, err := format.Source(buffer.Bytes())
		if err != nil {
			log.Fatalln(err)
		}
		coderFile.Write(data)
	}()

	mod.Members = normMember
	genInitCode(mod)

	mod.Members = gValMember
	genGlobalConstCode(mod)

	os.Chdir(oldWd)
}
