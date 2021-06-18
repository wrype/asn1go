package uper

import (
	"encoding/hex"
	"encoding/json"
	"fmt"
	"testing"

	"github.com/wrype/asn1go/common"
)

var SeqMember common.Member

var SSMember common.Member

var SeqExtMember common.Member

func init() {
	json.Unmarshal([]byte(`{
                "members": [
                    {
                        "expr_type": "INTEGER",
                        "meta_type": "AMT_TYPE",
                        "constraints": {
                            "type": "SIZE_RANG_AND_EXT",
                            "min": 1,
                            "max": 32
                        },
                        "tag": {
                            "tag_mode": "TM_DEFAULT",
                            "tag_class": "",
                            "tag_value": ""
                        },
                        "Identifier": "Range",
                        "marker": {
                            "flags": "EM_NOMARK"
                        }
                    },
                    {
                        "expr_type": "SS",
                        "meta_type": "ASN_TYPE",
                        "Identifier": "Ss",
                        "marker": {
                            "flags": "EM_OPTIONAL"
                        }
                    }
                ],
                "expr_type": "SEQUENCE",
                "meta_type": "AMT_TYPE",
                "constraints": null,
                "tag": {
                    "tag_mode": "TM_DEFAULT",
                    "tag_class": "",
                    "tag_value": ""
                },
                "Identifier": "ROOT"
            }`), &SeqMember)

	json.Unmarshal([]byte(`{
                "expr_type": "INTEGER",
                "meta_type": "AMT_TYPE",
                "constraints": {
                    "type": "SIZE_RANG",
                    "min": 1,
                    "max": 32
                },
                "tag": {
                    "tag_mode": "TM_DEFAULT",
                    "tag_class": "",
                    "tag_value": ""
                },
                "Identifier": "SS"
            }`), &SSMember)

	json.Unmarshal([]byte(`{
                "members": [
                    {
                        "expr_type": "INTEGER",
                        "meta_type": "AMT_TYPE",
                        "constraints": {
                            "type": "SIZE_RANG_AND_EXT",
                            "min": 1,
                            "max": 32
                        },
                        "tag": {
                            "tag_mode": "TM_DEFAULT",
                            "tag_class": "",
                            "tag_value": ""
                        },
                        "Identifier": "Range",
                        "marker": {
                            "flags": "EM_NOMARK"
                        }
                    },
                    {
                        "Identifier": "...",
                        "expr_type": "A1TC_EXTENSIBLE",
                        "meta_type": "AMT_TYPE"
                    }
                ],
                "expr_type": "SEQUENCE",
                "meta_type": "AMT_TYPE",
                "constraints": null,
                "tag": {
                    "tag_mode": "TM_DEFAULT",
                    "tag_class": "",
                    "tag_value": ""
                },
                "Identifier": "ROOT"
            }`), &SeqExtMember)
}

/**
  SEQ ::= SEQUENCE
  {
     range     INTEGER(1..32,...),
	 ss		   SS OPTIONAL
  }

  SS ::=   INTEGER(1..32)
*/

type SEQ struct {
	Range int64
	Ss    *SS
}

type SS int64

func (s *SS) UperEncode() *common.BitBuffer {
	return UperEncodeInteger(SSMember, int64(*s))
}

func (s *SS) UperDecode(b *common.BitBuffer) {
	*s = SS(UperDecodeInteger(SSMember, b))
}

func (a *SEQ) UperEncode() *common.BitBuffer {
	return UperEncodeSequence(SeqMember, a)
}

func (a *SEQ) UperDecode(b *common.BitBuffer) {
	UperDecodeSequence(SeqMember, b, a)
}

func TestSeqEncode(t *testing.T) {
	ss := SS(10)
	a := SEQ{Range: 10, Ss: &ss}
	ab, ae := a.UperEncode().Result()
	t.Logf("TestSeqEncode %2x %v", ab, ae)
	if fmt.Sprintf("%2x", ab) == "9290" {
		t.Log("Pass")
	} else {
		t.Fatalf("TestSeqEncode Error %2x != 9290", ab)
	}

	b := SEQ{Range: 32}
	bb, be := b.UperEncode().Result()
	t.Logf("TestSeqEncode %2x %v", bb, be)
	if fmt.Sprintf("%2x", bb) == "3e" {
		t.Log("Pass")
	} else {
		t.Fatalf("TestSeqEncode Error %2x != 3e", bb)
	}
	//b := SEQ{}
	//b.UperDecode(common.NewBitBufferFromBytes(a.UperEncode().Bytes()))
	//t.Log(b)
}

func TestSeqDecode(t *testing.T) {
	a := &SEQ{}
	buf, _ := hex.DecodeString("9290") //10,10
	bs := common.NewBitBufferFromBytes(buf)
	a.UperDecode(bs)
	if a.Range == 10 && a.Ss != nil && *a.Ss == 10 && bs.Error() == nil {
		t.Log("Pass")
	} else {
		t.Fatalf("TestSeqDecode Error %+v %v", a, bs.Error())
	}

	b := &SEQ{}
	buf2, _ := hex.DecodeString("3e") //32,nil
	bs2 := common.NewBitBufferFromBytes(buf2)
	b.UperDecode(bs2)
	if b.Range == 32 && b.Ss == nil && bs2.Error() == nil {
		t.Log("Pass")
	} else {
		t.Fatalf("TestSeqDecode Error %+v %v", b, bs2.Error())
	}

}

/**
  ROOT ::= SEQUENCE
  {
     range     INTEGER(1..32,...) ,
     ...
  }
*/
type SeqExt struct {
	Range int64
}

