package uper

import (
	"encoding/hex"
	"encoding/json"
	"fmt"
	"testing"

	"github.com/wrype/asn1go/common"
)

/**
SeqExt ::= SEQUENCE
{
   range     INTEGER(1..32,...) ,
   ...
}

SeqExtOf ::= SEQUENCE (SIZE(1..8)) OF SeqExt
*/
var SeqExtOfMember common.Member

/**
  SEQ ::= SEQUENCE
  {
     range     INTEGER(1..32,...),
	 ss		   SS OPTIONAL
  }

  SS ::= INTEGER(1..32)

  SeqExtFixed ::= SEQUENCE (SIZE(2)) OF SEQ
*/
var SeqOfFixedMember common.Member

func init() {
	json.Unmarshal([]byte(`{
                "constraints": {
                    "type": "SIZE_RANG",
                    "min": 1,
                    "max": 8
                },
                "expr_type": "SEQUENCE_OF",
                "meta_type": "AMT_TYPE",
                "members": [
                    {
                        "expr_type": "SeqExt",
                        "meta_type": "ASN_TYPE",
                        "Identifier": "",
                        "tag": {
                            "tag_mode": "TM_DEFAULT",
                            "tag_class": "",
                            "tag_value": ""
                        }
                    }
                ],
                "tag": {
                    "tag_mode": "TM_DEFAULT",
                    "tag_class": "",
                    "tag_value": ""
                },
                "Identifier": "SeqExtOf"
            }`), &SeqExtOfMember)

	json.Unmarshal([]byte(`{
                "constraints": {
                    "type": "SIZE_FIXED",
                    "value": 2
                },
                "expr_type": "SEQUENCE_OF",
                "meta_type": "AMT_TYPE",
                "members": [
                    {
                        "expr_type": "SEQ",
                        "meta_type": "ASN_TYPE",
                        "Identifier": "",
                        "tag": {
                            "tag_mode": "TM_DEFAULT",
                            "tag_class": "",
                            "tag_value": ""
                        }
                    }
                ],
                "tag": {
                    "tag_mode": "TM_DEFAULT",
                    "tag_class": "",
                    "tag_value": ""
                },
                "Identifier": "SeqExtFixed"
            }`), &SeqOfFixedMember)
}

type SeqExtList []SeqExt

func (s *SeqExtList) UperEncode() *common.BitBuffer {
	return UperEncodeSequenceOf(SeqExtOfMember, s)
}

func (s *SeqExtList) UperDecode(b *common.BitBuffer) {
	UperDecodeSequenceOf(SeqExtOfMember, b, s)
}

type SeqOfFixedList []SEQ

func (s *SeqOfFixedList) UperEncode() *common.BitBuffer {
	return UperEncodeSequenceOf(SeqOfFixedMember, s)
}

func (s *SeqOfFixedList) UperDecode(b *common.BitBuffer) {
	UperDecodeSequenceOf(SeqOfFixedMember, b, s)
}

func TestSequenceOfEncode(t *testing.T) {
	as := SeqExtList{{Range: 10}, {Range: 12}}
	asb, ase := as.UperEncode().Result()
	if fmt.Sprintf("%2X", asb) != "224580" || ase != nil {
		t.Fatalf("TestSequenceOfEncode Error: %2x err:%v\n", asb, ase)
	} else {
		t.Log("Pass")
	}

	ss := SS(22)
	bs := SeqOfFixedList{{Range: 20}, {Range: 11, Ss: &ss}}
	bsb, bse := bs.UperEncode().Result()
	if fmt.Sprintf("%2X", bsb) != "272AA0" || bse != nil {
		t.Fatalf("TestSequenceOfEncode Error: %2x err:%v\n", bsb, bse)
	} else {
		t.Log("Pass")
	}
}

func TestSequenceOfDecode(t *testing.T) {
	as := &SeqExtList{}
	buf, _ := hex.DecodeString("224580")
	bitBuffer := common.NewBitBufferFromBytes(buf)
	as.UperDecode(bitBuffer)
	if bitBuffer.Error() == nil && len(*as) == 2 && (*as)[0].Range == 10 && (*as)[1].Range == 12 {
		t.Log("Pass")
	} else {
		t.Fatalf("TestSequenceOfDecode Error: %+v err:%v\n", *as, bitBuffer.Error())
	}

	bs := &SeqOfFixedList{}
	buf2, _ := hex.DecodeString("272AA0")
	bitBuffer2 := common.NewBitBufferFromBytes(buf2)
	bs.UperDecode(bitBuffer2)
	if bitBuffer2.Error() == nil && len(*bs) == 2 && (*bs)[0].Range == 20 && (*bs)[0].Ss == nil && (*bs)[1].Range == 11 && (*bs)[1].Ss != nil && *(*bs)[1].Ss == 22 {
		t.Log("Pass")
	} else {
		t.Fatalf("TestSequenceOfDecode Error: %+v err:%v\n", *bs, bitBuffer.Error())
	}
}
