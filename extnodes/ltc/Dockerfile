FROM exchangeunion/litecoind

ENTRYPOINT [ "litecoind", \ 
    "-disablewallet", \
    "-prune=550", \
    "-maxmempool=20", \
    "-zmqpubsequence=tcp://0.0.0.0:5003" \
]