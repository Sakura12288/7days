package RPC

import (
	"fmt"
	"html/template"
	"net/http"
)

const debugText = `<html>
	<body>
	<title>GeeRPC Services</title>
	{{range .}}
	<hr>
	Service {{.Name}}
	<hr>
		<table>
		<th align=center>Method</th><th align=center>Calls</th>
		{{range $name, $mtype := .Method}}
			<tr>
			<td align=left font=fixed>{{$name}}({{$mtype.ArgType}}, {{$mtype.ReplyType}}) error</td>
			<td align=center>{{$mtype.NumCalls}}</td>
			</tr>
		{{end}}
		</table>
	{{end}}
	</body>
	</html>`

var debug = template.Must(template.New("rpc debug").Parse(debugText))

type DebugServer struct {
	*Server
}

type debugService struct {
	Name   string
	Method map[string]*methodType
}

func (d DebugServer) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	var services []debugService
	d.services.Range(func(name, value interface{}) bool {
		svc := value.(*service)
		services = append(services, debugService{
			Name:   name.(string),
			Method: svc.method,
		})
		return true
	})
	err := debug.Execute(w, services)
	if err != nil {
		_, _ = fmt.Fprintln(w, "解析有问题"+err.Error())
		return
	}
}
