FROM golang:1.21.6-bullseye
ARG ENTRYPOINT

RUN apt-get update && apt-get install -y netcat

# TODO: explore the debugging tool in the future
# RUN go install github.com/go-delve/delve/cmd/dlv@latest

RUN go install github.com/githubnemo/CompileDaemon@latest

# TODO: explore the debugging tool in the future
# RUN go install gotest.tools/gotestsum@latest

RUN go install github.com/golang-migrate/migrate/v4/cmd/migrate@latest

ADD https://raw.githubusercontent.com/eficode/wait-for/v2.1.0/wait-for /usr/local/bin/wait-for
COPY $ENTRYPOINT /entrypoint.sh

RUN chmod +rx /usr/local/bin/wait-for /entrypoint.sh
