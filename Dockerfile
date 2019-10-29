FROM golang:1.13.0
WORKDIR /src
COPY . .
RUN CGO_ENABLED=0 go build -o "githashcrash" "cmd/githashcrash/main.go"

FROM scratch
COPY --from=0 "/src/githashcrash" "/githashcrash"
ENTRYPOINT ["/githashcrash"]
