FROM golang:latest AS build-env
RUN mkdir -p /go/src/github.com/moqmar/formaility
WORKDIR /go/src/github.com/moqmar/formaility
COPY *.go /go/src/github.com/moqmar/formaility
RUN go get github.com/moqmar/formaility
RUN CGO_ENABLED=0 GOOS=linux go build -a -ldflags '-extldflags "-static" -s' -installsuffix cgo -o formaility -v .

FROM alpine:latest as network
RUN apk --no-cache add tzdata zip ca-certificates
WORKDIR /usr/share/zoneinfo
RUN zip -r -0 /zoneinfo.zip .

# Put everything together
FROM scratch
COPY --from=build-env /go/src/github.com/moqmar/formaility/formaility /
COPY --from=network /zoneinfo.zip /
COPY --from=network /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/

ENV GIN_MODE=release
WORKDIR /
EXPOSE 8080

ENTRYPOINT [ "/formaility" ]
