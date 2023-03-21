#!/usr/bin/env bash
# SPDX-License-Identifier: Apache-2.0
# Copyright © 2022 Roberto Hidalgo <chinampa@un.rob.mx>

root="$(dirname "$MILPA_COMMAND_REPO")"
cd "$root" || @milpa.fail "could not cd into $root"
@milpa.log info "Running unit tests"
args=()

if [[ "${MILPA_OPT_COVERAGE}" ]]; then
  cover_dir="$root/coverage/unit"
  rm -rf "$cover_dir"
  mkdir -p "$cover_dir"
  args=( -test.gocoverdir="$cover_dir" --coverpkg=./... )
fi
gotestsum --format short -- ./... "${args[@]}" || exit 2
@milpa.log complete "Unit tests passed"
