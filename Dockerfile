# Start from a Debian image with the latest version of Go installed
# and a workspace (GOPATH) configured at /go.
FROM golang

# Fetching gouuid dependencies
RUN go get github.com/nu7hatch/gouuid

# Copy the local package files to the container's workspace.
ADD . /go/src/github.com/hamzamac/chord

#Install the main project
RUN go install github.com/hamzamac/chord

# Run the outyet command by default when the container starts.
ENTRYPOINT ["/go/bin/"]

# Document that the service listens on port 8080.
EXPOSE 8080
