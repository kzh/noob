FROM golang:alpine AS build

WORKDIR /go/src/github.com/kzh/noob/
ADD . /go/src/github.com/kzh/noob/

RUN go build -o main ./frontend/

FROM alpine

WORKDIR /app/

COPY --from=build /go/src/github.com/kzh/noob/main ./frontend
ADD ./frontend/templates/ ./templates/
ADD ./frontend/static/ ./static/

CMD ["./frontend"]
