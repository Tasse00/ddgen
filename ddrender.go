package main

import (
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
	flag.StringVar(&renderArgs.renderType, "t", "", fmt.Sprintf("render type, one of %s", strings.Join(renderer.GlobalRendererRepository.GetRenderers(), ",")))
	flag.StringVar(&renderArgs.renderParams, "p", "", "additional params that render type need")
}

func main() {

	flag.Parse()

	if renderArgs.h {
		flag.Usage()
		return
	}

	// validate renderType
	if !utils.ContainsString(renderer.GlobalRendererRepository.GetRenderers(), renderArgs.renderType) {
		log.Fatalf("renderType must be one of %s", strings.Join(renderer.GlobalRendererRepository.GetRenderers(), ","))
		return
	}

	ddr := renderer.CreateRenderFromData(renderArgs.datFile, renderArgs.renderType, renderArgs.renderParams, renderArgs.outFile)
	ddr.Render()
}
