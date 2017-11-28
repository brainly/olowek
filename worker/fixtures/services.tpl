{{- range $id, $app := .Apps }}
  {{- with $tasks := $app.Tasks }}
upstream {{ $app.Name }} {
    {{- range $tasks }}
      server {{ .Host }}:{{ index .Ports 0 }};
    {{- end }}
}
  {{- end }}

server {
  server_name {{ $app.Name }}.example.com;
  listen 80;

  location / {
    {{- with $tasks := $app.Tasks }}
      proxy_next_upstream error timeout http_502 http_503;
      proxy_pass http://{{ $app.Name }};
    {{- else }}
      return 503;
    {{- end }}
  }
}
{{- end }}
