/**
采用DBInspector的数据进行渲染
1. html 渲染
2. markdown 渲染
*/

package renderer

import (
	"ddgen/inspector"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
)

// 一次渲染操作
type DDRender struct {
	RenderType string
	DatFile    string
	OutFile    string
	Params     string
	Inspector  *inspector.DBInspector
}

func CreateRenderFromData(datFile string, renderType string, params string, outfile string) DDRender {
	ddr := DDRender{
		RenderType: renderType,
		DatFile:    datFile,
		Inspector:  nil,
		OutFile:    outfile,
		Params:     params,
	}

	var dbi inspector.DBInspector

	bDat, err := ioutil.ReadFile(datFile)
	if err != nil {
		panic(err)
	}

	err = json.Unmarshal(bDat, &dbi)
	if err != nil {
		panic(err)
	}

	ddr.Inspector = &dbi

	return ddr
}

func (ddr *DDRender) Render() {
	renderer, ok := GlobalRendererRepository.get(ddr.RenderType)
	if !ok {
		panic(fmt.Sprintln("invalid render type", ddr.RenderType))
	}
	log.Printf("use renderer '%s'", renderer.GetRendererId())

	err := renderer.Render(ddr.Inspector, ddr.Params, ddr.OutFile)
	if err != nil {
		panic(err)
	}
	log.Println("OK.")
}
