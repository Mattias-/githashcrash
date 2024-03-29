FROM golang:1.20.1
WORKDIR /src
COPY . .
RUN CGO_ENABLED=0 go build -o "githashcrash" "cmd/githashcrash/main.go"

FROM alpine:latest
RUN apk add git
COPY --from=0 "/src/githashcrash" "/githashcrash"
CMD ["/githashcrash"]
