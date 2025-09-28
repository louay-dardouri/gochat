#!/bin/sh

set -e
/usr/local/bin/chat-server &
ssh-keygen -A
exec /usr/sbin/sshd -D -e
