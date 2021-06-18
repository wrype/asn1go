package uper

import (
	"encoding/hex"
	"encoding/json"
	"fmt"
	"testing"

	"github.com/wrype/asn1go/common"
)

/**
MsgDayII ::= SEQUENCE {
		messageId INTEGER (0..32767) ,
		value OCTET STRING (SIZE(1..16)),
        aa INTEGER(0..7),
        ...,
        ss INTEGER(0..7) OPTIONAL,
        gg INTEGER(0..7),
        ...
}
*/
type MsgDayII struct {
	MessageId int64  `json:"messageId"`
	Val       []byte `json:"value"`
	Aa        int64  `json:"aa"`
	Ss        *int64 `json:"ss"`
	Gg        *int64 `json:"gg"`
}

var MsgDayIIMember common.Member

func (t *MsgDayII) UperEncode() *common.BitBuffer {
	return UperEncodeSequence(MsgDayIIMember, t)
}

func (t *MsgDayII) UperDecode(buf *common.BitBuffer) {
	if buf.Error() != nil {
		return
	}
	UperDecodeSequence(MsgDayIIMember, buf, t)
}

func init() {
	def := `{"members":[{"expr_type":"INTEGER","meta_type":"AMT_TYPE","constraints":{"type":"SIZE_RANG","value":0,"min":0,"max":32767,"base":0,"exponent":0},"Identifier":"MessageId","marker":{"flags":"EM_NOMARK"},"value":0},{"expr_type":"OCTET_STRING","meta_type":"AMT_TYPE","constraints":{"type":"SIZE_RANG","value":0,"min":1,"max":16,"base":0,"exponent":0},"Identifier":"Val","marker":{"flags":"EM_NOMARK"},"value":0},{"expr_type":"INTEGER","meta_type":"AMT_TYPE","constraints":{"type":"SIZE_RANG","value":0,"min":0,"max":7,"base":0,"exponent":0},"Identifier":"Aa","marker":{"flags":"EM_NOMARK"},"value":0},{"expr_type":"A1TC_EXTENSIBLE","meta_type":"AMT_TYPE","Identifier":"...","value":0},{"expr_type":"INTEGER","meta_type":"AMT_TYPE","constraints":{"type":"SIZE_RANG","value":0,"min":0,"max":7,"base":0,"exponent":0},"Identifier":"Ss","marker":{"flags":"EM_NOMARK"},"value":0},{"expr_type":"INTEGER","meta_type":"AMT_TYPE","constraints":{"type":"SIZE_RANG","value":0,"min":0,"max":7,"base":0,"exponent":0},"Identifier":"Gg","marker":{"flags":"EM_NOMARK"},"value":0},{"expr_type":"A1TC_EXTENSIBLE","meta_type":"AMT_TYPE","Identifier":"...","value":0}],"expr_type":"SEQUENCE","meta_type":"AMT_TYPE","Identifier":"MsgDayII","value":0}`
	json.Unmarshal([]byte(def), &MsgDayIIMember)
}

func TestUperEncodeSequenceExt(t *testing.T) {
	aa := int64(7)
	ss := int64(7)
	gg := int64(3)
	val, _ := hex.DecodeString("ABCD")
	/**
	rec1value MsgDayII ::=
	{
	  messageId 99,
	  value 'ABCD'H,
	  aa 7,
	  ss 7,
	  gg 3
	}
	*/
	msg := MsgDayII{
		MessageId: 99,
		Val:       val,
		Aa:        aa,
		Ss:        &ss,
		Gg:        &gg,
	}
	buf, err := msg.UperEncode().Result()
	if err != nil || fmt.Sprintf("%2X", buf) != "80631ABCDE0701E00160" {
		t.Fatalf("UperEncodeSequenceExt Error:%v %2X != 80631ABCDE0701E00160 \n", err, buf)
	} else {
		t.Log("Pass")
	}

	/**
	rec1value MsgDayII ::=
	{
	  messageId 10,
	  value '1234ABCD'H,
	  aa 5,
	  gg 7
	}
	*/
	msg.MessageId = 10
	msg.Aa = 5
	msg.Ss = nil
	msg.Gg = &aa
	msg.Val, _ = hex.DecodeString("1234ABCD")
	buf, err = msg.UperEncode().Result()
	if err != nil || fmt.Sprintf("%2X", buf) != "800A31234ABCDA0501E0" {
		t.Fatalf("UperEncodeSequenceExt Error:%v %2X != 800A31234ABCDA0501E0 \n", err, buf)
	} else {
		t.Log("Pass")
	}
}

func TestUperDecodeSequenceExt(t *testing.T) {
	msg := &MsgDayII{}
	buffer, _ := hex.DecodeString("80631ABCDE0701E00160")
	bitBuffer := common.NewBitBufferFromBytes(buffer)
	msg.UperDecode(bitBuffer)
	if msg.MessageId != 99 || msg.Aa != 7 || fmt.Sprintf("%2X", msg.Val) != "ABCD" && msg.Gg == nil || msg.Ss == nil || *msg.Gg != 3 || *msg.Ss != 7 || bitBuffer.Error() != nil {
		j, _ := json.Marshal(msg)
		t.Fatalf("TestUperDecodeSequenceExt Error:%v %s %2X", bitBuffer.Error(), string(j), msg.Val)
	} else {
		t.Log("Pass")
	}

	msg = &MsgDayII{}
	buffer, _ = hex.DecodeString("800A31234ABCDA0900F000")
	bitBuffer = common.NewBitBufferFromBytes(buffer)
	msg.UperDecode(bitBuffer)
	if msg.MessageId != 10 || msg.Aa != 5 || fmt.Sprintf("%2X", msg.Val) != "1234ABCD" && msg.Gg == nil || msg.Ss != nil || *msg.Gg != 7 || bitBuffer.Error() != nil {
		j, _ := json.Marshal(msg)
		t.Fatalf("TestUperDecodeSequenceExt Error:%v %s %2X", bitBuffer.Error(), string(j), msg.Val)
	} else {
		t.Log("Pass")
	}

	/**
	MsgDayII ::= SEQUENCE {
			messageId INTEGER (0..32767) ,
			value OCTET STRING (SIZE(1..16)),
	        aa INTEGER(0..7),
	        ...,
	        ss INTEGER(0..7) OPTIONAL,
	        gg INTEGER(0..7),
			qq  INTEGER(0..7),
			...
		}
	{
	  "value":"1234EFABCD",
	  "messageId":11
	  "gg":6,
	  "aa":5,
	  "qq":6
	}
	*/
	msg = &MsgDayII{}
	buffer, _ = hex.DecodeString("800B41234EFABCDA0980E000E000")
	bitBuffer = common.NewBitBufferFromBytes(buffer)
	msg.UperDecode(bitBuffer)
	if msg.MessageId != 11 || msg.Aa != 5 || fmt.Sprintf("%2X", msg.Val) != "1234EFABCD" && msg.Gg == nil || msg.Ss != nil || *msg.Gg != 6 || bitBuffer.Error() != nil {
		j, _ := json.Marshal(msg)
		t.Fatalf("TestUperDecodeSequenceExt Error:%v %s %2X", bitBuffer.Error(), string(j), msg.Val)
	} else {
		t.Log("Pass")
	}
}
