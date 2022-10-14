// Copyright 2022 Liam White
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package generate

import (
	"fmt"
	"io/ioutil"
	"path"
	"regexp"

	"gopkg.in/yaml.v3"

	"google.golang.org/protobuf/compiler/protogen"
)

var configMatch = regexp.MustCompile(`\+terraform-gen:config:([^\/]+\.yaml|[^\/]+\.yml)`)

type config struct {
	InjectedFields map[string]injectedField `yaml:"injectedFields,omitempty"`
}

type injectedField struct {
	// Has to be in format types.<Type> where type is one of these constants
	// https://pkg.go.dev/github.com/hashicorp/terraform-plugin-framework@v0.14.0/types?utm_source=gopls#pkg-constants
	Type     string `yaml:"type,omitempty"`
	Required bool   `yaml:"required,omitempty"`
	Computed bool   `yaml:"computed,omitempty"`
	Optional bool   `yaml:"optional,omitempty"`
}

func loadConfig(m *protogen.Message) config {
	dir := path.Dir(m.Location.SourceFile)
	filename := getFileName(m.Comments.Leading)
	cfg := config{}
	if len(filename) > 0 {
		location := path.Join(dir, filename)
		contents, err := ioutil.ReadFile(location)
		if err != nil {
			panic(fmt.Sprintf("unable to read contents of '%s': %v", location, err))
		}

		if err := yaml.Unmarshal(contents, &cfg); err != nil {
			panic(fmt.Sprintf("unable to unmarshal contents of '%s': %v", location, err))
		}
	}
	return cfg
}

func getFileName(c protogen.Comments) string {
	match := configMatch.FindAllStringSubmatch(string(c), 1)
	if len(match) != 1 || len(match[0]) != 2 {
		return ""
	}
	return match[0][1]
}
