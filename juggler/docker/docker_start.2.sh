#!/bin/sh

# Wait for redis to be available (containers named `redis{p,c}`, port `6379`).
while ! nc -z redisp 6379; do sleep 1; done
while ! nc -z redisc 6379; do sleep 1; done

# Exec, so that docker signals (sent to PID 1) are sent to the correct process.
exec "$@"

