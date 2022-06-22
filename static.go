/*
处理静态sql
*/
package fgl

import (
	"database/sql"
	"log"
	"strconv"
	"strings"

	"github.com/jmoiron/sqlx"
	_ "github.com/sijms/go-ora/v2"
)

func NewBuild() *build {
	return &build{
		head:      "",
		unique:    false,
		cols:      make([]string, 0),
		tables:    make([]string, 0),
		joins:     make(map[string][]join),
		joinvalue: make([]any, 0),
		where:     "",
		values:    make([]any, 0),
		group:     "",
		order:     "",
		having:    "",
		session:   nil,
	}
}

func (d *build) Insert() *string {
	d.head = insert
	defer new(d)
	return insertString(d)
}
func (d *build) InsertDo(db *DataBase) int64 {
	d.head = insert
	defer new(d)
	return exec(db, insertString(d), d.values...)
}

func (d *build) Delete() *string {
	d.head = delete
	defer new(d)
	return deleteString(d)
}
func (d *build) DeleteDo(db *DataBase) int64 {
	d.head = delete
	defer new(d)
	return exec(db, deleteString(d), d.values...)
}

func (d *build) Update() *string {
	d.head = update
	defer new(d)
	return updateString(d)
}
func (d *build) UpdateDo(db *DataBase) int64 {
	d.head = update
	defer new(d)
	return exec(db, updateString(d), d.values...)
}

func (d *build) Select() *string {
	d.head = select_
	defer new(d)
	return selectString(d)
}
func (d *build) SelectDo(db *DataBase) [][]any {
	d.head = select_
	defer new(d)
	return query(db, selectString(d), d.joinvalue...)
}
func (d *build) Create() *string {
	d.head = create
	defer new(d)
	return createString(d)
}
func (d *build) CreateDo(db *DataBase) int64 {
	d.head = create
	defer new(d)
	return exec(db, createString(d), d.values...)
}

func (d *build) Index() *string {
	d.head = index
	defer new(d)
	return indexSring(d)
}

func (d *build) IndexDo(db *DataBase) int64 {
	d.head = index
	defer new(d)
	return exec(db, indexSring(d), d.values...)
}

// 处理的sql，可以调用多次
/*
	示例： ima01,im02 as ima,ima03,max(ima04)
*/
func (db *build) Col(cols ...string) *build {
	if db.cols == nil {
		db.cols = make([]string, 1)
	}
	db.cols = append(db.cols, cols...)
	return db
}

// 需要处理的表，可多次调用
/*
	示例：ima_file,imb_file as b
*/
func (db *build) From(s string) *build {
	db.tables = append(db.tables, s)
	return db
}

// 查询条件中所有值用:1,:2,:3...表示
// 可调用多个
/*
	示例： ima01 like :1 and ima02 between :2 and :3
*/
func (db *build) Where(s string) *build {
	db.where += " " + s
	return db
}

// 汇总字段
/*
	示例: ima01,ima02
*/
func (db *build) Group(s string) *build {
	db.group = s
	return db
}

// 排序字段
/*
 汇总字段
 示例: ima01,ima02
*/
func (db *build) Order(s string) *build {
	db.order = s
	return db
}

// 是否使用unique
func (db *build) Unique() *build {
	db.unique = true
	return db
}

// 传入的参数，只能调用一次，后面会覆盖前面的值
/*
	示例：[]any{1212,"12123"}
*/
func (db *build) Values(v ...any) *build {
	db.values = v
	return db
}

// 聚合条件，只能调用一次
/*
	示例：" sum(ima04) > 0 and sum(ima05) >0 "
*/
func (db *build) Having(s string) *build {
	db.having = s
	return db
}

// 左连接，可调用多次
/*
	示例："sfb_file"," sfb05 = ima01 and sfb81 >= :1 ",time.Now()
*/
func (db *build) Join(from, table, where string, v ...any) *build {
	if db.joins == nil {
		// 初始化
		db.joins = make(map[string][]join, 1)
	}
	var _, ok = db.joins[from]
	if !ok {
		db.joins[from] = make([]join, 0)
	}
	db.joins[from] = append(db.joins[from], join{
		table:  table,
		where:  where,
		values: v,
	})
	return db
}

