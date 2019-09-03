package markdown

import (
	"ddgen/inspector"
	"ddgen/renderer"
	"fmt"
	"html/template"
	"os"
	"strings"
)

type Renderer struct {
	rendererId string
}

const Template = `
# 数据字典

{{- range $si, $schema := .}}
## 数据库: {{$schema.SchemaName}}
	{{- range $ti, $table := $schema.Tables }}
### 表: {{$table.TableName}}
{{ MakeTableTitle $table}}
{{ MakeTableInline $table }}
		{{- range $ci, $col := $table.Columns }}
{{ MakeTableRow $col }}
		{{- end }}
	{{- end }}
{{- end }}
`

func (r Renderer) GetRendererId() string {
	return r.rendererId
}

func (r Renderer) Render(dbi *inspector.DBInspector, params string, outfile string) error {
	tmpl := template.New("renderer")

	tmpl, err := tmpl.Funcs(template.FuncMap{
		"MakeTableTitle": func(tb inspector.TableDesc) string {
			return strings.Join(inspector.ColumnDesc{}.GetRenderLabels(), "|")
		},

		"MakeTableInline": func(tb inspector.TableDesc) string {
			empty := make([]string, len(inspector.ColumnDesc{}.GetRenderLabels()))
			for idx := range empty {
				empty[idx] = "---"
			}
			return strings.Join(empty, "|")
		},
		"MakeTableRow": func(cd inspector.ColumnDesc) string {
			fields := inspector.ColumnDesc{}.GetRenderFields()
			values := cd.GetRenderValues(fields)
			var strValues []string
			for _, v := range values {
				strValues = append(strValues, strings.ReplaceAll(fmt.Sprintf("%s", v), "\n", ""))
			}

			return strings.Join(strValues, "|")
		},
	}).Parse(Template)

	if err != nil {
		return err
	}

	f, err := os.Create(outfile)
	if err != nil {
		return nil
	}
	defer func() {
		err := f.Close()
		if err != nil {
			panic(err)
		}
	}()
	return tmpl.Execute(f, dbi.Schemas)
}

func init() {
	renderer.GlobalRendererRepository.Register(Renderer{rendererId: "md"})
}
