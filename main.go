package main

import (
	"ddgen/common"
	"ddgen/inspector"
	_ "ddgen/inspector/mysql5.7"
	"ddgen/renderer"
	_ "ddgen/renderer/markdown"
	_ "ddgen/renderer/office_word"
	"ddgen/utils"
	"errors"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"github.com/urfave/cli"
	"log"
	"os"
	"sort"
	"strings"
)

func inspect(c *cli.Context) error {

	inspectorId := c.String("inspector")
	params := c.String("params")
	source := c.String("source")
	schema := c.String("schema")
	outfile := c.String("out")
	validDbDrivers := inspector.GlobalRendererRepository.GetComponentIds()
	if !utils.ContainsString(validDbDrivers, strings.ToLower(inspectorId)) {
		return errors.New(fmt.Sprintf("inspector must in %s", strings.Join(validDbDrivers, ",")))
	}

	ins, err := inspector.GlobalRendererRepository.Get(inspectorId)
	if err != nil {
		return err
	}

	log.Printf("user inspector: %s with params: %s", ins.GetComponentId(), params)
	ss, err := ins.(inspector.Inspector).Inspect(source, schema, params)
	if err != nil {
		return err
	}

	err = ss.SaveToFile(outfile)
	if err != nil {
		return err
	}
	log.Println("OK.")
	return nil
}

func render(c *cli.Context) {
	rendererId := c.String("renderer")
	datFile := c.String("dat")
	outFile := c.String("out")
	params := c.String("params")

	// validate renderType
	if !utils.ContainsString(renderer.GlobalRendererRepository.GetComponentIds(), rendererId) {
		log.Fatalf("renderType must be one of %s", strings.Join(renderer.GlobalRendererRepository.GetComponentIds(), ","))
		return
	}

	ss := common.SchemaSpec{}
	err := ss.LoadFromFile(datFile)
	if err != nil {
		log.Printf("open dat file %s failed", datFile)
		panic(err)
	}

	ren, err := renderer.GlobalRendererRepository.Get(rendererId)
	if err != nil {
		panic(err)
	}

	err = ren.(renderer.Renderer).Render(ss, outFile, params)
	if err != nil {
		log.Printf("render failed")
		panic(err)
	}
	log.Println("OK.")
}

func export(c *cli.Context) error {
	rendererId := c.String("renderer")
	renderParams := c.String("renderParams")
	inspectorId := c.String("inspector")
	inspectParams := c.String("inspectParams")
	source := c.String("source")
	schema := c.String("schema")
	outfile := c.String("out")

	validDbDrivers := inspector.GlobalRendererRepository.GetComponentIds()
	if !utils.ContainsString(validDbDrivers, strings.ToLower(inspectorId)) {
		return errors.New(fmt.Sprintf("inspector must in %s", strings.Join(validDbDrivers, ",")))
	}

	ins, err := inspector.GlobalRendererRepository.Get(inspectorId)
	if err != nil {
		return err
	}

	log.Printf("user inspector: %s with params: %s", ins.GetComponentId(), inspectParams)
	ss, err := ins.(inspector.Inspector).Inspect(source, schema, inspectParams)
	if err != nil {
		return err
	}

	// render
	// validate renderType
	if !utils.ContainsString(renderer.GlobalRendererRepository.GetComponentIds(), rendererId) {
		return errors.New(fmt.Sprintf("renderType must be one of %s", strings.Join(renderer.GlobalRendererRepository.GetComponentIds(), ",")))
	}

	ren, err := renderer.GlobalRendererRepository.Get(rendererId)
	if err != nil {
		panic(err)
	}

	err = ren.(renderer.Renderer).Render(*ss, outfile, renderParams)
	if err != nil {
		log.Printf("render failed")
		panic(err)
	}
	log.Println("OK.")
	return nil
}

func main() {

	app := cli.NewApp()
	app.Name = "DDGen"
	app.Usage = "inspect database and export it"

	app.Commands = []cli.Command{
		{
			Name:  "export",
			Usage: "inspect database and export it",
			Flags: []cli.Flag{
				cli.StringFlag{
					Name:  "inspector, i",
					Usage: "the inspector of database",
				},
				cli.StringFlag{
					Name:  "source, s",
					Usage: "database source",
				},

				cli.StringFlag{
					Name:  "schema",
					Usage: "the schema to export",
				},

				cli.StringFlag{
					Name:  "renderer, r",
					Usage: "schema spec renderer",
				},

				cli.StringFlag{
					Name:  "inspectParams, ip",
					Usage: "parameters passed to inspector",
				},

				cli.StringFlag{
					Name:  "renderParams, rp",
					Usage: "parameters passed to renderer",
				},
				cli.StringFlag{
					Name:     "out, o",
					Usage:    "output file",
					Required: true,
				},
			},
			Action: export,
		},
		{
			Name:  "inspect",
			Usage: "only export schema spec data",
			Flags: []cli.Flag{
				cli.StringFlag{
					Name:     "inspector, i",
					Usage:    "the inspector of database",
					Required: true,
				},
				cli.StringFlag{
					Name:     "source, s",
					Usage:    "database source",
					Required: true,
				},

				cli.StringFlag{
					Name:     "schema",
					Usage:    "the schema to export",
					Required: true,
				},
				cli.StringFlag{
					Name:  "out, o",
					Usage: "output data filepath",
					Value: "dat.json",
				},
				cli.StringFlag{
					Name:  "params, p",
					Usage: "parameters passed to inspector",
				},
			},

			Action: inspect,
		},
		{
			Name:  "render",
			Usage: "render with schema spec data",
			Flags: []cli.Flag{
				cli.StringFlag{
					Name:     "dat, d",
					Usage:    "schema spec data file",
					Required: true,
				},
				cli.StringFlag{
					Name:     "renderer, r",
					Usage:    "schema spec renderer",
					Required: true,
				},
				cli.StringFlag{
					Name:  "params, p",
					Usage: "parameters passed to renderer",
				},
				cli.StringFlag{
					Name:  "out, o",
					Usage: "output data filepath",
					Value: "dat.json",
				},
			},
			Action: render,
		},
	}

	sort.Sort(cli.FlagsByName(app.Flags))
	sort.Sort(cli.CommandsByName(app.Commands))

	err := app.Run(os.Args)
	if err != nil {
		log.Fatal(err)
	}
}
