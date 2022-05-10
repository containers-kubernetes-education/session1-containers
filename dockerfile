# Layer we are looking to build from
FROM scratch

# Copy the compiled binary, data and assests
COPY bin/server-linux-amd64 /server
COPY assets /assets
COPY data/names.json /data/names.json

# Runtime command
ENTRYPOINT ["./server"]