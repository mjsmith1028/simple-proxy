FROM golang:1.15-alpine AS builder

WORKDIR /home/go/

COPY main.go /home/go/

RUN CGO_ENABLED=0 go build -o /bin/proxy

# Start building the final image
FROM scratch

COPY --from=builder /bin/proxy /bin/proxy

ENTRYPOINT ["/bin/proxy"]