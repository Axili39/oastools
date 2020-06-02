
package oasmodel

type XwsRpcOperation struct {
	Name string `yaml:"name"`
	Schema *SchemaOrRef `yaml:"schema"`
}

type XwsRpcService struct {
	Server struct {
		Name string		`yaml:"name"`
		OpenPath string	`yaml:"openPath"` //deprecated
		Interface []XwsRpcOperation	`yaml:"interface,omitempty"`
	}	`yaml:"server"`
	Client struct {
		Name string		`yaml:"name"`
		Interface []XwsRpcOperation	`yaml:"interface,omitempty"`
	}
}
