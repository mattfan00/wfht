# syntax=docker/dockerfile:1 
FROM node:21-alpine as tailwind

WORKDIR /tailwind

COPY ui/static/tailwind.css ./static/tailwind.css
COPY ui/package.json ./
RUN npm install
RUN npm run build 


FROM golang:1.21-alpine as app
# install gcc needed for cgo
RUN apk add gcc musl-dev 

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download && go mod verify

COPY . .
COPY --from=tailwind /tailwind/static/index.css ./ui/static/index.css
RUN go build -o wfht ./cmd/wfht


FROM alpine:3.19

WORKDIR /app

COPY --from=app /app/ui/static ./ui/static
COPY --from=app /app/ui/views ./ui/views
COPY --from=app /app/wfht ./

ENTRYPOINT [ "./wfht" ]
