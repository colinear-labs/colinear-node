FROM ethereum/client-go:v1.10.8

ENTRYPOINT [ "geth", \
    "--syncmode", "light", \
    "--snapshot=false", \
    "--ws", \
    "--ws.addr", "0.0.0.0", \
    "--ws.port", "5001", \
    "--ws.rpcprefix", "/wsrpc", \
    "--ws.origins", "127.0.0.1" \
]