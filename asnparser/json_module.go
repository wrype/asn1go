package asnparser

type Module struct {
	ModuleName String     `json:"moduleName"`
	Flags      String     `json:"flags"` //MSF_AUTOMATIC_TAGS
	Members    MemberList `json:"members,omitempty"`
}

type Member struct {
	Members     MemberList  `json:"members,omitempty"`     //包含的字段
	ExprType    String      `json:"expr_type"`             // AMT_TYPE 类型 INTEGER,BOOLEAN,CHOICE,ENUMERATED,SEQUENCE,IA5String,OCTET_STRING,BIT_STRING,REAL,SEQUENCE_OF
	MetaType    String      `json:"meta_type"`             // AMT_TYPE（原始 ASN 类型） ASN_TYPE(自定义类型) AMT_VALUE(数字类型如 ENUMERATED {solid(1), liquid(2), gas(3)})
	Constraints *Constraint `json:"constraints,omitempty"` //约束
	Identifier  String      `json:"Identifier"`            //字段名称，注意标识符为 ...(IDENTIFY_EXTENSIBLE) 时需特殊处理扩展的情况
	Marker      *Marker     `json:"marker,omitempty"`      //"marker": { "flags": "EM_OPTIONAL" }
	Value       Int64       `json:"value"`                 //meta_type 为 AMT_VALUE 需要处理该字段,若没设置时,value 为 -1
	//Tag         *Tag        `json:"tag,omitempty"`         //tag 暂时没用
}

type Constraint struct {
	Type String `json:"type"` //SIZE_FIXED SIZE_RANG SIZE_FIXED_AND_EXT SIZE_RANG_AND_EXT REAL_WITH_COMPONENTS
	//SIZE_FIXED SIZE_FIXED_AND_EXT 定长时包含 Value
	Value Int64 `json:"value"`
	//SIZE_RANG SIZE_RANG_AND_EXT 变长时包含值范围
	Min Int64 `json:"min"`
	Max Int64 `json:"max"`
	//Type 为 REAL_WITH_COMPONENTS 时的约束
	Base     Int64 `json:"base"`     //10
	Exponent Int64 `json:"exponent"` //-2
}

type Marker struct {
	Flags String `json:"flags"`
}

type Tag struct {
	TagMode  String `json:"tag_mode"`
	TagClass String `json:"tag_class"`
	TagValue Int64  `json:"tag_value"`
}

type Import struct {
	ModuleName String `json:"moduleName"`
	Form       String `json:"form"`
}

type ImportList []Import
type ModuleList []Module
type MemberList []Member
type String string
type Int64 int64
