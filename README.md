# Docker static

Simple static server written in Go with basic prometheus metrics

Inspired by: [PierreZ/goStatic](https://github.com/PierreZ/goStatic)

## Usage

```
$ /docker-static [args] 

-auth string
    Basic authorization in form of username:password
-fallback string
    Default file that will be served
-maxAge int
    Cache-Control header value (default 3600)
-monitoring
    Enable or disable monitoring (default true)
-monitoringPort int
    Monitoring port to listen on (default 9600)
-path string
    Path on which files will be found (default "/public")
-port int
    Port to listen on (default 3000)
```

## Monitoring

By default prometheus monitoring and a basic health route are enabled on a separate configurable port
The `GET /health` route will always return `200 OK` and based on response timings you can decide to replace or scale up.
Or just use is to let Docker now that the server is alive.

`GET /metrics` exposes metrics to be consumed by prometheus. Apart from the default 
[go runtime metrics](https://github.com/prometheus/client_golang/blob/master/prometheus/go_collector.go)
this server exposes authorization success and failure counts, request count per path, method, response status
 and lastly response timings per path in seconds.

## Deploying

This project does not support HTTPS and it's expected that a reverse proxy is setup to handle that. Logging will all be done
to stdout and stderr with the expectation that a different system will handle log collection and processing

## Docker usage

Simple run:
```bash
$ docker run -v ./path/to/public:/public dirkdev98/docker-static -fallback index.html
```

Multistage build with for example a React app
```Dockerfile
# Build React app, customize as needed
FROM node:10 as build-deps
WORKDIR /usr/src/app
COPY package.json yarn.lock ./
RUN yarn
COPY . ./
RUN yarn build

# Build final image with only the builded files
FROM dirkdev98/docker-static
COPY --from=build-deps /usr/src/app/build /public
CMD ["/docker-static", "-fallback", "index.html", "-maxAge", "7200"]
```