package db

import (
	"database/sql"
	"fmt"
	"github.com/Centny/gwf/dbutil"
	"github.com/Centny/gwf/log"
)

func migrateAddTable(db *sql.DB, table, def string) error {
	log.I("begin to add table %s", table)
	exist, err := tableExist(db, table)
	if err != nil {
		return err
	}
	if exist {
		log.I("table %s already exist, no need to add", table)
		return nil
	}

	if _, err := db.Exec(def); err != nil {
		return err
	}
	log.I("add table %s success", table)
	return nil
}

func tableExist(db *sql.DB, table string) (bool, error) {
	str, err := dbutil.DbQueryStr(db, `show tables like ?`, table)
	return str != "", err
}

func tableColumnExist(db *sql.DB, table, column string) (bool, error) {
	dbname, err := dbutil.DbQueryStr(db, `SELECT DATABASE();`)
	if err != nil {
		return false, err
	}
	q := `SELECT count(*) FROM information_schema.COLUMNS WHERE TABLE_SCHEMA = ? AND TABLE_NAME = ? AND COLUMN_NAME = ?`
	cnt, err := dbutil.DbQueryI(db, q, dbname, table, column)
	return cnt > 0, err
}

func migrateAddColumn(db *sql.DB, table, column, typ, defaul, comment string) error {
	log.I("begin to add column[%s] to the table[%s]...", column, table)

	exists, err := tableColumnExist(db, table, column)
	if err != nil {
		log.E("get table column exists fail %v", err.Error())
		return err
	}
	if exists == false {
		_, err := db.Exec(fmt.Sprintf(`ALTER TABLE %s ADD %s %s default %s comment '%s';`, table, column, typ, defaul, comment))
		if err != nil {
			log.E("add column[%s] fail %v", column, err.Error())
			return err
		}
		log.I("have added column[%s] to the table[%s]", column, table)

		//set the old data of this field to default value
		_, err = db.Exec(fmt.Sprintf(`update %s set %s = ?`, table, column), "")
		return err

	} else {
		log.I("column[%s] already exists no need to add", column)
	}
	return nil
}

func RunTx(fn func(*sql.Tx) error) error {
	tx, err := DbConn().Begin()
	if err != nil {
		return err
	}
	err = fn(tx)
	if err != nil {
		return fmt.Errorf("err:%v rollbackErr:%v", err, tx.Rollback())
	}
	return tx.Commit()
}

func Init(driver string, db string) error {
	t_dbMgr.Init(driver, db, "rcp")
	t_dbMgr.SetMaxOpenConn("rcp", MaxOpenConns)
	return CheckDb(DbConn())
}

type ModelField struct {
	reflect.StructField
	Column  string
	Null    bool
	Autoinc bool
}

// Represents meta info about a model
type ModelInfo struct {
	Type           reflect.Type
	TableName      string
	Fields         []*ModelField
	FieldsSimple   string
	FieldsPrefixed string
	FieldsInsert   string
	Placeholders   string
}

// Global cache
var allModelInfos = map[string]*ModelInfo{}

func GetModelInfo(i interface{}) *ModelInfo {
	t := reflect.TypeOf(i)
	return GetModelInfoFromType(t)
}

func IsScanner(i interface{}) bool {
	_, ok := i.(sql.Scanner)
	return ok
}

// Call this once after each struct type declaration
func GetModelInfoFromType(modelType reflect.Type) *ModelInfo {
	if modelType.Kind() == reflect.Ptr {
		modelType = modelType.Elem()
	}
	if modelType.Kind() != reflect.Struct {
		return nil
	}
	if modelType.Implements(reflect.TypeOf((*sql.Scanner)(nil)).Elem()) {
		return nil
	}

	modelName := modelType.Name()

	// Check cache
	if allModelInfos[modelName] != nil {
		return allModelInfos[modelName]
	}

	// Construct
	m := &ModelInfo{}
	allModelInfos[modelName] = m
	m.Type = modelType
	m.TableName = strings.ToLower(modelName)

	// Fields
	numFields := m.Type.NumField()
	for i := 0; i < numFields; i++ {
		field := m.Type.Field(i)
		if field.Tag.Get("m2s") != "" {
			column, null, autoinc := parseDBTag(field.Tag.Get("m2s"))
			m.Fields = append(m.Fields, &ModelField{field, column, null, autoinc})
		}
	}

	// Simple & Prefixed
	fieldNames := []string{}
	fieldInsertNames := []string{}
	ph := []string{}
	for _, field := range m.Fields {
		fieldName, _, _ := parseDBTag(field.Tag.Get("m2s"))
		fieldNames = append(fieldNames, fieldName)
		if !field.Autoinc {
			fieldInsertNames = append(fieldInsertNames, fieldName)
			ph = append(ph, "?")
		}
	}

	m.FieldsSimple = strings.Join(fieldNames, ", ")
	m.FieldsPrefixed = m.TableName + "." + strings.Join(fieldNames, ", "+m.TableName+".")
	m.FieldsInsert = strings.Join(fieldInsertNames, ", ")
	m.Placeholders = strings.Join(ph, ", ")

	return m
}

// Helper
func parseDBTag(tag string) (fieldName string, null bool, autoinc bool) {
	s := strings.Split(tag, ",")
	fieldName = s[0]
	for _, ss := range s[1:] {
		if ss == "null" {
			null = true
		}
		if ss == "autoinc" {
			autoinc = true
		}
	}
	return
}
func MustExec(q string, args ...interface{}) {
	_, err := DbConn().Exec(q, args...)
	if err != nil {
		panic(err)
	}
}
