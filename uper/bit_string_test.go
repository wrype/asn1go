package uper

import (
	"encoding/json"
	"fmt"
	"testing"

	"github.com/wrype/asn1go/common"
)

var vehicleEventFlagsmember common.Member

type vehicleEventFlag uint64

var laneSharingmember common.Member

type laneSharing uint64

func init() {
	//value vehicleEventFlag ::= '1001011001'B
	json.Unmarshal([]byte(`{
		"expr_type": "BIT_STRING",
		"meta_type": "AMT_TYPE",
		"constraints": {
	  		"type": "SIZE_FIXED",
	  		"value": 10
		},
		"tag": {
	  		"tag_mode": "TM_DEFAULT",
	  		"tag_class": "",
	  		"tag_value": ""
		},
		"Identifier": "LaneSharing"
	  }`), &vehicleEventFlagsmember)

	//value laneSharing ::= '1001011001101'B
	json.Unmarshal([]byte(`{
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
        "Identifier": "VehicleEventFlags"
      }`), &laneSharingmember)
}

func (v vehicleEventFlag) UperEncode() *common.BitBuffer {
	return UperEncodeBitString(vehicleEventFlagsmember, uint64(v))
}

func (v vehicleEventFlag) UperDecode(b *common.BitBuffer) uint64 {
	return UperDecodeBitString(vehicleEventFlagsmember, b)
}

func (l laneSharing) UperEncode() *common.BitBuffer {
	return UperEncodeBitString(laneSharingmember, uint64(l))
}

func (l laneSharing) UperDecode(b *common.BitBuffer) uint64 {
	return UperDecodeBitString(laneSharingmember, b)
}

/* func TestBS_SIZE_FIXED(t *testing.T) {
	var data1 vehicleEventFlag = 601
	bitbuf := data1.UperEncode()
	bs := bitbuf.Bytes()
	t.Logf("value vehicleEventFlag ::= 1001011001; data.uper is %4x\n", bs)
	if fmt.Sprintf("%4x", bs) == "9640" {
		t.Log("fixed encode pass")
	} else {
		t.Fatal("vehicleEventFlag encode error!")
	}

	var data2 vehicleEventFlag
	tempdata2 := data2.UperDecode(bitbuf)
	t.Logf("buffer %2x; Decode vehicleEventFlag is %d\n", bs, tempdata2)
	if tempdata2 == 601 {
		t.Log("fixed decode pass")
	} else {
		t.Fatal("vehicleEventFlag decode error!")
	}
} */

func TestBS_SIZE_FIXED_AND_EXT(t *testing.T) {
	var data1 laneSharing = 5
	bitbuf := data1.UperEncode()
	bs := bitbuf.Bytes()
	fmt.Println(bs)
	t.Logf("value laneSharing ::= 000101; data.uper is %2x\n", bs)
	if fmt.Sprintf("%2x", bs) == "50" {
		t.Log("fixed encode pass")
	} else {
		t.Fatal("laneSharing encode error!")
	}

	var data2 laneSharing
	tempdata2 := data2.UperDecode(bitbuf)
	t.Logf("buffer %2x; Decode laneSharing is %d\n", bs, tempdata2)
	if tempdata2 == 5 {
		t.Log("fixed decode pass")
	} else {
		t.Fatal("laneSharing decode error!")
	}
}

/**
DriveBehavior ::= BIT STRING {
	goStraightForward(0),
	laneChangingToLeft(1),
	laneChangingToRight(2),
	rampIn(3),
	rampOut(4),
	intersectionStraightThrough(5),
	intersectionTurnLeft(6),
	intersectionTurnRight(7),
	intersectionUTurn(8),
	stop-and-go(9),
	stop(10),
	slow-down(11) --,
	--parking(12)
} (SIZE(12,...))
*/

var DriveBehaviorMember common.Member

type DriveBehavior uint64

const (
	DriveBehavior_GoStraightForward           DriveBehavior = 1 << 0
	DriveBehavior_LaneChangingToLeft          DriveBehavior = 1 << 1
	DriveBehavior_LaneChangingToRight         DriveBehavior = 1 << 2
	DriveBehavior_RampIn                      DriveBehavior = 1 << 3
	DriveBehavior_RampOut                     DriveBehavior = 1 << 4
	DriveBehavior_IntersectionStraightThrough DriveBehavior = 1 << 5
	DriveBehavior_IntersectionTurnLeft        DriveBehavior = 1 << 6
	DriveBehavior_IntersectionTurnRight       DriveBehavior = 1 << 7
	DriveBehavior_IntersectionUTurn           DriveBehavior = 1 << 8
	DriveBehavior_Stop_and_go                 DriveBehavior = 1 << 9
	DriveBehavior_Stop                        DriveBehavior = 1 << 10
	DriveBehavior_Slow_down                   DriveBehavior = 1 << 11
	DriveBehavior_Parking                     DriveBehavior = 1 << 12
)

