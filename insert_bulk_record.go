package main

import (
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"database/sql"
	"log"
)

func insert_bulk_record(ip string, vlan string, hostname string, mac string, dbtype string, constring string) {
        var dhcpidtype int = 0
        fmt.Println("Connecting to database", dbtype)
        db, err := sql.Open(dbtype, constring)
        defer db.Close()
        fmt.Println("Inserting record with: IP: ", ip, "\nVLAN: ", vlan, "\nMAC address: ", mac, "\nHostname: ", hostname, "\nDHCP ID Type: ", dhcpidtype)
        stmt, err := db.Prepare("INSERT INTO hosts (dhcp_identifier,dhcp_identifier_type,dhcp4_subnet_id,ipv4_address,hostname) VALUES (UNHEX(?),?,?,INET_ATON(?),?);")
        defer db.Close()
        if err != nil {
                panic(err.Error())
        }
        fmt.Println("Using this value for MAC address:", mac)
        res,err := stmt.Exec(mac, dhcpidtype, vlan, ip, hostname)
        if err!=nil{
                log.Fatal(err)
        }
        log.Println(res)
        //db.Close()

}
