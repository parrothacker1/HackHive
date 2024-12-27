FROM node:20-alpine AS base
WORKDIR /app
RUN apk add --no-cache pnpm
COPY ./package.json ./pnpm-lock.yaml ./
RUN pnpm install --production
COPY ./src ./src
ENTRYPOINT [ "pnpm", "start" ]
