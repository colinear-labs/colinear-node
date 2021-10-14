FROM zquestz/bitcoin-cash-node

ENTRYPOINT [ "bitcoind", \ 
    "-disablewallet", \
    "-prune=550", \
    "-maxmempool=20", \
    "-zmqpubrawtx=tcp://0.0.0.0:5002" \
]