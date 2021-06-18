package uper

import (
	"math/bits"

	"github.com/wrype/asn1go/common"
)

func UperEncodeInteger(member common.Member, val int64) *common.BitBuffer {
	b := common.NewBitBuffer()
	if member.Constraints != nil {
		if member.Constraints.IsOutOfRange(val) {
			if member.Constraints.HasExt() {
				if val < 0 {
					//TODO 负数的情况？
					return b.SetErrorText("INTEGER 不支持负数的扩展编码")
				}
				b.PushBool(true)
				valBitSize := bits.Len64(uint64(val))
				valBitSize += 8 - valBitSize%8
				b.PushInteger(uint64(valBitSize/8), 8)
				return b.PushInteger(uint64(val), valBitSize)
			} else {
				return b.SetErrorTextf("%v存在非法的数值：INTEGER(%v) out of range", member.Identifier, val)
			}
		}
		if member.Constraints.HasExt() {
			//有扩展标记的需要一个 bit 代表是否选择了扩展
			b.PushBit(false)
		}
		//实际的数值，从 0 开始
		valuec := val - member.Constraints.Min
		b.PushInteger(uint64(valuec), member.Constraints.RangeBitLength())
	} else {
		b.SetErrorTextf("Integer %v没有对应的约束", member.Identifier)
	}
	return b
}

func UperDecodeInteger(member common.Member, b *common.BitBuffer) int64 {
	if member.Constraints != nil {
		if member.Constraints.HasExt() {
			//读取扩展标记
			ext := b.ShiftBit()
			if ext {
				valByteSize := b.ShiftInteger(8)
				valBuffer := b.ShiftBytes(int(valByteSize * 8))
				val := common.NewBitBufferFromBytes(valBuffer).ShiftInteger(int(valByteSize * 8))
				//TODO 负数的情况
				return int64(val)
			} else {
				valuec := b.ShiftInteger(member.Constraints.RangeBitLength())
				return int64(valuec) + member.Constraints.Min
			}
		} else {
			valuec := b.ShiftInteger(member.Constraints.RangeBitLength())
			return int64(valuec) + member.Constraints.Min
		}
	} else {
		b.SetErrorTextf("Integer %v没有对应的约束", member.Identifier)
	}
	return 0
}
