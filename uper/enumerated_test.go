package uper

import (
	"encoding/hex"
	"encoding/json"
	"fmt"
	"testing"

	"github.com/wrype/asn1go/common"
)

//var SeqMember common.Member
var SlotStatusMember common.Member
var AngularVConfidenceMember common.Member

func init() {
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
                        "Identifier": "prec100deg",
                        "value": 1
                    },
                    {
                        "expr_type": "A1TC_UNIVERVAL",
                        "meta_type": "AMT_VALUE",
                        "Identifier": "prec10deg",
                        "value": 2
                    },
                    {
                        "expr_type": "A1TC_UNIVERVAL",
                        "meta_type": "AMT_VALUE",
                        "Identifier": "prec5deg",
                        "value": 3
                    },
                    {
                        "expr_type": "A1TC_UNIVERVAL",
                        "meta_type": "AMT_VALUE",
                        "Identifier": "prec1deg",
                        "value": 4
                    },
                    {
                        "expr_type": "A1TC_UNIVERVAL",
                        "meta_type": "AMT_VALUE",
                        "Identifier": "prec0-1deg",
                        "value": 5
                    },
                    {
                        "expr_type": "A1TC_UNIVERVAL",
                        "meta_type": "AMT_VALUE",
                        "Identifier": "prec0-05deg",
                        "value": 6
                    },
                    {
                        "expr_type": "A1TC_UNIVERVAL",
                        "meta_type": "AMT_VALUE",
                        "Identifier": "prec0-01deg",
                        "value": 7
                    }
                ],
                "expr_type": "ENUMERATED",
                "meta_type": "AMT_TYPE",
                "constraints": null,
                "tag": {
                    "tag_mode": "TM_DEFAULT",
                    "tag_class": "",
                    "tag_value": ""
                },
                "Identifier": "AngularVConfidence"
            }`), &AngularVConfidenceMember)
	json.Unmarshal([]byte(`{
		"members": [
			{
				"expr_type": "A1TC_UNIVERVAL",
				"meta_type": "AMT_VALUE",
				"Identifier": "unknown",
				"value": 0
			},
			{
				"expr_type": "A1TC_UNIVERVAL",
				"meta_type": "AMT_VALUE",
				"Identifier": "available",
				"value": 1
			},
			{
				"expr_type": "A1TC_UNIVERVAL",
				"meta_type": "AMT_VALUE",
				"Identifier": "occupied",
				"value": 2
			},
			{
				"Identifier": "...",
				"expr_type": "A1TC_EXTENSIBLE",
				"meta_type": "AMT_TYPE"
			}
		],
		"expr_type": "ENUMERATED",
		"meta_type": "AMT_TYPE",
		"constraints": null,
		"tag": {
			"tag_mode": "TM_DEFAULT",
			"tag_class": "",
			"tag_value": ""
		},
		"Identifier": "SlotStatus"
	}`), &SlotStatusMember)
}

type SlotStatus int64
type AngularVConfidence int64

func (e *SlotStatus) UperEncode() *common.BitBuffer {
	return UperEncodeEnumerated(SlotStatusMember, int64(*e))
}

func (e *SlotStatus) UperDecode(b *common.BitBuffer) {
	*e = SlotStatus(UperDecodeEnumerated(SlotStatusMember, b))
}

func (e *AngularVConfidence) UperEncode() *common.BitBuffer {
	return UperEncodeEnumerated(AngularVConfidenceMember, int64(*e))
}

func (e *AngularVConfidence) UperDecode(b *common.BitBuffer) {
	*e = AngularVConfidence(UperDecodeEnumerated(AngularVConfidenceMember, b))
}

func TestEnumWithoutExt(t *testing.T) {
	for i := 0; i < 8; i++ {
		t1 := AngularVConfidence(i)
		b1 := t1.UperEncode().Bytes()
		t.Logf("value AngularVConfidenceMember ::= %v; data.uper is %2x\n", i, b1)
		ans := fmt.Sprintf("%x0", uint8(i*2))
		if fmt.Sprintf("%2x", b1) == ans {
			t.Log("Pass")
		} else {
			t.Fatalf("SpeedRang Uper Encode Error! %2x(want %s)", b1, ans)
		}

		var t2 AngularVConfidence
		t2.UperDecode(common.NewBitBufferFromBytes(b1))
		t.Logf("buffer %2x; Decode AngularVConfidence is %d\n", b1, t2)
		if t2 == AngularVConfidence(i) {
			t.Log("Pass")
		} else {
			t.Fatal("AngularVConfidence Uper Decode Error!")
		}
	}
}

func TestEnumWithoutExtNoInEnum(t *testing.T) {
	for i := 8; i < 10; i++ {
		t1 := AngularVConfidence(i)
		b1, err := t1.UperEncode().Result()
		t.Logf("value AngularVConfidenceMember ::= %v; data.uper is error", i)
		if err != nil {
			t.Log("Pass")
		} else {
			t.Fatalf("SpeedRang Uper Encode Error!")
		}

		var t2 AngularVConfidence
		bitBuffer := common.NewBitBufferFromBytes(b1)
		t2.UperDecode(bitBuffer)
		t.Logf("buffer %2x; Decode AngularVConfidence error", b1)
		if bitBuffer.Error() != nil {
			t.Log("Pass")
		} else {
			t.Fatal("AngularVConfidence Uper Decode Error!")
		}
	}
}

func TestEnumWithExt(t *testing.T) {
	for i := 0; i < 3; i++ {
		t1 := SlotStatus(i)
		b1 := t1.UperEncode().Bytes()
		t.Logf("value SlotStatus ::= %v; data.uper is %2x\n", i, b1)
		ans := fmt.Sprintf("%v0", i*2)
		if fmt.Sprintf("%2x", b1) == ans {
			t.Log("Pass")
		} else {
			t.Fatalf("SpeedRang Uper Encode Error! %2x(want %s)", b1, ans)
		}

		var t2 SlotStatus
		t2.UperDecode(common.NewBitBufferFromBytes(b1))
		t.Logf("buffer %2x; Decode SlotStatus is %d\n", b1, t2)
		if t2 == SlotStatus(i) {
			t.Log("Pass")
		} else {
			t.Fatal("SlotStatus Uper Decode Error!")
		}
	}
}

func TestEnumWithExtNoInEnum(t *testing.T) {
	for i := 3; i < 5; i++ {
		t1 := SlotStatus(i)
		b1, err := t1.UperEncode().Result()
		t.Logf("value SlotStatus ::= %v; data.uper is error", i)
		if err != nil {
			t.Log("Pass")
		} else {
			t.Fatalf("SpeedRang Uper Encode Error!")
		}

		var t2 SlotStatus
		bitBuffer := common.NewBitBufferFromBytes(b1)
		t2.UperDecode(bitBuffer)
		t.Logf("buffer %2x; Decode SlotStatus error", b1)
		if bitBuffer.Error() != nil {
			t.Log("Pass")
		} else {
			t.Fatal("SlotStatus Uper Decode Error!")
		}
	}
}

func TestEnumWithExtNosupport(t *testing.T) {
	var t2 SlotStatus
	buf, _ := hex.DecodeString("80") //
	bitBuffer := common.NewBitBufferFromBytes(buf)
	t2.UperDecode(bitBuffer)
	t.Logf("buffer %2x; Decode SlotStatus error: %s", buf, bitBuffer.Error())
	if bitBuffer.Error() != nil {
		t.Log("Pass")
	} else {
		t.Fatal("SlotStatus Uper Decode Error!")
	}
}
