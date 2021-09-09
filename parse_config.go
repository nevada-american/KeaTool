package main

import (
        "fmt"
        "bufio"
        "os"
        "strings"
        "log"
        "regexp"

)

func parse_config(config string) (string, string, string, string) {

        var cfglines []string
        var cfgrecord []string
        var cfgvalues []string
        var s string
        var u string
        var p string
        var connect []string
        file, err := os.Open(config)
                if err != nil {
                        log.Fatal(err)
                }

        defer file.Close()

        parser:= bufio.NewScanner(file)

        for parser.Scan() {
                matched, _ := regexp.MatchString("#.*", parser.Text())
                if (matched) {
                        continue
                }
                cfglines = append(cfglines, parser.Text())
        }
        for f := 0; f <= (len(cfglines)-1); f++ {
                cfgrecord = strings.SplitN(cfglines[f], ":", -1)
                cfgvalues = append(cfgvalues, cfgrecord[1])
        }
        s = strings.Replace(cfgvalues[0], " ", "", -1)
        u = strings.Replace(cfgvalues[1], " ", "", -1)
        p = strings.Replace(cfgvalues[2], " ", "", -1)
        if (debug) {
                fmt.Println(s,"\n",u,"\n",p,"\n")
        }
        file.Close()
        connect = []string{u,p}
        constring := strings.Join(connect,":")
        more := []string{constring,s}
        constring = strings.Join(more,"@tcp(")
        constring = (constring+")/kea?tls=true")
        //if (debug) {
        //      fmt.Println("In parse_config() function!!",constring)
        //}
        return s, u, p, constring
}
