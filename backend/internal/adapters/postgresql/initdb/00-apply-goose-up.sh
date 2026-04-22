#!/bin/sh
set -eu

echo "Applying goose Up migrations from /migrations"

for file in /migrations/*.sql; do
  echo "Running migration: $(basename "$file")"
  awk '
    /^-- \+goose Up/ { in_up = 1; next }
    /^-- \+goose Down/ { in_up = 0; exit }
    /^-- \+goose StatementBegin/ { next }
    /^-- \+goose StatementEnd/ { next }
    in_up { print }
  ' "$file" | psql -v ON_ERROR_STOP=1 -U "$POSTGRES_USER" -d "$POSTGRES_DB"
done

echo "Finished applying goose Up migrations"