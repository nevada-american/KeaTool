package main

import (
        "fmt"
	"bytes"
	"os"
	"net"
	"sort"
        _ "github.com/go-sql-driver/mysql"
        "database/sql"

)

func vlan_search(vlan string, allips []string, dbtype string, constring string) {
	var a int
	//freeIPs := make([]net.IP, 0)
	foundIPs := make([]string, 0)
	byteIPs := make([]net.IP, 0)
	allbyteIPs := make([]net.IP, 0)
	//db, err := sql.Open("mysql", "keaadmin:password@tcp(wizard-dev.dri.edu)/kea")
	db, err := sql.Open(dbtype, constring)
        defer db.Close()
	if (debug) {
        	fmt.Println("VLAN search string is", vlan)
	}
        results, err := db.Query("SELECT INET_NTOA(ipv4_address), HEX(dhcp_identifier), dhcp4_subnet_id, hostname FROM hosts where dhcp4_subnet_id =?", vlan)
        for results.Next() {
                        rows++
                        var tag kea_entry
                        err = results.Scan(&tag.ipntoa, &tag.dehexmac, &tag.hname, &tag.dhcp_id)
                                if err != nil {
                                        panic(err.Error())

                                }
			//if (debug) {
				fmt.Println("---------------------------------------------------------------")
                        	fmt.Println("|",tag.ipntoa,"|",tag.dehexmac,"|",tag.hname,"|",tag.dhcp_id,"|")
				fmt.Println("---------------------------------------------------------------")
			//}
			foundIPs = append(foundIPs, tag.ipntoa)
			byteIPs = append(byteIPs, net.ParseIP(tag.ipntoa))
			//if (verbose) {
		//		for x :=0; x <= (len(foundIPs)-1 ); x++ {
		 //       		fmt.Println(foundIPs[x])
		//		}
		//	}
	 }
	 if (debug) {
         	a = len(foundIPs)
	 	fmt.Println("Found", a, "IPs in use.")
	 }

        //if (verbose) {
        //      fmt.Println("----------------------------------------------------")
        //       fmt.Println("Rows returned: ",rows)
        //		for x := 0; x <= (len(foundIPs) -1); x++ {
	//		fmt.Println(foundIPs[x])
	//	}
	// }

      	 sort.Slice(byteIPs, func(i, j int) bool {
		return bytes.Compare(byteIPs[i], byteIPs[j])  < 0
	 })
//	 if (debug) {
//	 	fmt.Println("IP address(es) in VLAN", vlan, "sorted in ascending order:")
//	 	fmt.Println("------------------------------------------------------")
//	 	for _, z := range byteIPs {
//	  	 	fmt.Printf("%s\n", z)
//	 	}
//		fmt.Println("------------------------------------------------------")
//	 }
	 // NOTE NOTE NOTE
	 // I think in here is where I would exclude the upper range of a subnet as reserved,
	 // per discussion with Ryan.
	 // Chew on this for a day and then try something!
	 for m := 0; m <= (len(allips) -1); m++ {
		allbyteIPs = append(allbyteIPs, net.ParseIP(allips[m]))
		//fmt.Println(len(allbyteIPs))
	 }
	 sort.Slice(allbyteIPs, func(b, d int) bool {
		return bytes.Compare(allbyteIPs[b], allbyteIPs[d]) < 0
	 })
	 for x := 0; x <= (len(byteIPs)-1); x++ {
		//for y := 0; y <= (len(allbyteIPs)-1); y++ {
		for y := 0; y <= x; y++ {
			compresult := bytes.Equal(allbyteIPs[y], byteIPs[x])
			// debugging stuff //
			//if (debug) {
				//fmt.Println("allByte and byte slice lengths follow:")
				//fmt.Println(len(allbyteIPs))
				//fmt.Println(byteIPs)
			//}
			if compresult {
			//	fmt.Println("This IP is in use:", byteIPs[x])
		 	}
			//if (!compresult) {
			//	fmt.Println("This IP is NOT in use:", allbyteIPs[y])
			//	freeIPs = append(freeIPs, allbyteIPs[y])
			//}
		}
	 }
	 total := cap(allbyteIPs)
	 if (debug) {
	 	fmt.Println(total)
	 }
	 q := len(byteIPs)
	 fmt.Println("The following IP addresses are IN USE:")
	 for z := 0; z < q; z++ {
		fmt.Println(byteIPs[z])
	 }
	 if ( q == 0 ) {
		 fmt.Println("There are no entries for VLAN", vlan)
		 os.Exit(99)
	 }
	 //fmt.Println("\n\nFirst free IP is:", byteIPs[q - 1])
	 ip := nextfree(byteIPs[q - 1])
	 subnet := get_subnet(vlan)
	 fmt.Println("The next free IP for subnet",subnet,"is:",ip)
	 fmt.Println("NOTE NOTE NOTE NOTE NOTE NOTE!!!!")
	 fmt.Println("DO NOT ASSIGN the top of the subnet, from .240-254!!!")
	 fmt.Println("EXCEPTION: VLAN 1008 - DO NOT ASSIGN 10.10.8.252 or 10.10.15.240-254")
	 db.Close()


}
