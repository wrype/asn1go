%{
package asnparser

import (
	"fmt"
	"encoding/json"
)

%}

//SymType
%union {
	str string //匹配到的字符
	//临时变量
	String String
	Int64 Int64
	ModuleList ModuleList
	Module Module
	Member Member
	MemberList MemberList
	Constraint Constraint
	Marker Marker
	ImportList ImportList
	Import Import
	Tag Tag
}

//ASN.1关键字
%token ABSENT
%token ABSTRACT_SYNTAX
%token ALL
%token ANY
%token APPLICATION
%token AUTOMATIC
%token BEGIN
%token BIT
%token BMPString
%token BOOLEAN
%token BY
%token CHARACTER
%token CHOICE
%token CLASS
%token COMPONENT
%token COMPONENTS
%token CONSTRAINED
%token CONTAINING
%token DEFAULT
%token DEFINITIONS
%token DEFINED
%token EMBEDDED
%token ENCODED
%token ENCODING_CONTROL
%token END
%token ENUMERATED
%token EXPLICIT
%token EXPORTS
%token EXTENSIBILITY
%token EXTERNAL
%token FALSE
%token FROM
%token GeneralizedTime
%token GeneralString
%token GraphicString
%token IA5String
%token IDENTIFIER
%token IMPLICIT
%token IMPLIED
%token IMPORTS
%token INCLUDES
%token INSTANCE
%token INSTRUCTIONS
%token INTEGER
%token ISO646String
%token MAX
%token MIN
%token MINUS_INFINITY
%token NULL
%token NumericString
%token OBJECT
%token ObjectDescriptor
%token OCTET
%token OF
%token OPTIONAL
%token PATTERN
%token PDV
%token PLUS_INFINITY
%token PRESENT
%token PrintableString
%token PRIVATE
%token REAL
%token RELATIVE_OID
%token SEQUENCE
%token SET
%token SIZE
%token STRING
%token SYNTAX
%token T61String
%token TAGS
%token TeletexString
%token TRUE
%token TYPE_IDENTIFIER
%token UNIQUE
%token UNIVERSAL
%token UniversalString
%token UTCTime
%token UTF8String
%token VideotexString
%token VisibleString
%token WITH
%token UTF8_BOM

%token ASSIGNMENT //ASSIGNMENT
%token '{' '}'
%token '[' ']'
%token '(' ')'
%token ',' ';'
%token ELLIPSISS // ELLIPSISS
%token ELLIPSIS  //ELLIPSIS

%token Number     //[-]?[0-9]+
%token Identifier //[a-zA-Z0-9][-a-zA-Z0-9]*

//指定类型
%type <String> Identifier
%type <Int64> Number
%type <Module> ModuleDefinition
%type <ModuleList> ModuleList
%type <String> optModuleDefinitionFlags
%type <String> ModuleDefinitionFlags
%type <String> ModuleDefinitionFlag
%type <MemberList> optModuleBody
%type <MemberList> ModuleBody
%type <MemberList> AssignmentList
%type <ImportList> optImports
%type <ImportList> Imports
%type <ImportList> ImportModuleList
%type <Import> ImportModule
%type <Member> Assignment
%type <Member> ValueAssignment
%type <Member> DataTypeReference
%type <Member> Type
%type <Member> TaggedType
%type <Member> UntaggedType
%type <Member> TypeDeclaration
%type <Member> ConcreteTypeDeclaration
%type <Marker> optMarker
%type <Marker> Marker
%type <Tag> optTag
%type <Tag> Tag
%type <Tag> TagTypeValue
%type <String> TagClass
%type <String> TagPlicit
%type <String> TypeRefName
%type <String> optIdentifier
%type <Member> BuiltinType
%type <MemberList> NamedNumberList
%type <Member> NamedNumber
%type <Int64> SignedNumber
%type <MemberList> Enumerations
%type <MemberList> UniverationList
%type <Member> UniverationElement
%type <MemberList> NamedBitList
%type <Member> NamedBit
%type <MemberList> AlternativeTypeLists
%type <Member> AlternativeType
%type <Member> ASN_TYPE
%type <Member> ALL_TYPE
%type <MemberList> optComponentTypeLists
%type <MemberList> ComponentTypeLists
%type <Member> ComponentType
%type <Constraint> optSizeOrConstraint
%type <Constraint> Constraint
%type <Constraint> ConstraintSpec
%type <MemberList> MemberConstraintList
%type <Member> MemberConstraint
%type <String> BasicTypeId
%type <String> BasicTypeId_UniverationCompatible
%type <String> BasicString
%type <Member> ExtensionAndException
%type <String> TagClass
%%
ParsedGrammar:
	UTF8_BOM ModuleList {
		ParsedGrammarlex.(*AsnLexer).result = &$2
	}
	| ModuleList {
		ParsedGrammarlex.(*AsnLexer).result = &$1
	}
