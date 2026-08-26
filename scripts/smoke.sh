#!/usr/bin/env sh
set -eu
base=${BASE_URL:-http://localhost:8080}
curl -fsS "$base/healthz" >/dev/null
body='{"tenant_id":"demo","protocol":"sum","protocol_version":"1","threshold":2,"participant_count":3,"input_commitment":"demo"}'
resp=$(curl -fsS -X POST "$base/api/v1/computations" -H 'Content-Type: application/json' -H 'Idempotency-Key: smoke-1' -d "$body")
id=$(printf '%s' "$resp" | sed -n 's/.*"id":"\([^"]*\)".*/\1/p')
test -n "$id"
round=$(curl -fsS -X POST "$base/api/v1/computations/$id/start")
rid=$(printf '%s' "$round" | sed -n 's/.*"id":"\([^"]*\)".*/\1/p')
test -n "$rid"
echo "smoke ok computation=$id round=$rid"
