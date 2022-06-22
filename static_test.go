package fgl

import (
	"testing"

	_ "github.com/sijms/go-ora/v2"
)

var (
	db = NewDB("oracle", "oracle://forewin:forewin@192.168.1.20:1521/topprod")
)

func TestSelect(t *testing.T) {

}
func TestInsert(t *testing.T) {

}

func TestUpdate(t *testing.T) {

}

func TestDelete(t *testing.T) {

}
