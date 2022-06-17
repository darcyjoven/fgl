/*
	事务处理
*/
package fgl

import (
	"log"
)

func (d *DataBase) Begin() {
	_, err := d.DB.Beginx()
	if err != nil {
		log.Println(err)
		return
	}
}

// 提交事务
func (d *Stmp) Commit() error {
	defer clearTx(d)
	if d.tx == nil {
		log.Println("未开启事务，不需要提交")
		return nil
	}
	return d.tx.Commit()
}

// 回滚事务
func (d *Stmp) Rollback() error {
	defer clearTx(d)
	if d.tx == nil {
		log.Println("未开启事务，不需要提交")
		return nil
	}
	return d.tx.Rollback()
}

// 检查是否再事务中,不再事务中，直接报错。不能继续执行
func (d *Stmp) MustBegin() *Stmp {
	if d.tx == nil {
		d.DB = nil
	}
	return d
}
