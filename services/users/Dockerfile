FROM node:20-alpine AS build
WORKDIR /app
RUN apk add --no-cache pnpm
COPY ./package.json ./pnpm-lock.yaml ./
RUN pnpm install --production
RUN mv ./node_modules ./prod_node_modules
RUN pnpm install
COPY . .
RUN npx tsc 

FROM node:20-alpine AS runtime
WORKDIR /root
COPY --from=build /app/dist /root/dist
COPY --from=build /app/prod_node_modules /root/node_modules
ENTRYPOINT [ "node", "/root/dist/src/index.js" ]
