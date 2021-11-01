FROM alpine

RUN apk add --no-cache curl zeromq=4.3.4 libsodium==1.0.18

# Copy binary from github releases page once that exists :P

ENTRYPOINT [ "/home/alpine/xnode" ]