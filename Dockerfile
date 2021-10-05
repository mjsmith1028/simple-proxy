FROM golang:1.15-alpine AS builder

WORKDIR /home/go/

COPY main.go /home/go/

RUN CGO_ENABLED=0 go build -o /bin/simple-proxy

# Start building the final image
FROM scratch

COPY --from=builder /bin/simple-proxy /bin/simple-proxy

USER 2000

ENTRYPOINT ["/bin/simple-proxy"]
