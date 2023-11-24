FROM golang:latest as Builder
WORKDIR /go/src/github.com/api
COPY . .
RUN go mod download
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o app ./cmd

FROM alpine:latest
RUN apk --no-cache add ca-certificates
WORKDIR /root/
COPY --from=Builder /go/src/github.com/api/app .
COPY --from=Builder /go/src/github.com/api/configs/base_prod.yaml .
CMD ["./app"]

