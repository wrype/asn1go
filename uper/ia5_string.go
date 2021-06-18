package uper

import (
	"bytes"

	"github.com/wrype/asn1go/common"
)

const (
	emptyStr   = ""
	bitsOfChar = 7
)

//IA5String 长度位要计算需要多少个bit，之后每7个比特代表一个ascii字符（去掉了前面一位0）
//扩展时，改变的只有长度位，固定长度 8bit
func UperEncodeIA5String(member common.Member, val string) *common.BitBuffer {
	extBuffer := common.NewBitBuffer() //标记位，代表是否扩展，不可扩展时不编码
	lenBuffer := common.NewBitBuffer() //长度位，定长时不编码，变长时根据是否是扩展编码成不同长度，范围外时固定 8bit，在范围内时按范围的最小 bit 数编码
	bitBuffer := common.NewBitBuffer() //每7个比特代表一个ascii字符
	if member.Constraints == nil {
		return lenBuffer.SetErrorTextf("IA5String %v没有对应的约束", member.Identifier)
	}
	strLen := int64(len(val))

	if member.Constraints.HasExt() {
		extBuffer.PushBool(member.Constraints.IsOutOfRange(strLen))
	}

	if !member.Constraints.IsRange() {
		if member.Constraints.Value != strLen {
			return lenBuffer.SetErrorTextf("IA5String %v非法的字符串，破坏定长约束", member.Identifier)
		}
	} else {
		if member.Constraints.IsOutOfRange(strLen) {
			if !member.Constraints.HasExt() {
				return lenBuffer.SetErrorTextf("IA5String %v非法的字符串，破坏定长约束", member.Identifier)
			} else {
				lenBuffer.PushInteger(uint64(strLen), 8)
			}
		} else {
			lenBuffer.PushInteger(uint64(strLen-member.Constraints.Min), member.Constraints.RangeBitLength())
		}
	}

	// 字符串内容
	for id := 0; id < int(strLen); id++ {
		if val[id] > 127 {
			return bitBuffer.SetErrorTextf("IA5String %v非法的字符串，存在无效字符", member.Identifier)
		}
		bitBuffer.PushByte(val[id], bitsOfChar)
	}
	return extBuffer.PushBitBuffer(lenBuffer.PushBitBuffer(bitBuffer))
}

func UperDecodeIA5String(member common.Member, b *common.BitBuffer) string {
	if member.Constraints == nil {
		b.SetErrorTextf("IA5String %v没有对应的约束", member.Identifier)
		return emptyStr
	}
	strLen := 0
	isOutOfRang := false
	if member.Constraints.HasExt() {
		isOutOfRang = b.ShiftBool()
	}

	if isOutOfRang {
		strLen = int(b.ShiftInteger(8))
	} else {
		if !member.Constraints.IsRange() {
			strLen = int(member.Constraints.Value)
		} else {
			strLen = int(int64(b.ShiftInteger(member.Constraints.RangeBitLength())) + member.Constraints.Min)
		}
	}

	var buffer bytes.Buffer
	for ; strLen > 0; strLen-- {
		char := b.ShiftInteger(bitsOfChar)
		err := buffer.WriteByte(byte(char))
		if err != nil {
			b.SetErrorTextf("IA5String %v buffer写入失败: %v", member.Identifier, err)
			return emptyStr
		}
	}
	return buffer.String()
}
