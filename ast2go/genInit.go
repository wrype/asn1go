package ast2go

import (
	"bytes"
	"fmt"
	"log"
	"os"
	"text/template"

	"github.com/wrype/asn1go/common"
	"github.com/wrype/go-render/render"
)

func genInitCode(mod *common.Module) {
	f, err := os.OpenFile("init.go", os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0666)
	if err != nil {
		log.Fatal(err)
	}
	defer f.Close()

	tmpl := template.New("initCode")
	tmpl.Funcs(template.FuncMap{
		"ToVar": func(v interface{}) string {
			return render.AsCode(v)
		},
	})
	tmpl, _ = tmpl.Parse(fmt.Sprintf(initTmpl))

	var buf bytes.Buffer
	err = tmpl.Execute(&buf, mod)
	if err != nil {
		log.Fatal(err)
	}
	code, _ := gShortener.Shorten(buf.Bytes())
	f.Write(code)
}
