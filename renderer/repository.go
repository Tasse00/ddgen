package renderer

import (
	"ddgen/common"
	"ddgen/utils"
)

// 渲染器
type Renderer interface {
	utils.Component
	Render(ss common.SchemaSpec, outfile string, params string) error
}

var GlobalRendererRepository = utils.CreateRepository("renderers")
