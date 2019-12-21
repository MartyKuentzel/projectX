# Dockerfile References: https://docs.docker.com/engine/reference/builder/

# Start from the latest golang base image
FROM golang:latest

# Add Maintainer Info
LABEL maintainer="Marty Kuentzel"

# Set the Current Working Directory inside the container
WORKDIR /app

# Copy go mod and sum files
COPY go.mod go.sum ./

# Download all dependencies. Dependencies will be cached if the go.mod and go.sum files are not changed
RUN go mod download

# Copy the source from the current directory to the Working Directory inside the container
COPY api api
COPY pkg pkg
COPY cmd cmd 
COPY third_party third_party
COPY my_wrapper_script.sh .

RUN wget https://dl.google.com/cloudsql/cloud_sql_proxy.linux.amd64 -O cloud_sql_proxy \
    && chmod +x cloud_sql_proxy my_wrapper_script.sh

# Build the Go app
RUN go build /app/cmd/server/main.go

# Expose port 8080 to the outside world
EXPOSE 8080

CMD ./my_wrapper_script.sh ${instances} credentials.json ${dbPw}