func (d *DriveBehavior) UperEncode() *common.BitBuffer {
	return UperEncodeBitString(DriveBehaviorMember, uint64(*d))
}

func (d *DriveBehavior) UperDecode(b *common.BitBuffer) {
	*d = DriveBehavior(UperDecodeBitString(DriveBehaviorMember, b))
}

func TestDriveBehaviorBitString(t *testing.T) {
	json.Unmarshal([]byte(`{
                "members": [
                    {
                        "expr_type": "A1TC_UNIVERVAL",
                        "meta_type": "AMT_VALUE",
                        "Identifier": "goStraightForward",
                        "value": 0
                    },
                    {
                        "expr_type": "A1TC_UNIVERVAL",
                        "meta_type": "AMT_VALUE",
                        "Identifier": "laneChangingToLeft",
                        "value": 1
                    },
                    {
                        "expr_type": "A1TC_UNIVERVAL",
                        "meta_type": "AMT_VALUE",
                        "Identifier": "laneChangingToRight",
                        "value": 2
                    },
                    {
                        "expr_type": "A1TC_UNIVERVAL",
                        "meta_type": "AMT_VALUE",
                        "Identifier": "rampIn",
                        "value": 3
                    },
                    {
                        "expr_type": "A1TC_UNIVERVAL",
                        "meta_type": "AMT_VALUE",
                        "Identifier": "rampOut",
                        "value": 4
                    },
                    {
                        "expr_type": "A1TC_UNIVERVAL",
                        "meta_type": "AMT_VALUE",
                        "Identifier": "intersectionStraightThrough",
                        "value": 5
                    },
                    {
                        "expr_type": "A1TC_UNIVERVAL",
                        "meta_type": "AMT_VALUE",
                        "Identifier": "intersectionTurnLeft",
                        "value": 6
                    },
                    {
                        "expr_type": "A1TC_UNIVERVAL",
                        "meta_type": "AMT_VALUE",
                        "Identifier": "intersectionTurnRight",
                        "value": 7
                    },
                    {
                        "expr_type": "A1TC_UNIVERVAL",
                        "meta_type": "AMT_VALUE",
                        "Identifier": "intersectionUTurn",
                        "value": 8
                    },
                    {
                        "expr_type": "A1TC_UNIVERVAL",
                        "meta_type": "AMT_VALUE",
                        "Identifier": "stop-and-go",
                        "value": 9
                    },
                    {
                        "expr_type": "A1TC_UNIVERVAL",
                        "meta_type": "AMT_VALUE",
                        "Identifier": "stop",
                        "value": 10
                    },
                    {
                        "expr_type": "A1TC_UNIVERVAL",
                        "meta_type": "AMT_VALUE",
                        "Identifier": "slow-down",
                        "value": 11
                    },
                    {
                        "expr_type": "A1TC_UNIVERVAL",
                        "meta_type": "AMT_VALUE",
                        "Identifier": "parking",
                        "value": 12
                    }
                ],
                "expr_type": "BIT_STRING",
                "meta_type": "AMT_TYPE",
                "constraints": {
                    "type": "SIZE_FIXED_AND_EXT",
                    "value": 12
                },
                "tag": {
                    "tag_mode": "TM_DEFAULT",
                    "tag_class": "",
                    "tag_value": ""
                },
                "Identifier": "DriveBehavior"
            }`), &DriveBehaviorMember)

	//rec1value DriveBehavior ::= '1001 0000 0000'B
	var driveBehavior = DriveBehavior_GoStraightForward | DriveBehavior_RampIn
	buf, err := driveBehavior.UperEncode().Result()
	if err != nil || fmt.Sprintf("%2X", buf) != "4800" {
		t.Fatalf("Bit String 编码错误 buf: %2X != 4800 err:%v", buf, err)
	} else {
		t.Log("Pass")
	}

	driveBehavior.UperDecode(common.NewBitBufferFromBytes(buf))
	if driveBehavior == DriveBehavior_GoStraightForward|DriveBehavior_RampIn {
		t.Log("Pass")
	} else {
		t.Fatalf("Bit String 解码错误 get:%v", driveBehavior)
	}

	// rec1value DriveBehavior ::= '1000 0000 0000 1'B
	var driveBehavior2 = DriveBehavior_GoStraightForward | DriveBehavior_Parking
	fmt.Println(driveBehavior2)
	buf2, err2 := driveBehavior2.UperEncode().Result()
	if err2 != nil || fmt.Sprintf("%2X", buf2) != "86C004" {
		t.Fatalf("Bit String 编码错误 buf: %2X != 86C004 err:%v", buf2, err2)
	} else {
		t.Log("Pass")
	}

	driveBehavior2.UperDecode(common.NewBitBufferFromBytes(buf2))
	if driveBehavior2 == DriveBehavior_GoStraightForward|DriveBehavior_Parking {
		t.Log("Pass")
	} else {
		t.Fatalf("Bit String 解码错误 get:%v", driveBehavior2)
	}
}
