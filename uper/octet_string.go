package uper

import (
	"github.com/wrype/asn1go/common"
)

func UperEncodeOctetString(member common.Member, val []byte) *common.BitBuffer {
	b := common.NewBitBuffer()
	if member.Constraints != nil {
		/* if member.Constraints.Max != member.Constraints.Min {
			if member.Constraints.IsOutOfRange(int64(len(val))) {
				if member.Constraints.HasExt() {
					b.SetErrorText("暂不支持选择扩展的数值")
				} else {
					b.SetErrorText("非法的数值")
				}
			}
		} else {
			if member.Constraints.Value != int64(len(val)) {
				if member.Constraints.HasExt() {
					b.SetErrorText("暂不支持选择扩展的数值")
				} else {
					b.SetErrorText("非法的数值")
				}
			}
		}
		if member.Constraints.HasExt() {
			if member.Constraints.Range() != 0 {
				b.PushBit(common.Bit(member.Constraints.IsOutOfRange(int64(len(val)))))
			} else {
				b.PushBit(common.Bit(int64(len(val)) != member.Constraints.Value))
			}
		} */
		//ext处理
		if member.Constraints.HasExt() {
			if member.Constraints.Range() != 0 {
				if member.Constraints.IsOutOfRange(int64(len(val))) {
					//变长ext编码，1+(len(val)->8位组)
					b.PushBit(common.Bit(true))
					realLen := 8
					/* for bits.Len64(uint64(len(val))) > realLen {
						realLen = realLen + 8
					} */
					b.PushInteger(uint64(len(val)), realLen)
				} else {
					b.PushBit(common.Bit(false))
				}
			} else {
				if member.Constraints.Value != int64(len(val)) {
					//变长ext编码，1+(len(val)->8位组)
					b.PushBit(common.Bit(true))
					realLen := 8
					/* for bits.Len64(uint64(len(val))) > realLen {
						realLen = realLen + 8
					} */
					b.PushInteger(uint64(len(val)), realLen)
				} else {
					b.PushBit(common.Bit(false))
				}
			}
		} else {
			if member.Constraints.Range() != 0 {
				if member.Constraints.IsOutOfRange(int64(len(val))) {
					b.SetErrorTextf("%v存在非法的数值：OCTET STRING len(%v) out of range", member.Identifier, len(val))
				}
			} else {
				if int64(len(val)) != member.Constraints.Value {
					b.SetErrorTextf("%v存在非法的数值：OCTET STRING len(%v) not eq %v", member.Identifier, len(val), member.Constraints.Value)
				}
			}
			if member.Constraints.Range() != 0 {
				/* templen := int64(member.Constraints.Range() + 1)
				lenbitlen := 0
				for templen != 0 {
					templen = templen >> 1
					lenbitlen++
				} */
				lenbitlen := member.Constraints.RangeBitLength()
				b.PushInteger(uint64(int64(len(val))-member.Constraints.Min), lenbitlen)
			}
		}

		b.PushBytes(val, len(val)*8)

	} else {
		b.PushInteger(uint64(len(val)), 8)
		b.PushBytes(val, len(val)*8)
	}
	return b
}

func UperDecodeOctetString(member common.Member, b *common.BitBuffer) []byte {
	result := make([]byte, 0)
	bitlen := uint64(0)
	lenbitlen := 0
	var ext common.Bit
	if member.Constraints != nil {
		if member.Constraints.HasExt() {
			ext = b.ShiftBit()
		}
		if ext {
			bitlen = b.ShiftInteger(8)
		} else {
			if member.Constraints.Range() != 0 {
				/* templen := int64(member.Constraints.Range()) + 1
				for templen != 0 {
					templen = templen >> 1
					lenbitlen++
				} */
				lenbitlen = member.Constraints.RangeBitLength()
				bitlen = b.ShiftInteger(lenbitlen) + uint64(member.Constraints.Min)
			} else {
				bitlen = uint64(member.Constraints.Value)
			}
		}
		result = b.ShiftBytes(int(bitlen) * 8)
	} else {
		dataLen := b.ShiftInteger(8)
		result = b.ShiftBytes(int(dataLen))
	}
	return result
}
