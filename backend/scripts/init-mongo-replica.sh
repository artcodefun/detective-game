#!/usr/bin/env bash

set -euo pipefail

cd "$(dirname "$0")/.."

if [[ ! -f .env ]]; then
  echo "Missing .env file" >&2
  exit 1
fi

set -a
source .env
set +a

if [[ ! -f mongo-keyfile ]]; then
  umask 077
  openssl rand -base64 512 > mongo-keyfile
fi

chmod 400 mongo-keyfile
sudo chown 999:999 mongo-keyfile

mongo_shell=(
  docker compose exec -T mongo mongosh --quiet
  --username "$MONGO_ROOT_USERNAME"
  --password "$MONGO_ROOT_PASSWORD"
  --authenticationDatabase admin
)

docker compose up -d mongo

until "${mongo_shell[@]}" --eval 'db.runCommand({ping: 1}).ok' | grep -qx '1'; do
  sleep 1
done

replica_set="$("${mongo_shell[@]}" --eval 'try { print(rs.status().set) } catch (error) { print("") }')"
if [[ "$replica_set" != "rs0" ]]; then
  "${mongo_shell[@]}" --eval 'rs.initiate({_id: "rs0", members: [{_id: 0, host: "mongo:27017"}]})'
fi

until "${mongo_shell[@]}" --eval 'db.hello().isWritablePrimary' | grep -qx 'true'; do
  sleep 1
done

echo "MongoDB replica set rs0 is ready."
