package inspector

import (
	"database/sql"
	"ddgen/utils"
	"fmt"
	"log"
	"reflect"
	"strings"
)

type ColumnDesc struct {
	ColumnName    sql.NullString `col:"COLUMN_NAME" render:"列"`
	ColumnType    sql.NullString `col:"COLUMN_TYPE" render:"类型"`
	ColumnDefault sql.NullString `col:"COLUMN_DEFAULT" render:"默认"`
	IsNullable    sql.NullString `col:"IS_NULLABLE" render:"Null"`
	ColumnKey     sql.NullString `col:"COLUMN_KEY" render:"Key"`
	Extra         sql.NullString `col:"EXTRA" render:"额外"`
	ColumnComment sql.NullString `col:"COLUMN_COMMENT" render:"评论"`

	OrdinalPosition sql.NullInt64  `col:"ORDINAL_POSITION"`
	Privileges      sql.NullString `col:"PRIVILEGES"`
}

func (cd ColumnDesc) GetRenderLabels() []string {
	vt := reflect.TypeOf(cd)
	var tableTitleValues []string
	for l := 0; l < vt.NumField(); l++ {
		f := vt.Field(l)
		renderText := f.Tag.Get("render")
		if len(renderText) > 0 {
			tableTitleValues = append(tableTitleValues, renderText)
		}
	}
	return tableTitleValues
}

func (cd ColumnDesc) GetRenderFields() []string {
	vt := reflect.TypeOf(cd)
	var renderFieldsIdx []string
	for l := 0; l < vt.NumField(); l++ {
		f := vt.Field(l)
		renderText := f.Tag.Get("render")
		if len(renderText) > 0 {
			renderFieldsIdx = append(renderFieldsIdx, f.Name)
		}
	}
	return renderFieldsIdx
}

const NullValue = "<nil>"
const EmptyValue = " "

func (cd ColumnDesc) GetRenderValues(fieldsOrder []string) []interface{} {
	v := reflect.ValueOf(cd)

	var renderFields []interface{}
	for _, fieldName := range fieldsOrder {

		fv := v.FieldByName(fieldName)
		fvt := fv.Type()
		switch fvt.String() {
		case "sql.NullString":
			d := fv.Interface().(sql.NullString)
			if d.Valid {
				if len(d.String) == 0 {
					renderFields = append(renderFields, EmptyValue)
				} else {
					renderFields = append(renderFields, d.String)
				}

			} else {
				renderFields = append(renderFields, NullValue)
			}
			break
		case "sql.NullBool":
			d := fv.Interface().(sql.NullBool)
			if d.Valid {
				renderFields = append(renderFields, d.Bool)
			} else {
				renderFields = append(renderFields, NullValue)
			}
			break
		case "sql.NullInt64":
			d := fv.Interface().(sql.NullInt64)
			if d.Valid {
				renderFields = append(renderFields, d.Int64)
			} else {
				renderFields = append(renderFields, NullValue)
			}
			break
		case "sql.NullFloat64":
			d := fv.Interface().(sql.NullFloat64)
			if d.Valid {
				renderFields = append(renderFields, d.Float64)
			} else {
				renderFields = append(renderFields, NullValue)
			}
			break
		default:
			log.Fatalln("unknown type", fvt.String())
		}
	}
	return renderFields
}

type TableDesc struct {
	TableName string
	Columns   []ColumnDesc
}

type SchemaDesc struct {
	SchemaName string
	Tables     []TableDesc
}

type DBInspector struct {
	// for database connection
	DriverName string
	DataSource string
	db         *sql.DB

	// desc data (for rendering DD)
	SchemasOnly []string
	Schemas     []SchemaDesc
}

func CreateDBInspector(driverName string, dataSource string) DBInspector {
	return DBInspector{
		DriverName:  driverName,
		DataSource:  dataSource,
		db:          nil,
		Schemas:     []SchemaDesc{},
		SchemasOnly: []string{},
	}
}

func (i *DBInspector) checkErr(err error) {
	log.Fatalln(err)
}

func (i *DBInspector) Initialize() {
	db, err := sql.Open(i.DriverName, i.DataSource)
	if err != nil {
		panic(err)
	}
	i.db = db
}

func (i *DBInspector) Destroy() {
	err := i.db.Close()
	if err != nil {
		panic(err)
	}
}

func (i *DBInspector) InspectSchemas() {
	rows, err := i.db.Query("show databases;")
	if err != nil {
		panic(err)
	}

	for rows.Next() {
		sd := SchemaDesc{
			SchemaName: "",
			Tables:     []TableDesc{},
		}

		err := rows.Scan(&sd.SchemaName)
		if err != nil {
			panic(err)
		}

		if !utils.ContainsString(i.SchemasOnly, sd.SchemaName) {
			continue
		}
		err = i.inspectSchema(&sd)
		if err != nil {
			i.checkErr(err)
			continue
		}
		i.Schemas = append(i.Schemas, sd)
	}
}

func (i *DBInspector) inspectSchema(s *SchemaDesc) error {
	log.Printf("* inspecting schema '%s'", s.SchemaName)

	tblQry := fmt.Sprintf("select table_name from information_schema.tables where table_schema='%s'", s.SchemaName)
	//log.Println(tblQry)
	rows, err := i.db.Query(tblQry)
	if err != nil {
		return err
	}
	for rows.Next() {
		td := TableDesc{
			TableName: "",
			Columns:   []ColumnDesc{},
		}

		err := rows.Scan(&td.TableName)

		if err != nil {
			panic(err)
		}
		log.Printf("** inspecting table '%s.%s'", s.SchemaName, td.TableName)

		err, colDatArr := i.queryAndAutoScan(
			fmt.Sprintf("select %%s from information_schema.columns where table_name='%s' and table_schema='%s'", td.TableName, s.SchemaName),
			reflect.TypeOf(ColumnDesc{}),
		)
		if err != nil {
			return err
		}

		for _, colDat := range colDatArr {
			td.Columns = append(td.Columns, colDat.(ColumnDesc))
		}

		s.Tables = append(s.Tables, td)
	}
	return nil
}

// QueryAndAutoScan 依据传入对带有col tag对结构体及sql来进行查询及自动填充
func (i *DBInspector) queryAndAutoScan(sqlFmt string, vt reflect.Type) (error, []interface{}) {

	count := vt.NumField()

	var dbCols []string
	var fields []string

	for i := 0; i < count; i++ {
		ft := vt.Field(i)
		dbColName := ft.Tag.Get("col")
		if len(dbColName) != 0 {
			dbCols = append(dbCols, dbColName)
			fields = append(fields, ft.Name)
		}
	}
	querySql := fmt.Sprintf(sqlFmt, strings.Join(dbCols, ","))
	rows, err := i.db.Query(querySql)

	if err != nil {
		i.checkErr(err)
		return err, nil
	}

	var datArr []interface{}

	for rows.Next() {
		rowData := reflect.New(vt)
		v := rowData.Elem()
		var sf []interface{}

		for _, vf := range fields {
			f := v.FieldByName(vf)
			sf = append(sf, f.Addr().Interface())
		}
		err := rows.Scan(sf...)

		if err != nil {
			i.checkErr(err)
		}
		datArr = append(datArr, v.Interface())
	}
	return nil, datArr
}
