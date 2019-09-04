package main

import (
	"ddgen/inspector"
	_ "ddgen/inspector/mysql5.7"
	"ddgen/utils"
	"flag"
	_ "github.com/go-sql-driver/mysql"
	"log"
	"strings"
)

var genArgs = struct {
	insId   string
	dbSrc   string
	outFile string
	schema  string
	h       bool
	params  string
}{}

func init() {
	flag.BoolVar(&genArgs.h, "h", false, "show this help")
	flag.StringVar(&genArgs.insId, "i", "mysql5.7", "inspector")
	flag.StringVar(&genArgs.dbSrc, "S", "", "database source")
	flag.StringVar(&genArgs.schema, "s", "", "schema to export dd")
	flag.StringVar(&genArgs.outFile, "o", "dat.json", "output data file (json)")
	flag.StringVar(&genArgs.params, "p", "", "params pass to inspector.")
}

func main() {

	flag.Parse()

	if genArgs.h {
		flag.Usage()
		return
	}

	validDbDrivers := inspector.GlobalRendererRepository.GetComponentIds()
	if !utils.ContainsString(validDbDrivers, strings.ToLower(genArgs.insId)) {
		log.Fatalf("inspector must in %s", strings.Join(validDbDrivers, ","))
		return
	}

	ins, err := inspector.GlobalRendererRepository.Get(genArgs.insId)
	if err != nil {
		panic(err)
	}

	log.Printf("user inspector: %s with params: %s", ins.GetComponentId(), genArgs.params)
	ss, err := ins.(inspector.Inspector).Inspect(genArgs.dbSrc, genArgs.schema, genArgs.params)
	if err != nil {
		panic(err)
	}

	err = ss.SaveToFile(genArgs.outFile)
	if err != nil {
		panic(err)
	}
}
