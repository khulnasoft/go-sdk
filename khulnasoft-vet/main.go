// Copyright 2020 The Khulnasoft Authors. All rights reserved.
// SPDX-License-Identifier: MIT

package main

import (
	"github.com/khulnasoft/go-sdk/khulnasoft-vet/checks"

	"golang.org/x/tools/go/analysis/unitchecker"
)

func main() {
	unitchecker.Main(
		checks.Imports,
		checks.License,
		checks.Migrations,
		checks.ModelsSession,
	)
}
