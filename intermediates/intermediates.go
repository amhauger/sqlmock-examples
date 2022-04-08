package intermediates

import (
	"time"

	"github.com/google/uuid"
	"github.com/jinzhu/gorm"
	"github.com/rs/zerolog/log"
)

type IntermediateDB struct {
	Connection *gorm.DB
}

type Value struct {
	ID   string
	Val  int
	Time time.Time
}

func (d IntermediateDB) InsertVal() error {
	v := Value{uuid.New().String(), 1, time.Now()}
	res := d.Connection.Save(&v)
	if err := res.Error; err != nil {
		return err
	}

	return nil
}

func (d IntermediateDB) GetValFromSqlByID() (*Value, error) {
	id := uuid.New().String()
	rows, err := d.Connection.Table("users").Where("id=?", id).Rows()
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	log.Info().Interface("rows", rows).Msg("dorp")
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
