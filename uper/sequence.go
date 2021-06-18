package uper

import (
	"reflect"

	"github.com/wrype/asn1go/common"
)

func UperEncodeSequence(member common.Member, seqPtr interface{}) *common.BitBuffer {
	maskBuffer := common.NewBitBuffer() //OPTIONAL掩码
	bitBuffer := common.NewBitBuffer()  //所有非扩展成员
	seqValue := reflect.ValueOf(seqPtr)
	if seqValue.Type().Kind() != reflect.Ptr {
		return maskBuffer.SetErrorTextf("UperEncodeSequence %v必须传入指针类型", seqValue.Type().String())
	} else if seqValue.IsNil() {
		return maskBuffer.SetErrorTextf("UperEncodeSequence %v必须传入非空指针类型", seqValue.Type().String())
	}
	seqValue = seqValue.Elem()
	for _, m := range member.Members {
		//成员可能是个扩展标记 ... 不做处理
		if m.IsExtLabel() {
			//扩展标记后面的成员都是扩展成员,额外处理
			break
		}
		if _, ok := seqValue.Type().FieldByName(m.Identifier); !ok {
			return bitBuffer.SetErrorTextf("%s 反射中没有找到对应的成员: %s", seqValue.Type().String(), m.Identifier)
		}
		//获取该成员的值
		val := seqValue.FieldByName(m.Identifier)
		//判断该成员是否是可选的
		if m.Marker.IsOptional() {
			if val.Type().Kind() == reflect.Ptr {
				if val.IsNil() {
					maskBuffer.PushBit(false)
					continue
				} else if val.Elem().Type().Kind() == reflect.Slice && val.Elem().IsNil() {
					maskBuffer.PushBit(false)
					continue
				} else {
					maskBuffer.PushBit(true)
				}
			} else if val.Type().Kind() == reflect.Slice && val.IsNil() {
				maskBuffer.PushBit(false)
				continue
			} else {
				maskBuffer.PushBit(true)
			}
		}
		bitBuffer.PushBitBuffer(UperEncode(m, val))
	}
	isExtBuffer := common.NewBitBuffer() //是否选择了扩展
	if member.HasExtLabel() {
		extBufferr := encodeExtMember(seqValue, member)
		if extBufferr == nil {
			//没有选择扩展成员，扩展成员都为 nil
			isExtBuffer.PushBit(false)
			return isExtBuffer.PushBitBuffer(maskBuffer).PushBitBuffer(bitBuffer)
		} else {
			//选择了扩展成员，有一个扩展成员不为 nil
			isExtBuffer.PushBit(true)
			return isExtBuffer.PushBitBuffer(maskBuffer).PushBitBuffer(bitBuffer).PushBitBuffer(extBufferr)
		}
	}
	return maskBuffer.PushBitBuffer(bitBuffer)
}

func UperDecodeSequence(member common.Member, bitBuffer *common.BitBuffer, seqPtr interface{}) {
	hasExtBuffer := false
	seqValue := reflect.ValueOf(seqPtr)
	if seqValue.Type().Kind() != reflect.Ptr {
		bitBuffer.SetErrorTextf("UperEncodeSequence %v必须传入指针类型", seqValue.Type().String())
		return
	}
	seqValue = seqValue.Elem()
	if member.HasExtLabel() {
		hasExtBuffer = bitBuffer.ShiftBool()
	}
	//处理 OPTIONAL 标记
	maskBuffer := common.NewBitBuffer()
	if member.HasOptionalMember() {
		maskBuffer.PushBits(bitBuffer.ShiftBits(member.OptionalCount()))
	}
	for _, m := range member.Members {
		//成员可能是个扩展标记 ... 不做处理, 后面的 buffer 都是扩展的，做特殊处理
		if m.IsExtLabel() {
			break
		}
		if m.Marker.IsOptional() {
			exist := maskBuffer.ShiftBool()
			if !exist {
				continue
			}
		}
		if seq, ok := seqValue.Type().FieldByName(m.Identifier); ok {
			//获取该成员的值,可能是个空指针
			val := seqValue.FieldByName(m.Identifier)
			if val.Type().Kind() == reflect.Ptr && val.IsNil() {
				val.Set(reflect.New(seq.Type.Elem()))
			}
			UperDecode(m, bitBuffer, val)
		} else {
			bitBuffer.SetErrorTextf("反射中没有找到对应的成员: %s, %s", seqValue.Type().String(), m.Identifier)
			return
		}
	}
	if hasExtBuffer {
		decodeExtMember(seqValue, member, bitBuffer)
	}
}

