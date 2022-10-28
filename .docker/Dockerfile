##
## Build
##

FROM golang:1.18-alpine AS build

ADD ./ /app

WORKDIR /app

RUN set -ex \
    && apk add --no-cache --virtual .build-deps build-base \
    && apk add --virtual rundeps $runDeps \
    && apk add tzdata \
    && cp /usr/share/zoneinfo/Asia/Singapore /etc/localtime \
    && apk del .build-deps \
    && apk del tzdata

RUN CGO_ENABLED=0 go build -a -ldflags '-extldflags "-static"' .

##
## Deploy
##

FROM gcr.io/distroless/base-debian10

#WORKDIR /

COPY --from=build /app/chatwebsocketgin /app/chatwebsocketgin

WORKDIR /app

EXPOSE 8080

USER nonroot:nonroot

ENTRYPOINT ["/app/chatwebsocketgin"]
