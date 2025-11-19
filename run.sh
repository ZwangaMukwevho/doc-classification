#!/usr/bin/env bash

# ========= CONFIG =========
APP="./my-go-binary"      # <-- change to your compiled Go binary path (e.g. ./main)
OUTPUT_LOG="./output.log" # Logs from the app's stdout/stderr
APP_LOG="./app.log"       # Logs from the supervisor (restarts, start, stop, etc.)
SLEEP_SECONDS=60          # Wait time before restart in seconds
# ==========================

set -u  # safer: error on unset vars (no `-e` so we can handle restart logic)

timestamp() {
  date +"%Y-%m-%d %H:%M:%S"
}

log_app() {
  # Log supervisor (this script) events
  echo "$(timestamp) [APP] $*" | tee -a "$APP_LOG" >&2
}

init_logs() {
  # Ensure log files exist (don’t truncate; just create if missing)
  touch "$OUTPUT_LOG" "$APP_LOG"
  log_app "Log files initialized: OUTPUT_LOG=$OUTPUT_LOG, APP_LOG=$APP_LOG"
}

run_app() {
  log_app "Starting application: $APP"

  # Run the binary, timestamping each line of its output into OUTPUT_LOG
  # PIPESTATUS[0] gives the exit code of $APP in the pipeline
  "$APP" "$@" 2>&1 | while IFS= read -r line; do
    echo "$(timestamp) [OUT] $line"
  done >> "$OUTPUT_LOG"

  local app_exit=${PIPESTATUS[0]}
  log_app "Application exited with code $app_exit"
  return "$app_exit"
}

main_loop() {
  log_app "Supervisor started."

  while true; do
    run_app "$@"
    exit_code=$?

    if [ "$exit_code" -eq 0 ]; then
      log_app "Application exited normally (code 0). Not restarting. Supervisor exiting."
      break
    else
      log_app "Application crashed (code $exit_code). Restarting in ${SLEEP_SECONDS}s..."
      sleep "$SLEEP_SECONDS"
    fi
  done
}

### Entry point
init_logs
main_loop "$@"
