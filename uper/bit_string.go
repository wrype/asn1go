package uper

import (
	"math/bits"

	"github.com/wrype/asn1go/common"
)

func UperEncodeBitString(member common.Member, value uint64) *common.BitBuffer {
	b := common.NewBitBuffer()
	if member.Constraints != nil {
		var lenbitLen int
		bitLen := int64(bits.Len64(value))
		if member.Constraints.IsRange() {
			if int(bitLen) < member.Constraints.RangeBitLength() {
				lenbitLen = member.Constraints.RangeBitLength()
			} else {
				lenbitLen = int(bitLen)
			}
		} else {
			if bitLen < member.Constraints.Value {
				lenbitLen = int(member.Constraints.Value)
			} else {
				lenbitLen = int(bitLen)
			}
		}
		temp := common.NewBitBuffer()
		for i := 0; i < lenbitLen; i++ {
			temp.PushBit((value & (1 << i) >> i) == 1)
		}
		newVal := temp.ShiftInteger(lenbitLen)
		bitLen = int64(bits.Len64(newVal))
		/* if member.Constraints.IsRange() {
			if bitLen < member.Constraints.Min || bitLen > member.Constraints.Max {
				if member.Constraints.HasExt() {
					b.SetErrorText("暂不支持选择扩展的数值")
				} else {
					b.SetErrorText("非法的数值")
				}
			}
		} else {
			if bitLen != member.Constraints.Value {
				if member.Constraints.HasExt() {
					b.SetErrorText("暂不支持选择扩展的数值")
				} else {
					b.SetErrorText("非法的数值")
				}
			}
		} */

		if member.Constraints.HasExt() {
			if member.Constraints.IsRange() {
				if member.Constraints.IsOutOfRange(bitLen) {
					//变长ext编码，1+(len(val)->8位组)
					b.PushBit(common.Bit(true))
					realLen := 8
					/* for bits.Len64(uint64(bitLen)) > realLen {
						realLen = realLen + 8
					} */
					b.PushInteger(uint64(bitLen), realLen)
				} else {
					b.PushBit(common.Bit(false))
				}
			} else {
				if member.Constraints.Value < bitLen {
					b.PushBit(common.Bit(true))
					realLen := 8
					/* for bits.Len64(uint64(bitLen)) > realLen {
						realLen = realLen + 8
					} */
					b.PushInteger(uint64(bitLen), realLen)
				} else {
					b.PushBit(common.Bit(false))
					bitLen = member.Constraints.Value
				}
			}
		} else {
			if member.Constraints.IsRange() {
				if member.Constraints.IsOutOfRange(int64(bitLen)) {
					b.SetErrorTextf("%v存在非法的数值：BITSTRING len(%v) out of range", member.Identifier, bitLen)
				}
			} else {
				bitLen = member.Constraints.Value
			}
			if member.Constraints.IsRange() {
				/* templen := int64(member.Constraints.Range() + 1)
				lenbitlen := 0
				for templen != 0 {
					templen = templen >> 1
					lenbitlen++
				} */
				b.PushInteger(uint64(int64(bitLen)-member.Constraints.Min), int(lenbitLen))
			}
		}

		/* if member.Constraints.HasExt() {
			if member.Constraints.IsRange() {
				b.PushBit(common.Bit(bitLen < member.Constraints.Min || bitLen > member.Constraints.Max))
			} else {
				b.PushBit(common.Bit(bitLen != member.Constraints.Value))
			}

		}
		if member.Constraints.IsRange() {
			b.PushInteger(uint64(bitLen-member.Constraints.Min), lenbitLen)
		} */
		b.PushInteger(newVal, int(bitLen))

	} else {
		b.SetErrorTextf("BitString %v没有对应的约束", member.Identifier)
	}
	return b
}

func UperDecodeBitString(member common.Member, b *common.BitBuffer) uint64 {
	bitLen := uint64(0)
	lenBitLen := 0
	var result uint64
	var ext common.Bit
	if member.Constraints != nil {
		if member.Constraints.HasExt() {
			ext = b.ShiftBit()
		}
		if ext {
			bitLen = b.ShiftInteger(8)
			lenBitLen = int(bitLen)
		} else {
			if member.Constraints.IsRange() {
				lenBitLen = member.Constraints.RangeBitLength()
				bitLen = b.ShiftInteger(lenBitLen) + uint64(member.Constraints.Min)

			} else {
				bitLen = uint64(member.Constraints.Value)
				lenBitLen = int(bitLen)
			}
		}

		result = b.ShiftInteger(int(bitLen))
		temp := common.NewBitBuffer()
		for i := 0; i < lenBitLen; i++ {
			temp.PushBit((result & (1 << i) >> i) == 1)
		}
		result = temp.ShiftInteger(lenBitLen)
	} else {
		b.SetErrorTextf("BitString %v没有对应的约束", member.Identifier)
	}
	return result
}
