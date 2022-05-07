FROM scratch

COPY bin/server-linux-amd64 /server
COPY assets /assets
COPY names.json /names.json

ENTRYPOINT ["./server"]