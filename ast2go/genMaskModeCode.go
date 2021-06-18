package ast2go

import (
	"bytes"
	"text/template"
)

type identifier = string
type maskTypeIdent = string
type maskType struct {
	ElemIdent identifier
	Type      maskTypeIdent
}
type maskTable map[identifier][]maskType
type maskTypeCache map[maskTypeIdent]struct{}
type maskTypeRecorder struct {
	choiceMaskTbl maskTable
	maskTypeList  maskTypeCache
}

func newMaskTypeRecorder() *maskTypeRecorder {
	return &maskTypeRecorder{
		choiceMaskTbl: make(maskTable),
		maskTypeList:  make(maskTypeCache),
	}
}

func genMaskModeCode(buffer *bytes.Buffer, recorder maskTypeRecorder) {
	enableMaskMode := func(tbl []maskType, maskTypeList maskTypeCache) bool {
		if len(tbl) == 0 {
			return false
		}
		for _, elem := range tbl {
			if _, ok := maskTypeList[elem.Type]; !ok {
				return false
			}
		}
		return true
	}
	for ident, elemList := range recorder.choiceMaskTbl {
		if !enableMaskMode(elemList, recorder.maskTypeList) {
			delete(recorder.choiceMaskTbl, ident)
		}
	}
	if len(recorder.choiceMaskTbl) == 0 {
		return
	}
	coderMask, _ := template.New("coderMaskMode").Parse(tmplMaskMode)
	coderMask.Execute(buffer, recorder.choiceMaskTbl)
}
