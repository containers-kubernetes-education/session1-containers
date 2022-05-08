FROM scratch

COPY bin/server-linux-amd64 /server
COPY assets /assets
COPY data/names.json /data/names.json

ENTRYPOINT ["./server"]