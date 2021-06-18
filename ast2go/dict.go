package ast2go

import (
	"encoding/csv"
	"log"
	"os"
	"strings"
)

type Dict struct {
	From string
	To   string
}

var gDict = []Dict{
	{From: "-", To: "_"},
	{From: "request", To: "req"},
	{From: "response", To: "resp"},
}

func GetVarMappingDict() [][]string {
	ret := make([][]string, 0)
	for _, elem := range gDict {
		ret = append(ret, []string{elem.From, elem.To})
	}
	return ret
}

func LoadDict(filename string) error {
	file, err := os.Open(filename)
	if err != nil {
		return err
	}
	defer file.Close()
	reader := csv.NewReader(file)
	data, err := reader.ReadAll()
	if err != nil {
		return err
	}
	gDict = make([]Dict, 0)
	for _, record := range data {
		gDict = append(gDict, Dict{
			From: record[0],
			To:   record[1],
		})
	}
	log.Printf("load dict: %+v\n", gDict)
	return nil
}

func getRecommendName(s string) string {
	s = strings.Title(s)
	for _, d := range gDict {
		s = strings.ReplaceAll(s, d.From, d.To)
	}
	s = Case2Camel(s, "")
	return s
}

//把名称转为驼峰形式
func Case2Camel(name string, trim string) string {
	name = strings.Title(name)
	name = strings.Replace(name, "_", "", -1)
	if trim != "" {
		name = strings.Replace(name, trim, "", -1)
	}
	return name
}
