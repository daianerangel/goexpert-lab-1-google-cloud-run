FROM golang:1.19 as builder
WORKDIR /app
COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -v -o cloudrun

FROM scratch
WORKDIR /app
COPY --from=builder /app/cloudrun .
ENTRYPOINT ["./cloudrun"]