func (a *SeqExt) UperEncode() *common.BitBuffer {
	return UperEncodeSequence(SeqExtMember, a)
}

func (a *SeqExt) UperDecode(b *common.BitBuffer) {
	UperDecodeSequence(SeqExtMember, b, a)
}

func TestSeqExtEncode(t *testing.T) {
	a := SeqExt{Range: 10}
	ab, ae := a.UperEncode().Result()
	t.Logf("TestSeqEncode %2x %v", ab, ae)
	if fmt.Sprintf("%2x", ab) == "12" {
		t.Log("Pass")
	} else {
		t.Fatalf("TestSeqEncode Error %2x != 12", ab)
	}
}

func TestSeqExtDecode(t *testing.T) {
	a := &SeqExt{}
	buf, _ := hex.DecodeString("12") //10
	bs := common.NewBitBufferFromBytes(buf)
	a.UperDecode(bs)
	if a.Range == 10 && bs.Error() == nil {
		t.Log("Pass")
	} else {
		t.Fatalf("TestSeqExtDecode Error %+v %v", a, bs.Error())
	}
}

/**
    SEQALL ::= SEQUENCE
	{
		range     INTEGER(1..10),
		data      OCTET STRING (SIZE(1..2047)),
		bit       PersonalAssistive OPTIONAL,
		...
	}

	PersonalAssistive ::= BIT STRING {
		unavailable (0),
		otherType (1),
		vision (2),
		hearing (3),
		movement (4),
		cognition (5)
	} (SIZE (6, ...))
*/

var PersonalAssistiveMember common.Member
var SeqAllMember common.Member

type PersonalAssistive uint64

type SeqAll struct {
	Range int64
	Data  []byte
	Bit   PersonalAssistive
}

func (p *PersonalAssistive) UperEncode() *common.BitBuffer {
	return UperEncodeBitString(PersonalAssistiveMember, uint64(*p))
}

func (s *SeqAll) UperEncode() *common.BitBuffer {
	return UperEncodeSequence(SeqAllMember, s)
}

func TestSeqAllEncode(t *testing.T) {
	json.Unmarshal([]byte(`{
                "members": [
                    {
                        "expr_type": "A1TC_UNIVERVAL",
                        "meta_type": "AMT_VALUE",
                        "Identifier": "unavailable",
                        "value": 0
                    },
                    {
                        "expr_type": "A1TC_UNIVERVAL",
                        "meta_type": "AMT_VALUE",
                        "Identifier": "otherType",
                        "value": 1
                    },
                    {
                        "expr_type": "A1TC_UNIVERVAL",
                        "meta_type": "AMT_VALUE",
                        "Identifier": "vision",
                        "value": 2
                    },
                    {
                        "expr_type": "A1TC_UNIVERVAL",
                        "meta_type": "AMT_VALUE",
                        "Identifier": "hearing",
                        "value": 3
                    },
                    {
                        "expr_type": "A1TC_UNIVERVAL",
                        "meta_type": "AMT_VALUE",
                        "Identifier": "movement",
                        "value": 4
                    },
                    {
                        "expr_type": "A1TC_UNIVERVAL",
                        "meta_type": "AMT_VALUE",
                        "Identifier": "cognition",
                        "value": 5
                    }
                ],
                "expr_type": "BIT_STRING",
                "meta_type": "AMT_TYPE",
                "constraints": {
                    "type": "SIZE_FIXED_AND_EXT",
                    "value": 6
                },
                "tag": {
                    "tag_mode": "TM_DEFAULT",
                    "tag_class": "",
                    "tag_value": ""
                },
                "Identifier": "PersonalAssistive"
            }`), &PersonalAssistiveMember)
	json.Unmarshal([]byte(`{
                "members": [
                    {
                        "expr_type": "INTEGER",
                        "meta_type": "AMT_TYPE",
                        "constraints": {
                            "type": "SIZE_RANG",
                            "min": 1,
                            "max": 10
                        },
                        "tag": {
                            "tag_mode": "TM_DEFAULT",
                            "tag_class": "",
                            "tag_value": ""
                        },
                        "Identifier": "Range",
                        "marker": {
                            "flags": "EM_NOMARK"
                        }
                    },
                    {
                        "expr_type": "OCTET_STRING",
                        "meta_type": "AMT_TYPE",
                        "constraints": {
                            "type": "SIZE_RANG",
                            "min": 1,
                            "max": 2047
                        },
                        "Identifier": "Data",
                        "marker": {
                            "flags": "EM_NOMARK"
                        }
                    },
                    {
                        "expr_type": "PersonalAssistive",
                        "meta_type": "ASN_TYPE",
                        "Identifier": "Bit",
                        "marker": {
                            "flags": "EM_OPTIONAL"
                        }
                    },
                    {
                        "Identifier": "...",
                        "expr_type": "A1TC_EXTENSIBLE",
                        "meta_type": "AMT_TYPE"
                    }
                ],
                "expr_type": "SEQUENCE",
                "meta_type": "AMT_TYPE",
                "constraints": null,
                "tag": {
                    "tag_mode": "TM_DEFAULT",
                    "tag_class": "",
                    "tag_value": ""
                },
                "Identifier": "SEQALL"
            }`), &SeqAllMember)

	a := &SeqAll{
		Range: 9,
		Data:  []byte("ABC"),
		Bit:   0x05, // unavailable, vision
	}
	ab, ae := a.UperEncode().Result()
	if fmt.Sprintf("%2X", ab) != "600120A121A8" || ae != nil {
		t.Fatalf("SeqAll Encode error buf: %2X  err: %v", ab, ae)
	} else {
		t.Log("Pass")
	}
}
