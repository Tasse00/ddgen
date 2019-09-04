package office_word

import (
	"fmt"
	"github.com/Tasse00/ddgen/common"
	"github.com/Tasse00/ddgen/renderer"
	"github.com/unidoc/unioffice/color"
	"github.com/unidoc/unioffice/document"
	"github.com/unidoc/unioffice/schema/soo/wml"
	"log"
)

type OfficeWordRenderer struct {
	rendererId string
}

func (r OfficeWordRenderer) GetComponentId() string {
	return r.rendererId
}

func (r OfficeWordRenderer) Render(ss common.SchemaSpec, outfile string, params string) error {

	doc := document.New()

	para := doc.AddParagraph()
	para.SetStyle("Title")
	para.AddRun().AddText("Data Dict")

	// 表title及字段顺序准备

	para = doc.AddParagraph()
	para.SetStyle("Heading1")
	para.AddRun().AddText(ss.Name)

	veryLightGray := color.RGB(240, 240, 240)

	for _, tbl := range ss.GetTables() {

		var tableTitleValues = tbl.GetDefaultSpecRenderFields()
		var columnRenderFields = tbl.GetDefaultSpecRenderFields()

		para = doc.AddParagraph()
		para.SetStyle("Heading2")
		para.AddRun().AddText(tbl.Name)

		table := doc.AddTable()
		table.Properties().SetWidthPercent(100)
		titleRow := table.AddRow()
		for _, title := range tableTitleValues {
			cell := titleRow.AddCell()
			cell.Properties().Margins().SetLeft(2)
			cell.Properties().Margins().SetBottom(2)
			cell.Properties().Margins().SetRight(2)
			cell.Properties().Margins().SetTop(2)
			cell.Properties().Borders().SetAll(wml.ST_BorderSingle, color.Gray, 1)
			cell.Properties().SetShading(wml.ST_ShdSolid, color.LightGray, color.Auto)
			r := cell.AddParagraph().AddRun()

			r.Properties().SetBold(true)
			r.Properties().SetSize(10)
			r.AddText(title)
		}

		for rIdx, cd := range tbl.GetColumns() {
			row := table.AddRow()
			for _, value := range cd.GetSpecRenderFieldsValue(columnRenderFields) {
				cell := row.AddCell()
				cell.Properties().Margins().SetLeft(2)
				cell.Properties().Margins().SetBottom(2)
				cell.Properties().Margins().SetRight(2)
				cell.Properties().Margins().SetTop(2)
				cell.Properties().Borders().SetAll(wml.ST_BorderSingle, color.Gray, 1)
				if rIdx%2 == 0 {
					cell.Properties().SetShading(wml.ST_ShdSolid, veryLightGray, color.Auto)
				}

				r := cell.AddParagraph().AddRun()
				r.Properties().SetSize(10)
				r.AddText(fmt.Sprintln(value))
			}
		}
	}
	log.Printf("save to file %s", outfile)
	return doc.SaveToFile(outfile)
}

func init() {
	ren := OfficeWordRenderer{rendererId: "office-word"}
	renderer.GlobalRendererRepository.Register(ren)
}
