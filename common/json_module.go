package common

import (
	"math/bits"
)

type Module struct {
	ModuleName string   `json:"moduleName"`
	Flags      string   `json:"flags"` //MSF_AUTOMATIC_TAGS
	Members    []Member `json:"members,omitempty"`
}

type Member struct {
	Members     []Member    `json:"members,omitempty"`     //包含的字段
	ExprType    ExprType    `json:"expr_type"`             // AMT_TYPE 类型 INTEGER,BOOLEAN,CHOICE,ENUMERATED,SEQUENCE,IA5String,OCTET_STRING,BIT_STRING,REAL,SEQUENCE_OF
	MetaType    string      `json:"meta_type"`             // AMT_TYPE（原始 ASN 类型） ASN_TYPE(自定义类型) AMT_VALUE(数字类型如 ENUMERATED {solid(1), liquid(2), gas(3)})
	Constraints *Constraint `json:"constraints,omitempty"` //约束
	Identifier  string      `json:"Identifier"`            //字段名称，注意标识符为 ...(IDENTIFY_EXTENSIBLE) 时需特殊处理扩展的情况
	Marker      *Marker     `json:"marker,omitempty"`      //"marker": { "flags": "EM_OPTIONAL" }
	Value       int         `json:"value"`                 //meta_type 为 AMT_VALUE 需要处理该字段,若没设置时,value 为 -1
}

type Constraint struct {
	Type Type `json:"type"` //SIZE_FIXED SIZE_RANG SIZE_FIXED_AND_EXT SIZE_RANG_AND_EXT REAL_WITH_COMPONENTS
	//SIZE_FIXED SIZE_FIXED_AND_EXT 定长时包含 Value
	Value int64 `json:"value"`
	//SIZE_RANG SIZE_RANG_AND_EXT 变长时包含值范围
	Min int64 `json:"min"`
	Max int64 `json:"max"`
	//Type 为 REAL_WITH_COMPONENTS 时的约束
	Base     int `json:"base"`     //10
	Exponent int `json:"exponent"` //-2
}

type Marker struct {
	Flags Flag `json:"flags"`
}

const (
	AMT_TYPE  = "AMT_TYPE"  //原始 ASN 类型 : {"expr_type":"INTEGER","meta_type":"AMT_TYPE","constraints":{"type":"SIZE_RANG","min":1,"max":255},"Identifier":"Speed"}
	ASN_TYPE  = "ASN_TYPE"  //自定义类型
	AMT_VALUE = "AMT_VALUE" //数字类型如 ENUMERATED {solid(1), liquid(2), gas(3)} 其中成员 solid 定义: {"expr_type":"A1TC_UNIVERVAL","meta_type":"AMT_VALUE","Identifier":"solid","value":1}
)

type ExprType string

const (
	AMT_TYPE_INTEGER      ExprType = "INTEGER"
	AMT_TYPE_BOOLEAN      ExprType = "BOOLEAN"
	AMT_TYPE_CHOICE       ExprType = "CHOICE"
	AMT_TYPE_ENUMERATED   ExprType = "ENUMERATED"
	AMT_TYPE_SEQUENCE     ExprType = "SEQUENCE"
	AMT_TYPE_IA5String    ExprType = "IA5String"
	AMT_TYPE_OCTET_STRING ExprType = "OCTET_STRING"
	AMT_TYPE_BIT_STRING   ExprType = "BIT_STRING"
	AMT_TYPE_REAL         ExprType = "REAL"
	AMT_TYPE_SEQUENCE_OF  ExprType = "SEQUENCE_OF"
)

type Type string

const (
	SIZE_FIXED           Type = "SIZE_FIXED"           //定长 SIZE(5)
	SIZE_RANG            Type = "SIZE_RANG"            //变长 SIZE(1..8)
	SIZE_FIXED_AND_EXT   Type = "SIZE_FIXED_AND_EXT"   //定长可扩展 SIZE(5,...) Max==Min
	SIZE_RANG_AND_EXT    Type = "SIZE_RANG_AND_EXT"    //变长可扩展 SIZE(1..7,...) //暂时没有这个
	REAL_WITH_COMPONENTS Type = "REAL_WITH_COMPONENTS" //Type 为 REAL_WITH_COMPONENTS 时包含 Base 和 Exponent 约束
)

type Flag string

const (
	EM_NOMARK   Flag = "EM_NOMARK"
	EM_OPTIONAL Flag = "EM_OPTIONAL"
)

const IDENTIFY_EXTENSIBLE = "..."

//判断约束是否是变长的
func (c *Constraint) IsRange() bool {
	if c != nil {
		if c.Type == SIZE_RANG || c.Type == SIZE_RANG_AND_EXT {
			return true
		}
	}
	return false
}

//获取约束的范围 0~range
func (c *Constraint) Range() uint64 {
	if c != nil {
		return uint64(c.Max - c.Min)
	}
	return 0
}

//获取该范围 uper 编码需要的比特长度
func (c *Constraint) RangeBitLength() int {
	return bits.Len64(c.Range())
}

//是否可扩展
func (c *Constraint) HasExt() bool {
	if c != nil {
		return c.Type == SIZE_RANG_AND_EXT || c.Type == SIZE_FIXED_AND_EXT
	}
	return false
}

//判断某个值是否超出了指定范围
func (c *Constraint) IsOutOfRange(val int64) bool {
	if c != nil {
		return val > c.Max || val < c.Min
	}
	return false
}

//判断该成员是否是个扩展标记
func (m Member) IsExtLabel() bool {
	return m.Identifier == IDENTIFY_EXTENSIBLE || m.ExprType == "A1TC_EXTENSIBLE"
}

//Sequence 判断是否有可扩展的标记 ...
func (m Member) HasExtLabel() bool {
	for _, member := range m.Members {
		if member.IsExtLabel() {
			return true
		}
	}
	return false
}

//扩展标记后面的成员都是可扩展成员
func (m Member) ExtMemberCount() int {
	count := 0
	isExtMember := false
	for _, member := range m.Members {
		if member.IsExtLabel() {
			isExtMember = true
		} else if isExtMember {
			count++
		}
	}
	return count
}

//获取非扩展的成员数量
func (m Member) MemberCount() int {
	count := 0
	for _, member := range m.Members {
		if !member.IsExtLabel() {
			count++
		} else {
			return count
		}
	}
	return count
}

//非扩展的成员数量需要的编码 bit_len
func (m Member) MemberCountBitLength() int {
	return bits.Len64(uint64(m.MemberCount() - 1))
}

//返回有可选标记(OPTIONAL)的子成员数,不包含扩展成员
func (m Member) OptionalCount() int {
	count := 0
	for _, member := range m.Members {
		if member.IsExtLabel() {
			return count
		}
		if member.Marker.IsOptional() {
			count++
		}
	}
	return count
}

//是否包含可选的子成员
func (m Member) HasOptionalMember() bool {
	return m.OptionalCount() != 0
}

//判断是否是可选的
func (m *Marker) IsOptional() bool {
	if m != nil {
		return m.Flags == EM_OPTIONAL
	}
	return false
}
