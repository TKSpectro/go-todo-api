FROM golang:1.19-alpine

# Set destination for COPY
WORKDIR /app

RUN apk update && \
    apk add curl git protobuf bash make openssh-client gcc && \
    rm -rf /var/cache/apk/*

# Download Go modules
COPY go.mod go.sum ./
RUN go mod download

# Copy the source code. Note the slash at the end, as explained in
# https://docs.docker.com/engine/reference/builder/#copy
COPY *.go ./
COPY . ./

# install ginkgo
RUN go get github.com/onsi/ginkgo/v2
RUN go install github.com/onsi/ginkgo/v2/ginkgo

# Run
CMD CGO_ENABLED=0 GOOS=linux GTA_ROOT_PATH=$(pwd) IS_TEST=true ginkgo -r