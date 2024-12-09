FROM golang:latest AS build_section
WORKDIR /src

# dependencies download
RUN --mount=type=cache,target=/go/pkg/mod/ \
    --mount=type=bind,source=go.sum,target=go.sum \
    --mount=type=bind,source=go.mod,target=go.mod \
    go mod download -x

# Build
RUN --mount=type=cache,target=/go/pkg/mod/ \
    --mount=type=bind,target=. \
    CGO_ENABLED=0 go build -o /usr/local/bin/cli ./main.go

FROM alpine:latest AS final

# Copy the executable from the "build" stage.
COPY --from=build_section /usr/local/bin/cli /usr/local/bin

# Expose the port that the application listens on.
#EXPOSE 8080

# What the container should run when it is started.
ENTRYPOINT [ "/usr/local/bin/cli" ]
