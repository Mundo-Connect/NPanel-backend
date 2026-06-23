#!/usr/bin/env sh
set -eu

ROOT_DIR=$(CDPATH= cd -- "$(dirname -- "$0")/.." && pwd)
cd "$ROOT_DIR"

COMPOSE=${DOCKER_COMPOSE:-"docker compose"}

if [ "${SKIP_GIT_PULL:-0}" != "1" ] && [ -d .git ]; then
  echo "==> Updating source tree"
  git pull --ff-only
fi

echo "==> Rebuilding NPanel backend image"
$COMPOSE build npanel

echo "==> Restarting NPanel backend container"
$COMPOSE up -d --no-deps --force-recreate npanel

echo "==> Current service status"
$COMPOSE ps

echo "==> Upgrade complete"
