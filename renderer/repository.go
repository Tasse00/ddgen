package renderer

import (
	"ddgen/inspector"
	"log"
)

// 渲染器
type Renderer interface {
	GetRendererId() string
	Render(dbi *inspector.DBInspector, outfile string, params string) error
}

type Repository struct {
	renders map[string]Renderer
}

var GlobalRendererRepository = Repository{
	renders: map[string]Renderer{},
}

// register 用于renderer注册自身
func (rr *Repository) Register(re Renderer) {

	if _, ok := rr.renders[re.GetRendererId()]; ok {
		log.Fatalf("Renderer '%s' already registed.", re.GetRendererId())
		return
	}
	rr.renders[re.GetRendererId()] = re

	log.Printf("detect renderer '%s'", re.GetRendererId())
}

func (rr *Repository) get(rendererId string) (Renderer, bool) {
	if renderer, ok := rr.renders[rendererId]; !ok {
		return nil, false
	} else {
		return renderer, true
	}
}

func (rr *Repository) GetRenderers() []string {
	ids := make([]string, len(rr.renders))
	idx := 0
	for k := range rr.renders {
		ids[idx] = k
		idx++
	}
	return ids
}
