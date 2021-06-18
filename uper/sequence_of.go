package uper

import (
	"reflect"

	"github.com/wrype/asn1go/common"
)

func UperEncodeSequenceOf(member common.Member, seqOfPtr interface{}) *common.BitBuffer {
	maskBuffer := common.NewBitBuffer()
	bitBuffer := common.NewBitBuffer()
	if member.Constraints == nil {
		return maskBuffer.SetErrorTextf("SEQUENCE OF %v没有对应的约束", member.Identifier)
	}
	seqListValue := reflect.ValueOf(seqOfPtr)
	if seqListValue.Type().Kind() != reflect.Ptr {
		return maskBuffer.SetErrorTextf("UperEncodeSequenceOf %v必须传入指针类型", seqListValue.Type().String())
	}
	seqListValue = seqListValue.Elem()
	//有扩展的时候
	if member.Constraints.HasExt() {
		if member.Constraints.IsOutOfRange(int64(seqListValue.Len())) {
			return maskBuffer.SetErrorTextf("UperEncodeSequenceOf %v不支持选择扩展", member.Identifier)
		} else {
			maskBuffer.PushBit(false)
		}
	}
	//处理数组为空或不在约束范围内时
	if member.Constraints.IsRange() { //变长有长度位
		if member.Constraints.IsOutOfRange(int64(seqListValue.Len())) {
			return maskBuffer.SetErrorTextf("%s 长度范围不符合约束(%d..%d) 当前长度是 %d", member.Identifier, member.Constraints.Min, member.Constraints.Max, seqListValue.Len())
		} else {
			//长度位
			maskBuffer.PushInteger(uint64(int64(seqListValue.Len())-member.Constraints.Min), member.Constraints.RangeBitLength())
		}
	} else { //定长没有长度位
		if member.Constraints.Value != int64(seqListValue.Len()) {
			return maskBuffer.SetErrorTextf("%s 长度范围不符合约束(%d) 当前长度是 %d", member.Identifier, member.Constraints.Value, seqListValue.Len())
		}
	}
	//判断数组内元素的类型，只处理自定义类型 ASN_TYPE
	if len(member.Members) != 1 || member.Members[0].MetaType != common.ASN_TYPE {
		return maskBuffer.SetErrorTextf("SEQUENCE OF %v只能是自定义类型", member.Identifier)
	}
	//遍历成员并编码
	for i := 0; i < seqListValue.Len(); i++ {
		seqValue := seqListValue.Index(i)
		if seqValue.Type().Kind() != reflect.Ptr {
			if seqValue.CanAddr() {
				seqValue = seqValue.Addr()
			} else {
				return maskBuffer.SetErrorTextf("无法获取 %v(SEQUENCE OF) 元素指针", seqValue.Type().String())
			}
		}
		if seqValue.Elem().Type().Implements(reflect.TypeOf(new(common.UperEncoder)).Elem()) {
			return maskBuffer.SetErrorTextf("%s 元素未实现 UperEncoder 接口", member.Identifier)
		}
		encoder := seqValue.Interface().(common.UperEncoder)
		bitBuffer.PushBitBuffer(encoder.UperEncode())
	}
	return maskBuffer.PushBitBuffer(bitBuffer)
}

func UperDecodeSequenceOf(member common.Member, bitBuffer *common.BitBuffer, seqOfPtr interface{}) {
	if member.Constraints == nil {
		bitBuffer.SetErrorTextf("SEQUENCE OF %v没有对应的约束", member.Identifier)
		return
	}
	seqListValue := reflect.ValueOf(seqOfPtr)
	if seqListValue.Type().Kind() != reflect.Ptr {
		bitBuffer.SetErrorTextf("UperEncodeSequenceOf %v必须传入指针类型", seqListValue.Type().String())
		return
	}
	//判断数组内元素的类型，只处理自定义类型 ASN_TYPE
	if len(member.Members) != 1 || member.Members[0].MetaType != common.ASN_TYPE {
		bitBuffer.SetErrorTextf("SEQUENCE OF %v只能是自定义类型", member.Identifier)
		return
	}
	seqListValue = seqListValue.Elem()
	//有扩展的时候
	if member.Constraints.HasExt() {
		ext := bitBuffer.ShiftBit() //读取扩展标记
		if ext {
			bitBuffer.SetErrorTextf("UperDecodeSequenceOf %v不支持选择扩展", member.Identifier)
			return
		}
	}
	//获取长度
	seqListLen := 0
	if member.Constraints.IsRange() {
		seqListLen = int(bitBuffer.ShiftInteger(member.Constraints.RangeBitLength()))
		seqListLen += int(member.Constraints.Min)
	} else {
		seqListLen = int(member.Constraints.Value)
	}
	//循环读取每个元素
	seqType := seqListValue.Type().Elem()
	for i := 0; i < seqListLen; i++ {
		seqValue := reflect.New(seqType)
		if !seqValue.Type().Implements(reflect.TypeOf(new(common.UperDecoder)).Elem()) {
			bitBuffer.SetErrorTextf("类型 %s 未实现 UperDecoder 接口", seqType.String())
			return
		}
		decoder := seqValue.Interface().(common.UperDecoder)
		decoder.UperDecode(bitBuffer)
		if bitBuffer.Error() == nil {
			seqListValue.Set(reflect.Append(seqListValue, seqValue.Elem()))
		} else {
			return
		}
	}
}
