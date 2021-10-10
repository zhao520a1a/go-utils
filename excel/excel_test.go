package excel

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"reflect"
	"testing"
)

var (
	TestFileUri    = "10/sea/37/f4/5914da82c3f80263989b93ae6ed2"
	TestFilePath   = "./testdata/simple.xlsx"
	StdSheetName   = "Standard"
	AdvSheetName   = "Advance"
	AdvSheetSuffix = ".suffix"
	DupSheetName   = "DuplicatedTitle"
)

var expectStandardList = []Standard{
	{
		ID:      1,
		Name:    "Andy",
		NamePtr: StrPtr("Andy"),
		Age:     1,
		Slice:   []int{1, 2},
		Temp: &Temp{
			Foo: "Andy",
		},
	},
	{
		ID:      2,
		Name:    "Leo",
		NamePtr: StrPtr("Leo"),
		Age:     2,
		Slice:   []int{2, 3, 4},
		Temp: &Temp{
			Foo: "Leo",
		},
	},
	{
		ID:      3,
		Name:    "Ben",
		NamePtr: StrPtr("Ben"),
		Age:     3,
		Slice:   []int{3, 4, 5, 6},
		Temp: &Temp{
			Foo: "Ben",
		},
	},
	{
		ID:      4,
		Name:    "Ming",
		NamePtr: StrPtr("Ming"),
		Age:     4,
		Slice:   []int{1},
		Temp: &Temp{
			Foo: "Ming",
		},
	},
}

var expectStandardPtrList = []*Standard{
	{
		ID:      1,
		Name:    "Andy",
		NamePtr: StrPtr("Andy"),
		Age:     1,
		Slice:   []int{1, 2},
		Temp: &Temp{
			Foo: "Andy",
		},
	},
	{
		ID:      2,
		Name:    "Leo",
		NamePtr: StrPtr("Leo"),
		Age:     2,
		Slice:   []int{2, 3, 4},
		Temp: &Temp{
			Foo: "Leo",
		},
	},
	{
		ID:      3,
		Name:    "Ben",
		NamePtr: StrPtr("Ben"),
		Age:     3,
		Slice:   []int{3, 4, 5, 6},
		Temp: &Temp{
			Foo: "Ben",
		},
	},
	{
		ID:      4,
		Name:    "Ming",
		NamePtr: StrPtr("Ming"),
		Age:     4,
		Slice:   []int{1},
		Temp: &Temp{
			Foo: "Ming",
		},
	},
}

var expectStandardMapList = []map[string]string{
	{
		"ID":              "1",
		"NameOf":          "Andy",
		"AgeOf":           "1",
		"Slice":           "1|2",
		"UnmarshalString": "{\"Foo\":\"Andy\"}",
	},
	{
		"ID":              "2",
		"NameOf":          "Leo",
		"AgeOf":           "2",
		"Slice":           "2|3|4",
		"UnmarshalString": "{\"Foo\":\"Leo\"}",
	},
	{
		"ID":              "3",
		"NameOf":          "Ben",
		"AgeOf":           "3",
		"Slice":           "3|4|5|6",
		"UnmarshalString": "{\"Foo\":\"Ben\"}",
	},
	{
		"ID":              "4",
		"NameOf":          "Ming",
		"AgeOf":           "4",
		"Slice":           "1",
		"UnmarshalString": "{\"Foo\":\"Ming\"}",
	},
}

var expectStandardSliceList = [][]string{
	{
		"1",
		"Andy",
		"1",
		"1|2",
		"{\"Foo\":\"Andy\"}",
	},
	{
		"2",
		"Leo",
		"2",
		"2|3|4",
		"{\"Foo\":\"Leo\"}",
	},
	{
		"3",
		"Ben",
		"3",
		"3|4|5|6",
		"{\"Foo\":\"Ben\"}",
	},
	{
		"4",
		"Ming",
		"4",
		"1",
		"{\"Foo\":\"Ming\"}",
	},
}

// defined a struct
type Standard struct {
	// use field name as default column name
	ID int
	// column means to map the column name
	Name string `xlsx:"column(NameOf)"`
	// you can map a column into more than one field
	NamePtr *string `xlsx:"column(NameOf)"`
	// omit `column` if only want to map to column name, it's equal to `column(AgeOf)`
	Age int `xlsx:"AgeOf"`
	// split means to split the string into slice by the `|`
	Slice []int `xlsx:"split(|)"`
	Temp  *Temp `xlsx:"column(UnmarshalString)"`
	// use '-' to ignore.
	WantIgnored string `xlsx:"-"`
}

// func (this Standard) GetXLSXSheetName() string {
// 	return "Some sheet name if need"
// }

type Temp struct {
	Foo string
}

