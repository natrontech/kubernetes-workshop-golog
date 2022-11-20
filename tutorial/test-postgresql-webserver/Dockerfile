# syntax=docker/dockerfile:1
FROM golang:1.19 as build
WORKDIR /go/src/github.com/natrongmbh/test-postgresql-webserver
COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o test-postgresql-webserver .

FROM alpine:latest  
RUN apk --no-cache add ca-certificates
WORKDIR /root/
COPY --from=build /go/src/github.com/natrongmbh/test-postgresql-webserver/test-postgresql-webserver ./
EXPOSE 8080
CMD ["./test-postgresql-webserver", "cluster"]