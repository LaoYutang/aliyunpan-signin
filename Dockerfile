FROM golang:1.20 AS builder
WORKDIR /data
COPY . .
RUN go mod download
RUN CGO_ENABLED=0 go build -o ./app .

FROM alpine:3.19
WORKDIR /data
COPY --from=builder /data/app .
CMD [ "./app" ]
