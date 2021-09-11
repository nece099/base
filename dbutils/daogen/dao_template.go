package daogen

var daogo_template string = `
package dao

import (
	"errors"

	"github.com/nece099/base/dbutils"
	"gorm.io/gorm"
)

func FindPage(
	db *gorm.DB,
	dest interface{},
	p *dbutils.Paging,
) error {

	if p == nil {
		return errors.New("paging is nil")
	} else {
		err := db.Count(&p.AllCount).Error
		if err != nil {
			return err
		}

		err = db.Limit(int(p.PageSize)).
			Offset(int(p.PageIndex * p.PageSize)).
			Find(dest).Error
		if err != nil {
			return err
		}
	}

	return nil
}

func RecoverTransaction(tx *gorm.DB) {
	err := recover()
	if err != nil {
		Log.Errorf("panic captured, will rollback tx, err = %v", err)
		tx.Rollback()
	}
}

func likeParams(s string) string {
	return "%" + s + "%"
}

`
