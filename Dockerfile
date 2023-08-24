ARG GO_VERSION=1.20

FROM golang:${GO_VERSION}-alpine AS builder

RUN go env -w GOPROXY=direct
RUN apk add --no-cache git
RUN apk --no-cache add ca-certificates && update-ca-certificates

WORKDIR /src

COPY ./go.mod -/go.sum ./
RUN go mod download

COPY ./ ./

RUN CGO_ENABLED=0 go build \
    -installsuffix 'static' \
    -o /curso-go-rest .

FROM scratch AS runner

COPY --from=builder /etc/ssl/certs/ca-certificates.crt /ect/ssl/certs/

COPY .env ./
COPY --from=builder /curso-go-rest /curso-go-rest

EXPOSE 5050

ENTRYPOINT [ "/curso-go-rest" ]
