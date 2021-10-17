package excel_test

import (
	"fmt"
	"reflect"
	"testing"

	"github.com/zhao520a1a/go-utils/excel"
)

type StandardFieldConfig excel.Standard

func (StandardFieldConfig) GetXLSXSheetName() string {
	return "Standard"
}

func (StandardFieldConfig) GetXLSXFieldConfigs() map[string]excel.FieldConfig {
	return map[string]excel.FieldConfig{
		"Name": {
			ColumnName: "NameOf",
		},
		"NameOf": {
			ColumnName: "NameOf",
		},
		"Age": {
			ColumnName: "AgeOf",
		},
		"Slice": {
			Split: "|",
		},
		"Temp": {
			ColumnName: "UnmarshalString",
		},
		"WantIgnored": {
			Ignore: true,
		},
	}
}

var expectStandardFieldConfigList = []StandardFieldConfig{
	{
		ID:      1,
		Name:    "Andy",
		NamePtr: excel.StrPtr("Andy"),
		Age:     1,
		Slice:   []int{1, 2},
		Temp: &excel.Temp{
			Foo: "Andy",
		},
	},
	{
		ID:      2,
		Name:    "Leo",
		NamePtr: excel.StrPtr("Leo"),
		Age:     2,
		Slice:   []int{2, 3, 4},
		Temp: &excel.Temp{
			Foo: "Leo",
		},
	},
	{
		ID:      3,
		Name:    "Ben",
		NamePtr: excel.StrPtr("Ben"),
		Age:     3,
		Slice:   []int{3, 4, 5, 6},
		Temp: &excel.Temp{
			Foo: "Ben",
		},
	},
	{
		ID:      4,
		Name:    "Ming",
		NamePtr: excel.StrPtr("Ming"),
		Age:     4,
		Slice:   []int{1},
		Temp: &excel.Temp{
			Foo: "Ming",
		},
	},
}

func TestReadStandardFieldConfigSimple(t *testing.T) {
	var stdList []StandardFieldConfig
	err := excel.UnmarshalXLSX(excel.TestFilePath, &stdList)
	if err != nil {
		t.Error(err)
		return
	}
	if !reflect.DeepEqual(stdList, expectStandardFieldConfigList) {
		t.Errorf("unexprect std list: %s", excel.MustJsonPrettyString(stdList))
	}
}

func TestReadStandardFieldConfig(t *testing.T) {
	conn := excel.NewConnector()
	err := conn.Open(excel.TestFilePath)
	if err != nil {
		t.Error(err)
		return
	}
	defer conn.Close()

	rd, err := conn.NewReader(excel.StdSheetName)
	if err != nil {
		t.Error(err)
		return
	}
	defer rd.Close()

	idx := 0
	for rd.Next() {
		var s StandardFieldConfig
		if err := rd.Read(&s); err != nil {
			fmt.Println(err)
			return
		}
		expectStd := expectStandardFieldConfigList[idx]
		if !reflect.DeepEqual(s, expectStd) {
			t.Errorf("unexpect std at %d = \n%s", idx, excel.MustJsonPrettyString(expectStd))
		}
		idx++
	}
}

func TestReadStandardFieldConfigIndex(t *testing.T) {
	conn := excel.NewConnector()
	err := conn.Open(excel.TestFilePath)
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
		var s StandardFieldConfig
		if err := rd.Read(&s); err != nil {
			fmt.Println(err)
			return
		}
		expectStd := expectStandardFieldConfigList[idx]
		if !reflect.DeepEqual(s, expectStd) {
			t.Errorf("unexpect std at %d = \n%s", idx, excel.MustJsonPrettyString(expectStd))
		}
		idx++
	}
}

func TestReadStandardFieldConfigAll(t *testing.T) {
	conn := excel.NewConnector()
	err := conn.Open(excel.TestFilePath)
	if err != nil {
		t.Error(err)
		return
	}
	defer conn.Close()

	var stdList []StandardFieldConfig
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
	if !reflect.DeepEqual(expectStandardFieldConfigList, stdList) {
		t.Errorf("unexpect stdlist: \n%s", excel.MustJsonPrettyString(stdList))
	}
}
