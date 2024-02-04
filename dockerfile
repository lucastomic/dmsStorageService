FROM golang:alpine AS build
WORKDIR /go/src/myapp
COPY . .
RUN go build -o /go/bin/myapp cmd/main.go


FROM alpine
COPY --from=build /go/bin/myapp /go/bin/myapp
ENV PROJECT_ROOT=/go/src/myapp
RUN mkdir files
ENTRYPOINT ["/go/bin/myapp"]
