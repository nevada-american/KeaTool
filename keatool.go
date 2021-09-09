// keatool v1
// Ed Mitchell (ed.mitchell@dri.edu) 2020

// To Do:
// Modify Record functionality
// Set a limit to avoid handing out .240-.254 on subnets (reserved for network equipment and firewalls)
// Newly Done:
// Added the go-getoptions library and implemented basic functionality for the debug option.
package main

import (
	//"strconv"
	"fmt"
	"bufio"
	"os"
	"strings"
	"bytes"
	"encoding/binary"
	"net"
	"database/sql"
	"sort"
	"log"
	_ "github.com/go-sql-driver/mysql"
	"github.com/DavidGamba/go-getoptions"
	//"regexp"

)


type kea_entry struct {
	ipntoa string
	dehexmac string
	host_id string
	dhcp_id string
	dtype string
	vlan string
	vlan6 string
	ipaddr string
	hname string
	d4class string
	d6class string
	d4next string
	d4sname string
	d4bfn string
	ucontext string
	authkey string
	}

type kea_lease_entry struct {
        address string
        hwaddr string
        expire string
        subnet_id string
        hostname string
        }

// declare our version variable first //

var version string = "1.0"

// declare some variables //

var vlan string
var hostname string
var ip string
var mac string
var rows int
var opt string
var debug bool
var verbose bool
var import_file string

// adding for config file //

var config_file string = "ktool.cfg"


// worker functions //

// IP and subnet functions //

func Hosts(cidr string) ([]string, error) {
        ip, ipnet, err := net.ParseCIDR(cidr)
        if err != nil {
                return nil, err
        }

        var ips []string
	var reserved int = 16
        for ip := ip.Mask(ipnet.Mask); ipnet.Contains(ip); inc(ip) {
                ips = append(ips, ip.String())
        }
        // remove network address and broadcast address
	// also, remove .254-.240
        return ips[1 : len(ips)-reserved], nil
}

func inc(ip net.IP) {
        for j := len(ip) - 1; j >= 0; j-- {
                ip[j]++
                if ip[j] > 0 {
                        break
                }
        }
}

func nextfree(ip net.IP) (net.IP) {
        for j := len(ip) - 1; j >= 0; j-- {
                ip[j]++
                if ip[j] > 0 {
                        break
                }
        }
	return ip
}


func ip2Long(ip string) uint32 {
         var long uint32
         binary.Read(bytes.NewBuffer(net.ParseIP(ip).To4()), binary.BigEndian, &long)
         return long
}


func mac_search(mac string, dbtype string, constring string) {
	db, err := sql.Open(dbtype, constring)
        defer db.Close()
	fmt.Println("MAC search string is", mac)
	size := len(mac)
	fmt.Println("Length of search string is", size, "characters.")
        //results, err := db.Query("SELECT INET_NTOA(ipv4_address), HEX(dhcp_identifier), dhcp4_subnet_id, hostname FROM hosts where HEX(dhcp_identifier) =?", mac)
	results, err := db.Query("SELECT INET_NTOA(ipv4_address), HEX(dhcp_identifier), dhcp4_subnet_id, hostname FROM hosts where HEX(dhcp_identifier) =?", mac)
        for results.Next() {
                        rows++
                        var tag kea_entry
                        err = results.Scan(&tag.ipntoa, &tag.dehexmac, &tag.hname, &tag.dhcp_id)
                                if err != nil {
                                        panic(err.Error())

                                }
			fmt.Println("---------------------------------------------------------------")	
                        fmt.Println("|",tag.ipntoa,"|",tag.dehexmac,"|",tag.hname,"|",tag.dhcp_id,"|")
			fmt.Println("---------------------------------------------------------------")
			if (debug) {
                       		fmt.Println("Rows returned: ",rows)
			}
                }

	db.Close()
}

