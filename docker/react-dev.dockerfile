FROM node:24-alpine

WORKDIR /app

COPY internal/presentation/http/web/package.json ./
COPY internal/presentation/http/web/package-lock.json ./

RUN npm install
