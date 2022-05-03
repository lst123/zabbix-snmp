package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"
	"time"

	g "github.com/gosnmp/gosnmp"
)

func main() {
	flag.Usage = func() {
		fmt.Printf("Usage:\n")
		fmt.Printf("   %s [-community=<community>] host [oid]\n", filepath.Base(os.Args[0]))
		fmt.Printf("     host      - the host to walk/scan\n")
		fmt.Printf("     oid       - the MIB/Oid defining a subtree of values\n\n")
		flag.PrintDefaults()
	}

	var community string
	flag.StringVar(&community, "community", "public", "the community string for device")
	flag.Parse()

	if len(flag.Args()) < 1 {
		flag.Usage()
		os.Exit(1)
	}

	target := flag.Args()[0]
	// var oid string
	oid := "1.3.6.1.2.1.2.2.1.3"
	if len(flag.Args()) > 1 {
		oid = flag.Args()[1]
	}

	g.Default.Target = target
	g.Default.Community = community
	g.Default.Timeout = time.Duration(10 * time.Second)
	err := g.Default.Connect()
	if err != nil {
		fmt.Printf("Connect err: %v\n", err)
		os.Exit(1)
	}
	defer g.Default.Conn.Close()

	m := map[string]string{"ifname": "1.3.6.1.2.1.31.1.1.1.1",
		"ifalias":      "1.3.6.1.2.1.31.1.1.1.18",
		"ifoperstatus": "1.3.6.1.2.1.2.2.1.8"}

	// Get ifTypes
	results, err := g.Default.BulkWalkAll(oid)
	if err != nil {
		log.Fatalf("Get() err: %v", err)
		os.Exit(1)
	}
	ifTypes := []string{}
	for _, variable := range results {
		o := strings.Split(variable.Name, ".")
		oidId := o[len(o)-1]
		if variable.Value == 6 {
			ifTypes = append(ifTypes, oidId)
		}
	}
	for _, iface := range ifTypes {
		oids := []string{}
		for _, v := range m {
			newOid := v + "." + iface
			oids = append(oids, newOid)
		}
		fmt.Println(oids)
	}
	// fmt.Println(ifTypes)
}

// func getSNMPValues(x g, oids []string) {
// 	result, err := x.Get(oids)
// 	if err != nil {
// 		log.Fatalf("Get() err: %v", err)
// 	}
// 	fmt.Println(result)
// }
