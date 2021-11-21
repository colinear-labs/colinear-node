FROM kylemanna/bitcoind

RUN apt-get update && apt-get install -y curl

ENTRYPOINT bitcoind \ 
    -disablewallet \
    -prune=550 \
    -maxmempool=20 \
    # -zmqpubsequence=tcp://0.0.0.0:5000 \
    -server \
    # auth --> user:pass
    -rpcauth=user:50c2fd98a724caf63afadca4a94ff94f\$18b83531c0eb837705dff39ad6cd5c5c012ca0d54a484f7cb0ce139667acb269 \ 
    -rpcport=5000 \
    -rpcallowip=127.0.0.1 \
    -blocknotify="curl -d 'name=btc&hash=%s' http://127.0.0.1:4999"