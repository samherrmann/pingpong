# pingpong

A simple Go app that issues HTTP GET requests to a specified list of servers to verify that the servers are operational. The results can either be observed through a web UI or in JSON format.

**This is only a exercise project app and is not intended to be used in production!**

## Usage

Example:
```
$ pingpong http://server-1-url [http://server-2-url http://server-N-url]
```
Using a web browser, navigate to `http://localhost:8080` to view the results in a web UI. Navigate to `http://localhost:8080/json` to view the results
in JSON format.

Run `pingpong -h` for additional options.