;

ModuleList:
	ModuleDefinition {
		$$ = append(make(ModuleList,0),$1)
	}
	| ModuleList ModuleDefinition {
		$$ = append($1,$2)
	}
;

/*
 * ASN module definition.
 * === EXAMPLE ===
 * MySyntax DEFINITIONS AUTOMATIC TAGS ::=
 * BEGIN
 * ...
 * END
 * === EOF ===
 */
ModuleDefinition:
	Identifier DEFINITIONS optModuleDefinitionFlags ASSIGNMENT
	BEGIN
	optModuleBody
	END {
		$$ = Module{ ModuleName: $1 , Flags: $3, Members: $6}
	}
;

optModuleDefinitionFlags:
	{ $$ = "MSF_NOFLAGS" }
	| ModuleDefinitionFlags {
		$$ = $1
	}
;

/*
 * Module flags.
 */
ModuleDefinitionFlags:
	ModuleDefinitionFlag {
		$$ = $1
	}
	| ModuleDefinitionFlags ModuleDefinitionFlag {
		$$ = $1 + "," + $2
	}
;

/*
 * Single module flag.
 */
ModuleDefinitionFlag:
	EXPLICIT TAGS {
		$$ = "MSF_EXPLICIT_TAGS"
	}
	| IMPLICIT TAGS {
		$$ = "MSF_IMPLICIT_TAGS"
	}
	| AUTOMATIC TAGS {
		$$ = "MSF_AUTOMATIC_TAGS"
	}
	| EXTENSIBILITY IMPLIED {
		$$ = "MSF_EXTENSIBILITY_IMPLIED"
	}
	/* EncodingReferenceDefault */
	| Identifier INSTRUCTIONS {
		/* X.680Amd1 specifies TAG and XER */
		if($1 == "TAG") {
		 	$$ = "MSF_TAG_INSTRUCTIONS"
		} else if($1 == "XER") {
		 	$$ = "MSF_XER_INSTRUCTIONS"
		} else {
			//console.error(`WARNING: $1 INSTRUCTIONS Unrecognized encoding reference`);
		 	$$ = "MSF_unk_INSTRUCTIONS"
		}
	}
;

optModuleBody: { $$= []Member{} }| ModuleBody ;

/*
 * ASN.1 Module body.
 */
ModuleBody:
	optImports AssignmentList {
		if(len($1)!=0){
			fmt.Println("不支持 Import !")
		}
		$$ = $2
	}
	| AssignmentList { $$ = $1 }
;

AssignmentList:
	Assignment {
		$$ = append(make([]Member,0),$1);
	}
	| AssignmentList Assignment {
		$$ = append($1,$2)
	}
;

// TODO 不支持模块导入导出
optImports: {$$=[]Import{}} | Imports  { $$ = $1 };

Imports:
	IMPORTS ImportModuleList ';' {
		$$ = $2
	}
;

ImportModuleList:
	ImportModule {
		$$ = append(make([]Import,0),$1)
	}
	| ImportModuleList ImportModule {
		$$ = append($1,$2)
	}
	| {
		$$ = []Import{}
	}
;

ImportModule:
	Identifier FROM Identifier {
		$$ = Import{ ModuleName: $1, Form: $3 }
	}
;

/*
 * One of the elements of ASN.1 module specification.
 */
