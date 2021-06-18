package ast2go

import (
	"fmt"
	"log"

	"github.com/wrype/asn1go/common"
)

const (
	jsonTag    = "`json:\"%s\"`\n"
	jsonTagOpt = "`json:\"%s,omitempty\"`\n"
)

type xctx struct {
	detectConst          bool
	detectGlobalConst    bool
	detectExtFlag        bool
	hitSlice             bool
	curIdent             string
	originEnumValMapping map[int]string
	recorder             *maskTypeRecorder
}

/**
 * @return: 标识符、类型。当处理AMT_VALUE时，返回标识符、""
 */
func parseType(ctx *xctx, member *common.Member) (string, string) {
	ident := getRecommendName(member.Identifier)
	var xtype string
	switch member.MetaType {
	case "AMT_TYPE":
		switch member.ExprType {
		case "INTEGER":
			xtype = "int64"
		case "BOOLEAN":
			xtype = "bool"
		case "CHOICE":
			xtype = "struct{\n"
			maskTypeList := make([]maskType, 0)
			for idx, field := range member.Members {
				ctx.hitSlice = false
				elemIdent, elemType := parseType(ctx, &member.Members[idx])
				if len(elemIdent) == 0 || len(elemType) == 0 {
					continue
				}
				maskTypeList = append(maskTypeList, maskType{
					ElemIdent: elemIdent,
					Type:      elemType,
				})
				tag := fmt.Sprintf(jsonTagOpt, field.Identifier)
				if !ctx.hitSlice {
					elemType = "*" + elemType
				}
				xtype = xtype + elemIdent + " " + elemType + tag
				member.Members[idx].Identifier = elemIdent
			}
			ctx.recorder.choiceMaskTbl[ident] = maskTypeList
			xtype += "}"
		case "ENUMERATED":
			xtype = "int64"
			ctx.curIdent = ident
			ctx.detectConst = true
			ctx.originEnumValMapping = make(map[int]string)
			for idx, mem := range member.Members {
				elemIdent, _ := parseType(ctx, &member.Members[idx])
				if len(elemIdent) == 0 {
					continue
				}
				ctx.originEnumValMapping[mem.Value] = mem.Identifier
				member.Members[idx].Identifier = elemIdent
			}
			ctx.recorder.maskTypeList[ident] = struct{}{}
		case "SEQUENCE":
			xtype = "struct{\n"
			for idx, field := range member.Members {
				ctx.hitSlice = false
				elemIdent, elemType := parseType(ctx, &member.Members[idx])
				if len(elemIdent) == 0 || len(elemType) == 0 {
					continue
				}
				var tag string
				if field.Marker.Flags == common.EM_OPTIONAL || ctx.detectExtFlag {
					if !ctx.hitSlice {
						elemType = "*" + elemType
					}
					tag = fmt.Sprintf(jsonTagOpt, field.Identifier)
				} else {
					tag = fmt.Sprintf(jsonTag, field.Identifier)
				}
				xtype = xtype + elemIdent + " " + elemType + tag
				member.Members[idx].Identifier = elemIdent
			}
			xtype += "}"
		case "IA5String":
			xtype = "string"
		case "OCTET_STRING":
			xtype = "[]byte"
			ctx.hitSlice = true
		case "BIT_STRING":
			xtype = "uint64"
			ctx.curIdent = ident
			ctx.detectConst = true
			for idx := range member.Members {
				elemIdent, _ := parseType(ctx, &member.Members[idx])
				if len(elemIdent) == 0 {
					continue
				}
				member.Members[idx].Identifier = elemIdent
			}
			ctx.recorder.maskTypeList[ident] = struct{}{}
		case "REAL":
			xtype = "float64"
		case "SEQUENCE_OF":
			_, tmpT := parseType(ctx, &member.Members[0])
			xtype = "[]" + tmpT
			ctx.hitSlice = true
		case "A1TC_EXTENSIBLE": //"..."
			ident = ""
			xtype = ""
			ctx.detectExtFlag = !ctx.detectExtFlag
		default:
			log.Fatal("unsupported AMT_TYPE: ", member.ExprType)
		}
	case "ASN_TYPE":
		xtype = getRecommendName(string(member.ExprType))
	case "AMT_VALUE":
		ident = ctx.curIdent + "_" + ident
		ident = getRecommendName(string(ident))
	case "AMT_VALUE_ASSIGMENT":
		_, xtype = parseType(ctx, &member.Members[0])
		ctx.detectGlobalConst = true
	}
	return ident, xtype
}
