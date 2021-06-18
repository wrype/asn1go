package uper

import (
	"fmt"
	"math"
	"strconv"
	"strings"

	"github.com/wrype/asn1go/common"
)

const (
	minus  uint8 = 0x2D
	point  uint8 = 0x2E
	CharE  uint8 = 0x45
	plus   uint8 = 0x2B
	prefix uint8 = 0x03
)

func UperEncodeReal(member common.Member, val float64) *common.BitBuffer {
	mask := common.NewBitBuffer()
	b := common.NewBitBuffer()
	var length uint8 = 0
	if val == 0 {
		b.PushByte(length, 8)
		return b
	} else {
		//处理正负值
		if val < 0 {
			val = -val
			b.PushByte(minus, 8)
			length++
		}
	}
	strarr := strings.Split(fmt.Sprintf("%v", val), ".")
	//只有整数部分
	if len(strarr) == 1 {
		//确认整数尾部是否有0
		var zeronum int = 0
		i := len(strarr[0])
		for ; i > 0; i-- {
			if strarr[0][i-1] == '0' {
				zeronum++
			} else {
				break
			}
		}
		//转化得到的应为 BaseNumber.E+zeronum
		BaseNumber := strarr[0][0:i]
		for _, char := range BaseNumber {
			b.PushByte(uint8(char), 8)
			length++
		}
		b.PushByte(point, 8)
		b.PushByte(CharE, 8)
		b.PushByte(plus, 8)
		length += 4
		zeronumstr := strconv.Itoa(zeronum)
		for _, char := range zeronumstr {
			b.PushByte(uint8(char), 8)
			length++
		}

		//往最前面塞入长度和前缀
		mask.PushByte(length, 8)
		mask.PushByte(prefix, 8)
		return mask.PushBitBuffer(b)
	} else {
		//有尾数，转化得到的应为 BaseNumber.E-zeronum
		zeronum := len(strarr[1])
		//底数部分
		num, _ := strconv.ParseInt(strarr[0]+strarr[1], 10, 64)
		BaseNumber := strconv.FormatInt(num, 10)
		for _, char := range BaseNumber {
			b.PushByte(uint8(char), 8)
			length++
		}
		b.PushByte(point, 8)
		b.PushByte(CharE, 8)
		b.PushByte(minus, 8)
		length += 4

		zeronumstr := strconv.Itoa(zeronum)
		for _, char := range zeronumstr {
			b.PushByte(uint8(char), 8)
			length++
		}

		//往最前面塞入长度和前缀
		mask.PushByte(length, 8)
		mask.PushByte(prefix, 8)
		return mask.PushBitBuffer(b)
	}
}

func UperDecodeReal(member common.Member, b *common.BitBuffer) float64 {
	//去掉前缀和长度
	_ = b.ShiftBytes(16)
	//读取底数部分
	baseStr := []byte{}
	for {
		char := b.ShiftByte(8)
		if point == char {
			break
		} else {
			baseStr = append(baseStr, char)
		}
	}
	//去掉E
	_ = b.ShiftBits(8)
	zerostr := []byte{}
	for {
		char := b.ShiftByte(8)
		if byte(0) == char {
			break
		}
		zerostr = append(zerostr, char)

	}
	baseNum, _ := strconv.Atoi(string(baseStr))
	zeroNum, _ := strconv.Atoi(string(zerostr))
	return float64(baseNum) * math.Pow10(zeroNum)
}
