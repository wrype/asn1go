package uper

import (
	"encoding/hex"
	"encoding/json"
	"fmt"
	"testing"

	"github.com/wrype/asn1go/common"
)

/**
IA5 ::= IA5String(SIZE(1..2,...))
*/
var IA5Member common.Member

type IA5 string

func init() {
	def := `{"expr_type":"IA5String","meta_type":"AMT_TYPE","constraints":{"type":"SIZE_RANG_AND_EXT","value":0,"min":1,"max":2,"base":0,"exponent":0},"Identifier":"IA5","value":0}`
	json.Unmarshal([]byte(def), &IA5Member)
}

func (t *IA5) UperEncode() *common.BitBuffer {
	return UperEncodeIA5String(IA5Member, string(*t))
}

func (t *IA5) UperDecode(buf *common.BitBuffer) {
	if buf.Error() != nil {
		return
	}
	*t = IA5(UperDecodeIA5String(IA5Member, buf))
}

func TestIA5StringExt(t *testing.T) {
	var str IA5 = "123"
	buf, err := str.UperEncode().Result()
	if err != nil || fmt.Sprintf("%2X", buf) != "81B164CC" {
		t.Fatalf("IA5String Ext Encode Error: %v %2X", err, buf)
	} else {
		t.Log("Pass")
	}

	buf, _ = hex.DecodeString("863180C99B56CDDC39610590")
	bitBuf := common.NewBitBufferFromBytes(buf)
	str.UperDecode(bitBuf)
	if bitBuf.Error() != nil || str != "1@23567890A2" {
		t.Fatalf("IA5String Ext Encode Error: %v %v", bitBuf.Error(), str)
	} else {
		t.Log("Pass")
	}
}
