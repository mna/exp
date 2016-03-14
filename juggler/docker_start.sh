#!/bin/sh

# Wait for redis to be available (container named `redis`, port `6379`).
while ! nc -z redis 6379; do sleep 1; done

# Exec, so that docker signals (sent to PID 1) are sent to the correct process.
exec "$@"

