FROM golang:1.23-alpine AS builder

RUN apk add --no-cache git ca-certificates

WORKDIR /build

COPY go.mod go.sum ./
RUN go mod download

COPY . .

ARG VERSION
ARG COMMIT
RUN CGO_ENABLED=0 go build -ldflags "-X neocut/internal/config.Commit=${COMMIT}" -o neocut ./cmd/neocut/

FROM alpine:3.20

RUN apk add --no-cache ca-certificates ffmpeg

WORKDIR /root/

COPY --from=builder /build/neocut /usr/local/bin/neocut

ENTRYPOINT ["neocut"]
CMD ["--help"]
