#!/bin/sh
# Starts clid in the background just long enough to run one `cli` command,
# then tears it down. Everything (the socket, the daemon process) lives
# only inside this container's lifetime.
set -e

clid start &
clid_pid=$!

sock="$HOME/cli.sock"
i=0
while [ ! -S "$sock" ]; do
    i=$((i + 1))
    if [ "$i" -gt 50 ]; then
        echo "clid did not come up (no socket at $sock after 10s)" >&2
        kill "$clid_pid" 2>/dev/null || true
        exit 1
    fi
    sleep 0.2
done

cli "$@"
status=$?

kill "$clid_pid" 2>/dev/null || true
exit "$status"
