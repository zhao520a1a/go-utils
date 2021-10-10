# 简介

## 引言

> 对于很多非开发人员来说， 写个excel要比写json或者yaml什么简单得多。这种场景下，读取特定格式（符合关系数据库特点的表格）的数据会比各种花式写入Excel的功能更为重要，于是写了这个简化库。

## 目标

期望创建一个阅读器库以轻松阅读类似excel，就像读取 DB 中数据一样。

## 示例

> 所有的示例可以移步查看 [excel/excel_test.go](excel/excel_test.go).

假设你有一个如下的 excel 文件（第一行是表头，其他行是内容）：

|ID|NameOf|AgeOf|Slice|UnmarshalString|
|-|-|-|-|-|
|1|Andy|1|1\|2|{"Foo":"Andy"}|
|2|Leo|2|2\|3\|4|{"Foo":"Leo"}|
|3|Ben|3|3\|4\|5\|6|{"Foo":"Ben"}|
|4|Ming|4|1|{"Foo":"Ming"}|

``` go

// defined a struct
type Standard struct {
	// 使用字段名称作为默认列名称
	ID      int
	// column表示映射列名
	Name    string `xlsx:"column(NameOf)"`
	// 可以将一列映射到多个字段中
	NamePtr *string `xlsx:"column(NameOf)"`
	// 如果只想映射到列名，则忽略“ column”，等价于“ column（AgeOf）”
	Age     int     `xlsx:"AgeOf"`
	// split表示通过`|`将字符串分割成切片
	Slice   []int `xlsx:"split(|)"`
	// *Temp 实现了 `encoding.BinaryUnmarshaler`
	Temp    *Temp `xlsx:"column(UnmarshalString)"`
	// 使用“-”去忽略映射信息
	Ignored string `xlsx:"-"`
}

type Temp struct {
	Foo string
}

func (this *Temp) UnmarshalBinary(d []byte) error {
	return json.Unmarshal(d, this)
}

func simpleUsage() {
	var stdList []Standard
	err := excel.UnmarshalXLSX("./testdata/simple.xlsx", &stdList)
	if err != nil {
		panic(err)
	}
}
```

Tips:

+ 空行将被跳过。
+ 大于len(TitleRow)的列将被跳过。
+ 只有空单元格可以填充默认值，如果一个单元格不能解析成一个字段，将返回一个错误。
+ 默认值也可以通过`encoding.BinaryUnmarshaler`来解读。
+ 如果没有标题行私有化，默认的列名如`'A', 'B', 'C', 'D' ......, 'XFC', 'XFD'`可以作为26数字系统的列名使用。
+ 当标题行有重复的标题，将返回错误`ErrDuplicatedTitles'。
+ 针对excel版本，支持.xlsx 不支持.xls。

## 进阶用法

### 自定义配置

Using a config as "excel.Config":

``` go
type Config struct {
	//如果sheet是字符串，将使用sheet作为工作表名称。
	// 如果sheet是int，将使用工作簿中的第i个sheet，注意隐藏的sheet会被计算在内。
	// 如果工作表是一个实现`GetXLSXSheetName()string'的对象，将使用其返回值。
	// 否则，将使用sheet作为结构并反映它的名称。
	// 如果sheet是一个切片，将像以前一样使用元素的类型来推断。
	Sheet interface{}
	// 指定作为标题的索引行，标题行之前的每一行都将被忽略，默认为0。
	TitleRowIndex int
	// 跳过标题后的n行，默认为0（不跳过），空行不计算在内。
	Skip int
	// 自动为sheet添加前缀。
	Prefix string
	// 自动为sheet添加后缀。
	Suffix string
}

```

## XLSX 标签使用

### column

映射到标题行的字段名，默认将使用字段名。

### default

当excel单元格中没有填入数值时，设置默认值，默认为0或""。

### split

分割一个字符串，并将其转换为一个切片。

### nil

当数据等于'nil(xxx)'中xxx时，将跳过扫描单元格中的值，不会赋值到struct字段。

### req

如果excel中不存在clomun标题，将返回错误。

## XLSX Field Config | 字段的解析配置

有时处理转义字符有点麻烦，所以实现`GetXLSXFieldConfigs() map[string]FieldConfig`的接口将比`tag`
更优先提供字段配置，更多信息请看测试文件[excel/field_config_test.go](excel/field_config_test.go)。

## 参考资料

这个库主要参考了[szyhf/go-excel](https://github.com/szyhf/go-excel) 的部分实现和读取逻辑。