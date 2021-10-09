FROM zquestz/bitcoin-cash-node

ENTRYPOINT [ "bitcoind", \ 
    "-disablewallet", \
    "-prune=10", \
    "-maxmempool=20", \
    "-server", \
    "-rpcport=5002", \
    "-rpcallowip=127.0.0.1" \
]