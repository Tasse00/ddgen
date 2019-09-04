package inspector

import (
	"ddgen/common"
	"errors"
	"fmt"
	"log"
	"strings"
)

// 解析器
type Inspector interface {
	GetInspectorId() string
	Inspect(dbSrc string, schema string, params string) (*common.SchemaSpec, error)
}

type Repository struct {
	inspectors map[string]Inspector
}

func (r *Repository) Register(ins Inspector) {
	insId := strings.ToLower(ins.GetInspectorId())
	if _, ok := r.inspectors[insId]; ok {
		log.Fatalf("Inspector '%s' already registed.", insId)
		return
	}
	r.inspectors[insId] = ins
	log.Printf("detect inspector '%s'", insId)
}

func (r *Repository) Get(insId string) (*Inspector, error) {
	if inspector, ok := r.inspectors[insId]; !ok {
		return nil, errors.New(fmt.Sprintf("invalid inspectorId %s", insId))
	} else {
		return &inspector, nil
	}
}

func (r *Repository) GetInspectorIds() []string {
	var iIds []string
	for iId := range r.inspectors {
		iIds = append(iIds, iId)
	}
	return iIds
}

var GlobalRendererRepository = Repository{
	inspectors: map[string]Inspector{},
}
