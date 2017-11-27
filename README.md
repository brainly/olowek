# Ołówek
<p align="center">
  <img style="float: right;" height="300" src="doc/olowek.png" alt="Ołówek logo"/>
</p>

Ołówek [ɔˈwuvɛk] - configure your Nginx with applications deployed on Mesos/Marathon.

**Features:**

* Custom templates with Go `text/template`
* Uses Marathon event stream for updates *(requires at least Marathon v0.9.0)*
* Automatic discovery of all applications running on Marathon

## Installation

1. Get pre-compiled binary from `releases` page
2. Edit configuration *(default location is `/etc/olowek/olowek.json`)*:

```json
{
  "scope": "internal",
  "marathon": "http://127.0.0.1:8080,127.0.0.1:8080",
  "nginx_config": "/etc/nginx/conf.d/services.conf",
  "nginx_template": "/etc/olowek/services.tpl",
  "nginx_cmd": "/usr/sbin/nginx"
}
```
3. Start olowek
