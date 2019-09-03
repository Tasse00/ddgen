package office_word

import (
	"ddgen/inspector"
	"ddgen/renderer"
	"fmt"
	"github.com/unidoc/unioffice/color"
	"github.com/unidoc/unioffice/document"
	"github.com/unidoc/unioffice/schema/soo/wml"
)

type OfficeWordRenderer struct {
	rendererId string
}

func (r OfficeWordRenderer) GetRendererId() string {
	return r.rendererId
}

func (r OfficeWordRenderer) Render(dbi *inspector.DBInspector, params string, outfile string) error {

	doc := document.New()

	para := doc.AddParagraph()
	para.SetStyle("Title")
	para.AddRun().AddText("Data Dict")

	// 表title及字段顺序准备

	var tableTitleValues = inspector.ColumnDesc{}.GetRenderLabels()
	var renderFieldsIdx = inspector.ColumnDesc{}.GetRenderFields()

	for _, schema := range dbi.Schemas {
		para := doc.AddParagraph()
		para.SetStyle("Heading1")
		para.AddRun().AddText(schema.SchemaName)

		veryLightGray := color.RGB(240, 240, 240)

		for _, tbl := range schema.Tables {
			para = doc.AddParagraph()
			para.SetStyle("Heading2")
			para.AddRun().AddText(tbl.TableName)

			table := doc.AddTable()
			table.Properties().SetWidthPercent(100)
			titleRow := table.AddRow()
			for _, title := range tableTitleValues {
				cell := titleRow.AddCell()
				cell.Properties().SetShading(wml.ST_ShdSolid, color.LightGray, color.Auto)
				cell.AddParagraph().AddRun().AddText(title)
			}

			for rIdx, cd := range tbl.Columns {
				row := table.AddRow()

				for _, value := range cd.GetRenderValues(renderFieldsIdx) {
					cell := row.AddCell()
					if rIdx%2 == 0 {
						cell.Properties().SetShading(wml.ST_ShdSolid, veryLightGray, color.Auto)
					}

					cell.AddParagraph().AddRun().AddText(fmt.Sprintln(value))
				}
			}
		}
	}
	return doc.SaveToFile(outfile)
}

func init() {
	renderer.GlobalRendererRepository.Register(OfficeWordRenderer{rendererId: "office-word"})
}
