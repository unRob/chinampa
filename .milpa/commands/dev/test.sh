#!/usr/bin/env bash
# SPDX-License-Identifier: Apache-2.0
# Copyright © 2022 Roberto Hidalgo <chinampa@un.rob.mx>

set -e
if [[ "$MILPA_OPT_COVERAGE" ]]; then
  rm -rf coverage
  milpa dev test unit --coverage
  milpa dev test coverage-report
else
  milpa dev test unit
fi
