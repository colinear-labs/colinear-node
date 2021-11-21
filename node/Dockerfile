# NOTE: STILL UNTESTED

FROM alpine

RUN apk update && apk add --no-cache go

COPY . .

ENTRYPOINT [ "go", "run", "main.go" ]

EXPOSE 4999