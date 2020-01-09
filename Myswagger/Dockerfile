FROM golang:1.12 as build

ENV GIN_MODE release
ENV GOPROXY https://goproxy.io
ENV GO111MODULE on

WORKDIR /go/cache

ADD go.mod .
RUN go mod tidy
RUN go mod download

WORKDIR /go/release

ADD . .

RUN GOOS=linux CGO_ENABLED=0 go build -ldflags="-s -w" -installsuffix cgo -o grpcswagger main.go

FROM scratch as prod

COPY --from=build /usr/share/zoneinfo/Asia/Shanghai /usr/share/zoneinfo/Asia/Shanghai
COPY --from=build /go/release/grpcswagger /

CMD ["/grpcswagger"]
