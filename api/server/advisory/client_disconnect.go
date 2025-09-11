package advisory

import (
	"github.com/nats-io/jsm.go/api/event"
)

// DisconnectEventMsgV1 is sent when a new connection previously defined from a
// ConnectEventMsg is closed.
//
// NATS Schema Type io.nats.server.advisory.v1.client_disconnect
type DisconnectEventMsgV1 struct {
	event.NATSEvent

	Server   ServerInfoV1 `json:"server"`
	Client   ClientInfoV1 `json:"client"`
	Sent     DataStatsV1  `json:"sent"`
	Received DataStatsV1  `json:"received"`
	Reason   string       `json:"reason"`
}

func init() {
	err := event.RegisterTextCompactTemplate("io.nats.server.advisory.v1.client_disconnect", `{{ .Time | ShortTime }} [Disconnection] {{ if .Client.User }}user: {{ .Client.User }}{{ end }} cid: {{ .Client.ID }} in account {{ .Client.Account }}: {{ .Reason }}`)
	if err != nil {
		panic(err)
	}

	err = event.RegisterTextExtendedTemplate("io.nats.server.advisory.v1.client_disconnect", `
[{{ .Time | ShortTime }}] [{{ .ID }}] Client Disconnection
{{ if .Reason }}
    Reason: {{ .Reason }}
{{- end }}
    Server: {{ .Server.Name }}
{{- if .Server.Cluster }}
   Cluster: {{ .Server.Cluster }}
{{- end }}

   Client:
                 ID: {{ .Client.ID }}
{{- if .Client.User }}
               User: {{ .Client.User }}
{{- end }}
{{- if .Client.Name }}
               Name: {{ .Client.Name }}
{{- end }}
            Account: {{ .Client.Account }}
    Library Version: {{ .Client.Version }}  Language: {{ with .Client.Lang }}{{ . }}{{ else }}Unknown{{ end }}
{{- if .Client.Host }}
               Host: {{ .Client.Host }}
{{- end }}
{{- if .Client.Jwt }}
         Issuer Key: {{ .Client.IssuerKey }}
           Name Tag: {{ .Client.NameTag }}
{{- if .Client.Tags }}
               Tags: {{ .Client.Tags | JoinStrings }}
{{- end }}
{{- end }}
{{- if .Client.Kind }}
        Client Kind: {{ .Client.Kind }}
{{- end }}
{{- if .Client.ClientType }}
        Client Type: {{ .Client.ClientType }}
{{- end }}

   Stats:
                  RTT: {{ .Client.RTT }}
      Client Received: {{ .Received.Msgs }} messages ({{ .Received.Bytes | IBytes }})
                   Gateways: {{ .Received.Gateways.Msgs }} messages ({{ .Received.Gateways.Bytes | IBytes }})
                     Routes: {{ .Received.Routes.Msgs }} messages ({{ .Received.Routes.Bytes | IBytes }})
                 Leaf Nodes: {{ .Received.Leafs.Msgs }} messages ({{ .Received.Leafs.Bytes | IBytes }})
     Client Published: {{ .Sent.Msgs }} messages ({{ .Sent.Bytes | IBytes }})
                   Gateways: {{ .Sent.Gateways.Msgs }} messages ({{ .Sent.Gateways.Bytes | IBytes }})
                     Routes: {{ .Sent.Routes.Msgs }} messages ({{ .Sent.Routes.Bytes | IBytes }})
                 Leaf Nodes: {{ .Sent.Leafs.Msgs }} messages ({{ .Sent.Leafs.Bytes | IBytes }})`)
	if err != nil {
		panic(err)
	}
}
