FROM golang:alpine as build

WORKDIR /go/src/garchive
COPY . .

RUN go install -i

FROM alpine

COPY --from=build /go/bin/garchive /go/bin/

WORKDIR /app
COPY index.html .
COPY static static/
EXPOSE 8080
ENTRYPOINT ["/go/bin/garchive"]
