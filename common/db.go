package common

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
