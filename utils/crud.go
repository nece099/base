package utils

import (
	"github.com/nece099/base/datasource"

	"github.com/jinzhu/gorm"
)

type TCrudService struct{}

var crudService = new(TCrudService)

func CurdServiceInstance() *TCrudService {
	return crudService
}

func (s *TCrudService) Create(mod interface{}) error {

	db := datasource.DataSourceInstance().Master()

	err := db.Create(mod).Error
	if err != nil {
		return err
	}

	return nil
}

func (s *TCrudService) Update(mod interface{}) error {

	db := datasource.DataSourceInstance().Master()

	err := db.Save(mod).Error
	if err != nil {
		return err
	}

	return nil
}

func (s *TCrudService) List(out interface{}, where string, whereParams ...interface{}) error {

	db := datasource.DataSourceInstance().Master()

	if len(whereParams) > 0 {
		db = db.Where(where, whereParams...)
	}

	err := db.Order(" id desc ").Find(out).Error
	if err != nil {
		return err
	}

	return nil
}

func (s *TCrudService) TxReadForUpdate(tx *gorm.DB, mod interface{}) error {

	err := tx.Set("gorm:query_option", "FOR UPDATE").Find(mod).Error
	if err != nil {
		tx.Rollback()
		return err
	}

	return nil
}

func (s *TCrudService) TxRead(tx *gorm.DB, mod interface{}) error {

	err := tx.Find(mod).Error
	if err != nil {
		tx.Rollback()
		return err
	}

	return nil
}

func (s *TCrudService) TxCreate(tx *gorm.DB, mod interface{}) error {

	err := tx.Create(mod).Error
	if err != nil {
		tx.Rollback()
		return err
	}

	return nil
}

func (s *TCrudService) TxUpdate(tx *gorm.DB, mod interface{}) error {

	err := tx.Save(mod).Error
	if err != nil {
		tx.Rollback()
		return err
	}

	return nil
}

func (s *TCrudService) TxDelete(tx *gorm.DB, mod interface{}) error {

	err := tx.Delete(mod).Error
	if err != nil {
		tx.Rollback()
		return err
	}

	return nil
}
