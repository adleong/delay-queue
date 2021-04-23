FROM golang:1.16.2-alpine as golang
WORKDIR /messages-build
COPY . .
RUN go build .

FROM alpine
COPY --from=golang /messages-build/messages /
ENTRYPOINT ["/messages"]
