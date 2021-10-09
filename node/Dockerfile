# NOTE: STILL UNTESTED

FROM alpine

RUN apk add --no-cache go=1.17.1
RUN go build 

COPY ./xnode /home/alpine

ENTRYPOINT [ "/home/alpine/xnode" ]