func (tmp *Temp) UnmarshalBinary(d []byte) error {
	return json.Unmarshal(d, tmp)
}

func init() {
	log.SetFlags(log.Llongfile)
}

func StrPtr(s string) *string {
	return &s
}

// UnmarshalXLSX unmarshal a sheet of XLSX file into a slice container.
// The sheet name will be inferred from element of container
// If container implement the function of GetXLSXSheetName()string, the return string will used.
// Oterwise will use the reflect struct name.
func UnmarshalXLSX(filePath string, container interface{}) error {
	conn := NewConnector()
	err := conn.Open(filePath)
	if err != nil {
		return err
	}

	rd, err := conn.NewReader(container)
	if err != nil {
		conn.Close()
		return err
	}

	err = rd.ReadAll(container)
	if err != nil {
		conn.Close()
		rd.Close()
		return err
	}
	conn.Close()
	rd.Close()
	return nil
}

func MustJsonPrettyString(i interface{}) string {
	if d, err := json.MarshalIndent(i, "", "\t"); err == nil {
		return string(d)
	} else {
		panic(err)
	}
}

func TestReadStandardSimple(t *testing.T) {
	var stdList []Standard
	err := UnmarshalXLSX(TestFilePath, &stdList)
	if err != nil {
		t.Error(err)
		return
	}
	if !reflect.DeepEqual(stdList, expectStandardList) {
		t.Errorf("unexprect std list: %s", MustJsonPrettyString(stdList))
	}
}

func TestReadStandard(t *testing.T) {
	conn := NewConnector()
	err := conn.Open(TestFilePath)
	if err != nil {
		t.Error(err)
		return
	}
	defer conn.Close()

	rd, err := conn.NewReader(StdSheetName)
	if err != nil {
		t.Error(err)
		return
	}
	defer rd.Close()

	idx := 0
	for rd.Next() {
		var s Standard
		if err := rd.Read(&s); err != nil {
			fmt.Println(err)
			return
		}
		expectStd := expectStandardList[idx]
		if !reflect.DeepEqual(s, expectStd) {
			t.Errorf("unexpect std at %d = \n%s", idx, MustJsonPrettyString(expectStd))
		}
		idx++
	}
}

func TestReadStandardIndex(t *testing.T) {
	conn := NewConnector()
	err := conn.Open(TestFilePath)
	if err != nil {
		t.Error(err)
		return
	}
	defer conn.Close()

	rd, err := conn.NewReader(2)
	if err != nil {
		t.Error(err)
		return
	}
	defer rd.Close()

	idx := 0
	for rd.Next() {
		var s Standard
		if err := rd.Read(&s); err != nil {
			fmt.Println(err)
			return
		}
		expectStd := expectStandardList[idx]
		if !reflect.DeepEqual(s, expectStd) {
			t.Errorf("unexpect std at %d = \n%s", idx, MustJsonPrettyString(expectStd))
		}
		idx++
	}
}

func TestReadStandardAll(t *testing.T) {
	conn := NewConnector()
	err := conn.Open(TestFilePath)
	if err != nil {
		t.Error(err)
		return
	}
	defer conn.Close()

	var std Standard
	var stdList []Standard

	rd, err := conn.NewReader(std)
	if err != nil {
		t.Error(err)
		return
	}
	defer rd.Close()

	err = rd.ReadAll(&stdList)
	if err != nil {
		t.Error(err)
		return
	}
	if !reflect.DeepEqual(expectStandardList, stdList) {
		t.Errorf("unexpect stdlist: \n%s", MustJsonPrettyString(stdList))
	}
}

func TestReadStandardPtrSimple(t *testing.T) {
	var stdList []*Standard
	err := UnmarshalXLSX(TestFilePath, &stdList)
	if err != nil {
		t.Error(err)
		return
	}
	if !reflect.DeepEqual(stdList, expectStandardPtrList) {
		t.Errorf("unexprect std list: %s", MustJsonPrettyString(stdList))
	}
}

func TestReadStandardPtrAll(t *testing.T) {
	conn := NewConnector()
	err := conn.Open(TestFilePath)
	if err != nil {
		t.Error(err)
		return
	}
	defer conn.Close()

	var stdList []*Standard
	rd, err := conn.NewReader(stdList)
	if err != nil {
		t.Error(err)
		return
	}
	defer rd.Close()

	err = rd.ReadAll(&stdList)
	if err != nil {
		t.Error(err)
		return
	}
	if !reflect.DeepEqual(expectStandardPtrList, stdList) {
		t.Errorf("unexpect stdlist: \n%s", MustJsonPrettyString(stdList))
	}
}

