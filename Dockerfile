# Use base golang image from Docker Hub
FROM golang:1.16 AS build

WORKDIR /do-you-be-me

# Install dependencies in go.mod and go.sum
COPY go.mod go.sum ./
RUN go mod download

# Copy rest of the application source code
# COPY . ./
COPY cmd ./cmd
COPY internal ./internal

# Compile the application to /app.
# Skaffold passes in debug-oriented compiler flags
ARG SKAFFOLD_GO_GCFLAGS
RUN echo "Go gcflags: ${SKAFFOLD_GO_GCFLAGS}"
RUN go build -gcflags="${SKAFFOLD_GO_GCFLAGS}" -mod=readonly -v -o /app ./cmd/do-you-be-me

# Now create separate deployment image
FROM gcr.io/distroless/base

# Definition of this variable is used by 'skaffold debug' to identify a golang binary.
# Default behavior - a failure prints a stack trace for the current goroutine.
# See https://golang.org/pkg/runtime/
ENV GOTRACEBACK=single

# Copy template & assets
WORKDIR /do-you-be-me
COPY --from=build /app ./app
COPY web/template/index.gohtml ./index.gohtml
COPY web/assets ./assets
COPY internal/dybmimport/data/syllables.csv ./data/syllables.csv

ENTRYPOINT ["./app"]
