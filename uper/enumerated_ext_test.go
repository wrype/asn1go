package uper

import (
	"encoding/hex"
	"encoding/json"
	"fmt"
	"testing"

	"github.com/wrype/asn1go/common"
)

type LightbarInUse int64

/**
CV2X2
DEFINITIONS AUTOMATIC TAGS::=
BEGIN
	LightbarInUse ::= ENUMERATED {
		unavailable (0), -- Not Equipped or unavailable
		notInUse (1), -- none active
		inUse (2),
		yellowCautionLights (3),
		schooldBusLights (4),
		arrowSignsActive (5),
		slowMovingVehicle (6),
		freqStops (7),
		...,
		qq(8),
		gg(9),
		...
	}

END
*/
var LightbarInUseMember common.Member

func init() {
	def := `{"members":[{"expr_type":"A1TC_UNIVERVAL","meta_type":"AMT_VALUE","Identifier":"LIU_Unavailable","value":0},{"expr_type":"A1TC_UNIVERVAL","meta_type":"AMT_VALUE","Identifier":"LIU_NotInUse","value":1},{"expr_type":"A1TC_UNIVERVAL","meta_type":"AMT_VALUE","Identifier":"LIU_InUse","value":2},{"expr_type":"A1TC_UNIVERVAL","meta_type":"AMT_VALUE","Identifier":"LIU_YellowCautionLights","value":3},{"expr_type":"A1TC_UNIVERVAL","meta_type":"AMT_VALUE","Identifier":"LIU_SchooldBusLights","value":4},{"expr_type":"A1TC_UNIVERVAL","meta_type":"AMT_VALUE","Identifier":"LIU_ArrowSignsActive","value":5},{"expr_type":"A1TC_UNIVERVAL","meta_type":"AMT_VALUE","Identifier":"LIU_SlowMovingVeh","value":6},{"expr_type":"A1TC_UNIVERVAL","meta_type":"AMT_VALUE","Identifier":"LIU_FreqStops","value":7},{"expr_type":"A1TC_EXTENSIBLE","meta_type":"AMT_TYPE","Identifier":"...","value":0},{"expr_type":"A1TC_UNIVERVAL","meta_type":"AMT_VALUE","Identifier":"LIU_Qq","value":8},{"expr_type":"A1TC_UNIVERVAL","meta_type":"AMT_VALUE","Identifier":"LIU_Gg","value":9},{"expr_type":"A1TC_EXTENSIBLE","meta_type":"AMT_TYPE","Identifier":"...","value":0}],"expr_type":"ENUMERATED","meta_type":"AMT_TYPE","Identifier":"LightbarInUse","value":0}`
	json.Unmarshal([]byte(def), &LightbarInUseMember)
}

const (
	LIU_Unavailable         LightbarInUse = 0
	LIU_NotInUse            LightbarInUse = 1
	LIU_InUse               LightbarInUse = 2
	LIU_YellowCautionLights LightbarInUse = 3
	LIU_SchooldBusLights    LightbarInUse = 4
	LIU_ArrowSignsActive    LightbarInUse = 5
	LIU_SlowMovingVeh       LightbarInUse = 6
	LIU_FreqStops           LightbarInUse = 7
	LIU_Qq                  LightbarInUse = 8
	LIU_Gg                  LightbarInUse = 9
)

func (t *LightbarInUse) UperEncode() *common.BitBuffer {
	return UperEncodeEnumerated(LightbarInUseMember, int64(*t))
}

func (t *LightbarInUse) UperDecode(buf *common.BitBuffer) {
	if buf.Error() != nil {
		return
	}
	*t = LightbarInUse(UperDecodeEnumerated(LightbarInUseMember, buf))
}

func (t *LightbarInUse) MarshalJSON() ([]byte, error) {
	switch *t {
	case LIU_Unavailable:
		return []byte("\"unavailable\""), nil
	case LIU_NotInUse:
		return []byte("\"notInUse\""), nil
	case LIU_InUse:
		return []byte("\"inUse\""), nil
	case LIU_YellowCautionLights:
		return []byte("\"yellowCautionLights\""), nil
	case LIU_SchooldBusLights:
		return []byte("\"schooldBusLights\""), nil
	case LIU_ArrowSignsActive:
		return []byte("\"arrowSignsActive\""), nil
	case LIU_SlowMovingVeh:
		return []byte("\"slowMovingVehicle\""), nil
	case LIU_FreqStops:
		return []byte("\"freqStops\""), nil
	case LIU_Qq:
		return []byte("\"qq\""), nil
	case LIU_Gg:
		return []byte("\"gg\""), nil
	default:
		return []byte("null"), nil
	}
}

func (t *LightbarInUse) UnmarshalJSON(b []byte) error {
	switch string(b) {
	case "\"unavailable\"":
		*t = LIU_Unavailable
	case "\"notInUse\"":
		*t = LIU_NotInUse
	case "\"inUse\"":
		*t = LIU_InUse
	case "\"yellowCautionLights\"":
		*t = LIU_YellowCautionLights
	case "\"schooldBusLights\"":
		*t = LIU_SchooldBusLights
	case "\"arrowSignsActive\"":
		*t = LIU_ArrowSignsActive
	case "\"slowMovingVehicle\"":
		*t = LIU_SlowMovingVeh
	case "\"freqStops\"":
		*t = LIU_FreqStops
	case "\"qq\"":
		*t = LIU_Qq
	case "\"gg\"":
		*t = LIU_Gg
	default:
		*t = -1
	}
	return nil
}

func TestEnumeratedExt(t *testing.T) {
	var a = LIU_Gg
	buf, err := a.UperEncode().Result()
	if err != nil || fmt.Sprintf("%2X", buf) != "81" {
		t.Fatalf("EnumeratedExt Encode Error:%v %2X", err, buf)
	} else {
		t.Log("Pass")
	}

	var b LightbarInUse
	buf, _ = hex.DecodeString("80")
	bitBuf := common.NewBitBufferFromBytes(buf)
	b.UperDecode(bitBuf)
	if bitBuf.Error() != nil || b != LIU_Qq {
		t.Fatalf("EnumeratedExt Decode Error:%v %d", bitBuf.Error(), b)
	} else {
		t.Log("Pass")
	}
}
