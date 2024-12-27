FROM rust:1.83-alpine AS base
WORKDIR /app
RUN apk add --no-cache build-base
COPY Cargo.toml Cargo.lock ./
RUN mkdir src && echo "// for_fetch" > src/main.rs
RUN cargo fetch

FROM base AS builder
WORKDIR /app
COPY ./src ./src
RUN cargo build --release --frozen

FROM gcr.io/distroless/cc AS runtime
WORKDIR /root
COPY --from=builder /app/target/release/challenges /root/challenges
ENTRYPOINT ["/root/challenges"]