func dynamic_vlan_search(vlan string, allips []string, dbtype string, constring string) {
	var a int
	freeIPs := make([]net.IP, 0)
	foundIPs := make([]string, 0)
	byteIPs := make([]net.IP, 0)
	allbyteIPs := make([]net.IP, 0)
	db, err := sql.Open(dbtype, constring)
        defer db.Close()
	if (debug) {
        	fmt.Println("VLAN search string is", vlan)
	}
	// FIX this query
	// SELECT INET_NTOA(address), HEX(hwaddr), subnet_id, hostname FROM lease4 where subnet_id = 2514
        results, err := db.Query("SELECT INET_NTOA(address), HEX(hwaddr), subnet_id, hostname FROM lease4 where subnet_id =?", vlan)
        for results.Next() {
                        rows++
                        var tag kea_entry
                        err = results.Scan(&tag.ipntoa, &tag.dehexmac, &tag.hname, &tag.dhcp_id)
                                if err != nil {
                                        panic(err.Error())

                                }
			if (debug) {
				fmt.Println("---------------------------------------------------------------")
                        	fmt.Println("|",tag.ipntoa,"|",tag.dehexmac,"|",tag.hname,"|",tag.dhcp_id,"|")
				fmt.Println("---------------------------------------------------------------")
			}
			foundIPs = append(foundIPs, tag.ipntoa)
			byteIPs = append(byteIPs, net.ParseIP(tag.ipntoa))
	 }
	 if (debug) {
         	a = len(foundIPs)
	 	fmt.Println("Found", a, "IPs in use.")
	 }

      	 sort.Slice(byteIPs, func(i, j int) bool {
		return bytes.Compare(byteIPs[i], byteIPs[j])  < 0
	 })
	 if (debug) {
	 	fmt.Println("IP address(es) in VLAN", vlan, "sorted in ascending order:")
	 	fmt.Println("------------------------------------------------------")
	 	for _, z := range byteIPs {
	  	 	fmt.Printf("%s\n", z)
	 	}
		fmt.Println("------------------------------------------------------")
	 }
	 for m := 0; m <= (len(allips) -1); m++ {
		allbyteIPs = append(allbyteIPs, net.ParseIP(allips[m]))
	 }
	 sort.Slice(allbyteIPs, func(b, d int) bool {
		return bytes.Compare(allbyteIPs[b], allbyteIPs[d]) < 0
	 })
	 for x := 0; x <= (len(byteIPs)-1); x++ {
		//for y := 0; y <= (len(allbyteIPs)-1); y++ {
		for y := 0; y <= x; y++ {
			compresult := bytes.Equal(allbyteIPs[y], byteIPs[x])
			if (!compresult) {
				//fmt.Println("This IP is NOT in use:", allbyteIPs[y])
				freeIPs = append(freeIPs, allbyteIPs[y])
			}
		}
	 }
	 total := cap(allbyteIPs)
	 if (debug) {
	 	fmt.Println(total)
	 }
	 q := len(byteIPs)
	 if ( q == 0 ) {
		 fmt.Println("There are no entries for VLAN", vlan)
		 os.Exit(99)
	 }
	 //fmt.Println("\n\nFirst free IP is:", byteIPs[q - 1])
	 ip := nextfree(byteIPs[q - 1])
	 subnet := get_subnet(vlan)
	 fmt.Println("The next free IP for subnet",subnet,"is:",ip)
	 db.Close()


}

func ip_search(ip string, dbtype string, constring string) {
	//db, err := sql.Open("mysql", "keaadmin:password@tcp(wizard-dev.dri.edu)/kea")
	db, err := sql.Open(dbtype, constring)
	defer db.Close()
	results, err := db.Query("SELECT INET_NTOA(ipv4_address), HEX(dhcp_identifier), dhcp4_subnet_id, hostname FROM hosts where INET_NTOA(ipv4_address) =?", ip)
	for results.Next() {
			rows++
			var tag kea_entry
			err = results.Scan(&tag.ipntoa, &tag.dehexmac, &tag.hname, &tag.dhcp_id)
			if err != nil {
				panic(err.Error())
			}
                        fmt.Println("---------------------------------------------------------------")
			fmt.Println("|",tag.ipntoa,"|",tag.dehexmac,"|",tag.hname,"|",tag.dhcp_id,"|")
			fmt.Println("---------------------------------------------------------------")
	}
	db.Close()
}

