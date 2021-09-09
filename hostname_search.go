package main

import (
        "fmt"
        _ "github.com/go-sql-driver/mysql"
        "database/sql"

)

func hname_search(hostname string, dbtype string, constring string) {
        db, err := sql.Open(dbtype, constring)
        defer db.Close()
        results, err := db.Query("SELECT INET_NTOA(ipv4_address), HEX(dhcp_identifier), dhcp4_subnet_id, hostname FROM hosts where hostname like ?", hostname)
        for results.Next() {
                        rows++
                        var tag kea_entry
                        err = results.Scan(&tag.ipntoa, &tag.dehexmac, &tag.hname, &tag.dhcp_id)
                                if err != nil {
                                        panic(err.Error())

                                }
                        fmt.Println("----------------------------------------------------------------")
                        fmt.Println("|",tag.ipntoa,"|",tag.dehexmac,"|", tag.hname,"|",tag.dhcp_id,"|")
                        fmt.Println("----------------------------------------------------------------")
                        if (debug) {
                                fmt.Println("Rows returned: ", rows)
                        }
                }
        db.Close()
}
