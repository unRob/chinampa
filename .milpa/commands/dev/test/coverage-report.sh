#!/usr/bin/env bash
# SPDX-License-Identifier: Apache-2.0
# Copyright © 2022 Roberto Hidalgo <chinampa@un.rob.mx>

runs=()
while IFS=  read -r -d $'\0'; do
    runs+=("$REPLY")
done < <(find coverage -type d -maxdepth 1 -mindepth 1 -print0)
packages="$(IFS=, ; echo "${runs[*]}")"


@milpa.log info "Building coverage report from runs: ${runs[*]}"
go tool covdata textfmt -i="$packages" -o coverage/coverage.cov || @milpa.fail "could not merge runs"
go tool cover -html=coverage/coverage.cov -o coverage/coverage.html || @milpa.fail "could not build reports"

@milpa.log complete "Coverage report built"
go tool covdata percent -i="$packages"
go tool cover -func=coverage/coverage.cov | tail -n 1
