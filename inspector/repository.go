package inspector

import (
	"github.com/Tasse00/ddgen/common"
	"github.com/Tasse00/ddgen/utils"
)

// 解析器
type Inspector interface {
	utils.Component
	Inspect(dbSrc string, schema string, params string) (*common.SchemaSpec, error)
}

var GlobalRendererRepository = utils.CreateRepository("inspectors")
