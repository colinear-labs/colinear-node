FROM kylemanna/bitcoind

ENTRYPOINT [ "bitcoind", \ 
    "-disablewallet", \
    "-prune=550", \
    "-maxmempool=20", \
    "-server", \
    "-rest", \
    # auth=user:pass
    "-rpcauth=user:50c2fd98a724caf63afadca4a94ff94f\$18b83531c0eb837705dff39ad6cd5c5c012ca0d54a484f7cb0ce139667acb269", \ 
    "-rpcport=5000" \
    # "-rpcallowip=127.0.0.1", \
]