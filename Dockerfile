# Use: docker run --rm --pid=host thajeztah/go-top
FROM golang:1.10-alpine AS build
WORKDIR /go/src/github.com/thaJeztah/go-top/
COPY . . 
RUN go install -v ./...

FROM scratch
COPY --from=build /go/bin/go-top /
CMD ["/go-top"]
