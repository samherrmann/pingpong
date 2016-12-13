# pingpong

A simple Go app that monitors the availability of network nodes.

**This is only a exercise project app and is not intended to be used in production!**

## Usage

### Launching
To launch `pingpong`, simply call the executable from a terminal:
```bash
$ pingpong
```

### Observing
Using a web browser, navigate to `http://localhost:8080` to view the results in a web UI. Navigate to `http://localhost:8080/json` to view the results
in JSON format.

### Help
Run `pingpong -h` for additional options.

### config.json
When running `pingpong` with no `config.json` file present in the same directory, it will auto-generate a sample file for you with the following content:

```json
{
    "pingpong": "http://localhost:8080",
    "pingpong host": "localhost"
}
```
The config file consists of a simple JSON object in which the keys represent the names of the network nodes and the values are the network addresses of the nodes. Modify this file to suit your needs.

### Protocols
If an address of a network node starts with prefix `http://` or `https://`, `pingpong` will use the HTTP protpcol to test the nodes' availability. If neither prefix is present, `pingpong` will use ping.

## Contributing
### Reference Dev Environment
* [Go 1.7](https://golang.org/)
* [VS Code](http://code.visualstudio.com/) (with [vscode-go](https://github.com/Microsoft/vscode-go) extension)

### Building


#### For Development
```bash
$ go generate
$ go build
```
* Only builds for the currently set `GOOS` and `GOARCH`
* Does **not** embed the software version in the build


#### For Distribution
Prerequisite: If you don't have `gowrap` installed, install it by following its [installation instructions](https://github.com/samherrmann/gowrap#installation).

If `gowrap` is installed, run it from the terminal:
```bash
$ gowrap
```
* Automatically creates builds for all the platforms listed in `gowrap.json`
* Embeds the software version in the build