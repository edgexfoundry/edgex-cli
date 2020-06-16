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
	"encoding/json"
	"fmt"
	"github.com/spf13/viper"
	"io"
)

type JsonFormatter struct{}

func (f *JsonFormatter) Write(obj interface{}) (err error) {
	json, err := json.Marshal(obj)
	if err != nil {
		return err
	}
	pw := viper.Get("writer").(io.WriteCloser)
	defer pw.Close()
	_, err = fmt.Fprintf(pw, "%s\n", json)
	return
}
