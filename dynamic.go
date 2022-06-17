/*
	动态SQL
	复杂SQL
*/
package fgl

import (
	"log"

	"github.com/jmoiron/sqlx"
)

// 创建prepare
func (d *DataBase) Prepare(sql string) *Stmp {

	stmt, err := d.DB.Preparex(sql)
	if err != nil {
		log.Println(err)
		return nil
	}
	return &Stmp{
		sql:  &sql,
		DB:   d.DB,
		stmt: stmt,
		tx:   nil,
	}
}

func (f *DataBase) Close() error {
	// defer new(f)
	return f.DB.Close()
}

// 取第一行
func (d *Stmp) Fetch(args []any, dest ...any) error {
	return d.stmt.QueryRow(args...).Scan(dest...)
}

// 一笔笔循环
func (d *Stmp) Next(args []any, dest ...any) error {
	getRows(d, args...)
	if d.rows.Next() {
		err := d.rows.Scan(dest...)
		if err != nil {
			return err
		}
	} else {
		return log.Output(1, "已经是最后一笔")
	}
	return nil
}

// 取第一行
func (d *Stmp) FetchAny(args ...any) []any {
	var r []any
	rows, err := d.stmt.Queryx(args...)
	if err != nil {
		log.Println(err)
		return nil
	}
	if rows.Next() {
		r, err = rows.SliceScan()
		if err != nil {
			log.Println(err)
			return nil
		}
	}
	if err != nil {
		log.Println(err)
		return nil
	}
	return r
}

// 循环用
func (d *Stmp) ForeachAny(args ...any) [][]any {
	var data [][]any
	rows, err := d.stmt.Queryx(args...)
	if err != nil {
		log.Println(err)
		return nil
	}
	for rows.Next() {
		temp, err := rows.SliceScan()
		if err != nil {
			log.Println(err)
			return nil
		}
		data = append(data, temp)
	}
	return data
}

// 一笔笔循环
func (d *Stmp) NextAny(args ...any) []any {
	getRows(d, args...)
	if d.rows.Next() {
		r, err := d.rows.SliceScan()
		if err != nil {
			log.Println(err)
			return nil
		}
		return r
	} else {
		log.Panicln("已经没有下一笔了")
		return nil
	}
}

// 上锁用,返回影响的行数
func (d *Stmp) Open(arg ...any) int64 {
	re, _ := d.stmt.Exec(arg...)
	r, err := re.RowsAffected()
	if err != nil {
		log.Println(err)
		return 0
	}
	return r
}

// 释放Prepare
func (d *Stmp) Free() error {
	return d.stmt.Close()
}

// Prepare
func (d *Stmp) Prepare(sql string) error {
	var (
		stmt *sqlx.Stmt
		err  error
	)
	d.stmt = nil
	if d.tx != nil {
		stmt, err = d.tx.Preparex(sql)
		if err != nil {
			log.Println(err)
			return err
		}
	} else {
		stmt, err = d.DB.Preparex(sql)
		if err != nil {
			log.Println(err)
			return err
		}
	}
	d.stmt = stmt
	return nil
}

// 获取rows
func getRows(d *Stmp, args ...any) *sqlx.Rows {
	var err error
	if d.rows == nil {
		d.rows, err = d.stmt.Queryx(args...)
		if err != nil {
			log.Println(err)
			return nil
		}
	}
	return d.rows
}

func clearTx(d *Stmp) {
	d.tx = nil
}