func insert_record(ip string, vlan string, hostname string, mac string, dbtype string, constring string) {
	var dhcpidtype int = 0
	// swap in 'err' for '_' below //
	var doit string = "y"
	if (debug) {
		var ans string = "y"
		fmt.Println(ans)
	}
	insReader := bufio.NewReader(os.Stdin)
	
	db, _ := sql.Open(dbtype, constring)
        defer db.Close()
	fmt.Println("Inserting record with: IP: ", ip, "\nVLAN: ", vlan, "\nMAC address: ", mac, "\nHostname: ", hostname, "\nDHCP ID Type: ", dhcpidtype)
	fmt.Print("Is this what you want?")
	ans, _ :=insReader.ReadString('\n')
	ans = strings.Replace(ans, "\n", "", -1)
	strings.ToLower(ans)
	ans_result := ans == doit 
	if ans_result {
		fmt.Println("Record accepted; adding to Kea hosts table.")
		stmt, err := db.Prepare("INSERT INTO hosts (dhcp_identifier,dhcp_identifier_type,dhcp4_subnet_id,ipv4_address,hostname) VALUES (UNHEX(?),?,?,INET_ATON(?),?);")
		defer db.Close()
	        if err != nil {
        		panic(err.Error())
   		}
		if (debug) {
			fmt.Println("Using this value for MAC address:", mac)
		}
		res,err := stmt.Exec(mac, dhcpidtype, vlan, ip, hostname)
		if err!=nil{
 			log.Fatal(err)
 		}
		if (debug) {
			log.Println(res)	
		}
		db.Close()
	} else {
	    fmt.Println("User has rejected the proposed record; exiting.")
	    os.Exit(99)
	}
}

//func insert_bulk_record(ip string, vlan string, hostname string, mac string, dbtype string, constring string) {
//	var dhcpidtype int = 0
//        //insReader := bufio.NewReader(os.Stdin)
//        //db, _ := sql.Open("mysql", "keaadmin:password@tcp(wizard-dev.dri.edu)/kea")
//	fmt.Println("Connecting to database", dbtype)
//	db, err := sql.Open(dbtype, constring)
//        defer db.Close()
//        fmt.Println("Inserting record with: IP: ", ip, "\nVLAN: ", vlan, "\nMAC address: ", mac, "\nHostname: ", hostname, "\nDHCP ID Type: ", dhcpidtype)
//        stmt, err := db.Prepare("INSERT INTO hosts (dhcp_identifier,dhcp_identifier_type,dhcp4_subnet_id,ipv4_address,hostname) VALUES (UNHEX(?),?,?,INET_ATON(?),?);")
//        defer db.Close()
//        if err != nil {
//                panic(err.Error())
//        }
//        fmt.Println("Using this value for MAC address:", mac)
//        res,err := stmt.Exec(mac, dhcpidtype, vlan, ip, hostname)
//        if err!=nil{
//                log.Fatal(err)
//        }
//        log.Println(res)
//        //db.Close()
//        
//}

