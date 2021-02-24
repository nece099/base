package paging

import "gorm.io/gorm"

type Paging struct {
	AllCount  int64
	PageIndex int64
	PageSize  int64
}

func GetPagingDB(db *gorm.DB, model interface{}, p *Paging) *gorm.DB {
	if p != nil && p.PageSize > 0 {
		return db.Model(model).Count(&p.AllCount).Offset(p.PageIndex * p.PageSize).Limit(p.PageSize)
	}

	return db
}
