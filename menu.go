package main

import (
        "fmt"
        "bufio"
        "os"
        "strings"

)
func show_menu() string {
        reader := bufio.NewReader(os.Stdin)
        fmt.Println("\n\nDesert Research Institute\n\nKeatool", version, "GPLv2, 2020")
        fmt.Println("\n\n---------Keatool---------\n\n")
        fmt.Print("Below options are for reservations/static entries\n")
        fmt.Print("---------------------------------------------\n")
        fmt.Print("1) Search by IP Address\n")
        fmt.Print("2) Search by MAC Address\n")
        fmt.Print("3) Search by Host Name (not case sensitive; partial name OK)\n")
        fmt.Print("4) Search by VLAN ID[long list]\n")
        fmt.Print("5) Add a new reservation (need IP, hostname, VLAN, and MAC address)\n")
        fmt.Print("6) Search for a free IP address\n")
        fmt.Print("7) Modify an existing record -- ** NOT IMPLEMENTED YET **\n")
        fmt.Print("8) Add records via file\n")
        fmt.Print("9) Get IP subnet and mask from VLAN lookup\n")
	fmt.Println("\n")
        fmt.Print("Below options of dynamic/pool entries aka Leases (not Reservations)\n")
        fmt.Print("-------------------------------------\n")
        fmt.Print("a) Search by IP Address\n")
        fmt.Print("b) Search by MAC Address - NOT Done Yet.\n")
        fmt.Print("c) Search by Host Name - NOT Done Yet.\n")
        fmt.Print("v) Search by VLAN ID - get subnet information, etc. - NOT Done Yet.\n")
        fmt.Print("------------------------------------\n")
        fmt.Print("d) D(elete) a record by IP address\n")
	fmt.Print("e) E(xport) a list of all leases, ordered by expiration\n")
	fmt.Print("m) D(elete) a record by MAC address\n")
	fmt.Print("E) E(xport) all reservations - intensive operation, use sparingly\n")
        fmt.Print("q) Quit\n")
        fmt.Print("Enter an option: ")
        opt, _ := reader.ReadString('\n')
        opt = strings.Replace(opt, "\n", "", -1)
        return opt
}
