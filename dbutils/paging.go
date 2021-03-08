package dbutils

import "fmt"

type Paging struct {
	AllCount  int64
	PageIndex int64
	PageSize  int64
}

func (p *Paging) PageSql(sql string) string {

	newSql := `
	select * from (` + sql + `) ps 
	limit ` + fmt.Sprintf("%v", p.PageSize) + ` offset ` + fmt.Sprintf("%v", p.PageSize*p.PageIndex)

	return newSql
}

func (p *Paging) CountSql(sql string) string {
	countSql := `select count(*) as all_count from (` + sql + `) ps`
	return countSql
}

var PAGE_FULL = &Paging{
	PageIndex: 0,
	PageSize:  999999999,
}
