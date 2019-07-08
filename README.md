# Docker static

[![Build Status](https://api.cirrus-ci.com/github/dirkdev98/docker-static.svg)](https://cirrus-ci.com/github/dirkdev98/docker-static)
[![](https://images.microbadger.com/badges/version/dirkdev98/docker-static.svg)](https://microbadger.com/images/dirkdev98/docker-static)
![Docker Pulls](https://img.shields.io/docker/pulls/dirkdev98/docker-static.svg)
[![](https://images.microbadger.com/badges/image/dirkdev98/docker-static.svg)](https://microbadger.com/images/dirkdev98/docker-static)



Simple static server written in Go

Inspired by: [PierreZ/goStatic](https://github.com/PierreZ/goStatic)

## Usage

```
$ /docker-static [args] 

-auth string
    Basic authorization in form of username:password
-fallback
    Automatically try to serve index.html if file is not found (default true)
-maxAge int
    Cache-Control header value (default 3600)
-path string
    Path on which files will be found (default "/public")
-port int
    Port to listen on (default 3000)
```

## Monitoring

The `GET /health` route will always return `200 OK` and based on response timings you can decide to replace or scale up.
Or just use is to let Docker now that the server is alive.


## Deploying

This project does not support HTTPS and it's expected that a reverse proxy is setup to handle that. Logging will all be done
to stdout and stderr with the expectation that a different system will handle log collection and processing

## Docker usage

Simple run:
```bash
$ docker run -v ./path/to/public:/public dirkdev98/docker-static
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
CMD ["/docker-static", "-maxAge", "7200"]
```

## License

MIT Copyright (c) 2019 Dirk de Visser
