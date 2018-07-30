# Use: docker run --rm --pid=host thajeztah/go-top
FROM golang:1.10-alpine AS build
RUN apk add --no-cache gcc musl-dev
WORKDIR /go/src/github.com/thaJeztah/go-top/
COPY . . 
RUN go build -ldflags "-linkmode external -extldflags -static" -o /go/bin/go-top -a cmd/go-top/main.go


FROM scratch
COPY --from=build /go/bin/go-top /
CMD ["/go-top"]
