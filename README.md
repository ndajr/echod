# echod
Echod is an HTTP server which echoes whatever JSON a client sent through the API appending `"echoed": true`. This is a demo project to demonstrate an API built in Go with logs, metrics, unit and integration tests. I decided to also include Envoy as an HTTPS proxy (for demonstration purposes) and Prometheus on the docker-compose.yaml. If you want to run just the service with docker, you can do so with `make dockerrun`

The API consists on the following endpoints:
- POST /api/echo (requires [Authentication](#authentication))
- PUT /api/echo (requires [Authentication](#authentication))
- GET /health
- GET /metrics

## Running
Pre-requisites:
- docker
- docker-compose
```bash
make run
curl -k https://localhost:10000/health
```
Note: Always use `curl -k` or `curl --insecure` to bypass the SSL validation since it is a self-signed certificate installed on Envoy. For production we would use something like Let's Encrypt and certbot to avoid this issue.

Running the tests ([Go](https://go.dev/dl/) is required):
```bash
make test
```

Running the lints:
Pre-requisites:
- [golangci-lint](https://golangci-lint.run/usage/install/#local-installation)
```bash
make lint
```

## Authentication
To use the `/api/echo` endpoint you must be authenticated using basic auth with one of the following users and password respectively:
| user          | password |
|---------------|----------|
| remarkablebob | 2rFtex   |
| dextroussheep | vSKjGK   |
| luckyslug     | VRsBMW   |
| pettyrabbit   | 2BpgwH   |

Example:
```bash
curl -k -u remarkablebob:2rFtex -X POST https://localhost:10000/api/echo --data '{"foo": "bar"}'
```
```json
{"echoed":true,"foo":"bar"}
```

## Documentation
Pre-requisites:
- go
- godoc (you can install by running `go get golang.org/x/tools/cmd/godoc`)
```bash
make doc
```

## Benchmark
Simple performance benchmark running both the server and [wrk](https://github.com/wg/wrk) locally:
```
$ wrk -t4 -c400 -d30s --script config/wrk/echo.lua http://localhost:3000/api/echo

Running 30s test @ http://localhost:3000/api/echo
  4 threads and 400 connections
  Thread Stats   Avg      Stdev     Max   +/- Stdev
    Latency     3.47ms    1.66ms  53.39ms   84.25%
    Req/Sec    18.04k     1.89k   24.36k    68.92%
  2153668 requests in 30.01s, 384.08MB read
  Socket errors: connect 153, read 19, write 21, timeout 0
Requests/sec:  71764.55
Transfer/sec:     12.80MB
```

## Monitoring
Echod exposes two metrics on Prometheus format:
- `http_requests_total`: Measures the total http requests received labeled by status code, method and path (Counter).
- `http_request_duration_seconds`: Measures the duration of each http request labeled by status code, method and path (Histogram)

Alerts:
- `InstanceDown`: triggered when the service is down for more than 5 minutes.
