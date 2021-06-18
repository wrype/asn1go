package uper

import (
	"encoding/hex"
	"encoding/json"
	"fmt"
	"testing"

	"github.com/wrype/asn1go/common"
)

// {
//     "members": [
//         {
//             "expr_type": "IA5String",
//             "meta_type": "AMT_TYPE",
//             "constraints": {
//                 "type": "SIZE_RANG",
//                 "min": 1,
//                 "max": 8
//             },
//             "Identifier": "Choice0IA5str"
//         },
//         {
//             "expr_type": "ChoiceStruct",
//             "meta_type": "ASN_TYPE",
//             "Identifier": "Choice1Struct"
//         },
//         {
//             "expr_type": "INTEGER",
//             "meta_type": "AMT_TYPE",
//             "constraints": {
//                 "type": "SIZE_RANG",
//                 "min": 1,
//                 "max": 16
//             },
//             "Identifier": "Choice2"
//         },
//         {
//             "expr_type": "A1TC_EXTENSIBLE",
//             "meta_type": "AMT_TYPE",
//             "Identifier": "..."
//         },
//         {
//             "expr_type": "INTEGER",
//             "meta_type": "AMT_TYPE",
//             "constraints": {
//                 "type": "SIZE_RANG",
//                 "min": 1,
//                 "max": 32
//             },
//             "Identifier": "Choice3"
//         }
//     ],
//     "expr_type": "CHOICE",
//     "meta_type": "AMT_TYPE",
//     "constraints": null,
//     "Identifier": "Choice"
// }
var ChoiceRootMember common.Member

var ChoiceSubMember common.Member

func init() {
	json.Unmarshal([]byte(`{
                    "members": [
                        {
                            "expr_type": "IA5String",
                            "meta_type": "AMT_TYPE",
                            "constraints": {
                                "type": "SIZE_RANG",
                                "min": 1,
                                "max": 64
                            },
                            "Identifier": "Choice0IA5str"
                        },
                        {
                            "expr_type": "ChoiceASNType",
                            "meta_type": "ASN_TYPE",
                            "Identifier": "Choice1ASN"
                        },
                        {
                            "expr_type": "INTEGER",
                            "meta_type": "AMT_TYPE",
                            "constraints": {
                                "type": "SIZE_RANG",
                                "min": 1,
                                "max": 16
                            },
                            "Identifier": "Choice2Int"
                        },
                        {
                            "expr_type": "A1TC_EXTENSIBLE",
                            "meta_type": "AMT_TYPE",
                            "Identifier": "..."
                        },
                        {
                            "expr_type": "INTEGER",
                            "meta_type": "AMT_TYPE",
                            "constraints": {
                                "type": "SIZE_RANG",
                                "min": 1,
                                "max": 32
                            },
                            "Identifier": "Choice3Int"
                        }
                    ],
                    "expr_type": "CHOICE",
                    "meta_type": "AMT_TYPE",
                    "constraints": null,
                    "Identifier": "Choice"
			}`), &ChoiceRootMember)
	json.Unmarshal([]byte(`{
                "expr_type": "INTEGER",
                "meta_type": "AMT_TYPE",
                "constraints": {
                    "type": "SIZE_RANG",
                    "min": 1,
                    "max": 32
                }
            }`), &ChoiceSubMember)
}

type ChoiceASNType int64

type Choice struct {
	Choice0IA5str *string
	Choice1ASN    *ChoiceASNType
	Choice2Int    *int64
	Choice3Int    *int64
}

func (c *ChoiceASNType) UperEncode() *common.BitBuffer {
	return UperEncodeInteger(ChoiceSubMember, int64(*c))
}

func (c *ChoiceASNType) UperDecode(b *common.BitBuffer) {
	*c = ChoiceASNType(UperDecodeInteger(ChoiceSubMember, b))
}

func (c *Choice) UperEncode() *common.BitBuffer {
	return UperEncodeChoice(ChoiceRootMember, c)
}

func (c *Choice) UperDecode(b *common.BitBuffer) {
	UperDecodeChoice(ChoiceRootMember, b, c)
}

func TestChoiceEncode(t *testing.T) {
	field := ChoiceASNType(10)
	ab, ae := field.UperEncode().Result()
	t.Logf("Field Encode %2x %v", ab, ae)

	choice := Choice{Choice1ASN: &field}
	ab, ae = choice.UperEncode().Result()
	t.Logf("Choice Encode %2x %v", ab, ae)
	if fmt.Sprintf("%2x", ab) == "29" {
		t.Log("Pass")
	} else {
		t.Fatal("Choice Encode Error!")
	}

	field2 := `'afsef3';'`
	choice = Choice{Choice0IA5str: &field2}
	ab, ae = choice.UperEncode().Result()
	t.Logf("Choice Encode %2x %v", ab, ae)
	if fmt.Sprintf("%2x", ab) == "04a7c39b9e5cccd3bb4e" {
		t.Log("Pass")
	} else {
		t.Fatal("Choice Encode Error!")
	}
}

func TestChoiceDecode(t *testing.T) {
	choice := Choice{}
	hexStr, _ := hex.DecodeString("04a7c39b9e5cccd3bb4e")
	bs := common.NewBitBufferFromBytes(hexStr)
	choice.UperDecode(bs)
	if *choice.Choice0IA5str == `'afsef3';'` {
		t.Log("Pass")
	} else {
		t.Fatal("Choice Decode Error! ", bs.Error())
	}

	choice = Choice{}
	hexStr, _ = hex.DecodeString("29")
	bs = common.NewBitBufferFromBytes(hexStr)
	choice.UperDecode(bs)
	if *choice.Choice1ASN == 10 {
		t.Log("Pass")
	} else {
		t.Fatal("Choice Decode Error! ", bs.Error())
	}
}
