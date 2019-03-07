FROM golang:1.12.0
WORKDIR /src
COPY . .
RUN CGO_ENABLED=0 go build

FROM scratch
COPY --from=0 /src/githashcrash /githashcrash
ENTRYPOINT ["/githashcrash"]
