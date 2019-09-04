package renderer

import (
	"github.com/Tasse00/ddgen/common"
	"github.com/Tasse00/ddgen/utils"
)

// 渲染器
type Renderer interface {
	utils.Component
	Render(ss common.SchemaSpec, outfile string, params string) error
}

var GlobalRendererRepository = utils.CreateRepository("renderers")
