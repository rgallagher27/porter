FROM golang:1.19-alpine as build

ENV CGO_ENABLED="0" GOOS="linux" GOARCH="amd64"

WORKDIR /app

COPY . .

RUN go build -o /porter_bin .

FROM scratch
USER 10001:10001

COPY --from=build /porter_bin /porter

ENTRYPOINT ["/porter"]
