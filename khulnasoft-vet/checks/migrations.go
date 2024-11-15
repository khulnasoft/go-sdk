// Copyright 2020 The Khulnasoft Authors. All rights reserved.
// SPDX-License-Identifier: MIT

package checks

import (
	"bufio"
	"errors"
	"os"
	"os/exec"
	"strings"

	"golang.org/x/tools/go/analysis"
)

var Migrations = &analysis.Analyzer{
	Name: "migrations",
	Doc:  "check migrations for black-listed packages.",
	Run:  checkMigrations,
}

var (
	migrationDepBlockList = []string{
		"github.com/khulnasoft/khulnasoft/models",
	}
	migrationImpBlockList = []string{
		"github.com/khulnasoft/khulnasoft/modules/structs",
	}
	migrationImpPrefix = "github.com/khulnasoft/khulnasoft/models/migrations/v"
	migrationFile      = "migrations.go"
)

func checkMigrations(pass *analysis.Pass) (interface{}, error) {
	if !strings.EqualFold(pass.Pkg.Path(), "github.com/khulnasoft/khulnasoft/models/migrations") {
		return nil, nil
	}

	if _, err := exec.LookPath("go"); err != nil {
		return nil, errors.New("go was not found in the PATH")
	}

	depsCmd := exec.Command("go", "list", "-f", `{{join .Deps "\n"}}`, "github.com/khulnasoft/khulnasoft/models/migrations")
	depsOut, err := depsCmd.Output()
	if err != nil {
		return nil, err
	}

	deps := strings.Split(string(depsOut), "\n")
	for _, dep := range deps {
		if stringInSlice(dep, migrationDepBlockList) {
			pass.Reportf(0, "github.com/khulnasoft/khulnasoft/models/migrations cannot depend on the following packages: %s", migrationDepBlockList)
			return nil, nil
		}
	}

	impsCmd := exec.Command("go", "list", "-f", `{{join .Imports "\n"}}`, "github.com/khulnasoft/khulnasoft/models/migrations")
	impsOut, err := impsCmd.Output()
	if err != nil {
		return nil, err
	}

	imps := strings.Split(string(impsOut), "\n")
	migrationCount := 0
	for _, imp := range imps {
		if stringInSlice(imp, migrationImpBlockList) {
			pass.Reportf(0, "github.com/khulnasoft/khulnasoft/models/migrations cannot import the following packages: %s", migrationImpBlockList)
			return nil, nil
		}

		if strings.HasPrefix(imp, migrationImpPrefix) {
			goFilesCmd := exec.Command("go", "list", "-f", `{{join .GoFiles "\n"}}`, imp)
			goFilesOut, err := goFilesCmd.Output()
			if err != nil {
				return nil, err
			}
			// the last item is empty we need to ignore it
			migrationCount += len(strings.Split(string(goFilesOut), "\n")) - 1
		}
	}

	mf, err := os.Open(migrationFile)
	if err != nil {
		return nil, err
	}
	defer mf.Close()

	l := bufio.NewScanner(mf)
	for l.Scan() {
		if strings.Contains(l.Text(), "NewMigration(\"") {
			migrationCount--
		}
	}
	if migrationCount != 0 {
		pass.Reportf(0, "migration files count does not match migrations lengths in %s", migrationFile)
	}
	return nil, nil
}

func stringInSlice(needle string, haystack []string) bool {
	for _, h := range haystack {
		if strings.EqualFold(needle, h) {
			return true
		}
	}
	return false
}
