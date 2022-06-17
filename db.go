package fgl

import (
	"log"

	"github.com/jmoiron/sqlx"
	_ "github.com/sijms/go-ora/v2"
)

type DataBase struct {
	DB   *sqlx.DB   // 数据库连接
	stmt *sqlx.Stmt // prepare后的结果
	tx   *sqlx.Tx   // 开启事务使用tx执行
	rows *sqlx.Rows // Next的行数
}
type build struct {
	head      string // sql类型 SELECT//UPDATE//DELTE...
	unique    bool
	cols      []string
	tables    []string
	joins     map[string][]join //[]*join
	joinvalue []any
	where     string // where 条件
	values    []any
	group     string
	order     string
	having    string
	session   *bool
}
type join struct {
	table  string
	where  string
	values []any
}

type Stmp struct {
	sql  *string
	DB   *sqlx.DB
	stmt *sqlx.Stmt
	tx   *sqlx.Tx
	rows *sqlx.Rows
}

const (
	select_ = "SELECT"
	update  = "UPDATE"
	delete  = "DELETE FROM"
	insert  = "INSERT INTO"
	create  = "CREATE GLOBAL TEMPORARY"
	drop    = "DROP TABLE"
	index   = "CREATE INDEX"
)

func NewDB(url string) *DataBase {
	var d DataBase
	db, err := sqlx.Connect("oracle", url)
	if err != nil {
		log.Println(err)
		return nil
	}
	d.DB = db
	return &d
}

func dbErr(db *DataBase, err error) {
	log.Panicln("err is ", err)
	if db.tx != nil {
		// 事务中出错，直接回滚
		db.tx.Rollback()
	}
}

/*
	func (*sqlx.DB) Prepare(){}

	1. 静态sql
	type build struct {
		*sqlx.DB
	}
	func (*build) Select(){}
	func (*build) SelectDo(){}
	func (*build) Open() *stmt{}
	// 这是未开启事务动态sql
	func (*build) Prepare() *stmt{}

	2. 动态sql
	type stmt struct{
		*sqlx.Stmt
		*sqlx.DB //can not be nil
		*sqlx.Tx //can be nil
		*sqlx.Rows
	}
	func (*stmt) Prepare() *stmt{}
	func (*stmt) Commit() *stmt{}
	func (*stmt) Rollback() *stmt{}

	3. 事务动态sql

*/
