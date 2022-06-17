package fgl

import (
	"fmt"
	"strings"
	"testing"

	"github.com/jmoiron/sqlx"
	_ "github.com/sijms/go-ora/v2"
)

func addPlus(s1 *string, s2 string) string {
	return *s1 + " " + s2
}

func addBuilder(s1 *string, s2 string) string {
	var str strings.Builder
	str.WriteString(*s1)
	str.WriteString(" ")
	str.WriteString(s2)
	return str.String()
}

func prepare() {
	var db, err = sqlx.Connect("oracle", "oracle://forewin:forewin@192.168.1.20:1521/topprod")
	defer db.Close()
	if err != nil {
		fmt.Println(err)
		return
	}
	for i := 0; i <= 1000; i++ {
		rows, err := db.Queryx("select * from ima_file where rownum<:1", 10)
		if err != nil {
			fmt.Println(err)
			return
		}
		for rows.Next() {
			_, err := rows.SliceScan()
			if err != nil {
				fmt.Println(err)
				return
			}
		}
	}

}

func noPrepare() {
	var db, err = sqlx.Connect("oracle", "oracle://forewin:forewin@192.168.1.20:1521/topprod")
	defer db.Close()
	if err != nil {
		fmt.Println(err)
		return
	}
	st, err := db.Preparex("select * from ima_file where rownum<:1")
	if err != nil {
		fmt.Println(err)
		return
	}
	for i := 0; i <= 1000; i++ {
		rows, err := st.Queryx(10)
		if err != nil {
			fmt.Println(err)
			return
		}
		for rows.Next() {
			_, err := rows.SliceScan()
			if err != nil {
				fmt.Println(err)
				return
			}
		}
	}

}

func BenchmarkPlus(b *testing.B) {
	var s = "where"
	for i := 0; i < b.N; i++ {
		addPlus(&s, "and")
	}
}
func BenchmarkBuilder(b *testing.B) {
	var s = "where"
	for i := 0; i < b.N; i++ {
		addBuilder(&s, "and")
	}
}

func BenchmarkPrepare(b *testing.B) {
	for i := 0; i < b.N; i++ {
		prepare()
	}
}
func BenchmarkNoPrepare(b *testing.B) {
	for i := 0; i < b.N; i++ {
		noPrepare()
	}
}
