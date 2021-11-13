# NOTE: STILL UNTESTED

FROM alpine

RUN apk add --no-cache go=1.17.1 zeromq=4.3.4 libsodium==1.0.18

COPY . .

ENTRYPOINT [ "go", "run", "main.go" ]

EXPOSE 4999