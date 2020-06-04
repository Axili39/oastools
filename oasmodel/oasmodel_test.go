package oasmodel

import (
	"io/ioutil"
	"os"
	"testing"

	"gopkg.in/yaml.v3"
)

func check(data string, t *testing.T) {
	var oa OpenAPI
	err := yaml.Unmarshal([]byte(data), &oa)
	if err != nil {
		t.Errorf("error unmarshalling : %v", err)
	}

	buf, err2 := yaml.Marshal(&oa)
	if err2 != nil {
		t.Errorf("Marshal: %v", err)
	}

	if string(buf) != data {
		t.Errorf("UnMarshal -> MArshal differs :\n%s\n\n", buf)
		// Helper : save the file

		err := ioutil.WriteFile("/tmp/dat1.yaml", buf, 0644)
		if err != nil {
			t.Errorf("error saving tmp file : %v \n", err)
		}

	}

}
func TestResponseWithRef(t *testing.T) {
	data :=
		`openapi: 3.0.0
info:
    title: test simple
    version: 1.0.0
paths:
    /test1:
        get:
            responses:
                "200":
                    $ref: '#/a/b/c'
`
	check(data, t)
}
func TestAllRequired(t *testing.T) {
	data :=
		`openapi: 3.0.0
info:
    title: test simple
    version: 1.0.0
paths:
    /test1:
        get:
            responses:
                "200":
                    description: pet response
                    content:
                        '*/*':
                            schema:
                                type: string
`
	check(data, t)
}

func TestAssets(t *testing.T) {
	var files []os.FileInfo
	root := "./assets/"
	files, err := ioutil.ReadDir(root)
	if err != nil {
		t.Errorf("error listings assets : %v", err)
		return
	}
	for i := range files {
		if !files[i].IsDir() {
			t.Logf("Checking : %s", files[i].Name())
			yamlFile, err := ioutil.ReadFile(root + files[i].Name())
			if err != nil {
				t.Errorf("error loading assets : %s", files[i].Name())
				return
			}
			check(string(yamlFile), t)
		}
	}

}
