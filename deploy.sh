#!/bin/bash

make

case $(uname -s) in
    # For MacOS
    Darwin)
        sudo cp ./bin/cli /usr/local/bin/cli
        sudo cp ./bin/clid /usr/local/bin/clid
        sudo cat <<EOF | sudo tee ~/Library/LaunchAgents/cn.binacs.cli.plist
<?xml version="1.0" encoding="UTF-8"?>
<!DOCTYPE plist PUBLIC "-//Apple//DTD PLIST 1.0//EN" "http://www.apple.com/DTDs/PropertyList-1.0.dtd">
<plist version="1.0">
    <dict>
        <key>Label</key>
        <string>cn.binacs.cli</string>
        <key>ProgramArguments</key>
        <array>
            <string>/usr/local/bin/clid</string>
            <string>start</string>
        </array>
        <key>KeepAlive</key>
        <true/>
    </dict>
</plist>
EOF
        launchctl stop cn.binacs.cli
        launchctl unload ~/Library/LaunchAgents/cn.binacs.cli.plist
        chmod 644 ~/Library/LaunchAgents/cn.binacs.cli.plist
        launchctl load -w ~/Library/LaunchAgents/cn.binacs.cli.plist
        launchctl start cn.binacs.cli
        ;;
    # For Linux
    *)
        sudo cp ./bin/cli /usr/sbin/cli
        sudo cp ./bin/clid /usr/sbin/clid
        # For linux
        cat <<EOF | sudo tee /etc/systemd/system/binacs-cli.service
[Unit]
Description=binacs-cli
Documentation=https://github.com/binacs/cli
[Service]
ExecStart=/usr/sbin/clid start
Restart=on-failure
RestartSec=5
[Install]
WantedBy=multi-user.target
EOF
        sudo systemctl daemon-reload
        sudo systemctl enable binacs-cli
        sudo systemctl start binacs-cli
        ;;
esac