FROM golang:alpine AS build

WORKDIR /go/src/github.com/kzh/noob/
ADD . /go/src/github.com/kzh/noob/

RUN go build -o main ./cmd/executor/

FROM alpine

WORKDIR /app/

COPY --from=build /go/src/github.com/kzh/noob/main ./executor
COPY ./cmd/executor/image.tar .

CMD ["./executor"]
