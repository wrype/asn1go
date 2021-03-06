# ASN1GO
> ASN.1 ==> Go Uper Codec

## Build & Install
```bash
# 编译
make
# 安装
make install
```

## Usage
```bash
# 从 asn 生成 go 代码
./asngo -f xxx.asn
# 生成一个临时的语法树 json 文件
./asngo -f xxx.asn -o tmp.json
# 生成一个临时的语法树 json 输出到终端
./asngo -f xxx.asn -o -
```

## UPER

### BIT STRING

>！约定不超过 64 个成员，注意代码生成的位序与编码相反

#### 编码流程

    1.有扩展标记时，第一个 bit 代表在约束范围内
    2.长度部分：定长时无需编码，变长时根据约束范围计算需要多少 bit 编码，可扩展且在范围外时，固定 8bit
    3.数据部分：具体的编码，注意位序（将原始数据按其二进制进行反转后（如原数据1010，反转后0101），按长度位（长度位表示bit数）编码数据）

### BOOLEAN

    无特殊情况，只有一个 bit，代表 true、false

### CHOICE

> 扩展的约束与 SEQUENCE 一致

#### 编码流程

    1.若有扩展标记，则编码一个 bit 代表是否选择了扩展成员
    2.若选择非扩展成员，进入步骤 3，否则进入步骤 5
    3.根据非扩展成员数量需要的最多 X bit，编码选择的成员索引（如：有三个非扩展成员，则需要两个 bit 表示位置，第二个非扩展成员不为空则编码为 01）
    4.接着的是该成员的 uper 编码，完成
    5.编码一个 bit 代表扩展成员是否超过 64 个
    6.编码 6bit 代表选择的扩展成员索引，从 0 开始，（如：选择了第一个扩展成员，则编码为 000000 ）
    7.编码 8bit 代表扩展成员 uper 编码的字节长度X
    8.扩展成员的 uper 编码，对齐到 X Byte，完成
    
### ENUMERATED

#### 编码流程

    1.若有扩展标记，则编码一个 bit 代表是否是扩展
    2.根据传进来的值找到对应的索引Index
    3.若在非扩展成员中找到的，则根据成员数量需要的 bit 数编码 Index
    4.若在扩展成员中找到的，则固定编码 7Bit Index
    
### IA5String

#### 编码流程

    1.若有扩展标记，则编码一个 bit 代表是否是扩展
    2.长度部分：若定长则不编码，若变长则根据范围需要的 bit 数编码长度，若是扩展且范围不在约束范围内的，则固定编码 8bit 长度位
    3.之后每7个比特代表一个ascii字符（去掉了前面一位0）
    
### INTEGER

>约定：必须是变长的，定长无太大意义

#### 编码流程

    1.若有扩展标记，则编码一个 bit 代表是否是扩展
    2.长度部分：若是扩展且范围不在约束范围内的，则固定编码 8bit 长度位，否则不编码
    3.数据部分：按照长度位编码数据
    
### OCTET STRING

#### 编码流程

    1.若有扩展标记，则编码一个 bit 代表是否是扩展
    2.长度部分：若是扩展且范围不在约束范围内的，则固定编码8bit长度位，否则不编码
    3.数据部分：按长度位（长度位表示数据字节数）编码数据

### REAL

#### 编码流程

    1.把数值转成科学计数法字符串形式，如 1.23 = '123.E-2'
    2.第一个 Byte 代表接下来的长度，第二个 Byte 代表解码方式，固定 0x03
    2.然后按照 ASCII 字符编码，每个字符 8bit

### SEQUENCE

#### 编码流程

    1.有扩展标记的情况，第一个 bit 代表是否有非空的扩展成员，无扩展标记时无需编码该 bit
    2.非扩展成员的 OPTIONAL 标记的掩码，代表该成员是否存在 （如有共有可选非扩展成员，且都不为空，则掩码为 111）
    3.非扩展成员的实际编码，按顺序拼接
    4.若无可扩展成员，完成编码。若有扩展成员且有一个成员不为空，则继续编码。
    5.一个 bit 代表扩展成员数量是否超过 64 个。（目前不处理超过 64 个扩展成员的情况）
    6.6bit 代表实际的扩展成员数量，从 0 开始（编码 000000 代表有一个扩展成员）。
    7.扩展成员的掩码（如有三个扩展成员，且都不为空，则掩码为 111）
    8.扩展成员的编码：首先对单个成员进行 uper 编码，对该 uper 编码补位成 X Byte，最后写入长度位 8 bit + X Byte
    
### 扩展的情况

INTEGER 中最多出现一个扩展标记(...)。且扩展标记后的数据不影响编码。
以下三种声明是等价的。约定只能使用第一种声明方式。即第一个为约束范围，第二个可能接着扩展标记。

```asn.1
VAL1 ::= INTEGER(0..7,...)
VAL2 ::= INTEGER(0..7,...,89)
VAL3 ::= INTEGER(0..7,...,8..9)

-- VAL1 == VAL2 == VAL3

-- 以下声明是不合法的
VAL4 ::= INTEGER(0..7,8..9,...)
VAL5 ::= INTEGER(0..7,...,8,...)
```

SEQUENCE CHOICE 等成员中最多出现两个扩展标记(...)。以下三种定义是等价的。
目前支持使用第二、三种方式定义扩展，即扩展标记后面的成员都是扩展成员。

```asn.1
MsgI ::= SEQUENCE {
  messageId INTEGER (0..32767) ,
  value OCTET STRING (SIZE(1..16)),
  ...,
  ss INTEGER(0..7),
  gg INTEGER(0..7),
  ...,
  aa INTEGER(0..7)
}

MsgII ::= SEQUENCE {
  messageId INTEGER (0..32767) ,
  value OCTET STRING (SIZE(1..16)),
  aa INTEGER(0..7),
  ...,
  ss INTEGER(0..7),
  gg INTEGER(0..7),
  ...
}

MsgIII ::= SEQUENCE {
  messageId INTEGER (0..32767) ,
  value OCTET STRING (SIZE(1..16)),
  aa INTEGER(0..7),
  ...,
  ss INTEGER(0..7),
  gg INTEGER(0..7)
}

-- MsgIII == MsgII == MsgI
```

    对于扩展的编码通常可能会有以下部分
    扩展部分：1bit 代表是否选择了扩展
    长度部分：8bit 代表扩展数据编码的字节数（X Byte）
    数据部分：X Byte == 8*X bit