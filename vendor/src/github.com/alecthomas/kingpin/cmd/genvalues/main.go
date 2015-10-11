package main

import (
	"encoding/json"
	"os/exec"
	"strings"
	"text/template"

	"os"
)

const (
	tmpl = `package kingpin

// This file is autogenerated by "go generate .". Do not modify.

{{range .}}
{{if not .NoValueParser}}
// -- {{.Type}} Value
type {{.Type}}Value {{.Type}}

func new{{.|Name}}Value(p *{{.Type}}) *{{.Type}}Value {
	return (*{{.Type}}Value)(p)
}

func (f *{{.Type}}Value) Set(s string) error {
	v, err := {{.Parser}}
	*f = {{.Type}}Value(v)
	return err
}

func (f *{{.Type}}Value) Get() interface{} { return {{.Type}}(*f) }

func (f *{{.Type}}Value) String() string { return {{.|Format}} }

// {{.|Name}} parses the next command-line value as {{.Type}}.
func (p *parserMixin) {{.|Name}}() (target *{{.Type}}) {
	target = new({{.Type}})
	p.{{.|Name}}Var(target)
	return
}

func (p *parserMixin) {{.|Name}}Var(target *{{.Type}}) {
	p.SetValue(new{{.|Name}}Value(target))
}

{{end}}
// {{.|Plural}} accumulates {{.Type}} values into a slice.
func (p *parserMixin) {{.|Plural}}() (target *[]{{.Type}}) {
	target = new([]{{.Type}})
	p.{{.|Plural}}Var(target)
	return
}

func (p *parserMixin) {{.|Plural}}Var(target *[]{{.Type}}) {
	p.SetValue(newAccumulator(target, func(v interface{}) Value { return new{{.|Name}}Value(v.(*{{.Type}})) }))
}

{{end}}
`
)

type Value struct {
	Name          string `json:"name"`
	NoValueParser bool   `json:"no_value_parser"`
	Type          string `json:"type"`
	Parser        string `json:"parser"`
	Format        string `json:"format"`
	Plural        string `json:"plural"`
}

func fatalIfError(err error) {
	if err != nil {
		panic(err)
	}
}

func main() {
	r, err := os.Open("values.json")
	fatalIfError(err)
	defer r.Close()

	v := []Value{}
	err = json.NewDecoder(r).Decode(&v)
	fatalIfError(err)

	valueName := func(v *Value) string {
		if v.Name != "" {
			return v.Name
		}
		return strings.Title(v.Type)
	}

	t, err := template.New("genvalues").Funcs(template.FuncMap{
		"Lower": strings.ToLower,
		"Format": func(v *Value) string {
			if v.Format != "" {
				return v.Format
			}
			return "fmt.Sprintf(\"%v\", *f)"
		},
		"Name": valueName,
		"Plural": func(v *Value) string {
			if v.Plural != "" {
				return v.Plural
			}
			return valueName(v) + "List"
		},
	}).Parse(tmpl)
	fatalIfError(err)

	w, err := os.Create("values_generated.go")
	fatalIfError(err)
	defer w.Close()

	err = t.Execute(w, v)
	fatalIfError(err)

	err = exec.Command("goimports", "-w", "values_generated.go").Run()
	fatalIfError(err)
}
