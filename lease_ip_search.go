package main

import (
        "fmt"
        _ "github.com/go-sql-driver/mysql"
        "database/sql"
	)

func dynamic_ip_search(search_ip string, dbtype string, constring string) {
        db, err := sql.Open(dbtype, constring)
        defer db.Close()
        results, err := db.Query("SELECT INET_NTOA(address),HEX(hwaddr),subnet_id,hostname,expire FROM lease4 where INET_NTOA(address) =?", search_ip)
        fmt.Println("IP: ", search_ip)
        for results.Next() {
                        rows++
                        //var tag kea_entry
                        var tag kea_lease_entry
                        err = results.Scan(&tag.address, &tag.hwaddr, &tag.subnet_id, &tag.hostname, &tag.expire)
                        if err != nil {
                         //       panic(err.Error())
                                  fmt.Println("Got error in results.Scan")
                        }
                        fmt.Println("-----------------------------------------------------------------------------------")
                        fmt.Println("|",tag.address,"|",tag.hwaddr,"|",tag.subnet_id,"|",tag.hostname,"|", tag.expire,"|")
                        fmt.Println("-----------------------------------------------------------------------------------")
        }
        db.Close()
}
