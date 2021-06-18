package main

import (
	"encoding/csv"
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"os"

	"github.com/wrype/asn1go/asnparser"
	"github.com/wrype/asn1go/ast2go"
	"github.com/wrype/asn1go/common"
)

type CmdFlag struct {
	help           bool
	asnFile        string
	outFile        string
	dictFile       string
	printInnerDict bool
}

var (
	Version    string
	CommitTag  string
	CommitTime string
	GoVersion  string
)

var cmdFlag CmdFlag

func init() {
	cmdFlag = CmdFlag{}
	flag.BoolVar(&cmdFlag.help, "h", false, "this help")
	flag.StringVar(&cmdFlag.asnFile, "f", "", "input a ASN file")
	flag.StringVar(&cmdFlag.outFile, "o", "go", "[go|xxx.json|stdout|-] output go codes or ast json")
	flag.StringVar(&cmdFlag.dictFile, "l", "", "load csv for variable mapping; case sensitive and will mapping in order of lines in csv")
	flag.BoolVar(&cmdFlag.printInnerDict, "p", false, "print inner dict of variable mapping, field: from,to")
}

func main() {
	flag.Parse()
	if cmdFlag.help {
		printVersion()
		flag.Usage()
		return
	}

	if cmdFlag.printInnerDict {
		dict := ast2go.GetVarMappingDict()
		writer := csv.NewWriter(os.Stdout)
		writer.WriteAll(dict)
		writer.Flush()
		return
	}

	if cmdFlag.asnFile == "" {
		printVersion()
		flag.Usage()
		return
	}

	if len(cmdFlag.dictFile) > 0 {
		if err := ast2go.LoadDict(cmdFlag.dictFile); err != nil {
			log.Panic(err)
		}
	}

	buf, err := ioutil.ReadFile(cmdFlag.asnFile)
	if err != nil {
		log.Panic(err)
	}
	lex := asnparser.NewAsnLexer(buf)
	asnparser.ParsedGrammarParse(lex)
	if lex.Err() != nil {
		log.Panic(lex.Err(), lex.Text())
	}
	if lex.Result() != nil {
		j, _ := json.Marshal(lex.Result())
		switch cmdFlag.outFile {
		case "go":
			ms := *lex.Result()
			for _, m := range ms {
				ast2go.GenPkg(ConvertModule(m))
			}
		case "-", "", "stdout":
			log.Println(string(j))
		default:
			ioutil.WriteFile(cmdFlag.outFile, j, 0644)
		}
	} else {
		panic("生成失败")
	}
}

func printVersion() {
	fmt.Println(fmt.Sprintf(`
Version:	%s
CommitTag:	%s
CommitTime:	%s
GoVersion:	%s
`, Version, CommitTag, CommitTime, GoVersion))
}

func ConvertModule(module asnparser.Module) *common.Module {
	var m common.Module
	buf, _ := json.Marshal(module)
	json.Unmarshal(buf, &m)
	return &m
}
