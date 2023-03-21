#!/usr/bin/env bash
# SPDX-License-Identifier: Apache-2.0
# Copyright © 2022 Roberto Hidalgo <chinampa@un.rob.mx>

cd "$(dirname "$MILPA_COMMAND_REPO")" || @milpa.fail "could not cd into base dir"

@milpa.log info "Linting go files"
golangci-lint run || exit 2
@milpa.log complete "Go files are up to spec"

