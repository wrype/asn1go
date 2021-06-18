package uper

import (
	"encoding/json"
	"fmt"
	"testing"

	"github.com/wrype/asn1go/common"
)

type Mode bool

var ModeMember common.Member

func init() {
	json.Unmarshal([]byte(`{
                "expr_type": "BOOLEAN",
                "meta_type": "AMT_TYPE",
                "constraints": null,
                "Identifier": "Mode"
            }`), &ModeMember)
}

func (m *Mode) UperEncode() *common.BitBuffer {
	return UperEncodeBoolean(ModeMember, bool(*m))
}

func (m *Mode) UperDecode(b *common.BitBuffer) bool {
	*m = Mode(UperDecodeBoolean(ModeMember, b))
	return bool(*m)
}

func TestBOOLEAN(t *testing.T) {
	var m Mode = true
	bs := m.UperEncode().Bytes()
	t.Logf("value Mode ::= true; data.uper is %2x\n", bs)
	if fmt.Sprintf("%2x", bs) == "80" {
		t.Log("Pass")
	} else {
		t.Fatal("Mode Uper Encode Error!")
	}

	var m2 Mode
	m2.UperDecode(common.NewBitBufferFromBytes(bs))
	t.Logf("buffer %2x; Decode Mode is %t\n", bs, m2)
	if m2 == true {
		t.Log("Pass")
	} else {
		t.Fatal("Mode Uper Decode Error!")
	}
}