func TestReadStandardPtrAllFromUri(t *testing.T) {
	conn := NewConnector()
	err := conn.OpenFromUri(TestFileUri)
	if err != nil {
		t.Error(err)
		return
	}
	defer conn.Close()

	var stdList []*Standard
	rd, err := conn.NewReader(stdList)
	if err != nil {
		t.Error(err)
		return
	}
	defer rd.Close()

	err = rd.ReadAll(&stdList)
	if err != nil {
		t.Error(err)
		return
	}
	if !reflect.DeepEqual(expectStandardPtrList, stdList) {
		t.Errorf("unexpect stdlist: \n%s", MustJsonPrettyString(stdList))
	}
}

func TestReadBinaryStandardPtrAll(t *testing.T) {
	xlsxData, err := ioutil.ReadFile(TestFilePath)
	if err != nil {
		t.Error(err)
		return
	}

	conn := NewConnector()
	err = conn.OpenBinary(xlsxData)
	if err != nil {
		t.Error(err)
		return
	}
	defer conn.Close()

	var stdList []*Standard

	rd, err := conn.NewReader(stdList)
	if err != nil {
		t.Error(err)
		return
	}
	defer rd.Close()

	err = rd.ReadAll(&stdList)
	if err != nil {
		t.Error(err)
		return
	}
	if !reflect.DeepEqual(expectStandardPtrList, stdList) {
		t.Errorf("unexpect stdlist: \n%s", MustJsonPrettyString(stdList))
	}
}

func TestReadStandardMap(t *testing.T) {
	conn := NewConnector()
	err := conn.Open(TestFilePath)
	if err != nil {
		t.Error(err)
		return
	}
	defer conn.Close()

	rd, err := conn.NewReader(StdSheetName)
	if err != nil {
		t.Error(err)
		return
	}
	defer rd.Close()

	idx := 0
	for rd.Next() {
		var m map[string]string
		if err := rd.Read(&m); err != nil {
			fmt.Println(err)
			return
		}

		expectStdMap := expectStandardMapList[idx]
		if !reflect.DeepEqual(m, expectStdMap) {
			t.Errorf("unexpect std at %d = \n%s", idx, MustJsonPrettyString(expectStdMap))
		}
		idx++
	}
}

func TestReadStandardSliceMap(t *testing.T) {
	conn := NewConnector()
	err := conn.Open(TestFilePath)
	if err != nil {
		t.Error(err)
		return
	}
	defer conn.Close()

	rd, err := conn.NewReader(StdSheetName)
	if err != nil {
		t.Error(err)
		return
	}
	defer rd.Close()

	var stdMapList []map[string]string
	err = rd.ReadAll(&stdMapList)
	if err != nil {
		t.Error(err)
		return
	}

	if !reflect.DeepEqual(expectStandardMapList, stdMapList) {
		t.Errorf("unexpect stdlist: \n%s", MustJsonPrettyString(stdMapList))
	}
}

func TestReadStandardSlice(t *testing.T) {
	conn := NewConnector()
	err := conn.Open(TestFilePath)
	if err != nil {
		t.Error(err)
		return
	}
	defer conn.Close()

	rd, err := conn.NewReader(StdSheetName)
	if err != nil {
		t.Error(err)
		return
	}
	defer rd.Close()

	idx := 0
	for rd.Next() {
		var l []string
		if err := rd.Read(&l); err != nil {
			fmt.Println(err)
			return
		}

		expectStdList := expectStandardSliceList[idx]
		if !reflect.DeepEqual(l, expectStdList) {
			t.Errorf("unexpect std at %d %s = \n%s", idx, MustJsonPrettyString(l), MustJsonPrettyString(expectStdList))
		}
		idx++
	}
}

func TestReadStandardSliceList(t *testing.T) {
	conn := NewConnector()
	err := conn.Open(TestFilePath)
	if err != nil {
		t.Error(err)
		return
	}
	defer conn.Close()

	rd, err := conn.NewReader(StdSheetName)
	if err != nil {
		t.Error(err)
		return
	}
	defer rd.Close()

	var stdList [][]string
	err = rd.ReadAll(&stdList)
	if err != nil {
		t.Error(err)
		return
	}

	if !reflect.DeepEqual(expectStandardSliceList, stdList) {
		t.Errorf("unexpect stdlist: \n%s", MustJsonPrettyString(stdList))
	}
}

// -- 高级配置

