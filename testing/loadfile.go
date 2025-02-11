/*
Copyright 2021 The AlaudaDevops Authors.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package testing

import (
	"bufio"
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"os"

	. "github.com/onsi/gomega"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	utilyaml "k8s.io/apimachinery/pkg/util/yaml"
	"sigs.k8s.io/yaml"
)

// MustLoadFileString loads a file as string
// will panic if if failes
// ONLY FOR TEST USAGE
func MustLoadFileString(file string, content *string) {
	*content = string(MustLoadFileBytes(file))
}

// MustLoadFileBytes loads a file as []bytes
// will panic if if failes
// ONLY FOR TEST USAGE
func MustLoadFileBytes(file string) []byte {
	content, err := os.ReadFile(file)
	if err != nil {
		panic(err)
	}
	return content
}

// LoadJSON loads json
func LoadJSON(file string, obj interface{}) (err error) {
	var data []byte
	if data, err = os.ReadFile(file); err != nil {
		return
	}
	err = json.Unmarshal(data, obj)
	return
}

// MustLoadJSON loads json or panics if the parse fails.
func MustLoadJSON(file string, obj interface{}) {
	err := LoadJSON(file, obj)
	if err != nil {
		panic(fmt.Sprintf("load json file failed, file path: %s, err: %s", file, err))
	}
}

// LoadMultiYamlOrJson loads multi yamls
func LoadMultiYamlOrJson[T any](file string, list *[]T) (err error) {
	if list == nil {
		return errors.New("list should not be nil")
	}
	var data []byte
	if data, err = os.ReadFile(file); err != nil {
		return
	}
	return LoadMultiYamlOrJsonFromBytes(data, list)
}

// LoadMultiYamlOrJsonFromBytes loads multi yamls
// For historical reasons, this method still supports JSON documents separated by ---
// However, --- is not a valid separator for JSON documents.
// To be compatible with the previous handling logic, we cannot directly use the k8s built-in multiple document unmarshalling method
// and need to read line by line to implement it.
func LoadMultiYamlOrJsonFromBytes[T any](data []byte, list *[]T) (err error) {

	docs := [][]byte{}
	var currentDoc = bytes.NewBuffer(make([]byte, 0, 4096))

	reader := bufio.NewReader(bytes.NewReader(data))
	for {
		line, err := reader.ReadBytes('\n')

		if err != nil && err != io.EOF {
			return err
		}

		if isSeparator(line) {
			if currentDoc.Len() > 0 {
				docCopy := make([]byte, currentDoc.Len())
				copy(docCopy, currentDoc.Bytes())
				docs = append(docs, docCopy)
				currentDoc.Reset()
			}
		} else {
			currentDoc.Write(line)
		}

		if err == io.EOF {
			if currentDoc.Len() > 0 {
				docCopy := make([]byte, currentDoc.Len())
				copy(docCopy, currentDoc.Bytes())
				docs = append(docs, docCopy)
				currentDoc.Reset()
			}
			break
		}
	}

	for _, doc := range docs {
		if len(bytes.TrimSpace(doc)) == 0 {
			continue
		}
		obj := new(T)
		err = utilyaml.NewYAMLOrJSONDecoder(bytes.NewReader(doc), len(doc)).Decode(obj)
		if err != nil {
			return
		}

		*list = append(*list, *obj)
	}

	return nil
}

func isSeparator(line []byte) bool {
	trimmed := bytes.TrimSpace(line)

	if !bytes.HasPrefix(trimmed, []byte("---")) {
		return false
	}

	rest := bytes.TrimSpace(trimmed[3:])
	return len(rest) == 0 || rest[0] == '#'
}

// MustLoadMultiYamlOrJson loads multi yamls or panics if the parse fails.
func MustLoadMultiYamlOrJson[T any](file string, list *[]T) {
	err := LoadMultiYamlOrJson(file, list)
	if err != nil {
		panic(fmt.Sprintf("load yaml file failed, file path: %s, err: %s", file, err))
	}
}

// LoadYAML loads yaml
func LoadYAML(file string, obj interface{}) (err error) {
	var data []byte
	if data, err = os.ReadFile(file); err != nil {
		return
	}
	err = yaml.Unmarshal(data, obj)
	return
}

// MustLoadYaml loads yaml or panics if the parse fails.
func MustLoadYaml(file string, obj interface{}) {
	err := LoadYAML(file, obj)
	if err != nil {
		panic(fmt.Sprintf("load yaml file failed, file path: %s, err: %s", file, err))
	}
}

// LoadObjectOrDie loads object from yaml and returns
func LoadObjectOrDie(g *WithT, file string, obj metav1.Object, patches ...func(metav1.Object)) metav1.Object {
	g.Expect(LoadYAML(file, obj)).To(Succeed(), "could not load file into metav1.Object")
	for _, p := range patches {
		p(obj)
	}
	return obj
}

// LoadObjectReferenceOrDie loads object reference from yaml and returns
func LoadObjectReferenceOrDie(g *WithT, file string, obj *corev1.ObjectReference, patches ...func(*corev1.ObjectReference)) *corev1.ObjectReference {
	g.Expect(LoadYAML(file, obj)).To(Succeed(), "could not load file into corev1.ObjectReference")
	for _, p := range patches {
		p(obj)
	}
	return obj
}

// MustLoadReturnObjectFromYAML loads and object from yaml file and returns as interface{}
// if any loading errors happen will panic
// TO BE USED IN TESTS, DO NOT USE IN PRODUCTION CODE
func MustLoadReturnObjectFromYAML(file string, obj interface{}) interface{} {
	MustLoadYaml(file, obj)
	return obj
}
