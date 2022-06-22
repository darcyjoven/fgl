package fgl

import (
	"log"
	"testing"

	_ "github.com/sijms/go-ora/v2"
)

func TestForeach(t *testing.T) {
	var db = NewDB("oracle", "oracle://forewin:forewin@192.168.1.20:1521/topprod")
	st := db.Prepare("select ima01,ima02 from ima_file where rownum < :1")
	log.Println(st.FetchAny(10))
	log.Println(st.FetchAny(20))
}
