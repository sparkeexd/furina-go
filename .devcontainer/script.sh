#!/bin/bash

# Set dev container directory as safe to fix dubious ownership warning.
if ! git config --global --get-all safe.directory | grep -Fxq "$0"; then
    git config --global --add safe.directory "$0"
fi