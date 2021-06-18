package ast2go

import (
	"bytes"
	"go/format"
	"log"
	"os"
	"text/template"

	"github.com/wrype/asn1go/common"
)

func genGlobalConstCode(mod *common.Module) {
	if mod == nil || len(mod.Members) == 0 {
		return
	}

	f, err := os.OpenFile("const.go", os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0666)
	if err != nil {
		log.Fatal(err)
	}
	defer f.Close()

	tmpl := template.New("globalConstCode")
	tmpl, _ = tmpl.Parse(tmplGlobalConst)

	var buf bytes.Buffer
	err = tmpl.Execute(&buf, mod)
	if err != nil {
		log.Fatal(err)
	}
	code, _ := format.Source(buf.Bytes())
	f.Write(code)
}
