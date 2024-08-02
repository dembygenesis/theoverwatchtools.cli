FROM golang:1.21.6-bullseye
ARG ENTRYPOINT

# Install CompileDaemon and golang-migrate
RUN go install github.com/githubnemo/CompileDaemon@latest
# RUN go install github.com/golang-migrate/migrate/v4

# Expose ports
EXPOSE 40000

# Copy the entrypoint script
COPY $ENTRYPOINT /entrypoint.sh

# Set the entrypoint and work directory
# ENTRYPOINT ["/entrypoint.sh"]
WORKDIR /work

# Use a simple command to keep the container running
CMD ["tail", "-f", "/dev/null"]
