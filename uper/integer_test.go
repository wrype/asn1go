package uper

import (
	"encoding/hex"
	"encoding/json"
	"fmt"
	"testing"

	"github.com/wrype/asn1go/common"
)

type SpeedRang int64

var SpeedRangMember common.Member

type SpeedRangExt int64

var SpeedRangExtMember common.Member

func init() {
	json.Unmarshal([]byte(`{
                "expr_type": "INTEGER",
                "meta_type": "AMT_TYPE",
                "constraints": {
                    "type": "SIZE_RANG",
                    "min": 1,
                    "max": 255
                },
                "Identifier": "Speed"
            }`), &SpeedRangMember)

	json.Unmarshal([]byte(`{
                "expr_type": "INTEGER",
                "meta_type": "AMT_TYPE",
                "constraints": {
                    "type": "SIZE_RANG_AND_EXT",
                    "min": 1,
                    "max": 255
                },
                "Identifier": "Speed"
            }`), &SpeedRangExtMember)
}

/**
  SpeedRang ::= INTEGER(1..255)
*/
func (s *SpeedRang) UperEncode() *common.BitBuffer {
	return UperEncodeInteger(SpeedRangMember, int64(*s))
}

func (s *SpeedRang) UperDecode(b *common.BitBuffer) {
	*s = SpeedRang(UperDecodeInteger(SpeedRangMember, b))
}

/**
  SpeedRangExt ::= INTEGER(1..255,...)
*/
func (s *SpeedRangExt) UperEncode() *common.BitBuffer {
	return UperEncodeInteger(SpeedRangExtMember, int64(*s))
}

func (s *SpeedRangExt) UperDecode(b *common.BitBuffer) {
	*s = SpeedRangExt(UperDecodeInteger(SpeedRangExtMember, b))
}

func TestSIZE_RANG(t *testing.T) {
	var s SpeedRang = 10
	bs := s.UperEncode().Bytes()
	t.Logf("value SpeedRang ::= 10; data.uper is %2x\n", bs)
	if fmt.Sprintf("%2x", bs) == "09" {
		t.Log("Pass")
	} else {
		t.Fatal("SpeedRang Uper Encode Error!")
	}

	var s2 SpeedRang
	s2.UperDecode(common.NewBitBufferFromBytes(bs))
	t.Logf("buffer %2x; Decode SpeedRang is %d\n", bs, s2)
	if s2 == 10 {
		t.Log("Pass")
	} else {
		t.Fatal("SpeedRang Uper Decode Error!")
	}
}

func TestSIZE_RANG_AND_EXT(t *testing.T) {
	var s SpeedRangExt = 10
	bs := s.UperEncode().Bytes()
	t.Logf("value SpeedRangExt ::= 10; data.uper is %2x\n", bs)
	if fmt.Sprintf("%2x", bs) == "0480" {
		t.Log("Pass")
	} else {
		t.Fatal("SpeedRangExt Uper Encode Error!")
	}

	var s2 SpeedRangExt
	s2.UperDecode(common.NewBitBufferFromBytes(bs))
	t.Logf("buffer %2x; Decode SpeedRangExt is %d\n", bs, s2)
	if s2 == 10 {
		t.Log("Pass")
	} else {
		t.Fatal("SpeedRangExt Uper Decode Error!")
	}

	s = 256
	bs, err := s.UperEncode().Result()
	if err != nil || fmt.Sprintf("%2x", bs) != "81008000" {
		t.Fatalf("SpeedRangExt Uper Encode Error %v", s.UperEncode())
	} else {
		t.Log("Pass")
	}

	buf, _ := hex.DecodeString("813A9800")
	bitBuf := common.NewBitBufferFromBytes(buf)
	s.UperDecode(bitBuf)
	if bitBuf.Error() != nil || s != 30000 {
		t.Fatalf("SpeedRangExt Uper Decode Error %v %d", bitBuf.Error(), s)
	} else {
		t.Log("Pass")
	}
}