func UperEncode(m common.Member, val reflect.Value) *common.BitBuffer {
	bitBuffer := common.NewBitBuffer()
	valIsPtr := val.Type().Kind() == reflect.Ptr
	if valIsPtr && val.IsNil() {
		return bitBuffer.SetErrorTextf("%s is nil ", m.Identifier)
	}
	switch m.MetaType {
	case common.AMT_TYPE: //基础类型，调用基础库编码
		switch m.ExprType {
		case common.AMT_TYPE_INTEGER:
			if valIsPtr {
				val = val.Elem()
			}
			bitBuffer.PushBitBuffer(UperEncodeInteger(m, val.Int()))
		case common.AMT_TYPE_BIT_STRING:
			if valIsPtr {
				val = val.Elem()
			}
			bitBuffer.PushBitBuffer(UperEncodeBitString(m, uint64(val.Int())))
		case common.AMT_TYPE_BOOLEAN:
			if valIsPtr {
				val = val.Elem()
			}
			bitBuffer.PushBitBuffer(UperEncodeBoolean(m, val.Bool()))
		case common.AMT_TYPE_CHOICE:
			if !valIsPtr {
				if val.CanAddr() {
					val = val.Addr()
				} else {
					return bitBuffer.SetErrorTextf("无法获取类型 %s 的指针", val.Type().String())
				}
			}
			bitBuffer.PushBitBuffer(UperEncodeChoice(m, val.Interface()))
		case common.AMT_TYPE_ENUMERATED:
			if valIsPtr {
				val = val.Elem()
			}
			bitBuffer.PushBitBuffer(UperEncodeEnumerated(m, val.Int()))
		case common.AMT_TYPE_IA5String:
			if valIsPtr {
				val = val.Elem()
			}
			bitBuffer.PushBitBuffer(UperEncodeIA5String(m, val.String()))
		case common.AMT_TYPE_OCTET_STRING:
			if valIsPtr {
				val = val.Elem()
			}
			bitBuffer.PushBitBuffer(UperEncodeOctetString(m, val.Bytes()))
		case common.AMT_TYPE_REAL:
			if valIsPtr {
				val = val.Elem()
			}
			bitBuffer.PushBitBuffer(UperEncodeReal(m, val.Float()))
		case common.AMT_TYPE_SEQUENCE:
			if !valIsPtr {
				if val.CanAddr() {
					val = val.Addr()
				} else {
					return bitBuffer.SetErrorTextf("无法获取类型 %s 的指针", val.Type().String())
				}
			}
			bitBuffer.PushBitBuffer(UperEncodeSequence(m, val.Interface()))
		case common.AMT_TYPE_SEQUENCE_OF:
			if !valIsPtr {
				if val.CanAddr() {
					val = val.Addr()
				} else {
					return bitBuffer.SetErrorTextf("无法获取类型 %s 的指针", val.Type().String())
				}
			}
			bitBuffer.PushBitBuffer(UperEncodeSequenceOf(m, val.Interface()))
		default:
			return bitBuffer.SetErrorTextf("不支持的ExprType %s ::= %s", m.Identifier, m.ExprType)
		}
	case common.ASN_TYPE: //自定义类型，调用 UperEncode 编码
		//如果不是一个指针时，获取该类型的指针
		if !valIsPtr {
			if val.CanAddr() {
				val = val.Addr()
			} else {
				return bitBuffer.SetErrorTextf("%v 必须是个指针类型实现 UperEncoder 接口", val.Type().String())
			}
		}
		if !val.Type().Implements(reflect.TypeOf(new(common.UperEncoder)).Elem()) {
			return bitBuffer.SetErrorTextf("%v 未实现 UperEncoder 接口", val.Type().String())
		}
		encoder := val.Interface().(common.UperEncoder)
		bs := encoder.UperEncode()
		bitBuffer.PushBitBuffer(bs)
	default:
		return bitBuffer.SetErrorTextf("不支持的 MetaType %s ::= %s", m.Identifier, m.MetaType)
	}
	return bitBuffer
}

