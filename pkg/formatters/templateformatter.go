/*******************************************************************************
 * Copyright 2020 VMWare.
 *
 * Licensed under the Apache License, Version 2.0 (the "License"); you may not use this file except
 * in compliance with the License. You may obtain a copy of the License at
 *
 * http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software distributed under the License
 * is distributed on an "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express
 * or implied. See the License for the specific language governing permissions and limitations under
 * the License.
 *******************************************************************************/
package formatters

import (
	"github.com/spf13/viper"
	"html/template"
	"io"
	"text/tabwriter"
)

type HtmlTempleteFormatter struct {
	Name     string
	Format   string
	FuncMaps template.FuncMap
}

func NewHtmlTempleteFormatter(name string, format string, funcMaps template.FuncMap) *HtmlTempleteFormatter {
	return &HtmlTempleteFormatter{
		Name:     name,
		Format:   format,
		FuncMaps: funcMaps,
	}
}

func (f *HtmlTempleteFormatter) Write(obj interface{}) (err error) {
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

