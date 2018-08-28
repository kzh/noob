FROM golang:alpine AS build

WORKDIR /go/src/github.com/k3v/noob/
ADD . /go/src/github.com/k3v/noob/

RUN go build -o main ./frontend/

FROM alpine

WORKDIR /app/

COPY --from=build /go/src/github.com/k3v/noob/main ./frontend
ADD ./frontend/templates/ ./templates/
ADD ./frontend/static/ ./static/

CMD ["./frontend"]