func UperDecode(m common.Member, bitBuffer *common.BitBuffer, val reflect.Value) {
	valIsPtr := val.Type().Kind() == reflect.Ptr
	if !val.CanSet() {
		bitBuffer.SetErrorTextf("%s 不允许被修改", m.Identifier)
		return
	}
	switch m.MetaType {
	case common.AMT_TYPE: //基础类型，调用基础库解码
		switch m.ExprType {
		case common.AMT_TYPE_INTEGER:
			if valIsPtr {
				v := UperDecodeInteger(m, bitBuffer)
				val.Set(reflect.ValueOf(&v))
			} else {
				val.Set(reflect.ValueOf(UperDecodeInteger(m, bitBuffer)))
			}
		case common.AMT_TYPE_BIT_STRING:
			if valIsPtr {
				v := UperDecodeBitString(m, bitBuffer)
				val.Set(reflect.ValueOf(&v))
			} else {
				val.Set(reflect.ValueOf(UperDecodeBitString(m, bitBuffer)))
			}
		case common.AMT_TYPE_BOOLEAN:
			if valIsPtr {
				v := UperDecodeBoolean(m, bitBuffer)
				val.Set(reflect.ValueOf(&v))
			} else {
				val.Set(reflect.ValueOf(UperDecodeBoolean(m, bitBuffer)))
			}
		case common.AMT_TYPE_CHOICE:
			if !valIsPtr {
				if val.CanAddr() {
					val = val.Addr()
				} else {
					bitBuffer.SetErrorTextf("无法获取类型 %s 的指针", val.Type().String())
					return
				}
			}
			UperDecodeChoice(m, bitBuffer, val.Interface())
		case common.AMT_TYPE_ENUMERATED:
			if valIsPtr {
				v := UperDecodeEnumerated(m, bitBuffer)
				val.Set(reflect.ValueOf(&v))
			} else {
				val.Set(reflect.ValueOf(UperDecodeEnumerated(m, bitBuffer)))
			}
		case common.AMT_TYPE_IA5String:
			if valIsPtr {
				v := UperDecodeIA5String(m, bitBuffer)
				val.Set(reflect.ValueOf(&v))
			} else {
				val.Set(reflect.ValueOf(UperDecodeIA5String(m, bitBuffer)))
			}
		case common.AMT_TYPE_OCTET_STRING:
			if valIsPtr {
				v := UperDecodeOctetString(m, bitBuffer)
				val.Set(reflect.ValueOf(&v))
			} else {
				val.Set(reflect.ValueOf(UperDecodeOctetString(m, bitBuffer)))
			}
		case common.AMT_TYPE_REAL:
			if valIsPtr {
				v := UperDecodeReal(m, bitBuffer)
				val.Set(reflect.ValueOf(&v))
			} else {
				val.Set(reflect.ValueOf(UperDecodeReal(m, bitBuffer)))
			}
		case common.AMT_TYPE_SEQUENCE:
			if !valIsPtr {
				if val.CanAddr() {
					val = val.Addr()
				} else {
					bitBuffer.SetErrorTextf("无法获取类型 %s 的指针", val.Type().String())
					return
				}
			}
			UperDecodeSequence(m, bitBuffer, val.Interface())
		case common.AMT_TYPE_SEQUENCE_OF:
			if !valIsPtr {
				if val.CanAddr() {
					val = val.Addr()
				} else {
					bitBuffer.SetErrorTextf("无法获取类型 %s 的指针", val.Type().String())
					return
				}
			}
			UperDecodeSequenceOf(m, bitBuffer, val.Interface())
		default:
			bitBuffer.SetErrorTextf("不支持的ExprType %s ::= %s", m.Identifier, m.ExprType)
		}
	case common.ASN_TYPE: //自定义类型，可能为指针
		//如果不是一个指针时，获取该类型的指针
		if !valIsPtr {
			if val.CanAddr() {
				val = val.Addr()
			} else {
				bitBuffer.SetErrorTextf("%v 必须是个指针类型实现 UperDecoder 接口", val.Type().String())
				return
			}
		}
		if !val.Type().Implements(reflect.TypeOf(new(common.UperDecoder)).Elem()) {
			bitBuffer.SetErrorTextf("%v 未实现 UperDecoder 接口", val.Type().String())
			return
		}
		if val.IsNil() {
			val.Set(reflect.New(val.Type().Elem()))
		}
		decoder := val.Interface().(common.UperDecoder)
		decoder.UperDecode(bitBuffer)
	default:
		bitBuffer.SetErrorTextf("不支持的 MetaType %s ::= %s", m.Identifier, m.MetaType)
	}
}

