FROM golang:1.21.0-alpine AS base

WORKDIR /app

RUN apk --no-cache update && apk --no-cache upgrade  && \
    apk add --no-cache bash ca-certificates git make && \
    update-ca-certificates && rm -rf /var/cache/apk*

COPY go.* ./
RUN  go mod download

COPY . ./
RUN  make

FROM scratch AS final

copy --from=base /app/bin/bot /app/bot
copy --from=base /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs

ENTRYPOINT ["/app/bot"]