Assignment:
	DataTypeReference {
		$$ = $1
	}
	| ValueAssignment { $$ = $1 }
	// /*
	//  * Value set definition
	//  * === EXAMPLE ===
	//  * EvenNumbers INTEGER ::= { 2 | 4 | 6 | 8 }
	//  * === EOF ===
	//  * Also ObjectClassSet.
	//  */
	// | ValueSetTypeAssignment {
	// 	$$ = asn1p_module_new();
	// 	checkmem($$);
	// 	assert($1.expr_type != A1TC_INVALID);
	// 	assert($1.meta_type != AMT_INVALID);
	// 	asn1p_module_member_add($$, $1);
	// }
	// | ENCODING_CONTROL capitalreference
	// 	{ asn1p_lexer_hack_push_encoding_control(); }
	// 		{
	// 	fprintf(stderr,
	// 		"WARNING: ENCODING-CONTROL %s "
	// 		"specification at %s:%d ignored\n",
	// 		$2, ASN_FILENAME, yylineno);
	// 	free($2);
	// 	$$ = 0;
	// }

	/*
	 * Erroneous attemps
	 */
	| BasicString {
		panic("Attempt to redefine a standard basic string type, please comment out or remove this type redefinition.")
	}
;

ValueAssignment:
	Identifier ALL_TYPE ASSIGNMENT Number {
		$$ = Member{
			Identifier:$1, 
			MetaType:"AMT_VALUE_ASSIGMENT", 
			ExprType:$2.ExprType,
			Constraints: &Constraint{Type:"AMT_VALUE", Value:$4},
			Members: append(make([]Member,0),$2),
		}
	}
;

/*
 * Data Type Reference.
 * === EXAMPLE ===
 * Type3 ::= CHOICE { a Type1,  b Type 2 }
 * === EOF ===
 */
DataTypeReference:
	/*
	 * Optionally tagged type definition.
	 */
	TypeRefName ASSIGNMENT Type {
		$$ = $3;
		$$.Identifier = $1;
	}
	// | TypeRefName ASSIGNMENT ObjectClass {
	// 	$$ = $3;
	// 	$$.Identifier = $1;
	// 	assert($$.expr_type == A1TC_CLASSDEF);
	// 	assert($$.meta_type == AMT_OBJECTCLASS);
	// }
	// /*
	//  * Parameterized <Type> declaration:
	//  * === EXAMPLE ===
	//  *   SIGNED { ToBeSigned } ::= SEQUENCE {
	//  *      toBeSigned  ToBeSigned,
	//  *      algorithm   AlgorithmIdentifier,
	//  *      signature   BIT STRING
	//  *   }
	//  * === EOF ===
	//  */
	// | TypeRefName '{' ParameterArgumentList '}' ASSIGNMENT Type {
	// 	$$ = $6;
	// 	$$.Identifier = $1;
	// 	$$.lhs_params = $3;
	// }
	// /* Parameterized CLASS declaration */
	// | TypeRefName '{' ParameterArgumentList '}' ASSIGNMENT ObjectClass {
	// 	$$ = $6;
	// 	$$.Identifier = $1;
	// 	$$.lhs_params = $3;
	// }
;

Type: TaggedType { $$ = $1 };

TaggedType:
    optTag UntaggedType {
		$$ = $2
		//$$.Tag = CopyTag($1)
    }
;

//TODO optManyConstraints
UntaggedType:
	TypeDeclaration optSizeOrConstraint {
		$$ = $1
		if $$.Constraints == nil || $$.Constraints.Type == "" {
			$$.Constraints = CopyConstraint($2)
		}
	}
;

TypeDeclaration:
	ConcreteTypeDeclaration { $$ = $1 }
    //| DefinedType //类型引用等
;

/*
 * 具体类型声明
 * TODO SET
 */
