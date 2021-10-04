FROM zquestz/bitcoin-cash-node

ENTRYPOINT [ "bitcoind", \ 
    "-disablewallet", \
    "-prune=10", \
    "-maxmempool=20", \
    "-server", \
    "-rpcport=5002" \
]