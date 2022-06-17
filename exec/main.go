package main

import _ "github.com/darcyjoven/fgl"

func main() {
	// db := fgl.NewDB("oracle://forewin:forewin@192.168.1.20:1521/topprod")
	// db.Prepare("select ima01 from ima_file where ima01 = :1")
	/*
		db := fgl.NewDB("xxx")
		updima := db.Begin()
		ima,err :=updima.Prepare("select * from ima_file where ima01 like :1")
		if err!=nil{
			log.Println(err)
			updima.Rollback()
		}
		imb:=updima.Prepare("xxx")
		_,err:=imb.Exec(12)
		if err!=nil{
			log.Println(err)
			updima.Rollback()
		}
		_,err:=imb.Exec(13)
		if err!=nil{
			log.Println(err)
			updima.Rollback()
		}
		updaima.Commit()

		intoa :=db.Prepare("xx")
		rows,err:=intoa.DB.Queryx("xx")
		for rows.Next(){
			temp:=rows.SliceScan()
		}
	*/
}