func main() {
        var s string
        var u string
        var p string
	var constring string
	var dbtype string = ("mysql")
	flagged := getoptions.New()
        flagged.Bool("d", false)
	flagged.Bool("v", false)
	flagged.Parse(os.Args[1:])
	if flagged.Called("d") {
		debug = true
		fmt.Println("Debug flag set; increasing output level.")
	}
	if flagged.Called("v") {
		verbose = true
		fmt.Println("Verbose flag set; using verbose output.")
	}
	var allips []string
	var txtlines []string
	// Note: splitrecord will, with SplitN, split our records on the "," delimiter, on a per-line basis.
	// Then the splitrecord []string slice will hold in indices 0-3 all the data we need for a SQL insert.
	// I hope.
	var splitrecord []string
	if (verbose) {
        	fmt.Println("Connecting to kea database; table hosts\n")
	}
	// Ask the User For Stuff and Things //
	reader := bufio.NewReader(os.Stdin)
	//fmt.Println("\n\nDesert Research Institute\n\nKeatool", version, "GPLv2, 2020")
	
	// adding for config file //

        s, u, p, constring = parse_config(config_file)
        if (debug) {
        	fmt.Println("In parse_config return!!",constring)
                fmt.Println(config_file)
		fmt.Println(s, u, p)
        }

	opt := show_menu()

	// Start the switch statement here //
	switch opt {

	case "1":
		fmt.Print("Enter an IP address: ")
		ip, _ := reader.ReadString('\n')
		ip = strings.Replace(ip, "\n", "", -1)
		ip_search(ip, dbtype, constring)
	case "2":
		fmt.Print("Enter a MAC address: ")
		mac, _ := reader.ReadString('\n')
		mac = strings.Replace(mac, "\n", "", -1)
		// Adding to handle ":" - why can't strings.Replace handle multiple characters?  //
		//mac = strings.Replace(mac, ":", "", -1)
		mac_search(mac, dbtype, constring)
	case "3":
		fmt.Print("Enter a hostname: ")
		hostname, _ := reader.ReadString('\n')
		hostname = strings.Replace(hostname, "\n", "%", -1)
		if (debug) {
			fmt.Println(hostname, dbtype, constring)
		}
		hname_search(hostname, dbtype, constring)
	case "4":
		fmt.Print("Enter a VLAN ID: ")
		vlan, _ := reader.ReadString('\n')
		vlan = strings.Replace(vlan, "\n", "", -1)
		subnet := get_subnet(vlan)
		allips,_ = Hosts(subnet)
		//fmt.Println(allips)
		if (debug) {
			fmt.Println("Subnet for VLAN",vlan,"is",subnet)
			//fmt.Println(allips)
		} else {
			fmt.Println("Working with VLAN",vlan)
		}
		vlan_search(vlan, allips, dbtype, constring)
	case "5":
		fmt.Print("Enter IP address: ")
		ip, _ := reader.ReadString('\n')
		ip = strings.Replace(ip, "\n", "", -1)
		//longIP := ip2Long(ip)
		fmt.Print("Enter MAC address: ")
		mac, _ := reader.ReadString('\n')
		mac = strings.Replace(mac, "\n", "", -1)
		mac = strings.Replace(mac, ":" ,"", -1)
		macbyte := []byte(mac)
		if (debug) {
			//fmt.Println("IP-as-long:", longIP)//
			fmt.Println("MAC address as bytes:",macbyte)
		}
		fmt.Print("Enter VLAN ID: ")
		vlan, _ := reader.ReadString('\n')
		vlan = strings.Replace(vlan, "\n", "", -1)
		fmt.Print("Enter hostname: ")
		hostname, _ := reader.ReadString('\n')
		hostname = strings.Replace(hostname, "\n", "", -1)
		insert_record(ip, vlan, hostname, mac, dbtype, constring)
	case "6":
		fmt.Print("Enter a subnet to search with mask (a.b.c.d/nn):")
		searchnet, _ := reader.ReadString('\n')
		searchnet = strings.Replace(searchnet, "\n", "", -1)
		fmt.Println(searchnet)
		hosts, _ := Hosts(searchnet)
		//fmt.Println(len(hosts))
		fmt.Printf("%T\n", hosts)
	case "7":
		fmt.Print("Enter a hostname to find and modify a record:")
		mHostname, _ := reader.ReadString('\n')
		mHostname = strings.Replace(mHostname, "\n", "", -1)
		if (debug) {
                        fmt.Println(mHostname)
                }
		os.Exit(0)
	case "8":
		fmt.Print("Enter a filename to import records:")
		import_file,_  = reader.ReadString('\n')
		import_file = strings.Replace(import_file, "\n", "", -1)
		file, err := os.Open(import_file)
		if err != nil {
			log.Fatalf("failed opening file: %s", err)
		}
		scanner := bufio.NewScanner(file)
		scanner.Split(bufio.ScanLines)
		for scanner.Scan() {
			txtlines = append(txtlines, scanner.Text())
		}
		file.Close()
		if (debug) {
			for _, one_line := range txtlines {
				fmt.Println(one_line)
			}
		}
		for f := 0; f <= (len(txtlines)-1); f++ {
			splitrecord = strings.SplitN(txtlines[f], ",", -1)
			splitrecord[1] = strings.Replace(splitrecord[1], ":" ,"", -1)
			insert_bulk_record(splitrecord[3], splitrecord[2], splitrecord[0], splitrecord[1], dbtype, constring)
			}
		os.Exit(0)
	case "9":
                fmt.Print("Enter a VLAN ID: ")
                vlan, _ := reader.ReadString('\n')
                vlan = strings.Replace(vlan, "\n", "", -1)
                subnet := get_subnet(vlan)
		_, ipv4Net, _ := net.ParseCIDR(subnet)
		mask := ipv4Net.Mask
                fmt.Println("Subnet for VLAN",vlan,"is",subnet)
		fmt.Print("Decimal-format subnet mask is: ")
		fmt.Println(fmt.Sprintf("%d.%d.%d.%d", mask[0], mask[1], mask[2], mask[3]))

    	case "a":
        	fmt.Print("Enter an IP address: ")
        	search_ip, _ := reader.ReadString('\n')
        	search_ip = strings.Replace(search_ip, "\n", "", -1)
        	fmt.Println("IP to find:", search_ip)
		dynamic_ip_search(search_ip, dbtype, constring)
        	//os.Exit(1)

	case "d":
		fmt.Print("Enter an IP address: ")
		del_ip, _ := reader.ReadString('\n')
		del_ip = strings.Replace(del_ip, "\n", "", -1)
		fmt.Println("IP to delete:", del_ip)
		delete_record(del_ip, dbtype, constring)

	case "e":
		fmt.Print("Pulling all leases and ordering by expiration(ascending).\n")
		get_lease_expirations(dbtype, constring)

	case "m":
		fmt.Print("Enter a MAC address [xx:xx:xx:xx:xx:xx]: ")
		del_mac,_ := reader.ReadString('\n')
		del_mac = strings.Replace(del_mac,"\n", "", -1)
		fmt.Println("MAC to delete:", del_mac)
		//delete_by_mac(del_mac, dbtype, constring)
		os.Exit(1)
	
	case "E":
		//fmt.Print("Exporting all reservation data...")
		export_all(dbtype, constring)

	case "v":
    		fmt.Print("Enter a VLAN ID: ")
       		//fmt.Print("Enter a VLAN ID: ")
        	vlan, _ := reader.ReadString('\n')
        	vlan = strings.Replace(vlan, "\n", "", -1)
        	subnet := get_subnet(vlan)
        	allips,_ = Hosts(subnet)
        	dynamic_vlan_search(vlan, allips, dbtype, constring)

	case "z":
		fmt.Print("Undocumented Test Routine.  Plug it in via keatool.go.")
		no_assign_list := get_no_assign("1008")		
		//var prefix = []string {"10.10.8."}
                for q := 240; q < 255; q++ {
			//t := strconv.Itoa(q)
			//bad_addr := prefix[0] + t 
                        //fmt.Println(bad_addr)
                        //no_assign_list = append(no_assign_list, bad_addr)
                        //xulu.Use(no_assign_list)
			//fmt.Println(no_assign_list)
                }
		for qq := range no_assign_list {
			fmt.Println(no_assign_list[qq])
		}
		//fmt.Println(no_assign_list)

	case "q":
		os.Exit(1)

	default:
		os.Exit(99)
	}

	rows = 0 // init the rows value //


}
