FROM python:3.11-alpine AS base 
RUN apk add --no-cache poetry

FROM base AS builder 
WORKDIR /app 
COPY ./pyproject.toml ./poetry.lock ./
RUN poetry config virtualenvs.in-project true
RUN poetry install --no-interaction --no-dev --no-ansi

FROM gcr.io/distroless/python3 AS runtime 
WORKDIR /root
COPY --from=builder /app/.venv ./.venv 
COPY ./src ./src
ENV VIRTUAL_ENV=/root/.venv
ENV PATH="$VIRTUAL_ENV/bin:$PATH"

ENTRYPOINT [ "python3 src/app.py" ]
