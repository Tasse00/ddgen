package main

import (
	"ddgen/common"
	"ddgen/renderer"
	_ "ddgen/renderer/markdown"
	_ "ddgen/renderer/office_word"
	"ddgen/utils"
	"flag"
	"fmt"
	"log"
	"strings"
)

var renderArgs = struct {
	datFile      string
	outFile      string
	renderType   string
	renderParams string
	h            bool
}{}

func init() {
	flag.BoolVar(&renderArgs.h, "h", false, "show this help")
	flag.StringVar(&renderArgs.datFile, "d", "./dat.json", "dd data filepath")
	flag.StringVar(&renderArgs.outFile, "o", "./out", "output filepath")
	flag.StringVar(&renderArgs.renderType, "t", "", fmt.Sprintf("render type, one of %s", strings.Join(renderer.GlobalRendererRepository.GetComponentIds(), ",")))
	flag.StringVar(&renderArgs.renderParams, "p", "", "additional params that render type need")
}

func main() {

	flag.Parse()

	if renderArgs.h {
		flag.Usage()
		return
	}

	// validate renderType
	if !utils.ContainsString(renderer.GlobalRendererRepository.GetComponentIds(), renderArgs.renderType) {
		log.Fatalf("renderType must be one of %s", strings.Join(renderer.GlobalRendererRepository.GetComponentIds(), ","))
		return
	}

	ss := common.SchemaSpec{}
	err := ss.LoadFromFile(renderArgs.datFile)
	if err != nil {
		log.Printf("open dat file %s failed", renderArgs.datFile)
		panic(err)
	}

	ren, err := renderer.GlobalRendererRepository.Get(renderArgs.renderType)
	if err != nil {
		panic(err)
	}

	err = ren.(renderer.Renderer).Render(ss, renderArgs.outFile, renderArgs.renderParams)
	if err != nil {
		log.Printf("render failed")
		panic(err)
	}
	log.Println("OK.")
}
