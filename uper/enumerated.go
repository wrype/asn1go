package uper

import (
	"sort"

	"github.com/wrype/asn1go/common"
)

//枚举可以出现负数
//枚举值可以无序
func UperEncodeEnumerated(member common.Member, val int64) *common.BitBuffer {
	maskBuffer := common.NewBitBuffer()
	bitBuffer := common.NewBitBuffer()

	if member.HasExtLabel() {
		maskBuffer.PushBit(false)
	}

	vList := make([]int, 0)  //非扩展成员列表
	vxList := make([]int, 0) //扩展成员列表
	isExtMember := false
	for _, m := range member.Members {
		if !m.IsExtLabel() {
			if isExtMember {
				vxList = append(vxList, m.Value)
			} else {
				vList = append(vList, m.Value)
			}
		} else {
			isExtMember = true
		}
	}
	sort.Ints(vList)
	sort.Ints(vxList)
	isExist := false
	for index, v := range vList {
		if int64(v) == val {
			bitBuffer.PushInteger(uint64(index), member.MemberCountBitLength())
			isExist = true
			break
		}
	}
	if !isExist {
		//不在非扩展成员中，继续查找扩展成员
		for index, v := range vxList {
			if int64(v) == val {
				return common.NewBitBuffer().PushBool(true).PushInteger(uint64(index), 7) //固定 7bit
			}
		}
		return bitBuffer.SetErrorTextf("Enumerated %v没有匹配数据", member.Identifier)
	}

	return maskBuffer.PushBitBuffer(bitBuffer)
}

func UperDecodeEnumerated(member common.Member, b *common.BitBuffer) int64 {
	isExt := false
	if member.HasExtLabel() {
		isExt = b.ShiftBool()
	}

	vList := make([]int, 0)  //非扩展成员列表
	vxList := make([]int, 0) //扩展成员列表
	isExtMember := false
	for _, m := range member.Members {
		if !m.IsExtLabel() {
			if isExtMember {
				vxList = append(vxList, m.Value)
			} else {
				vList = append(vList, m.Value)
			}
		} else {
			isExtMember = true
		}
	}
	sort.Ints(vList)
	sort.Ints(vxList)

	if isExt {
		//扩展的情况
		index := int(b.ShiftInteger(7))
		if index >= len(vxList) {
			b.SetErrorTextf("Enumerated %v超出数据范围", member.Identifier)
			return -1
		}
		return int64(vxList[index])
	} else {
		//非扩展的情况
		index := int(b.ShiftInteger(member.MemberCountBitLength()))
		if index < 0 || index >= len(vList) {
			b.SetErrorTextf("Enumerated %v超出数据范围", member.Identifier)
			return -1
		}
		return int64(vList[index])
	}
}
