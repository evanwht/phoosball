FROM node AS ui

RUN mkdir /ui
WORKDIR /ui

COPY ui/package.json /ui
COPY ui/package-lock.json /ui
RUN npm install

COPY ui /ui

RUN npm run build

# build go binary
FROM golang:alpine AS app

RUN apk update && apk add --no-cache git
WORKDIR $GOPATH/src/app

COPY server ./
RUN go get -d -v ./...
RUN go build -o /go/bin/phoos_server phoos_server.go

# build final image
FROM alpine

RUN mkdir /public && mkdir /db && mkdir /db/prod

COPY --from=ui /ui/build /public
COPY --from=app /go/bin/phoos_server .
COPY db/prod /db/prod

EXPOSE 3032
ENTRYPOINT [ "./phoos_server" ]