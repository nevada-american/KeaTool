package main

import (
        "fmt"
        _ "github.com/go-sql-driver/mysql"
        "database/sql"
	"os"
	"strings"

)
func get_lease_expirations(dbtype string, constring string) {
	var line []string
        fmt.Println("Connecting to database: kea, table: lease4...")
        db, err := sql.Open(dbtype, constring)
        defer db.Close()
	outFile, err := os.Create("active-leases.txt")
	defer outFile.Close()	
        //results, err := db.Query("SELECT INET_NTOA(ipv4_address), HEX(dhcp_identifier), dhcp4_subnet_id, hostname FROM lease4")
        results, err := db.Query("SELECT INET_NTOA(address),HEX(hwaddr),expire,subnet_id,hostname from lease4 order by expire asc")
        if err != nil {
                panic(err.Error())
        }
	fmt.Println("Dumping active leases to active-leases.txt...")
        for results.Next() {
                var tag kea_lease_entry
                if err := results.Scan(&tag.address, &tag.hwaddr, &tag.expire, &tag.subnet_id, &tag.hostname); err != nil {
                        panic(err)
                }
                //fmt.Println("-----------------------------------------------------------------------------------")
                //fmt.Println("|", tag.address,"|",tag.hwaddr,"|",tag.expire,"|",tag.subnet_id,"|",tag.hostname,"|")
                //fmt.Println("-----------------------------------------------------------------------------------")
		line = []string{tag.address,tag.hwaddr,tag.expire,tag.subnet_id,tag.hostname}
		linestr := strings.Join(line,",")
		outFile.WriteString(linestr)
                outFile.WriteString("\n")
                rows++
        }
        fmt.Println("Records returned: ", rows, "(are actives leases with known expirations)")
        db.Close()
}
