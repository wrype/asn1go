package common

import (
	"errors"
	"fmt"
)

type Bit bool

type BitBuffer struct {
	bits []Bit
	err  error
}

func NewBitBuffer() *BitBuffer {
	return &BitBuffer{
		bits: make([]Bit, 0, 100),
		err:  nil,
	}
}

func NewBitBufferFromBytes(buf []byte) *BitBuffer {
	return NewBitBuffer().PushBytes(buf, len(buf)*8)
}

func (b *BitBuffer) PushBool(bit bool) *BitBuffer {
	return b.PushBit(Bit(bit))
}

func (b *BitBuffer) PushBit(bit Bit) *BitBuffer {
	b.bits = append(b.bits, bit)
	return b
}

//value: 0 or 1
func (b *BitBuffer) PushBit2(value byte) *BitBuffer {
	return b.PushBit(value != 0)
}

func (b *BitBuffer) PushBits(bit []Bit) *BitBuffer {
	b.bits = append(b.bits, bit...)
	return b
}

func (b *BitBuffer) PushBitBuffer(bitBuffer *BitBuffer) *BitBuffer {
	if bitBuffer == nil {
		return b.SetErrorText("PushBitBuffer arg is a nil value")
	}
	if bitBuffer.Error() != nil {
		return b.SetError(bitBuffer.Error())
	}
	b.bits = append(b.bits, bitBuffer.Bits()...)
	return b
}

func (b *BitBuffer) PushByte(byte2 byte, bit_len int) *BitBuffer {
	if bit_len > 8 || bit_len < 1 {
		return b.SetError(errors.New("PushByte bit_len invalid"))
	}
	return b.PushInteger(uint64(byte2), bit_len)
}

func (b *BitBuffer) PushBytes(buf []byte, bit_len int) *BitBuffer {
	for i, b2 := range buf {
		if i+1 == len(buf) && bit_len < len(buf)*8 {
			//最后一个字节长度不是 8
			b.PushByte(b2, bit_len%8)
		} else {
			b.PushByte(b2, 8)
		}
	}
	return b
}

//添加一个指定比特长度的整数
func (b *BitBuffer) PushInteger(value uint64, bit_len int) *BitBuffer {
	if bit_len > 64 || bit_len < 1 {
		return b.SetError(errors.New("PushInteger bit_len invalid"))
	}
	for i := bit_len - 1; i >= 0; i-- {
		if value&(1<<i) == 1<<i {
			b.bits = append(b.bits, true)
		} else {
			b.bits = append(b.bits, false)
		}
	}
	return b
}

//读取一个比特
func (b *BitBuffer) ShiftBit() Bit {
	if b.Length() == 0 {
		b.SetError(errors.New("当前比特流数据为空"))
		return false
	}
	b2 := b.bits[0]
	b.bits = b.bits[1:]
	return b2
}

func (b *BitBuffer) ShiftBool() bool {
	return bool(b.ShiftBit())
}

func (b *BitBuffer) ShiftBits(bit_len int) []Bit {
	if b.Length() < bit_len {
		b.SetError(errors.New("当前比特流数据过短"))
		return []Bit{}
	}
	b2 := b.bits[:bit_len]
	b.bits = b.bits[bit_len:]
	return b2
}

//根据比特数读取一个整数
func (b *BitBuffer) ShiftInteger(bit_len int) uint64 {
	if bit_len > 64 || bit_len < 1 {
		b.SetError(errors.New("ShiftInteger bit_len invalid"))
		return 0
	}
	var val uint64
	b2 := b.ShiftBits(bit_len)
	for i, bit := range b2 {
		if bit {
			val |= 1 << (bit_len - i - 1)
		}
	}
	return val
}

//根据比特数读取一个bye
func (b *BitBuffer) ShiftByte(bit_len int) byte {
	if bit_len > 8 || bit_len < 1 {
		b.SetError(errors.New("ShiftByte bit_len invalid"))
		return byte(0)
	}
	v := b.ShiftInteger(bit_len)
	return byte(v)
}

//根据比特数读取一个[]bye
func (b *BitBuffer) ShiftBytes(bit_len int) []byte {
	var bs []byte
	for i := 8; i <= bit_len; i += 8 {
		bs = append(bs, b.ShiftByte(8))
	}
	//最后一个比特长度不是 8
	if bit_len%8 != 0 {
		bs = append(bs, b.ShiftByte(bit_len%8))
	}
	return bs
}

//当前的比特长度
func (b *BitBuffer) Length() int {
	return len(b.bits)
}

func (b *BitBuffer) Bits() []Bit {
	return b.bits
}

//计算当前字节长度
func (b *BitBuffer) ByteLength() int {
	if b.Length()%8 == 0 {
		return b.Length() / 8
	}
	return b.Length()/8 + 1
}

//返回比特数组
func (b *BitBuffer) Bytes() []byte {
	var bs []byte
	//i 代表第 i 个 bye
	for i := 0; i < b.ByteLength(); i++ {
		var buf byte
		//j 代表每个 byte 的第 j 个 bit
		for j := 0; j < 8; j++ {
			if i*8+j < b.Length() && b.bits[i*8+j] {
				buf |= 1 << (7 - j)
			}
		}
		bs = append(bs, buf)
	}
	return bs
}

//外部可设置编码是否出错
func (b *BitBuffer) SetError(err error) *BitBuffer {
	//log.Panic(err)
	if b.err == nil {
		b.err = err
	}
	return b
}

//外部可设置编码是否出错
func (b *BitBuffer) SetErrorText(txt string) *BitBuffer {
	return b.SetError(errors.New(txt))
}

//外部可设置编码是否出错
func (b *BitBuffer) SetErrorTextf(format string, args ...interface{}) *BitBuffer {
	return b.SetErrorText(fmt.Sprintf(format, args...))
}

//最后需要通过该接口判断是否出错
func (b *BitBuffer) Error() error {
	return b.err
}

func (b *BitBuffer) Result() ([]byte, error) {
	return b.Bytes(), b.Error()
}

func (b *BitBuffer) String() string {
	if b.Error() != nil {
		return fmt.Sprintf("%v", b.Error())
	}
	str := ""
	for _, bit := range b.Bits() {
		if bit {
			str += "1"
		} else {
			str += "0"
		}
	}
	return str
}

//返回逆转后的 []Bit
func (b *BitBuffer) ReverseBits() []Bit {
	bits := make([]Bit, 0)
	for _, bit := range b.Bits() {
		bits = append(bits, bit)
	}
	return bits
}

func NewBitBufferFromBitString(bitString string) *BitBuffer {
	bits := make([]Bit, 0)
	for i, _ := range bitString {
		if bitString[i] == '0' {
			bits = append(bits, false)
		} else {
			bits = append(bits, true)
		}
	}
	return NewBitBuffer().PushBits(bits)
}
