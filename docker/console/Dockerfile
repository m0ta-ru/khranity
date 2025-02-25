# declare incomining arguments
ARG GO_VERSION
# start from a Alpine image with the specified version of Go installed
# and a workspace (GOPATH) configured at /go. AS golang
FROM golang:${GO_VERSION}-alpine AS golang
# set label for future deleting
LABEL stage=temp
# set environment variables
ENV GO_VERSION=${GO_VERSION}
# update packages and install tzdata for further usage
RUN apk update --no-cache && apk add --no-cache tzdata
# set default directory
WORKDIR /build
# copy the local package files before compiling it
COPY . .
# download and verify all required dependencies
RUN go mod download && go mod verify
# build executable file
RUN go build -o /build/exec ./cmd/console/*.go

# start another blank image from the latest version of an Alpine image
FROM alpine
# set label
LABEL stage=run
# contact maintainer with any issues you encounter
LABEL maintainer.name="m0ta" \
      maintainer.email="dev@m0ta.ru"
# update packages and install ca-certificates for further usage
RUN apk update --no-cache && apk add --no-cache ca-certificates
# set default directory
WORKDIR /exec
# copy necessary timezone into the production image
COPY --from=golang /usr/share/zoneinfo /usr/share/zoneinfo
# copy the executable file from the first stage
# into the production image
COPY --from=golang /build/exec /exec/exec
# run executable file as default command (entrypoint)
ENTRYPOINT ["./exec"]