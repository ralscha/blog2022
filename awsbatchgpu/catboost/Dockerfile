FROM nvidia/cuda:12.9.1-runtime-ubuntu24.04
ARG DEBIAN_FRONTEND=noninteractive

ENV PYTHONDONTWRITEBYTECODE=1
ENV PYTHONUNBUFFERED=1

RUN apt-get update && apt-get install -y --no-install-recommends \
    python3 \
    python3-pip \
    python3-venv \
    curl \
    && rm -rf /var/lib/apt/lists/*

RUN curl -sSL https://install.python-poetry.org | python3 -
ENV PATH="${PATH}:/root/.local/bin"

WORKDIR /app
COPY pyproject.toml poetry.lock /app/
RUN poetry install --no-interaction --no-ansi
COPY demo /app/demo

CMD ["poetry", "run", "python", "/app/demo"]
