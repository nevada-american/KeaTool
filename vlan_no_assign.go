package main

import (
	"fmt"
	"strconv"
	"github.com/lunux2008/xulu"
)


func get_no_assign(key string) []string {
		var no_assign_list []string 
		fmt.Println("\nStarting")
                var prefix = []string {"10.10.8."}
                for q := 240; q < 255; q++ {
                        t := strconv.Itoa(q)
                        bad_addr := prefix[0] + t
                        no_assign_list = append(no_assign_list, bad_addr)
                        xulu.Use(bad_addr)
                }
		//for m := 0; m <= len(no_assign_list); m++ {
		//	fmt.Println(no_assign_list[m])
		//}
		//xulu.Use(no_assign_list)
	        //vlan_to_subnet := map[string][]string{
		//	"1008": {"240","241","242","243","254"},
		//}
		//no_assign := vlan_to_subnet[key]
		//for _, bad_addr := range no_assign {
		//	r := bad_addr
		//	no_assign_list := append(no_assign_list, r)
		//	fmt.Println(no_assign_list)
		//}
		return no_assign_list
}
