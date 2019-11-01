FROM golang:latest as builder
WORKDIR /ipchecker
ADD . .
RUN CGO_ENABLED=0 go build -o ./bin/main ./cmd/checker_service.go
FROM scratch
COPY --from=builder /ipchecker/bin/main /
EXPOSE 8080
CMD ["/main"]