package common

type UperEncoder interface {
	UperEncode() *BitBuffer
}

type UperDecoder interface {
	UperDecode(buffer *BitBuffer)
}
