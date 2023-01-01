package main

import (
	"bytes"
	"embed"
	"flag"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"text/template"

	"google.golang.org/protobuf/compiler/protogen"
)

//go:embed fs
var FS embed.FS

var Version = "development"

var (
	implPath *string
	version  *string
)

func main() {
	var flags flag.FlagSet

	implPath = flags.String("impl", "", "Go import path to the handlers implementation")
	version = flags.String("version", "", "Semver version of the service")

	protogen.Options{
		ParamFunc: flags.Set,
	}.Run(func(gen *protogen.Plugin) error {
		if *implPath == "" {
			return fmt.Errorf("please set impl option")
		}

		for _, f := range gen.Files {
			if !f.Generate {
				continue
			}

			generateService(gen, f)
		}

		return nil
	})
}

func generatePromStats(gen *protogen.Plugin, file *protogen.File) {
	filename := fmt.Sprintf("%s_nmfw_stats.pb.go", file.GeneratedFilenamePrefix)
	g := gen.NewGeneratedFile(filename, file.GoImportPath)

	buf := bytes.NewBuffer([]byte{})
	t := template.New("prometheus.go.templ")
	f, err := FS.ReadFile("fs/prometheus.go.templ")
	if err != nil {
		panic(fmt.Sprintf("could not read template: %v", err))
	}

	p, err := t.Parse(string(f))
	if err != nil {
		panic(fmt.Sprintf("could not parse template: %v", err))
	}

	err = p.Execute(buf, map[string]any{
		"file": file,
	})
	if err != nil {
		panic(fmt.Sprintf("could not execute template: %v", err))
	}

	g.P(buf.String())
}

func generateService(gen *protogen.Plugin, file *protogen.File) {
	generatePromStats(gen, file)

	for _, s := range file.Services {
		filename := fmt.Sprintf("%s_nmfw_svc_%s.pb.go", file.GeneratedFilenamePrefix, strings.ToLower(s.GoName))
		g := gen.NewGeneratedFile(filename, file.GoImportPath)

		buf := bytes.NewBuffer([]byte{})
		t := template.New("service.go.templ")
		f, err := FS.ReadFile("fs/service.go.templ")
		if err != nil {
			panic(fmt.Sprintf("could not read template: %v", err))
		}

		t.Funcs(map[string]any{
			"toLower":          strings.ToLower,
			"QualifiedGoIdent": g.QualifiedGoIdent,
		})

		p, err := t.Parse(string(f))
		if err != nil {
			panic(fmt.Sprintf("could not parse template: %v", err))
		}

		err = p.Execute(buf, map[string]any{
			"file":             file,
			"service":          s,
			"version":          version,
			"generatorVersion": Version,
		})
		if err != nil {
			panic(fmt.Sprintf("could not execute template: %v", err))
		}

		g.P(buf.String())

		generateServiceCommand(gen, file, s)

		for _, m := range s.Methods {
			fmt.Fprintf(os.Stderr, "Implement %s.%sHandler of type %s.%sHandler\n", *implPath, m.GoName, strings.ReplaceAll(file.GoImportPath.String(), `"`, ``), m.GoName)
		}
	}
}

func generateServiceCommand(gen *protogen.Plugin, file *protogen.File, svc *protogen.Service) {
	ln := strings.ToLower(svc.GoName)
	g := gen.NewGeneratedFile(filepath.Join(filepath.Dir(file.GeneratedFilenamePrefix), ln, fmt.Sprintf("%s.go", ln)), "main")

	buf := bytes.NewBuffer([]byte{})
	t := template.New("cmd.go.templ")
	f, err := FS.ReadFile("fs/cmd.go.templ")
	if err != nil {
		panic(fmt.Sprintf("could not read template: %v", err))
	}

	t.Funcs(map[string]any{
		"toLower": strings.ToLower,
	})

	p, err := t.Parse(string(f))
	if err != nil {
		panic(fmt.Sprintf("could not parse template: %v", err))
	}

	err = p.Execute(buf, map[string]any{
		"file":             file,
		"service":          svc,
		"implPath":         *implPath,
		"version":          version,
		"generatorVersion": Version,
	})
	if err != nil {
		panic(fmt.Sprintf("could not execute template: %v", err))
	}

	g.P(buf.String())
}
