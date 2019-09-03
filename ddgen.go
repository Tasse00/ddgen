package main

import (
	"ddgen/inspector"
	"ddgen/utils"
	"encoding/json"
	"flag"
	_ "github.com/go-sql-driver/mysql"
	"io/ioutil"
	"log"
	"strings"
)

var genArgs = struct {
	dbDri   string
	dbSrc   string
	outFile string
	schemas string
	h       bool
}{}

func init() {
	flag.BoolVar(&genArgs.h, "h", false, "show this help")
	flag.StringVar(&genArgs.dbDri, "D", "mysql", "database driver")
	flag.StringVar(&genArgs.dbSrc, "S", "", "database source")
	flag.StringVar(&genArgs.schemas, "s", "", "schemas to export dd")
	flag.StringVar(&genArgs.outFile, "o", "dat.json", "output data file (json)")
}

func main() {

	flag.Parse()

	if genArgs.h {
		flag.Usage()
		return
	}

	validDbDrivers := []string{"mysql"}
	if !utils.ContainsString(validDbDrivers, strings.ToLower(genArgs.dbDri)) {
		log.Fatalf("db drivers must in %s", strings.Join(validDbDrivers, ","))
		return
	}

	//dbSrc := "root:root@tcp(127.0.0.1:32768)/information_schema?charset=utf8"

	dbi := inspector.CreateDBInspector(genArgs.dbDri, genArgs.dbSrc)

	// 排除schema
	dbi.SchemasOnly = strings.Split(genArgs.schemas, ",")

	dbi.Initialize()
	defer dbi.Destroy()

	dbi.InspectSchemas()

	b, err := json.Marshal(dbi)

	if err != nil {
		panic(err)
	}

	err = ioutil.WriteFile(genArgs.outFile, b, 0644)
	if err != nil {
		panic(err)
	}

}
