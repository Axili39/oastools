package {{.PackageName}}
// GENERATED DO NOT EDIT

import (
	"fmt"
	"net/http"
	"github.com/Axili39/gowapi"
)
{{- range $val := .WSHandlers }}
type {{$val.Name}}Server struct {

}
{{- end }}

type MyRPCServer struct {
{{- range $val := .WSHandlers }}
	{{$val.Name}} {{$val.Name}}Server
{{- end}}
}

// General HTTP Handlers
{{- range $val := .HTTPHandlers }}
func (s* MyRPCServer) {{$val.Operation}}(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w,"reached {{$val.Operation}} url=%s ws:%v\n", r.URL.Path)
}
{{- end }}

//-------------------------------------------------------------------------------------------------
// WebSocket RPC Services
{{- range $val := .WSHandlers }}
func (s* MyRPCServer)Get{{$val.Name}}() gowapi.WSHandler {
	return &s.{{$val.Name}}
}
{{- end }}

{{- range $val := .WSHandlers }}

//-{{$val.Name}} ---------------------------------------------------------------------------
// PreUpgrade : This Operation is called by gowapi when a new connection to 
// {{$val.Path}} is established.
// we need to validate request
func (s *{{$val.Name}}Server) PreUpgrade(w http.ResponseWriter,r *http.Request) bool {
	return true
}
// PostUpgrade : After connection has been upgraded to a websocket, gowapi give us the peer 
// reference. We should build a Peer context and memorize it.
func (s *{{$val.Name}}Server) PostUpgrade(c *gowapi.Conn) gowapi.WSUser {
	cli := {{$val.Name}}CreatePeer(c,s)
	return cli
}

// RemoveUser : Called when connection is closed
func (s*{{$val.Name}}Server) RemoveUser(u gowapi.WSUser) {
	// if client has been memorized in PostUpgrade Func , delete it ex : delete(s.clients, u)
}
{{- range $op := $val.ServerInterface.Ops}}
func (s *{{$val.Name}}Server) op{{$op.Name}}(*{{$op.PbMessageName}}) {

}
{{- end }}

{{- end }}