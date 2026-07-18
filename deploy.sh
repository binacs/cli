#!/bin/bash
# Installs cli/clid for the current user, no sudo required.
#
# Prefers a prebuilt binary from the latest GitHub release matching this
# machine's OS/arch; falls back to a local `go build` only if no release
# asset can be fetched and a Go toolchain is available.
set -euo pipefail

REPO="binacs/cli"
BIN_DIR="${BIN_DIR:-$HOME/.local/bin}"
SRC_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"

case "$(uname -m)" in
    x86_64) ARCH="amd64" ;;
    arm64|aarch64) ARCH="arm64" ;;
    *) echo "Unsupported architecture: $(uname -m)" >&2; exit 1 ;;
esac

case "$(uname -s)" in
    Darwin) OS="darwin" ;;
    Linux) OS="linux" ;;
    *) echo "Unsupported OS: $(uname -s)" >&2; exit 1 ;;
esac

TMP_DIR="$(mktemp -d)"
trap 'rm -rf "$TMP_DIR"' EXIT

ASSET="cli_${OS}_${ARCH}.tar.gz"
if curl -fsSL "https://github.com/${REPO}/releases/latest/download/${ASSET}" -o "$TMP_DIR/${ASSET}" 2>/dev/null; then
    echo "Downloaded prebuilt ${ASSET} from the latest release"
    tar -C "$TMP_DIR" -xzf "$TMP_DIR/${ASSET}"
elif command -v go >/dev/null 2>&1; then
    echo "No prebuilt release found; building from source with 'go build'"
    make -C "$SRC_DIR" build
    cp "$SRC_DIR/bin/cli" "$TMP_DIR/cli"
    cp "$SRC_DIR/bin/clid" "$TMP_DIR/clid"
else
    echo "Error: no prebuilt release available for ${OS}/${ARCH} and no local Go toolchain to build from source." >&2
    exit 1
fi

mkdir -p "$BIN_DIR"
install -m 755 "$TMP_DIR/cli" "$BIN_DIR/cli"
install -m 755 "$TMP_DIR/clid" "$BIN_DIR/clid"
echo "Installed cli/clid to $BIN_DIR"

case ":$PATH:" in
    *":$BIN_DIR:"*) ;;
    *) echo "Note: $BIN_DIR is not on your PATH. Add it, e.g.: export PATH=\"$BIN_DIR:\$PATH\"" ;;
esac

case "$OS" in
    darwin)
        mkdir -p "$HOME/Library/LaunchAgents"
        cat > "$HOME/Library/LaunchAgents/cn.binacs.cli.plist" <<EOF
<?xml version="1.0" encoding="UTF-8"?>
<!DOCTYPE plist PUBLIC "-//Apple//DTD PLIST 1.0//EN" "http://www.apple.com/DTDs/PropertyList-1.0.dtd">
<plist version="1.0">
    <dict>
        <key>Label</key>
        <string>cn.binacs.cli</string>
        <key>ProgramArguments</key>
        <array>
            <string>$BIN_DIR/clid</string>
            <string>start</string>
        </array>
        <key>KeepAlive</key>
        <true/>
    </dict>
</plist>
EOF
        launchctl unload "$HOME/Library/LaunchAgents/cn.binacs.cli.plist" 2>/dev/null || true
        launchctl load -w "$HOME/Library/LaunchAgents/cn.binacs.cli.plist"
        echo "clid installed as a per-user LaunchAgent (cn.binacs.cli)."
        ;;
    linux)
        mkdir -p "$HOME/.config/systemd/user"
        cat > "$HOME/.config/systemd/user/binacs-cli.service" <<EOF
[Unit]
Description=binacs-cli
Documentation=https://github.com/binacs/cli

[Service]
ExecStart=$BIN_DIR/clid start
Restart=on-failure
RestartSec=5

[Install]
WantedBy=default.target
EOF
        systemctl --user daemon-reload
        systemctl --user enable --now binacs-cli
        echo "clid installed as a user systemd service (binacs-cli)."
        echo "Note: to have it start on boot without an active login session, run: loginctl enable-linger \$USER"
        ;;
esac

echo "Done. Run '$BIN_DIR/cli --help' to get started."
