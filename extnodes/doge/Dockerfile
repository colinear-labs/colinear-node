FROM casperstack/dogecoin

ENTRYPOINT [ "dogecoind", \
    "-disablewallet", \
    "-prune=2200", \
    "-maxmempool=20", \
    "-server", \
    "-rpcport=5003", \
    "-rpcallowip=127.0.0.1" \
]