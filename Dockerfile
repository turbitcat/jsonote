FROM golang:1.19 AS builder

WORKDIR /go/src/jsonote
COPY . .
RUN CGO_ENABLED=0 go build -o /go/bin/jsonote

FROM gcr.io/distroless/static-debian11
COPY --from=builder /go/bin/jsonote /

ENV JSONOTE_PATH=/data
EXPOSE 8088

CMD ["/jsonote"]