ConcreteTypeDeclaration:
	BuiltinType { $$ = $1 }
	| CHOICE '{' AlternativeTypeLists '}' {
		$$ = Member{ Members:$3, ExprType: "CHOICE", MetaType: "AMT_TYPE" }
	}
	| SEQUENCE '{' optComponentTypeLists '}' {
		$$ = Member{ Members:$3, ExprType: "SEQUENCE", MetaType: "AMT_TYPE" }
	}
	// | SET '{' optComponentTypeLists '}' {
	// 	$$ = $3;
	// 	$$.expr_type = 'SET';
	// 	$$.meta_type = 'AMT_TYPE';
	// }
	| SEQUENCE optSizeOrConstraint OF optTag ALL_TYPE {
		$$ = Member{
			Identifier:"",
			Members:append(make([]Member,0),$5),
			ExprType: "SEQUENCE_OF",
			MetaType: "AMT_TYPE",
		}
		//$$.Tag = CopyTag($4)
		$$.Constraints = CopyConstraint($2)
	}
	// | SET optSizeOrConstraint OF optIdentifier optTag MaybeIndirectTypeDeclaration {
	// 	$$ = {}
	// 	$$.constraints = $2;
	// 	$$.expr_type = 'SET_OF';
	// 	$$.meta_type = 'AMT_TYPE';
	// 	$6.Identifier = $4;
	// 	$6.tag = $5;
	// 	asn1p_expr_add($$, $6);
	// }
	// | ANY {
	// 	$$ = {}
	// 	$$.expr_type = 'ASN_TYPE_ANY';
	// 	$$.meta_type = 'AMT_TYPE';
	// }
	// | INSTANCE OF ComplexTypeReference {
	// 	$$ = {}
	// 	$$.reference = $3;
	// 	$$.expr_type = 'A1TC_INSTANCE';
	// 	$$.meta_type = 'AMT_TYPE';
	// }
;



optMarker: { $$= Marker{ Flags: "EM_NOMARK" } }
	| Marker { $$ = $1 }
;

Marker:
	OPTIONAL {
		$$= Marker{Flags:"EM_OPTIONAL"}
	}
;

/* -------------------------------------------------------------------
 * SET definition.
 * === EXAMPLE ===
 * Person ::= SET {
 * 	name [0] PrintableString (SIZE(1..20)),
 * 	country [1] PrintableString (SIZE(1..20)) DEFAULT default-country,
 * }
 * EMBEDDEDPDV ::= [UNIVERSAL 11] IMPLICIT SEQUENCE {...}
 * === EOF ===
 */

optTag:
	{ $$ = Tag{ TagMode:"TM_DEFAULT", TagClass:"TC_CONTEXT_SPECIFIC", TagValue:0 } }
	| Tag { $$ = $1 }
;

Tag:
	TagTypeValue TagPlicit {
		$$ = $1
		$$.TagMode = $2
	}
;

TagTypeValue:
	'[' TagClass Number ']' {
		$$ = Tag{ TagMode:"TM_DEFAULT", TagClass:$2, TagValue:$3 } 
	}
;

TagClass:
	{ $$ = "TC_CONTEXT_SPECIFIC" }
	| UNIVERSAL { $$ = "TC_UNIVERSAL" }
	| APPLICATION { $$ = "TC_APPLICATION" }
	| PRIVATE { $$ = "TC_PRIVATE" }
;

TagPlicit:
	{ $$ = "TM_DEFAULT" }
	| IMPLICIT { $$ = "TM_IMPLICIT" }
	| EXPLICIT { $$ = "TM_EXPLICIT" }
;

TypeRefName:
	Identifier { $$ = $1 }
;

optIdentifier:
	{ $$ = "" }
	| Identifier {
		$$ = $1;
	}
;

/* ----------------------------------------------------
 * 一些基础类型
 */

BuiltinType:
	BasicTypeId {
		$$ = Member{ExprType:$1,MetaType:"AMT_TYPE"}
	}
	| INTEGER '{' NamedNumberList '}' {
		$$ = Member{ExprType:"INTEGER",MetaType:"AMT_TYPE",Members:$3}
	}
	| ENUMERATED '{' Enumerations '}' {
		$$ = Member{ExprType:"ENUMERATED",MetaType:"AMT_TYPE",Members:$3}
	}
	| BIT STRING '{' NamedBitList '}' {
	    	$$ = Member{ExprType:"BIT_STRING",MetaType:"AMT_TYPE",Members:$4}
	}
;

/*
 * INTEGER { a(0),b(1) }
 */
NamedNumberList:
	NamedNumber {
		$$ = append(make([]Member,0),$1)
	}
	| NamedNumberList ',' NamedNumber {
		$$ = append($1,$3)
	}
;

NamedNumber:
	Identifier '(' SignedNumber ')' {
		$$ = Member{Identifier:$1,ExprType:"A1TC_UNIVERVAL",MetaType:"AMT_VALUE",Value:$3}
	}
;

SignedNumber:
	Number { $$ = $1 }
;

// 枚举类型 a(0),b(1)
Enumerations:
	UniverationList {
		for i,v:=range $1{
			if v.Value==-999 {
				$1[i].Value = Int64(i)
			}
		}
		$$ = $1
	}
