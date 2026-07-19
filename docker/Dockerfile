FROM node:24-alpine AS frontend

WORKDIR /build/internal/presentation/http/web

COPY internal/presentation/http/web/package.json internal/presentation/http/web/package-lock.json ./
RUN npm ci

COPY internal/presentation/http/web/index.html ./
COPY internal/presentation/http/web/vite.config.ts ./
COPY internal/presentation/http/web/tsconfig.json ./
COPY internal/presentation/http/web/tsconfig.app.json ./
COPY internal/presentation/http/web/tsconfig.node.json ./
COPY internal/presentation/http/web/public ./public
COPY internal/presentation/http/web/src ./src

RUN npm run build

FROM golang:1.26.1 AS release

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY cmd ./cmd
COPY internal ./internal
COPY --from=frontend /build/internal/presentation/http/web/dist ./internal/presentation/http/web/dist

RUN go build -o /usr/local/bin/lifeline ./cmd/app/main.go

CMD ["lifeline"]
