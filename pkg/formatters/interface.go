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
)

type FormatWriter interface {
	Write(obj interface{}) error
}

func NewFormatter (template string, funcMaps template.FuncMap) FormatWriter {
	if viper.GetBool("verbose") {
		return &JsonFormatter{}
	}
	return NewHtmlTempleteFormatter(template, funcMaps)
}
