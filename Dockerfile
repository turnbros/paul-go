FROM golang:1.16-alpine as paul-builder

WORKDIR /paul

# Download dependencies
COPY go.mod ./
COPY go.sum ./
RUN go mod download

# Copy in the code and build scripts
COPY cmd cmd
COPY internal internal
COPY pkg pkg
COPY scripts scripts
RUN ./scripts/build.sh

FROM alpine:3.15
RUN adduser -h "/paul" -u 3240 -g "Paul" -D paul
COPY --from=paul-builder --chown=paul /paul/dist/* /usr/local/bin
USER paul
EXPOSE 8443