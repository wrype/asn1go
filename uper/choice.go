package uper

import (
	"reflect"

	"github.com/wrype/asn1go/common"
)

//ptr 是具体类型的指针
func UperEncodeChoice(member common.Member, ptr interface{}) *common.BitBuffer {
	buf := common.NewBitBuffer()
	elem := reflect.ValueOf(ptr)
	if elem.Type().Kind() != reflect.Ptr {
		return buf.SetErrorText(member.Identifier + "必须传入指针类型")
	}
	elem = elem.Elem()

	selectMemberIndex := -1
	var selectMember *common.Member
	isSelectExtMember := false
	// 成员选择和扩展标记统计
	index := 0
	for _, mem := range member.Members {
		//成员可能是个扩展标记 ... 不做处理
		//扩展标记后面的成员都是扩展成员
		if mem.IsExtLabel() {
			isSelectExtMember = true
			index = 0
			continue
		}
		//获取该成员的值并判断是否为选择项
		field := elem.FieldByName(mem.Identifier)
		switch field.Type().Kind() {
		case reflect.Ptr, reflect.Slice:
		default:
			return buf.SetErrorText(mem.Identifier + "无效成员类型")
		}
		if !field.IsNil() {
			selectMemberIndex = index
			selectMember = &mem
			break
		}
		index++
	}
	if selectMember == nil {
		return buf.SetErrorTextf("Choice %v必须有一个不为空的成员", member.Identifier)
	}
	//preamble编码
	if member.HasExtLabel() {
		buf.PushBool(isSelectExtMember)
	}
	if !isSelectExtMember {
		//索引号编码
		buf.PushInteger(uint64(selectMemberIndex), member.MemberCountBitLength())
		//内容编码
		field := elem.FieldByName(selectMember.Identifier)
		buf.PushBitBuffer(UperEncode(*selectMember, field))
	} else {
		//扩展成员编码
		//1bit 代表扩展成员数量是否大于 64 个
		//6bit 扩展第几个成员，从 0 开始
		if member.ExtMemberCount() > 64 {
			//TODO 大于 64 个扩展成员的处理
			return buf.SetErrorTextf("Choice %v不支持大于 64 个扩展成员的编码", member.Identifier)
		} else {
			buf.PushBit(false)                            //小于 64 个扩展成员
			buf.PushInteger(uint64(selectMemberIndex), 6) //扩展成员索引，固定 6bit
			field := elem.FieldByName(selectMember.Identifier)
			valBuf, err := UperEncode(*selectMember, field).Result()
			if err != nil {
				return buf.SetError(err)
			}
			buf.PushInteger(uint64(len(valBuf)), 8) //内容长度，固定 8bit
			buf.PushBytes(valBuf, len(valBuf)*8)    //内容
		}
	}
	return buf
}

//ptr 是具体类型的指针
func UperDecodeChoice(member common.Member, buf *common.BitBuffer, ptr interface{}) {
	elem := reflect.ValueOf(ptr)
	if elem.Type().Kind() != reflect.Ptr {
		buf.SetErrorText(member.Identifier + "必须传入指针类型")
		return
	}
	elem = elem.Elem()

	if member.HasExtLabel() {
		if buf.ShiftBool() {
			//选择了扩展成员
			isOver64 := buf.ShiftBool()
			if isOver64 {
				buf.SetErrorTextf("Choice %v不支持大于 64 个扩展成员的编码", member.Identifier)
				return
			}
			extMemberIndex := int(buf.ShiftInteger(6)) + 1 + member.MemberCount() //6bit 读取扩展成员位置
			extMemberValueLength := buf.ShiftInteger(8)                           //8bit 扩展成员编码长度
			extBuf := common.NewBitBufferFromBytes(buf.ShiftBytes(int(extMemberValueLength * 8)))
			mem := member.Members[extMemberIndex]
			if fieldStruct, ok := elem.Type().FieldByName(mem.Identifier); ok {
				field := elem.FieldByName(mem.Identifier)
				// 为空指针的成员字段(指针)创建对象
				switch fieldStruct.Type.Kind() {
				case reflect.Ptr:
					field.Set(reflect.New(fieldStruct.Type.Elem()))
				case reflect.Slice:
					field.Set(reflect.MakeSlice(field.Type(), 0, 0))
				default:
					buf.SetErrorText(mem.Identifier + "无效成员类型")
					return
				}

				UperDecode(mem, extBuf, field)
				if extBuf.Error() != nil {
					buf.SetError(extBuf.Error())
				}
			} else {
				//未在定义的扩展成员，忽略掉
			}
			return
		}
	}

	// 成员编号
	memId := buf.ShiftInteger(member.MemberCountBitLength())
	mem := member.Members[memId]
	// 是否存在相应成员字段
	if fieldStruct, ok := elem.Type().FieldByName(mem.Identifier); ok {
		field := elem.FieldByName(mem.Identifier)
		// 为空指针的成员字段(指针)创建对象
		switch fieldStruct.Type.Kind() {
		case reflect.Ptr:
			field.Set(reflect.New(fieldStruct.Type.Elem()))
		case reflect.Slice:
			field.Set(reflect.MakeSlice(field.Type(), 0, 0))
		default:
			buf.SetErrorText(mem.Identifier + "无效成员类型")
			return
		}

		UperDecode(mem, buf, field)
	} else {
		buf.SetErrorText("反射中没有找到对应的成员: " + mem.Identifier)
		return
	}
	// if _, ok := elem.Type().FieldByName(mem.Identifier); ok {
	// 	fieldType := elem.FieldByName(mem.Identifier).Type()
	// field := reflect.ValueOf(reflect.New(fieldType.Elem()))
	// 	elem.FieldByName(mem.Identifier).Set(reflect.New(fieldType.Elem()))
	// 	UperDecode(mem, buf, elem.FieldByName(mem.Identifier).Elem())

}
