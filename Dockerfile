FROM golang as build

WORKDIR /src

COPY ./ /src

RUN GOPROXY=https://goproxy.cn GOOS=linux GOARCH=amd64 CGO_ENABLED=0 go build -o api

FROM alpine

COPY --from=build /src/api /usr/bin/api

WORKDIR /app

COPY ./.env.common /app/.env.common

RUN chmod +x /usr/bin/api

CMD /usr/bin/api
