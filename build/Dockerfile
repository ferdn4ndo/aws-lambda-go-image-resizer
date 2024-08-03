FROM golang@sha256:f94174c5262af3d8446833277aa27af400fd1a880277d43ec436df06ef3bb8ab

LABEL maintaner="Fernando Constantino <const.fernando@gmail.com>"

# Install GCC + git + SSL ca certificates
# GCC is required to build the sqlite3 dependency
# Git is required for fetching the dependencies
RUN apk update && apk add git gcc musl-dev curl sqlite coreutils && rm -rf /var/cache/apk/*

# Create appuser
ENV USER=appuser
ENV UID=1000

# See https://stackoverflow.com/a/55757473/12429735RUN
RUN adduser \
    --disabled-password \
    --gecos "" \
    --home "/nonexistent" \
    --shell "/sbin/nologin" \
    --no-create-home \
    --uid "${UID}" \
    "${USER}"

WORKDIR /go/src/github.com/ferdn4ndo/userver-logger-api/

COPY ./ /go/src/github.com/ferdn4ndo/userver-logger-api/

RUN go mod download gorm.io/driver/sqlite github.com/go-chi/chi/v5 github.com/go-chi/docgen github.com/go-chi/render gorm.io/gorm

RUN go mod verify

RUN GOOS=linux GOARCH=amd64 go build \
    -ldflags='-w -s -extldflags "-static"' \
    -o /go/bin/server . \
    && chmod +x /go/bin/server

EXPOSE 5555

HEALTHCHECK --interval=30s --timeout=10s --retries=6 --start-period=5s CMD curl -s -f http://localhost:5555/health

ENTRYPOINT ["/go/bin/server"]
