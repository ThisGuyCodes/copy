#!/bin/bash
set -euo pipefail

sudo chown -R vscode:vscode /mnt/mise-data
sudo chown -R vscode:vscode /mnt/go-caches

WORKSPACE_DIR="${WORKSPACE_DIR:-${PWD}}"
MISE_FILE="$WORKSPACE_DIR/mise.toml"

if [[ -f "$MISE_FILE" ]]; then
  /usr/local/bin/mise trust "$MISE_FILE"
  /usr/local/bin/mise install

  go mod download

  pre-commit install
else
  echo "WARN: $MISE_FILE not found. Skipping mise install."
fi

if command -v zsh >/dev/null 2>&1; then
  sudo chsh -s "$(command -v zsh)" "$USER" || true
fi