func encodeExtMember(seqValue reflect.Value, member common.Member) *common.BitBuffer {
	maskBuffer := common.NewBitBuffer()
	bitBuffer := common.NewBitBuffer()
	if member.ExtMemberCount() <= 0 {
		return nil
	} else if member.ExtMemberCount() > 64 {
		//TODO 大于 64 个扩展成员的处理
		return maskBuffer.PushBit(true).SetErrorTextf("%v encodeExtMember 最多支持 64 个扩展成员", member.Identifier)
	} else {
		//小于等于 64 个成员，固定 6bit 长度位
		maskBuffer.PushBit(false).PushInteger(uint64(member.ExtMemberCount()-1), 6) //从 0 开始
	}
	isExtMember := false
	for i, m := range member.Members {
		if m.IsExtLabel() {
			isExtMember = true
			continue
		}
		//扩展标记后面的都是扩展成员
		if isExtMember {
			if m.IsExtLabel() {
				if len(member.Members) != i+1 {
					return maskBuffer.SetErrorTextf("%v encodeExtMember 第二个扩展标记后不允许还有成员！", member.Identifier)
				} else {
					break
				}
			}
			if _, ok := seqValue.Type().FieldByName(m.Identifier); !ok {
				return maskBuffer.SetErrorTextf("%s 反射中没有找到对应的成员: %s", seqValue.Type().String(), m.Identifier)
			}
			//获取该成员的值
			val := seqValue.FieldByName(m.Identifier)
			//扩展成员都是可选的
			if val.Type().Kind() != reflect.Ptr {
				return maskBuffer.SetErrorTextf("扩展成员 %s.%s 不是一个指针！", seqValue.Type().String(), m.Identifier)
			}
			maskBuffer.PushBit(common.Bit(!val.IsNil())) //掩码
			if !val.IsNil() {
				//选择了扩展
				valBytes, err := UperEncode(m, val.Elem()).Result()
				if err != nil {
					return maskBuffer.SetError(err)
				}
				bitBuffer.PushInteger(uint64(len(valBytes)), 8) //长度位，固定 8bit
				bitBuffer.PushBytes(valBytes, len(valBytes)*8)
			}
		}
	}
	//扩展成员都为空，直接返回 nil
	if len(bitBuffer.Bits()) == 0 {
		return nil
	}
	return maskBuffer.PushBitBuffer(bitBuffer)
}

func decodeExtMember(seqValue reflect.Value, member common.Member, bitBuffer *common.BitBuffer) {
	var maskBuffer *common.BitBuffer //扩展成员掩码
	over64 := bitBuffer.ShiftBool()
	if over64 {
		//TODO 处理大于 64 个扩展成员的情况
		bitBuffer.SetErrorTextf("%v decodeExtMember 不支持大于 64 个扩展成员的解码", member.Identifier)
	} else {
		memCount := bitBuffer.ShiftInteger(6) + 1 //小于 64 个扩展成员，固定读取 6bit 扩展成员数量
		maskBuffer = common.NewBitBuffer().PushBits(bitBuffer.ShiftBits(int(memCount)))
	}
	isExtMember := false
	for i, m := range member.Members {
		if m.IsExtLabel() {
			isExtMember = true
			continue
		}
		//扩展标记后面的都是扩展成员
		if isExtMember {
			if m.IsExtLabel() {
				if len(member.Members) != i+1 {
					bitBuffer.SetErrorTextf("%v decodeExtMember 第二个扩展标记后不允许还有成员！", member.Identifier)
				} else {
					break
				}
			}
			//判断该扩展成员是否存在
			if maskBuffer.ShiftBool() {
				valByteSize := bitBuffer.ShiftInteger(8) //长度
				valBitBuffer := common.NewBitBuffer().PushBits(bitBuffer.ShiftBits(int(valByteSize * 8)))
				UperDecode(m, valBitBuffer, seqValue.FieldByName(m.Identifier))
			}
		}
	}
	//还有未定义的成员，新版本的数据，忽略掉
	for len(maskBuffer.Bits()) > 0 {
		if maskBuffer.ShiftBool() {
			valByteSize := bitBuffer.ShiftInteger(8) //长度
			bitBuffer.ShiftBits(int(valByteSize * 8))
		}
	}
}
