FROM kylemanna/bitcoind

ENTRYPOINT [ "bitcoind", \ 
    "-disablewallet", \
    "-prune=550", \
    "-maxmempool=20", \
    "-server", \
    "-rpcport=5000", \
    "-rpcallowip=127.0.0.1" \
]