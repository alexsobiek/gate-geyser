FROM golang:1.23.5-alpine3.21 AS build

WORKDIR /src

COPY ./ /src/

RUN go build -o /bin/gate server/main.go

FROM alpine:3.21

COPY --from=build /bin/gate /bin/gate

CMD ["/bin/gate" "-c", "/etc/gate/config.yml"]