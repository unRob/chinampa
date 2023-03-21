#!/usr/bin/env bash
# SPDX-License-Identifier: Apache-2.0
# Copyright © 2022 Roberto Hidalgo <chinampa@un.rob.mx>

milpa dev lint || @milpa.fail "linter found errors"
milpa dev test unit || @milpa.fail "tests failed"
