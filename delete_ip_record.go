package main

import (
        "fmt"
        "log"
	_ "github.com/go-sql-driver/mysql"
	"database/sql"

)

func delete_record(del_ip string, dbtype string, constring string) {
        fmt.Println("Deleting IP:", del_ip)
        db, err := sql.Open(dbtype, constring)
        defer db.Close()
        stmt, err := db.Prepare("DELETE FROM hosts where (ipv4_address) = (INET_ATON(?)) ;")
        res,err := stmt.Exec(del_ip)
                if err!=nil{
                        log.Fatal(err)
                }
                if (debug) {
                        log.Println(res)
                }
                db.Close()
                fmt.Println(res)
                //return res
}
