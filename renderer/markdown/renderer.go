package markdown

import (
	"ddgen/common"
	"ddgen/renderer"
	"fmt"
	"html/template"
	"log"
	"os"
	"strings"
)

type Renderer struct {
	rendererId string
}

const Template = `
# 数据字典

## 数据库: {{.Name}}
{{- range $ti, $table := .GetTables }}
### 表: {{$table.Name}}
{{ MakeTableTitle $table}}
{{ MakeTableInline $table }}
	{{- range $ci, $col := $table.GetColumns }}
{{ MakeTableRow $table $col }}
	{{- end }}
{{- end }}
`

func (r Renderer) GetComponentId() string {
	return r.rendererId
}

func (r Renderer) Render(ss common.SchemaSpec, outfile string, params string) error {
	tmpl := template.New("renderer")

	tmpl, err := tmpl.Funcs(template.FuncMap{
		"MakeTableTitle": func(ts common.TableSpec) string {
			return strings.Join(ts.GetDefaultSpecRenderFields(), "|")
		},

		"MakeTableInline": func(ts common.TableSpec) string {
			empty := make([]string, len(ts.GetDefaultSpecRenderFields()))
			for idx := range empty {
				empty[idx] = "---"
			}
			return strings.Join(empty, "|")
		},
		"MakeTableRow": func(ts common.TableSpec, cs common.ColumnSpec) string {
			fields := ts.GetDefaultSpecRenderFields()
			values := cs.GetSpecRenderFieldsValue(fields)

			var strValues []string
			for _, v := range values {
				nowrapVal := strings.ReplaceAll(fmt.Sprintf("%s", v), "\n", "")
				if len(nowrapVal) == 0 {
					nowrapVal = " "
				}
				strValues = append(strValues, nowrapVal)
			}
			return strings.Join(strValues, "|")
		},
	}).Parse(Template)

	if err != nil {
		return err
	}

	f, err := os.Create(outfile)
	if err != nil {
		log.Printf("create file %s failed", outfile)
		return nil
	}
	defer func() {
		err := f.Close()
		if err != nil {
			log.Println(err)
		}
	}()
	return tmpl.Execute(f, ss)
}

func init() {
	renderer.GlobalRendererRepository.Register(Renderer{rendererId: "md"})
}