;

UniverationList:
	UniverationElement {
		$$ = append(make([]Member,0),$1)
	}
	| UniverationList ',' UniverationElement {
		$$ = append($1,$3)
	}
;

UniverationElement:
	Identifier {
		$$ = Member{Identifier:$1,ExprType:"A1TC_UNIVERVAL",MetaType:"AMT_VALUE",Value:-999}
	}
	| Identifier '(' SignedNumber ')' {
		$$ = Member{Identifier:$1,ExprType:"A1TC_UNIVERVAL",MetaType:"AMT_VALUE",Value:$3}
	}
	| SignedNumber {
		$$ = Member{Identifier:"",ExprType:"A1TC_UNIVERVAL",MetaType:"AMT_VALUE",Value:$1}
	}
	| ExtensionAndException {$$ =$1}
;

//BIT STRING 内容 	entrance(0),exit(1),
NamedBitList:
	NamedBit {
		$$ = append(make([]Member,0),$1)
	}
	| NamedBitList ',' NamedBit {
		$$ = append($1,$3)
	}
;

NamedBit:
	Identifier '(' Number ')' {
		$$ = Member{Identifier:$1,ExprType:"A1TC_UNIVERVAL",MetaType:"AMT_VALUE",Value:$3}
	}
;


// CHOICE 内容

AlternativeTypeLists:
	AlternativeType {
		$$ = append(make([]Member,0),$1)
	}
	| AlternativeTypeLists ',' AlternativeType {
		$$ = append($1,$3)
	}
;

AlternativeType:
	Identifier ALL_TYPE {
		$$ = $2;
		$$.Identifier = $1;
	}
	| ExtensionAndException {
		$$ = $1;
	}
;

ASN_TYPE:
	//不检查自定义类型
	Identifier {
		$$ = Member{Identifier:"",ExprType:$1,MetaType:"ASN_TYPE"}
	}
;

ALL_TYPE:
	Type { $$ = $1 } | ASN_TYPE { $$ = $1 }
;

// SEQUENCE 内容
optComponentTypeLists:
	{ $$ = []Member{} }
	| ComponentTypeLists { $$ = $1 };

ComponentTypeLists:
	ComponentType {
		$$ = append(make([]Member,0),$1)
	}
	| ComponentTypeLists ',' ComponentType {
		$$ = append($1,$3)
	}
	;

ComponentType:
	Identifier Type optMarker {
		$$ = $2
		$$.Identifier = $1
		$$.Marker = CopyMarker($3)
	}
	| Identifier ASN_TYPE optMarker {
		$$ = $2
		$$.Identifier = $1
		$$.Marker = CopyMarker($3)
	}
	| ExtensionAndException {
		$$ = $1
	}
;

optSizeOrConstraint:
	{ $$ = Constraint{} }
	| Constraint { $$ = $1 }
;

/*
	DataRateKByte ::= REAL (WITH COMPONENTS{
			mantissa (0..65536),
			base(10),
			exponent(-3)
	})
*/
// (SIZE(1..63)) (SIZE(63)) (-31..31) SIZE(8,...) (0..127,...)
Constraint:
    	'(' ConstraintSpec ')' {
		$$ = $2
    	}
	| '(' WITH COMPONENTS '{' MemberConstraintList '}' ')' {
		$$ = Constraint{}
		$$.Type = "REAL_WITH_COMPONENTS"
		for i :=0;i<len($5);i++ {
			constraint := $5[i]
			if(constraint.Identifier=="mantissa"){
				$$.Min=constraint.Constraints.Min
				$$.Max=constraint.Constraints.Max
			}else if(constraint.Identifier == "base"){
				$$.Base=constraint.Constraints.Value
			}else if(constraint.Identifier == "exponent"){
				$$.Exponent=constraint.Constraints.Value
			}
		}
	}
;

