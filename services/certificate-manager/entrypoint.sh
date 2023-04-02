#!/bin/sh

# Start the CompileDaemon for hot-reloading
CompileDaemon --build="/app/build.sh" --command="/app/start_debugger.sh" --graceful-kill
