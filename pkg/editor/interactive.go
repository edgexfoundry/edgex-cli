/*******************************************************************************
 * Copyright 2020 Dell Inc.
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

// Package editor contains functionality which allows the user to provide information in the editor of their choice.
package editor

import (
	"bytes"
	"encoding/json"
	"html/template"
	"io/ioutil"
	"os"
	"os/exec"
	"runtime"
)

const (
	UnixDefaultEditor       = "vi"
	WindowsDefaultEditor    = "notepad"
	EnvironmentVariableName = "EDITOR"
	InteractiveModeLabel    = "interactive mode"
)

// OpenFileInEditor opens filename in a text editor.
func OpenFileInEditor(filename string) error {
	editor := os.Getenv(EnvironmentVariableName)
	if editor == "" {

		if runtime.GOOS == "windows" {
			editor = WindowsDefaultEditor
		} else {
			editor = UnixDefaultEditor
		}
	}

	editorExecutable, err := exec.LookPath(editor)
	if err != nil {
		return err
	}

	cmd := exec.Command(editorExecutable, filename)
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	return cmd.Run()
}

// CaptureInputFromEditor opens a temporary file in a text editor and returns
// the written bytes on success or an error on failure. It handles deletion
// of the temporary file behind the scenes.
func CaptureInputFromEditor(template []byte) ([]byte, error) {
	file, err := ioutil.TempFile(os.TempDir(), "*")
	if err != nil {
		return []byte{}, err
	}

	filename := file.Name()
	defer func() { _ = os.Remove(filename) }()

	_, err = file.Write(template)
	if err != nil {
		return nil, err
	}

	if err = file.Close(); err != nil {
		return []byte{}, err
	}

	if err = OpenFileInEditor(filename); err != nil {
		return []byte{}, err
	}

	updatedBytes, err := ioutil.ReadFile(filename)
	if err != nil {
		return []byte{}, err
	}

	return updatedBytes, nil
}

// OpenInteractiveEditor opens users default editor populated with a JSON representation of a structure
func OpenInteractiveEditor(o interface{}, temp string, funcMap template.FuncMap) ([]byte, error) {
	t := template.New("Template")
	if funcMap != nil {
		t.Funcs(funcMap)
	}
	dsJsonTemplate, err := t.Parse(temp)
	if err != nil {
		return nil, err
	}
	buff := bytes.NewBuffer([]byte{})
	err = dsJsonTemplate.Execute(buff, o)
	if err != nil {
		return nil, err
	}

	return CaptureInputFromEditor(buff.Bytes())
}

// isLastElementOfSlice is a function which is used in HTML templates to determine the last element in a slice
func IsLastElementOfSlice(index int, lenght int) bool {
	return index == lenght-1
}

// EscapeHTML is a function which marshal `s` applying Intent and returns template.HTML. This makes hmtl.Template to
// escape HTML codes
func EscapeHTML(s interface{}) (template.HTML, error) {
	data, err := json.MarshalIndent(s, "", "    ")
	if err != nil {
		return "", err
	}
	return template.HTML(data), nil
}
