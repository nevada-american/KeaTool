package main

import (
)

func get_offlimits_addr(key string) string {
	        vlan_to_prefix := map[string]string{
			"1560": "10.15.61.",
			"1008": "10.10.8.",
			"1900": "10.19.1.",
			"1800": "10.18.1.",
			"1732": "10.17.33.",
			"1730": "10.17.31.",
			"1722": "10.17.23.",
			"1720": "10.17.21.",
			"1710": "10.17.11.",
			"1700": "10.17.1.",
			"1604": "10.16.5.",
			"1602": "10.16.3.",
			"1600": "10.16.1.",
			"1552": "10.15.53.",
			"1550": "10.15.51.",
			"1540": "10.15.41.",
			"1530": "10.15.31.",
			"1520": "10.15.21.",
			"1516": "10.15.17.",
			"1514": "10.15.15.",
			"1510": "10.15.11.",
			"1504": "10.15.5.",
			"1502": "10.15.3.",
			"1500": "10.15.1.",
			"1100": "10.11.15.",
			"2008": "10.20.15.",
			"2050": "10.20.50.",
			"2090": "10.20.90.",
			"2100": "10.21.7.",
			"2400": "10.24.7.",
			"2408": "10.24.15.",
			"2500": "10.25.1.",
			"2510": "10.25.11.",
			"2514": "10.25.15.",
			"2516": "10.25.17.",
			"2520": "10.25.21.",
			"2530": "10.25.31.",
			"2540": "10.25.41.",
			"2550": "10.25.51.",
			"2552": "10.25.52.",
			"2560": "10.25.61.",
			"2600": "10.26.1.",
			"2602": "10.26.3.",
			"2700": "10.27.1.",
			"2710": "10.27.11.",
			"2720": "10.27.21.",
			"2722": "10.27.23.",
			"2730": "10.27.31.",
			"2800": "10.28.1.",
			"2900": "10.29.1.",
			"2910": "10.29.11.",
		}
		prefix := vlan_to_prefix[key]
		return prefix 
}
