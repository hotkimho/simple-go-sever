FROM alpine:3.17.1

COPY ./simple-go-server /simple-go-server
RUN chmod +x /simple-go-server

ENTRYPOINT ["./simple-go-server"]