// 是否建立中间表，不调用此函数就不是中间表。
// true 代表是会话级中间表，false代表是事务级中间表
func (db *build) Session(b bool) *build {
	*db.session = b
	return db
}

// where 条件组成 检查
// select/update/delete/join 都可能调用
func whereStatements(s *string) *string {
	return s
}

// delete 删除
func deleteString(db *build) *string {

	var s = db.head
	if len(db.tables) != 1 {
		log.Println("删除只能一次删除一个表")
		return nil
	}
	s = s + " " + db.tables[0]
	s = s + " WHERE " + *whereStatements(&db.where)

	return &s
}

// update 更新

func updateString(db *build) *string {
	var s = db.head
	if len(db.tables) != 1 {
		log.Println("更新只能一次一个表")
		return nil
	}
	s = s + " " + db.tables[0]
	s = s + " SET " + updateColStatements(db.cols)
	s = s + " WHERE " + *whereStatements(&db.where)
	return &s
}

// INSERT 插入
func insertString(db *build) *string {
	var (
		s strings.Builder
		r string
	)
	s.WriteString(db.head)
	s.WriteString(" ")
	// table
	if len(db.tables) != 1 {
		log.Println("更新只能一次删除一个表")
		return nil
	}
	s.WriteString(db.tables[0])
	// columns
	if len(db.cols) > 0 {
		// 有指定栏位
		s.WriteString("(")
		s.WriteString(colStatements(db.cols))
		s.WriteString(")")
	}
	// values,只有insert用，所以不另外新增函数
	s.WriteString("VALUES (")
	for i := 0; i < len(db.values); i++ {
		s.WriteString(":")
		s.WriteString(strconv.Itoa(i + 1))
		s.WriteString(",")
	}
	r = s.String()[:s.Len()-1] + ")"
	return &r
}

// select
func selectString(db *build) *string {
	var (
		s strings.Builder
		r string
	)
	//SELECT
	s.WriteString(db.head)
	// Col
	s.WriteString(" ")
	s.WriteString(colStatements(db.cols))
	// FROM Table
	s.WriteString(" FROM ")
	sql, v := tableStatements(db.tables, db.joins)
	s.WriteString(sql)
	db.joinvalue = v
	db.joinvalue = append(db.joinvalue, db.values...)

	// where
	if db.where != "" {
		s.WriteString(" WHERE ")
		s.WriteString(*whereStatements(&db.where))
	}
	// Group
	if db.group != "" {
		s.WriteString(" GROUP BY ")
		s.WriteString(db.group)
	}
	// Having
	if db.having != "" {
		s.WriteString(" HAVING ")
		s.WriteString(db.having)
	}
	// Order
	if db.order != "" {
		s.WriteString(" ORDER BY ")
		s.WriteString(db.order)
	}
	r = s.String()
	return &r
}

// 创建表
func createString(db *build) *string {
	var (
		r string
	)
	// Create
	r = r + db.head + "("
	//col
	r = r + colStatements(db.cols)
	r = r + ")"
	if db.session != nil {
		if *db.session {
			r = r + " ON COMMIT DELETE ROWS"
		} else {
			r = r + " ON COMMIT PRESERVE ROWS"
		}
	}
	r = r + ")"
	return &r
}

// 新建索引
func indexSring(db *build) *string {
	var (
		r string
	)
	return &r
}

// Update 的set字段拼接
func updateColStatements(col []string) string {
	var s strings.Builder
	if len(col) == 0 {
		// 所有字段都加载
		// TODO:
		/*
			SELECT column_name FROM User_Tab_Cols
			 WHERE table_name ='IMA_FILE'
			 ORDER BY column_name
		*/
		return ""
	}
	for _, v := range col {
		s.WriteString(v)
		s.WriteString("=:1,")
	}
	return s.String()[:s.Len()-1]
}

// 字段拼接
// 格式 "ima01,ima02,ima03"
func colStatements(col []string) string {
	if len(col) == 0 {
		return ""
	}
	var s strings.Builder
	for _, v := range col {
		s.WriteString(v)
		s.WriteString(",")
	}
	return s.String()[:s.Len()-1]
}

