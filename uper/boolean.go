package uper

import (
	"github.com/wrype/asn1go/common"
)

func UperEncodeBoolean(member common.Member, value bool) *common.BitBuffer {
	return common.NewBitBuffer().PushBit(common.Bit(value))
}

func UperDecodeBoolean(member common.Member, b *common.BitBuffer) bool {
	return bool(b.ShiftBit())
}