var expectAdvanceList = []Advance{
	{
		ID:      1,
		Name:    "Andy",
		NamePtr: StrPtr("Andy"),
		Age:     1,
		Slice:   []int{1, 2},
		Temp: &Temp{
			Foo: "Andy",
		},
	},
	{
		ID:      2,
		Name:    "Leo",
		NamePtr: StrPtr("Leo"),
		Age:     2,
		Slice:   []int{2, 3, 4},
		Temp: &Temp{
			Foo: "Leo",
		},
	},
	{
		ID:      3,
		Name:    "",
		NamePtr: StrPtr("Ben"),
		Age:     180, //  using default
		Slice:   []int{3, 4, 5, 6},
		Temp: &Temp{
			Foo: "Ben",
		},
	},
	{
		ID:      4,
		Name:    "Ming",
		NamePtr: StrPtr("Ming"),
		Age:     4,
		Slice:   []int{1},
		Temp: &Temp{
			Foo: "Default",
		},
	},
}

type Advance struct {
	// use field name as default column name
	ID int
	// column means to map the column name, and skip cell that value equal to "Ben"
	Name string `xlsx:"column(NameOf);nil(Ben);req();"`
	// you can map a column into more than one field
	NamePtr *string `xlsx:"column(NameOf);req();"`
	// omit `column` if only want to map to column name, it's equal to `column(AgeOf)`
	// use 180 as default if cell is empty.
	Age int `xlsx:"column(AgeOf);default(180);req();"`
	// split means to split the string into slice by the `|`
	Slice []int `xlsx:"split(|);req();"`
	// use default also can marshal to struct
	Temp *Temp `xlsx:"column(UnmarshalString);default({\"Foo\":\"Default\"});req();"`
	// use '-' to ignore.
	WantIgnored string `xlsx:"-"`
	// By default, required tag req is not set
	NotRequired string
}

func TestRead(t *testing.T) {
	// file
	conn := NewConnector()
	err := conn.Open(TestFilePath)
	if err != nil {
		t.Error(err)
	}
	defer conn.Close()

	rd, err := conn.NewReaderByConfig(&Config{
		Sheet:         AdvSheetName,
		TitleRowIndex: 1,
		Skip:          1,
		Prefix:        "",
		Suffix:        AdvSheetSuffix,
	})
	if err != nil {
		t.Error(err)
	}
	defer rd.Close()

	idx := 0
	for rd.Next() {
		var a Advance
		err := rd.Read(&a)
		if err != nil {
			t.Error(err)
			return
		}
		expect := expectAdvanceList[idx]
		if !reflect.DeepEqual(expect, a) {
			t.Errorf("unexpect advance at %d = \n%s", idx, MustJsonPrettyString(a))
		}

		idx++
	}
}

func TestReadAll(t *testing.T) {
	// see the Advancd.suffix sheet in simple.xlsx
	conn := NewConnector()
	err := conn.Open(TestFilePath)
	if err != nil {
		t.Error(err)
	}
	rd, err := conn.NewReaderByConfig(&Config{
		Sheet:         AdvSheetName,
		TitleRowIndex: 1,
		Skip:          1,
		Prefix:        "",
		Suffix:        AdvSheetSuffix,
	})
	if err != nil {
		t.Error(err)
	}
	defer conn.Close()

	var slc []Advance
	err = rd.ReadAll(&slc)
	if err != nil {
		t.Error(err)
		return
	}
	if !reflect.DeepEqual(slc, expectAdvanceList) {
		t.Errorf("unexpect advance list: \n%s", MustJsonPrettyString(slc))
	}

}

// -- 重复Title
var expectDuplicatedTitleSliceList = [][]string{
	{
		"Value1",
		"EmptyTitleValue1",
		"Value2",
	},
	{
		"Value3",
		"EmptyTitleValue2",
		"Value2",
	},
}

type DuplicatedTitle struct {
	DuplicatedTitle string
}

func TestReadDuplicatedTitlePtrAll(t *testing.T) {
	conn := NewConnector()
	err := conn.Open(TestFilePath)
	if err != nil {
		t.Error(err)
		return
	}
	defer conn.Close()

	var stdList []*DuplicatedTitle
	rd, err := conn.NewReader(stdList)
	if err != nil {
		t.Error(err)
		return
	}
	defer rd.Close()

	err = rd.ReadAll(&stdList)
	if err != ErrDuplicatedTitles {
		t.Errorf("expect ErrDuplicatedTitles but got: %+v", err)
		return
	}
}

func TestReadDuplicatedTitleSliceList(t *testing.T) {
	conn := NewConnector()
	err := conn.Open(TestFilePath)
	if err != nil {
		t.Error(err)
		return
	}
	defer conn.Close()

	rd, err := conn.NewReader(DupSheetName)
	if err != nil {
		t.Error(err)
		return
	}
	defer rd.Close()

	var stdList [][]string
	err = rd.ReadAll(&stdList)
	if err != nil {
		t.Error(err)
		return
	}

	if !reflect.DeepEqual(expectDuplicatedTitleSliceList, stdList) {
		t.Errorf("unexpect stdlist: \n%s", MustJsonPrettyString(stdList))
	}
}
