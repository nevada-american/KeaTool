package main

import (
        "fmt"
	_ "github.com/go-sql-driver/mysql"
	"database/sql"
	"os"
	//"encoding/csv"
	"strings"
)

func export_all(dbtype string, constring string) {
	var line []string
	//var line2 string
	// set up the CSV file //
	//file, err := os.Create("export.csv")
	f2, err := os.Create("export_raw.csv")
	//checkError("Cannot create file", err)
	//defer file.Close()
	defer f2.Close()
	// set up the handle for writing //
	//writer := csv.NewWriter(file)
	//defer writer.Flush()
	// prep and connect to database //
        db, err := sql.Open(dbtype, constring)
        defer db.Close()
        //results, err := db.Query("SELECT * FROM hosts")
	results, err := db.Query("SELECT INET_NTOA(ipv4_address), HEX(dhcp_identifier), dhcp4_subnet_id, hostname FROM hosts")
	fmt.Println("Dumping all reservations.")
        for results.Next() {
                        rows++
                        var tag kea_entry
                        err = results.Scan(&tag.ipntoa, &tag.dehexmac, &tag.dhcp_id, &tag.hname)
                                if err != nil {
                                        panic(err.Error())

                                }
                        fmt.Println("----------------------------------------------------------------")
                        fmt.Println("|",tag.ipntoa,"|",tag.dehexmac,"|", tag.hname,"|",tag.dhcp_id,"|")
                        fmt.Println("----------------------------------------------------------------")
			line = []string{tag.ipntoa,tag.dehexmac,tag.hname,tag.dhcp_id}
        		linestr := strings.Join(line,",")
			//writer.Write(line)
			f2.WriteString(linestr)
			f2.WriteString("\n")
			//fmt.Println(linestr)
                }

	if (debug) {
		fmt.Println("Records returned: ", rows)
	}
        db.Close()
	//file.Close()
	f2.Close()
	f2.Close()
}
