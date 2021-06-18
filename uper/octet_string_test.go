package uper

import (
	"encoding/json"
	"fmt"
	"testing"

	"github.com/wrype/asn1go/common"
)

var RTCMmessageMember common.Member

type RTCMmessage []byte

var contractSerialNumerMember common.Member

type contractSerialNumer []byte

var RangeExtMember common.Member

type RangeExtbs []byte

func init() {
	//value RTCMmessage ::= '12345678'H
	json.Unmarshal([]byte(`{
        "expr_type": "OCTET_STRING",
        "meta_type": "AMT_TYPE",
        "constraints": {
          "type": "SIZE_RANG",
          "min": 1,
          "max": 2047
        },
        "tag": {
          "tag_mode": "TM_DEFAULT",
          "tag_class": "",
          "tag_value": ""
        },
        "Identifier": "RTCMmessage"
	  }`), &RTCMmessageMember)

	//value contractSerialNumer ::= '1234567812345678'H
	json.Unmarshal([]byte(`{
		"expr_type": "OCTET_STRING",
		"meta_type": "AMT_TYPE",
		"constraints": {
		  "type": "SIZE_FIXED",
		  "value": 8
		},
		"tag": {
		  "tag_mode": "TM_DEFAULT",
		  "tag_class": "",
		  "tag_value": ""
		},
		"Identifier": "contractSerialNumer",
		"marker": {
		  "flags": "EM_NOMARK"
		}
	  }`), &contractSerialNumerMember)

	//value RTCMmessage ::= ''H
	json.Unmarshal([]byte(`{
        "expr_type": "OCTET_STRING",
        "meta_type": "AMT_TYPE",
        "constraints": {
          "type": "SIZE_RANG_AND_EXT",
          "min": 1,
          "max": 3
        },
        "tag": {
          "tag_mode": "TM_DEFAULT",
          "tag_class": "",
          "tag_value": ""
        },
        "Identifier": "RTCMmessage"
	  }`), &RangeExtMember)

}

func (r RTCMmessage) UperEncode() *common.BitBuffer {
	return UperEncodeOctetString(RTCMmessageMember, r)
}

func (r RTCMmessage) UperDecode(b *common.BitBuffer) []byte {
	return UperDecodeOctetString(RTCMmessageMember, b)
}

func (c contractSerialNumer) UperEncode() *common.BitBuffer {
	return UperEncodeOctetString(contractSerialNumerMember, c)
}

func (c contractSerialNumer) UperDecode(b *common.BitBuffer) []byte {
	return UperDecodeOctetString(contractSerialNumerMember, b)
}

func (r RangeExtbs) UperEncode() *common.BitBuffer {
	return UperEncodeOctetString(RangeExtMember, r)
}

func (r RangeExtbs) UperDecode(b *common.BitBuffer) []byte {
	return UperDecodeOctetString(RangeExtMember, b)
}

func TestSIZE_RANGE(t *testing.T) {
	var data1 RTCMmessage
	data1 = append(data1, 0x12)
	data1 = append(data1, 0x34)
	data1 = append(data1, 0x56)
	data1 = append(data1, 0x78)
	bitbuf := data1.UperEncode()
	bs := bitbuf.Bytes()
	t.Logf("value RTCMmessage ::= 0x12345678; data.uper is %8x\n", bs)
	if bs[0] != 0 {
		t.Fatal("contractSerialNumer Uper Encode Error!")
	} else if bs[1] != 98 {
		t.Fatal("contractSerialNumer Uper Encode Error!")
	} else if bs[2] != 70 {
		t.Fatal("contractSerialNumer Uper Encode Error!")
	} else if bs[3] != 138 {
		t.Fatal("contractSerialNumer Uper Encode Error!")
	} else if bs[4] != 207 {
		t.Fatal("contractSerialNumer Uper Encode Error!")
	} else if bs[5] != 0 {
		t.Fatal("contractSerialNumer Uper Encode Error!")
	} else {
		t.Log("range encode Pass")
	}
	var data2 RTCMmessage
	data2 = data2.UperDecode(bitbuf)
	t.Logf("buffer %v; Decode contractSerialNumer is %8x\n", bs, data2)
	if fmt.Sprintf("%8x", data2) != "12345678" {
		t.Fatal("contractSerialNumer Uper Encode Error!")
	} else {
		t.Log("range decode Pass")
	}
}

func TestSIZE_FIXED(t *testing.T) {
	var data1 contractSerialNumer
	data1 = append(data1, 0x12)
	data1 = append(data1, 0x34)
	data1 = append(data1, 0x56)
	data1 = append(data1, 0x78)
	data1 = append(data1, 0x12)
	data1 = append(data1, 0x34)
	data1 = append(data1, 0x56)
	data1 = append(data1, 0x78)
	bitbuf := data1.UperEncode()
	bs := bitbuf.Bytes()
	t.Logf("value contractSerialNumer ::= 0x1234567812345678; data.uper is %16x\n", bs)
	if fmt.Sprintf("%16x", bs) != "1234567812345678" {
		t.Fatal("contractSerialNumer Uper Encode Error!")
	} else {
		t.Log("Pass")
	}

	var data2 contractSerialNumer
	data2 = data2.UperDecode(bitbuf)
	t.Logf("buffer %v; Decode contractSerialNumer is %16x\n", bs, data2)
	if fmt.Sprintf("%16x", data2) != "1234567812345678" {
		t.Fatal("contractSerialNumer Uper Decode Error!")
	} else {
		t.Log("Pass")
	}
}

//测试ext处理
func TestSIZE_RANGE_AND_EXT(t *testing.T) {
	var data1 RangeExtbs
	data1 = append(data1, 0x78)
	data1 = append(data1, 0x34)
	data1 = append(data1, 0x12)
	data1 = append(data1, 0x78)
	bitbuf := data1.UperEncode()
	bs := bitbuf.Bytes()
	t.Logf("value RangeExtbs ::= 0x78341278; data.uper is %8x\n", bs)
	if bs[0] != 130 {
		t.Fatal("contractSerialNumer Uper Encode Error!")
	} else if bs[1] != 60 {
		t.Fatal("contractSerialNumer Uper Encode Error!")
	} else if bs[2] != 26 {
		t.Fatal("contractSerialNumer Uper Encode Error!")
	} else if bs[3] != 9 {
		t.Fatal("contractSerialNumer Uper Encode Error!")
	} else if bs[4] != 60 {
		t.Fatal("contractSerialNumer Uper Encode Error!")
	} else if bs[5] != 0 {
		t.Fatal("contractSerialNumer Uper Encode Error!")
	} else {
		t.Log("range ext encode Pass")
	}
	var data2 RangeExtbs
	data2 = data2.UperDecode(bitbuf)
	t.Logf("buffer %v; Decode contractSerialNumer is %8x\n", bs, data2)
	if fmt.Sprintf("%8x", data2) != "78341278" {
		t.Fatal("contractSerialNumer Uper Encode Error!")
	} else {
		t.Log("range decode Pass")
	}
}
