FROM bitcoinabc/bitcoin-abc:0.24.3

USER root

RUN apt-get update && apt-get install -y curl

USER bitcoin

ENTRYPOINT bitcoind \ 
    -disablewallet \
    -prune=550 \
    -maxmempool=20 \
    -server \
    # auth --> user:pass
    -rpcauth=user:50c2fd98a724caf63afadca4a94ff94f\$18b83531c0eb837705dff39ad6cd5c5c012ca0d54a484f7cb0ce139667acb269 \ 
    -rpcport=5002 \
    -rpcallowip=127.0.0.1 \
    -blocknotify="curl -d 'name=bch&hash=%s' http://127.0.0.1:4999"
