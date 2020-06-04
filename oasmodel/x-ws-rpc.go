package oasmodel

type XwsRPCOperation struct {
	Name   string       `yaml:"name"`
	Schema *SchemaOrRef `yaml:"schema"`
}

type XwsRPCService struct {
	Server struct {
		Name      string            `yaml:"name"`
		OpenPath  string            `yaml:"openPath"` //deprecated
		Interface []XwsRPCOperation `yaml:"interface,omitempty"`
	} `yaml:"server"`
	Client struct {
		Name      string            `yaml:"name"`
		Interface []XwsRPCOperation `yaml:"interface,omitempty"`
	}
}
