FROM golang:alpine as build

WORKDIR /go/src/app
COPY . .

RUN go install -i

FROM alpine

COPY --from=build /go/bin/garchive /go/bin/

WORKDIR /app
COPY index.html .
COPY static static/
ENTRYPOINT ["/go/bin/garchive", "--port", "80"]
