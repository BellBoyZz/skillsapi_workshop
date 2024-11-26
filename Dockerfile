FROM cgr.dev/chainguard/go:latest as build

WORKDIR /app

RUN --mount=type=cache,target=/go/pkg/mod/ \
    --mount=type=bind,source=./go.sum,target=go.sum \
    --mount=type=bind,source=./go.mod,target=go.mod \
    go mod download

COPY . .

RUN go build \
    -ldflags "-linkmode external -extldflags -static" \
    -o api

FROM cgr.dev/chainguard/static:latest

LABEL version="X.Y.Z"

COPY --from=build /app/api .

EXPOSE ${PORT}

CMD ["/api"]
