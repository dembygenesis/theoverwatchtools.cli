FROM golang:1.21.6-bullseye
ARG ENTRYPOINT

# Install necessary packages including netcat
RUN apt-get update && apt-get install -y netcat

# TODO: explore the debugging tool in the future
# RUN go install github.com/go-delve/delve/cmd/dlv@latest

# Compiler that is what I want
RUN go install github.com/githubnemo/CompileDaemon@latest

# TODO: explore the debugging tool in the future
# RUN go install gotest.tools/gotestsum@latest

# Hello nancy
#RUN go install github.com/golang-migrate/migrate/v4@latest

# This is so we can run migrations commands well?? HMmm..
# so I can run both inside, and outside? What's easier for us...
# whats easier is running just inside so it's completely isolated
# Okay so my vision is... I can run the migrations inside the container,
# and then I can run the server outside the container
# So I can run the migrations inside the container, and then I can run the server outside the container
RUN go install github.com/golang-migrate/migrate/v4/cmd/migrate@latest

# Not sure if I can expose my fucking port ahaha
# EXPOSE 40000 40000

ADD https://raw.githubusercontent.com/eficode/wait-for/v2.1.0/wait-for /usr/local/bin/wait-for
COPY $ENTRYPOINT /entrypoint.sh

RUN chmod +rx /usr/local/bin/wait-for /entrypoint.sh
# ENTRYPOINT ["/entrypoint.sh"]
# WORKDIR /work
