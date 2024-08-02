FROM golang:1.21.6-bullseye
ARG ENTRYPOINT

# dont need this
# RUN apt-get update && apt-get install -y rsync iproute2 haproxy

# TODO: explore the debugging tool in the future
# RUN go install github.com/go-delve/delve/cmd/dlv@latest

# Compiler that is what I want
RUN go install github.com/githubnemo/CompileDaemon@latest

# TODO: explore the debugging tool in the future
# RUN go install gotest.tools/gotestsum@latest

# Hello nancy
#RUN go install github.com/golang-migrate/migrate/v4@latest
RUN go install github.com/golang-migrate/migrate/v4/cmd/migrate@latest

# EXPOSE 40000 40000

#COPY $ENTRYPOINT /entrypoint.sh
#ENTRYPOINT ["/entrypoint.sh"]
WORKDIR /work
