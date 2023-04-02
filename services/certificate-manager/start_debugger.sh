#!/bin/sh

# Start the Delve debugger with the updated binary
dlv exec --continue --accept-multiclient --api-version=2 --headless --listen=:2345 -- /certificate-manager &

# Store the PID of the Delve debugger
DLV_PID=$!

# When the script is terminated, kill the Delve debugger process
trap "kill $DLV_PID" TERM EXIT

# Wait for the Delve debugger process to finish
wait $DLV_PID