// select 表格部分
func tableStatements(tabs []string, joins map[string][]join) (string, []any) {
	if len(tabs) == 0 {
		log.Println("没有table栏位")
		return "", nil
	}
	var (
		s strings.Builder
		r []any
	)
	for _, v := range tabs {
		// 判断是否有left join
		s.WriteString(v)
		var j, ok = joins[v]
		if ok {
			// 至少有一个
			sql, v := joinStatements(j)
			s.WriteString(sql)
			if len(v) != 0 {
				r = append(r, v...)
			}
		}
		s.WriteString(",")
	}
	return s.String()[:s.Len()-1], r
}

// left join 部分
// 同时处理value 值，防止二次排序
func joinStatements(joins []join) (string, []any) {
	var (
		s strings.Builder
		r []any
	)
	for _, v := range joins {
		s.WriteString(" LEFT JOIN ")
		s.WriteString(v.table)
		s.WriteString(" ON ")
		s.WriteString(v.where)
		// chuli values 值
		r = append(r, v.values...)
	}
	return s.String(), r
}

// 执行sql，返回影响的行数
func exec(db *DataBase, s *string, arg ...any) int64 {
	var (
		re  sql.Result
		err error
	)
	if db.tx == nil {
		// 不是事务中，用DB运行
		re, err = db.DB.Exec(*s, arg...)
	} else {
		re, err = db.tx.Exec(*s, arg...)
	}
	if err != nil {
		dbErr(db, err)
		return 0
	}
	r, err := re.RowsAffected()
	if err != nil {
		log.Println(err)
		return 0
	}
	return r
}

// 查询sql，返回所有结果的any二维数组
func query(db *DataBase, sql *string, arg ...any) [][]any {
	var (
		data [][]any
		rows *sqlx.Rows
		err  error
	)
	if db.tx != nil {
		rows, err = db.tx.Queryx(*sql, arg...)
	} else {
		rows, err = db.DB.Queryx(*sql, arg...)
	}
	if err != nil {
		dbErr(db, err)
		return data
	}
	for rows.Next() {
		temp, err := rows.SliceScan()
		if err != nil {
			log.Println(err)
			return data
		}
		data = append(data, temp)
	}
	return data
}

// DB 也可直接执行静态SQL
func (db *DataBase) Insert() *string {
	return db.builder.Insert()
}
func (db *DataBase) InsertDo() int64 {
	return db.builder.InsertDo(db)
}
func (db *DataBase) Delete() *string {
	return db.builder.Delete()
}
func (db *DataBase) DeleteDo() int64 {
	return db.builder.DeleteDo(db)
}
func (db *DataBase) Update() *string {
	return db.builder.Update()
}
func (db *DataBase) UpdateDo() int64 {
	return db.builder.UpdateDo(db)
}
func (db *DataBase) Select() *string {
	return db.builder.Select()
}
func (db *DataBase) SelectDo() [][]any {
	return db.builder.SelectDo(db)
}
func (db *DataBase) Create() *string {
	return db.builder.Create()
}
func (db *DataBase) CreateDo() int64 {
	return db.builder.CreateDo(db)
}
func (db *DataBase) Index() *string {
	return db.builder.Index()
}
func (db *DataBase) IndexDo() int64 {
	return db.builder.IndexDo(db)
}
func (db *DataBase) Col(cols ...string) *DataBase {
	db.builder.Col(cols...)
	return db
}
func (db *DataBase) From(s string) *DataBase {
	db.builder.From(s)
	return db
}
func (db *DataBase) Where(s string) *DataBase {
	db.builder.Where(s)
	return db
}
func (db *DataBase) Group(s string) *DataBase {
	db.builder.Group(s)
	return db
}
func (db *DataBase) Order(s string) *DataBase {
	db.builder.Order(s)
	return db
}
func (db *DataBase) Unique() *DataBase {
	db.builder.Unique()
	return db
}
func (db *DataBase) Values(v ...any) *DataBase {
	db.builder.Values(v...)
	return db
}
func (db *DataBase) Having(s string) *DataBase {
	db.builder.Having(s)
	return db
}
func (db *DataBase) Join(from string, table string, where string, v ...any) *DataBase {
	db.builder.Join(from, table, where, v...)
	return db
}
func (db *DataBase) Session(b bool) *DataBase {
	db.builder.Session(b)
	return db
}
