# wmt (Web Multi-Tool)

A collection of web dev utilities, primarily for testing web services with templated requests.

## Installation

### Requirements
- go 1.19+, gcc, Cgo
- Linux: X11 dev package, e.g. libx11-dev, xorg-dev or libX11-devel

```sh
go install github.com/jbchouinard/wmt@latest
```

## Example
```sh
# wmt opt set editor "code -w -n"
# wmt opt set editor "subl -w -n" 
# wmt opt set editor "vim"

wmt env use dev
wmt env set baseUrl http://localhost:8080/v1

wmt request add createPerson POST /person --template
# Edit, save and close editor:
# {
#   firstName: "{{.firstName}}"
#   lastName: "{{.lastName}}"
# }
wmt request header createPerson Content-Type application/json
wmt request do createPerson -p firstName=Jane -p lastName=Smith

wmt request add getPerson GET "/person/{{.id}}"
wmt request do getPerson -p id=1
```

## Env Commands

## Template Commands

`wmt` uses go templates. For basic variable substitution, use `{{.foo}}`.

Text templates do not have any security features. HTML templates (created with the `--html` flag)
are safe against HTML code injection.

See the go documentation for advanced features:
- Text: https://pkg.go.dev/text/template
- HTML: https://pkg.go.dev/html/template

### list
```sh
wmt template list
```

### add
```sh
wmt template add <name> [--html]
```

### edit
```sh
wmt template edit <name>
```

### delete
```sh
wmt template delete <name>
```

### eval
```sh
wmt template eval <name> [-p <param>=<value>...]
```

## Options

### opt

List all options:
```sh
wmt opt
```

Get an option value:
```sh
wmt opt <key>
```

Set an option value (set to _ to clear the option):
```sh
wmt opt <key> <value>
```

| Key       | Values | Default | Details                       |
------------|--------|---------|-------------------------------|
| clipboard | yes/no | no      | Enable clipboard integration? |
| history   | yes/no | yes     | Save history?                 |
| editor    | *      | nano    | Text editor to spawn.         |

If using a GUI text editor, the text editor command should not return until
the window is closed. For example, use `code -w` for VS Code.

## Other Commands

### uuid
```sh
wmt uuid [--v4]
```

Generate a UUID (default: V1).


### history
```sh
wmt history <command>
```

Show history for a command.

### purge
```sh
wmt purge [--keep-days n]
```
Purge history, keeping n days (default: 7);

## Files

`wmt` tries to stores its files in `<HOME>/.wmt`. If that fails, it tries `<CURRENTDIR>/.wmt`.

The location can be changed by setting the `WMT_DIR` environment variable.

## License

Copyright 2022 Jerome Boisvert-Chouinard, under MIT License.
