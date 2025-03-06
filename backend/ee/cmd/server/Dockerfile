FROM golang:1.22 AS build
WORKDIR /go/src/github.com/trysourcetool/sourcetool/backend/ee
COPY go.mod .
COPY go.sum .
RUN go mod download
COPY . .
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o main -a -installsuffix cgo -ldflags="-s -w" cmd/server/*.go

FROM scratch
WORKDIR /root
COPY --from=build /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/
COPY --from=build /go/src/github.com/trysourcetool/sourcetool/backend/ee/main main
EXPOSE 8080
CMD ["/root/main"]