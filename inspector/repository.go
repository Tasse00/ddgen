package inspector

import (
	"ddgen/common"
	"ddgen/utils"
)

// 解析器
type Inspector interface {
	utils.Component
	Inspect(dbSrc string, schema string, params string) (*common.SchemaSpec, error)
}

var GlobalRendererRepository = utils.CreateRepository("inspectors")
