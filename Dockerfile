FROM golang:1.13 AS builder
ENV CGO_ENABLED 0
WORKDIR /go/src/app
ADD . .
RUN go build -mod vendor -o /deployer

FROM scratch
COPY --from=builder /deployer /deployer
CMD ["/deployer"]