ConstraintSpec:
	SIZE  '(' Number ')' {
		$$ = Constraint{
			Type: "SIZE_FIXED",
			Value: $3,
		}
	}
	| SIZE '(' Number ELLIPSIS Number ')' {
		$$ = Constraint{
			Type: "SIZE_RANG",
			Min:$3,
			Max:$5,
		}
	}
	| SIZE '(' Number ',' ELLIPSISS ')' {
		$$ = Constraint{
			Type: "SIZE_FIXED_AND_EXT",
			Value: $3,
		}
	}
	| SIZE '(' Number ELLIPSIS Number ',' ELLIPSISS ')' {
		$$ = Constraint{
			Type: "SIZE_RANG_AND_EXT",
                        Min:$3,
                        Max:$5,
		}
	}
	| Number ELLIPSIS Number ',' ELLIPSISS{
		$$ = Constraint{
			Type: "SIZE_RANG_AND_EXT",
			Min:$1,
			Max:$3,
		}
	}
	| Number ELLIPSIS Number {
		$$ = Constraint{
			Type: "SIZE_RANG",
			Min:$1,
			Max:$3,
		}
	}
;

MemberConstraintList:
	MemberConstraint {
		$$ = append(make(MemberList,0),$1)
	}
	| MemberConstraintList ',' MemberConstraint {
		$$ = append($1,$3)
	}
;

MemberConstraint:
	Identifier '(' ConstraintSpec ')' {
		$$ = Member{ }
		$$.Identifier=$1
		$$.Constraints=CopyConstraint($3)
	}
	| Identifier '(' Number ')' {
		$$ = Member{ }
		$$.Identifier=$1
		$$.Constraints=&Constraint{
			Type: "SIZE_FIXED",
			Value: $3,
		}
	}
;

//可扩展
ExtensionAndException:
	ELLIPSISS {
		$$ = Member{ Identifier:"...", ExprType: "A1TC_EXTENSIBLE", MetaType:"AMT_TYPE" }
	}
;

BasicTypeId:
	BOOLEAN { $$ = "BOOLEAN" }
	| NULL { $$ = "NULL" }
	| REAL { $$ = "REAL" }
	| OCTET STRING { $$ = "OCTET_STRING" }
	| OBJECT IDENTIFIER { $$ = "OBJECT_IDENTIFIER" }
	| RELATIVE_OID { $$ = "RELATIVE_OID" }
	| EXTERNAL  { $$ = "EXTERNAL" }
	| EMBEDDED PDV { $$ = "EMBEDDED_PDV" }
	| CHARACTER STRING { $$ = "CHARACTER_STRING" }
	| UTCTime  { $$ = "UTCTime" }
	| GeneralizedTime  { $$ = "GeneralizedTime" }
	| BasicString { $$ = $1 }
	| BasicTypeId_UniverationCompatible { $$ = $1 }
;

BasicString:
    	BMPString { $$ = "BMPString" }
	| GeneralString { $$ = "GeneralString" }
	| GraphicString { $$ = "GraphicString" }
	| IA5String  { $$ = "IA5String" }
	| ISO646String  { $$ = "ISO646String" }
	| NumericString { $$ = "NumericString" }
	| PrintableString  { $$ = "PrintableString" }
	| T61String  { $$ = "T61String" }
	| TeletexString { $$ = "TeletexString" }
	| UniversalString { $$ = "UniversalString" }
	| UTF8String { $$ = "UTF8String" }
	| VideotexString { $$ = "VideotexString" }
	| VisibleString { $$ = "VisibleString" }
	| ObjectDescriptor { $$ = "ObjectDescriptor" }
;
;

/*
 * A type identifier which may be used with "{ a(1), b(2) }" clause.
 */
BasicTypeId_UniverationCompatible:
	INTEGER { $$ = "INTEGER" }
	| ENUMERATED { $$ = "ENUMERATED" }
	| BIT STRING { $$ = "BIT_STRING" }
;



%%

//指针类型都要复制一份

//复制一个Constraint
func CopyConstraint(c Constraint) *Constraint {
	if c.Type == "" {
		return nil
	}
	var copy Constraint
	buf, _ := json.Marshal(c)
	json.Unmarshal(buf, &copy)
	return &copy
}

//复制一个Marker
func CopyMarker(c Marker) *Marker {
	if c.Flags == "" {
		return nil
	}
	var copy Marker
	buf, _ := json.Marshal(c)
	json.Unmarshal(buf, &copy)
	return &copy
}

//复制一个 Tag
func CopyTag(c Tag) *Tag {
	if c.TagMode == "" {
		return nil
	}
	var copy Tag
	buf, _ := json.Marshal(c)
	json.Unmarshal(buf, &copy)
	return &copy
}