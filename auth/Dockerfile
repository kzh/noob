FROM golang:alpine AS build

WORKDIR /go/src/github.com/k3v/noob/
ADD . /go/src/github.com/k3v/noob/

RUN go build -o main ./auth/

FROM alpine

WORKDIR /app/

COPY --from=build /go/src/github.com/k3v/noob/main ./auth

CMD ["./auth"]
