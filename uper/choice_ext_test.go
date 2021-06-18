package uper

import (
	"encoding/hex"
	"encoding/json"
	"fmt"
	"testing"

	"github.com/wrype/asn1go/common"
)

/**
Description ::= CHOICE{
		id INTEGER (0..7),
        val BOOLEAN,
        ...,
        aa INTEGER(0..7),
        bb INTEGER(0..7),
        ...
	}
*/
var DespMember common.Member

func init() {
	def := `{"members":[{"expr_type":"INTEGER","meta_type":"AMT_TYPE","constraints":{"type":"SIZE_RANG","value":0,"min":0,"max":7,"base":0,"exponent":0},"Identifier":"Id","value":0},{"expr_type":"BOOLEAN","meta_type":"AMT_TYPE","Identifier":"Val","value":0},{"expr_type":"A1TC_EXTENSIBLE","meta_type":"AMT_TYPE","Identifier":"...","value":0},{"expr_type":"INTEGER","meta_type":"AMT_TYPE","constraints":{"type":"SIZE_RANG","value":0,"min":0,"max":7,"base":0,"exponent":0},"Identifier":"Aa","value":0},{"expr_type":"INTEGER","meta_type":"AMT_TYPE","constraints":{"type":"SIZE_RANG","value":0,"min":0,"max":7,"base":0,"exponent":0},"Identifier":"Bb","value":0},{"expr_type":"A1TC_EXTENSIBLE","meta_type":"AMT_TYPE","Identifier":"...","value":0}],"expr_type":"CHOICE","meta_type":"AMT_TYPE","Identifier":"Desp","value":0}`
	json.Unmarshal([]byte(def), &DespMember)
}

type Desp struct {
	Id  *int64 `json:"id,omitempty"`
	Val *bool  `json:"val,omitempty"`
	Aa  *int64 `json:"aa,omitempty"`
	Bb  *int64 `json:"bb,omitempty"`
}

func (t *Desp) UperEncode() *common.BitBuffer {
	return UperEncodeChoice(DespMember, t)
}

func (t *Desp) UperDecode(buf *common.BitBuffer) {
	if buf.Error() != nil {
		return
	}
	UperDecodeChoice(DespMember, buf, t)
}

func TestUperEncodeChoiceExt(t *testing.T) {
	bb := int64(7)
	desp := Desp{
		Bb: &bb,
	}
	buf, err := desp.UperEncode().Result()
	if fmt.Sprintf("%2X", buf) != "8101E0" || err != nil {
		t.Fatalf("TestUperEncodeChoiceExt Error: %v %2X != 8101E0", err, buf)
	} else {
		t.Log("Pass")
	}

	desp = Desp{
		Aa: &bb,
	}
	buf, err = desp.UperEncode().Result()
	if fmt.Sprintf("%2X", buf) != "8001E0" || err != nil {
		t.Fatalf("TestUperEncodeChoiceExt Error: %v %2X != 8001E0", err, buf)
	} else {
		t.Log("Pass")
	}
}

func TestUperDecodeChoiceExt(t *testing.T) {
	desp := &Desp{}
	buf, _ := hex.DecodeString("8101E0")
	bitBuf := common.NewBitBufferFromBytes(buf)
	desp.UperDecode(bitBuf)
	if desp.Bb == nil || *desp.Bb != 7 || bitBuf.Error() != nil {
		t.Fatalf("TestUperDecodeChoiceExt Error:%v %+v", bitBuf.Error(), desp.Bb)
	} else {
		t.Log("Pass")
	}

	desp = &Desp{}
	buf, _ = hex.DecodeString("8001E0")
	bitBuf = common.NewBitBufferFromBytes(buf)
	desp.UperDecode(bitBuf)
	if desp.Aa == nil || *desp.Aa != 7 || bitBuf.Error() != nil {
		t.Fatalf("TestUperDecodeChoiceExt Error:%v %+v", bitBuf.Error(), desp.Aa)
	} else {
		t.Log("Pass")
	}

	/**
	定义中不存在的扩展成员的情况
	Description ::= CHOICE{
			id INTEGER (0..7),
			val BOOLEAN,
			...,
			aa INTEGER(0..7),
			bb INTEGER(0..7),
			gg INTEGER(0..7),
			...
	}
	{
	  "gg":7
	}
	*/
	desp = &Desp{}
	buf, _ = hex.DecodeString("8201E0")
	bitBuf = common.NewBitBufferFromBytes(buf)
	desp.UperDecode(bitBuf)
	if desp.Aa != nil || desp.Bb != nil || desp.Val != nil || desp.Id != nil || bitBuf.Error() != nil {
		t.Fatalf("TestUperDecodeChoiceExt Error:%v %+v", bitBuf.Error(), desp)
	} else {
		t.Log("Pass")
	}
}
