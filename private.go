package fgl

// 静态sql配置清空
func new(d *build) {
	d.head = ""
	d.unique = false
	d.cols = make([]string, 0)
	d.tables = make([]string, 0)
	d.joins = make(map[string][]join)
	d.joinvalue = make([]any, 0)
	d.where = ""
	d.values = make([]any, 0)
	d.group = ""
	d.order = ""
	d.having = ""
	d.session = nil
}
