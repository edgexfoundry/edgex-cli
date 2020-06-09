package formatter

import (
	"html/template"
	"io"
	"text/tabwriter"
	
	"github.com/spf13/viper"
)

type FormatWriter interface {
	Write(obj interface{}) error
}
//TODO Add VerboseFormatter instead
//EmptyFomatter does not write anything. Usually used when --verbose flag is present.
type EmptyFormatter struct{}

func (f *EmptyFormatter) Write(obj interface{}) (err error) {
	return nil
}

type Formatter struct {
	Name     string
	Format   string
	FuncMaps template.FuncMap
}

func New(name string, format string, funcMaps template.FuncMap) *Formatter {
	return &Formatter{
		Name:     name,
		Format:   format,
		FuncMaps: funcMaps,
	}
}

func (f *Formatter) Write(obj interface{}) (err error) {
	pw := viper.Get("writer").(io.WriteCloser)
	w := new(tabwriter.Writer)
	w.Init(pw, 0, 8, 1, '\t', 0)
	tmpl := template.New(f.Name)
	if f.FuncMaps != nil {
		tmpl.Funcs(f.FuncMaps)
	}
	tmpl, err = tmpl.Parse(f.Format)
	if err != nil {
		return err
	}
	err = tmpl.Execute(w, obj)
	if err != nil {
		return err
	}
	w.Flush()
	return nil
}
