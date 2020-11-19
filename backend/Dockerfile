FROM golang:alpine AS builder

RUN apk add --no-cache build-base

RUN adduser -S go && addgroup -S go
RUN mkdir /src/ && chown -R go:go /src
USER go:go
WORKDIR /src

COPY --chown=go:go ./go.* /src/
RUN go mod download

COPY --chown=go:go ./ /src/

RUN go build -o /src/backend .

############################################
FROM alpine:3.12

COPY --from=builder /src/backend /backend

RUN mkdir /data && chown -R nobody:nobody /data
VOLUME /data
USER nobody:nobody

ENTRYPOINT ["/backend"]
