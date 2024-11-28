FROM cgr.dev/chainguard/go:latest as build

WORKDIR /app

COPY go.mod go.sum ./

RUN --mount=type=cache,target=/go/pkg/mod \
  --mount=type=cache,target=/root/.cache/go-build \
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
