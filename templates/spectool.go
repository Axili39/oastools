package main

import(
	"../../pkg/oatool"
	"../{{.Package}}"
)

{{ $package := .Package }}

/* Package */
func GetObjByName(node string) interface{} {
	switch node {
	{{ range $val := .Components }}
		case "{{ $val }}":
			var obj {{ $package }}.{{ $val }}
			return &obj
	{{end}}
	}
	return nil
}

func main() {
	oatool.MainOAFileTool(GetObjByName)
}



