package dbMgr

import (
	"bytes"
	"database/sql"
	"github.com/Centny/gwf/dbutil"
	"github.com/Centny/gwf/log"
	"reflect"
	"strings"
)

type DbModel interface {
	GetTableName() string
}

type dbfield struct {
	pk         bool
	column     string
	fieldName  string
	fieldValue interface{}
}

func Insert(o DbModel) (int64, error) {
	dfs := parseToDbfield(o)
	s := bytes.Buffer{}
	s.WriteString("insert into ")
	s.WriteString(o.GetTableName())
	s.WriteString("(")
	s.WriteString(selectColumn(dfs))
	s.WriteString(")values(")
	s.WriteString(placeholders(dfs))
	s.WriteString(")")
	log.D("sql:%v", s.String())
	log.D("args:%v", args(dfs))
	return dbutil.DbInsert(DbConn(), s.String(), args(dfs))
}

func InsertTx(tx *sql.Tx, o interface{}, tableName string) (int64, error) {
	dfs := parseToDbfield(o)
	s := bytes.Buffer{}
	s.WriteString("insert into ")
	s.WriteString(tableName)
	s.WriteString("(")
	s.WriteString(selectColumn(dfs))
	s.WriteString(")values(")
	s.WriteString(placeholders(dfs))
	s.WriteString(")")
	log.D("sql:%v", s.String())
	log.D("args:%v", args(dfs))
	return dbutil.DbInsert2(tx, s.String(), args(dfs))
}

func args(fields []*dbfield) []interface{} {
	var args []interface{}
	for _, v := range fields {
		if v.pk {
		} else {
			args = append(args, v.fieldValue)
		}
	}
	return args
}

func placeholders(fields []*dbfield) string {
	var fs []string
	for _, v := range fields {
		if v.pk {
			fs = append(fs, "null")
		} else {
			fs = append(fs, "?")
		}
	}
	return strings.Join(fs, ",")
}

func selectColumn(fields []*dbfield) string {
	var fs []string
	for _, v := range fields {
		fs = append(fs, v.column)
	}
	return strings.Join(fs, ",")
}

func parseToDbfield(o interface{}) []*dbfield {

	objType := reflect.TypeOf(o)
	objVal := reflect.ValueOf(o)

	if objType.Kind() == reflect.Ptr {
		objType = objType.Elem()
	}
	if objVal.Kind() == reflect.Ptr {
		objVal = objVal.Elem()
	}

	fields := []*dbfield{}
	for i := 0; i < objType.NumField(); i++ {
		eachField := objType.Field(i)
		tag, pk := parseField(eachField)
		if tag != "" {
			var fd dbfield
			fd.fieldName = eachField.Name
			fd.column = tag
			fd.pk = pk
			fd.fieldValue = objVal.Field(i).Interface()
			fields = append(fields, &fd)
		}
	}
	return fields
}

func parseField(eachField reflect.StructField) (string, bool) {
	tag := eachField.Tag.Get("m2s")
	tags := strings.Split(tag, ",")
	if len(tags) == 2 {
		if tags[1] == "autoinc" {
			return tags[0], true
		} else {
			return tags[0], false
		}
	}
	if len(tags) == 1 {
		return tags[0], false
	}
	return "", false
}

func DbQueryS(db *sql.DB, res interface{}, query string, args ...interface{}) error {

	if reflect.TypeOf(res).Elem().Kind() == reflect.Slice {
		return dbQueryStructSlice(db, res, query, args...)
	}
	return dbQueryStruct(db, res, query, args...)
}

func dbQueryStructSlice(db *sql.DB, res interface{}, query string, args ...interface{}) error {
	rval := reflect.ValueOf(res).Elem()

	//this is the slice type
	rtype := rval.Type().Elem()

	if rtype.Kind() == reflect.Ptr {
		rtype = rtype.Elem()
	}

	rows, err := db.Query(query, args...)
	if err != nil {
		return err
	}
	cols, err := rows.Columns()
	if err != nil {
		return err
	}
	var scanFieldIdx []int
	for _, col := range cols {
		for i := 0; i < rtype.NumField(); i++ {
			if rtype.Field(i).Tag.Get("m2s") == col {
				scanFieldIdx = append(scanFieldIdx, i)
			}
		}
	}
	for rows.Next() {
		pls := []interface{}{}
		rtval := reflect.New(rtype)
		for _, idx := range scanFieldIdx {
			pls = append(pls, rtval.Elem().Field(idx).Addr().Interface())
		}
		if err := rows.Scan(pls...); err != nil {
			return err
		}
		rval = reflect.Append(rval, rtval)
	}

	reflect.Indirect(reflect.ValueOf(res)).Set(rval)

	return nil
}

func dbQueryStruct(db *sql.DB, res interface{}, query string, args ...interface{}) error {
	rows, err := db.Query(query, args...)
	if err != nil {
		return err
	}
	cols, err := rows.Columns()
	if err != nil {
		return err
	}
	rval := reflect.ValueOf(res).Elem()

	rtype := rval.Type()
	pls := []interface{}{}
	for _, col := range cols {
		for i := 0; i < rtype.NumField(); i++ {
			if rtype.Field(i).Tag.Get("m2s") == col {
				pls = append(pls, rval.Field(i).Addr().Interface())
			}
		}
	}
	log.D("query:%v \n pls:%v", query, pls)
	if rows.Next() {
		return rows.Scan(pls...)
	}
	return sql.ErrNoRows
}
