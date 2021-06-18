package uper

import (
	"encoding/json"
	"fmt"
	"strings"
	"testing"

	"github.com/wrype/asn1go/common"
)

type load float64

var LoadMember common.Member

func init() {
	json.Unmarshal([]byte(`{
		"expr_type": "REAL",
        "meta_type": "AMT_TYPE",
        "constraints": {
            "type": "REAL_WITH_COMPONENTS",
            "min": 0,
            "max": 100,
            "base": 10,
            "exponent": -2
        },
        "tag": {
            "tag_mode": "TM_DEFAULT",
            "tag_class": "",
            "tag_value": ""
        },
        "Identifier": "Load"
            }`), &LoadMember)

}

func (s *load) UperEncode() *common.BitBuffer {
	return UperEncodeReal(LoadMember, float64(*s))
}

func (s *load) UperDecode(b *common.BitBuffer) {
	*s = load(UperDecodeReal(LoadMember, b))
}

func TestLoad(t *testing.T) {
	var s load = 123.456
	bs := s.UperEncode().Bytes()
	t.Logf("value load ::=123.456; data.uper is %2x\n", bs)
	str := strings.ToLower("0B033132333435362E452D33")
	if fmt.Sprintf("%2x", bs) == str {
		t.Log("Pass")
	} else {
		t.Fatal("load Uper Encode Error!")
	}
	var s2 load
	s2.UperDecode(common.NewBitBufferFromBytes(bs))
	t.Logf("buffer %2x; Decode load is %x\n", bs, s2)
	if s2 == 123.456 {
		t.Log("Pass")
	} else {
		t.Fatal("load Uper Decode Error!")
	}
}
