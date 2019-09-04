package mysql5_7

import (
	"database/sql"
	"ddgen/common"
	"ddgen/inspector"
	"fmt"
	"log"
)

const Driver = "mysql"

type Inspector struct {
	id       string
	emptyTag string
	db       *sql.DB
}

func (i Inspector) GetComponentId() string {
	return i.id
}

func (i *Inspector) createEngine(dbSrc string) error {
	db, err := sql.Open(Driver, dbSrc)
	if err != nil {
		return err
	}
	i.db = db
	return nil
}

func (i *Inspector) destroyEngine() {
	if i.db != nil {
		err := i.db.Close()
		if err != nil {
			log.Fatalln(err)
		}
	}
}

func (i Inspector) Inspect(dbSrc string, schema string, params string) (*common.SchemaSpec, error) {
	err := i.createEngine(dbSrc)
	if err != nil {
		return nil, err
	}
	defer i.destroyEngine()

	ss := common.SchemaSpec{
		Name:                schema,
		DefaultCharacterSet: "",
		DefaultCollation:    "",
	}

	querySql := fmt.Sprintf("select DEFAULT_CHARACTER_SET_NAME, DEFAULT_COLLATION_NAME from schemata where SCHEMA_NAME='%s'", schema)
	var defaultCharacterSet, defaultCollation sql.NullString
	err = i.db.QueryRow(querySql).Scan(&defaultCharacterSet, &defaultCollation)

	if defaultCharacterSet.Valid {
		ss.DefaultCharacterSet = defaultCharacterSet.String
	} else {
		ss.DefaultCharacterSet = i.emptyTag
	}

	if defaultCollation.Valid {
		ss.DefaultCollation = defaultCollation.String
	} else {
		ss.DefaultCollation = i.emptyTag
	}

	if err != nil {
		return nil, err
	}

	log.Printf("* inspected schema %s", schema)

	// inspect tables

	rows, err := i.db.Query(fmt.Sprintf("select TABLE_NAME, TABLE_COMMENT, TABLE_COLLATION from tables  where TABLE_SCHEMA='%s'", schema))
	if err != nil {
		return nil, err
	}
	for rows.Next() {
		ts := common.TableSpec{
			Name:      "",
			Comment:   "",
			Collation: "",
		}
		var name, comment, collation sql.NullString
		err := rows.Scan(&name, &comment, &collation)
		if err != nil {
			return nil, err
		}

		if name.Valid {
			ts.Name = name.String
		} else {
			ts.Name = i.emptyTag
		}

		if comment.Valid {
			ts.Comment = comment.String
		} else {
			ts.Comment = i.emptyTag
		}

		if collation.Valid {
			ts.Collation = collation.String
		} else {
			ts.Collation = i.emptyTag
		}

		// inspect columns in table
		log.Printf("-* inspected table %s", ts.Name)

		rows, err := i.db.Query(fmt.Sprintf("select COLUMN_NAME, COLUMN_DEFAULT, IS_NULLABLE, CHARACTER_SET_NAME, COLLATION_NAME, COLUMN_TYPE, COLUMN_KEY,EXTRA,COLUMN_COMMENT from columns where TABLE_SCHEMA='%s' and TABLE_NAME='%s'", schema, ts.Name))
		if err != nil {
			log.Printf("query table %s columns failed %s", ts.Name, err)
			continue
		}

		scanErr := false
		for rows.Next() {
			var name, cDefault, nullable, characterSet, collation, cType, key, extra, comment sql.NullString
			err := rows.Scan(&name, &cDefault, &nullable, &characterSet, &collation, &cType, &key, &extra, &comment)
			if err != nil {
				scanErr = true
				log.Printf("scan table %s columns failed %s", ts.Name, err)
				break
			}

			cs := common.ColumnSpec{
				Name:         i.fieldValueToStr(name),
				Type:         i.fieldValueToStr(cType),
				Comment:      i.fieldValueToStr(comment),
				Default:      i.fieldValueToStr(cDefault),
				Nullable:     i.fieldValueToStr(nullable),
				Key:          i.fieldValueToStr(key),
				Extra:        i.fieldValueToStr(extra),
				CharacterSet: i.fieldValueToStr(characterSet),
				Collation:    i.fieldValueToStr(collation),
			}

			err = ts.AppendColumn(&cs)
			if err != nil {
				log.Println(err)
				scanErr = true
				break
			}

			log.Printf("--* inspected column %s", cs.Name)
		}
		if scanErr {
			log.Printf("inspect table %s columns failed", ts.Name)
			continue
		}

		// store
		err = ss.AppendTable(&ts)
		if err != nil {
			log.Printf("append table %s failed %s", ts.Name, err)
			continue
		}
	}
	return &ss, err
}

func (i *Inspector) fieldValueToStr(f interface{}) string {

	switch f.(type) {
	case sql.NullString:
		if f.(sql.NullString).Valid {
			return f.(sql.NullString).String
		} else {
			return i.emptyTag
		}
	case sql.NullFloat64:
		if f.(sql.NullFloat64).Valid {
			return fmt.Sprint(f.(sql.NullFloat64).Float64)
		} else {
			return i.emptyTag
		}
	case sql.NullInt64:
		if f.(sql.NullInt64).Valid {
			return fmt.Sprint(f.(sql.NullInt64).Int64)
		} else {
			return i.emptyTag
		}
	case sql.NullBool:
		if f.(sql.NullBool).Valid {
			return fmt.Sprint(f.(sql.NullBool).Bool)
		} else {
			return i.emptyTag
		}
	default:
		log.Println("unknown field type", f)
		return "<unknown>"
	}
}

func init() {
	inspector.GlobalRendererRepository.Register(Inspector{id: "mysql5.7", emptyTag: "<empty>"})
}
