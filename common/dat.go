package common

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
)

// ColumnSpec 列定义
type ColumnSpec struct {
	Name     string
	Type     string
	Comment  string
	Nullable string
	Default  string
	Key      string
	Extra    string

	// 字符集
	CharacterSet string
	Collation    string
}

// TableSpec 表定义
type TableSpec struct {
	Name      string
	Comment   string
	Collation string
	Columns   []ColumnSpec
}

// GetColumns 获取当前表的所有字段定义
func (ts *TableSpec) GetColumns() []ColumnSpec {
	return ts.Columns
}

// GetColumn 以字段名获取当前表的某一字段定义
func (ts *TableSpec) GetColumn(field string) (*ColumnSpec, error) {
	for _, cs := range ts.Columns {
		if cs.Name == field {
			return &cs, nil
		}
	}
	return nil, errors.New(fmt.Sprintf("column named %s not existed in table %s", field, ts.Name))
}

// AppendColumn 为表定义添加列定义
func (ts *TableSpec) AppendColumn(cs *ColumnSpec) error {
	for _, ecs := range ts.Columns {
		if ecs.Name == cs.Name {
			return errors.New(fmt.Sprintf("column named %s already existed in table %s", cs.Name, ts.Name))
		}
	}
	ts.Columns = append(ts.Columns, *cs)
	return nil
}

// SchemaSpec 库定义
type SchemaSpec struct {
	Name                string
	DefaultCharacterSet string
	DefaultCollation    string
	Tables              []TableSpec
}

// GetTables 获取当前库的所有表定义
func (ss *SchemaSpec) GetTables() []TableSpec {
	return ss.Tables
}

// GetTable 以表名名获取当前库的某一表定义
func (ss *SchemaSpec) GetTable(name string) (*TableSpec, error) {
	for _, ts := range ss.Tables {
		if ts.Name == name {
			return &ts, nil
		}
	}
	return nil, errors.New(fmt.Sprintf("table named %s not existed in schema %s", name, ss.Name))
}

// AppendTable 为库定义添加表定义
func (ss *SchemaSpec) AppendTable(ts *TableSpec) error {
	if ss.Tables == nil {
		ss.Tables = []TableSpec{}
	}
	for _, ets := range ss.Tables {
		if ets.Name == ts.Name {
			return errors.New(fmt.Sprintf("table named %s already existed in schema %s", ts.Name, ss.Name))
		}
	}
	ss.Tables = append(ss.Tables, *ts)
	return nil
}

// LoadFromFile 从json文件中加载"库定义描述"
func (ss *SchemaSpec) LoadFromFile(filename string) error {
	bDat, err := ioutil.ReadFile(filename)
	if err != nil {
		return err
	}
	return json.Unmarshal(bDat, ss)
}

// SaveToFile 将库定义描述存储到文件
func (ss *SchemaSpec) SaveToFile(filename string) error {
	b, err := json.Marshal(ss)

	if err != nil {
		return err
	}
	return ioutil.WriteFile(filename, b, 0644)
}
