package basics

import (
	"database/sql"

	"github.com/jinzhu/gorm"
)

type BasicsDB struct {
	Connection *gorm.DB
}

type Value struct {
	Val int
}

func (d BasicsDB) GetEverything() (*sql.Rows, error) {
	return d.Connection.Table("values").Rows()
}

func (d BasicsDB) GetValFromSql(val int) (*Value, error) {
	rows, err := d.Connection.Table("values").Where("val=?", val).Rows()
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		val := &Value{}
		err = d.Connection.ScanRows(rows, val)
		if err != nil {
			return nil, err
		}
		if val.Val != 0 {
			return val, nil
		}
	}

	return nil, nil
}

func (d BasicsDB) ReturnNothing() {}
