package uper

import (
	"encoding/hex"
	"encoding/json"
	"fmt"
	"testing"

	"github.com/wrype/asn1go/common"
)

type IA5Rang string

var IA5RangMember common.Member

type IA5Fixed string

var IA5FixedMember common.Member

type DescriptiveName string

/**
DescriptiveName ::= IA5String (SIZE(1..63))
*/
var DescriptiveNameMember common.Member

func init() {
	json.Unmarshal([]byte(`{
                "expr_type": "IA5String",
                "meta_type": "AMT_TYPE",
                "constraints": {
                    "type": "SIZE_RANG",
                    "min": 1,
                    "max": 64
                },
                "Identifier": "textString"
			}`), &IA5RangMember)
	json.Unmarshal([]byte(`{
				"expr_type": "IA5String",
				"meta_type": "AMT_TYPE",
				"constraints": {
					"type": "SIZE_FIXED",
					"value": 10
				},
				"Identifier": "name"
				}`), &IA5FixedMember)

	def := `{"expr_type":"IA5String","meta_type":"AMT_TYPE","constraints":{"type":"SIZE_RANG","value":0,"min":1,"max":63,"base":0,"exponent":0},"Identifier":"DescriptiveName","value":0}`
	json.Unmarshal([]byte(def), &DescriptiveNameMember)
}

func (s *IA5Rang) UperEncode() *common.BitBuffer {
	return UperEncodeIA5String(IA5RangMember, string(*s))
}

func (s *IA5Rang) UperDecode(b *common.BitBuffer) {
	*s = IA5Rang(UperDecodeIA5String(IA5RangMember, b))
}

func (s *IA5Fixed) UperEncode() *common.BitBuffer {
	return UperEncodeIA5String(IA5FixedMember, string(*s))
}

func (s *IA5Fixed) UperDecode(b *common.BitBuffer) {
	*s = IA5Fixed(UperDecodeIA5String(IA5FixedMember, b))
}

func (s *DescriptiveName) UperEncode() *common.BitBuffer {
	return UperEncodeIA5String(DescriptiveNameMember, string(*s))
}

func (s *DescriptiveName) UperDecode(b *common.BitBuffer) {
	*s = DescriptiveName(UperDecodeIA5String(DescriptiveNameMember, b))
}

func TestIA5_SIZE_RANG(t *testing.T) {
	var s IA5Rang = "'afsef3';'"
	bs := s.UperEncode().Bytes()
	t.Logf("value IA5Rang ::= \"%s\"; data.uper is %2x\n", s, bs)
	if fmt.Sprintf("%2X", bs) == "253E1CDCF2E6669DDA70" {
		t.Log("Pass")
	} else {
		t.Fatal("IA5Rang Uper Encode Error!")
	}

	var s2 IA5Rang
	s2.UperDecode(common.NewBitBufferFromBytes(bs))
	t.Logf("buffer %2x; Decode IA5Rang is \"%s\"\n", bs, s2)
	if s2 == s {
		t.Log("Pass")
	} else {
		t.Fatal("IA5Rang Uper Decode Error!")
	}

	//rec1value DescriptiveName ::= "link0"
	var s3 DescriptiveName = "link0"
	ab, ae := s3.UperEncode().Result()
	if fmt.Sprintf("%2X", ab) != "13669DDAD800" || ae != nil {
		t.Fatalf("DescriptiveName Encode error buf: %2X  err: %v", ab, ae)
	} else {
		t.Log("Pass")
	}

	var s4 DescriptiveName
	bufs4, _ := hex.DecodeString("13669DDAD800")
	s4b := common.NewBitBufferFromBytes(bufs4)
	s4.UperDecode(s4b)
	if string(s4) != "link0" || s4b.Error() != nil {
		t.Fatalf("DescriptiveName Encode error str:%s err:%v", s4, s4b.Error())
	} else {
		t.Log("Pass")
	}
}

func TestIA5_SIZE_FIXED(t *testing.T) {
	var s IA5Fixed = "'afsef3';'"
	bs := s.UperEncode().Bytes()
	t.Logf("value IA5Fixed ::= \"%s\"; data.uper is %2x\n", s, bs)
	if fmt.Sprintf("%2X", bs) == "4F87373CB999A7769C" {
		t.Log("Pass")
	} else {
		t.Fatal("IA5Fixed Uper Encode Error!")
	}

	var s2 IA5Fixed
	s2.UperDecode(common.NewBitBufferFromBytes(bs))
	t.Logf("buffer %2x; Decode IA5Fixed is \"%s\"\n", bs, s2)
	if s2 == s {
		t.Log("Pass")
	} else {
		t.Fatal("IA5Fixed Uper Decode Error!")
	}
}
