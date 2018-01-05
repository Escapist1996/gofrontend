// Copyright 2011 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// Package vet implements the ``go vet'' command.
package vet

import (
	"cmd/go/internal/base"
	"cmd/go/internal/load"
	"cmd/go/internal/work"
)

var CmdVet = &base.Command{
	Run:         runVet,
	CustomFlags: true,
	UsageLine:   "vet [-n] [-x] [build flags] [vet flags] [packages]",
	Short:       "report likely mistakes in packages",
	Long: `
Vet runs the Go vet command on the packages named by the import paths.

For more about vet and its flags, see 'go doc cmd/vet'.
For more about specifying packages, see 'go help packages'.

The -n flag prints commands that would be executed.
The -x flag prints commands as they are executed.

The build flags supported by go vet are those that control package resolution
and execution, such as -n, -x, -v, -tags, and -toolexec.
For more about these flags, see 'go help build'.

See also: go fmt, go fix.
	`,
}

func runVet(cmd *base.Command, args []string) {
	vetFlags, pkgArgs := vetFlags(args)

	work.BuildInit()
	work.VetFlags = vetFlags

	pkgs := load.PackagesForBuild(pkgArgs)
	if len(pkgs) == 0 {
		base.Fatalf("no packages to vet")
	}

	var b work.Builder
	b.Init()

	root := &work.Action{Mode: "go vet"}
	for _, p := range pkgs {
		root.Deps = append(root.Deps, b.VetAction(work.ModeBuild, work.ModeBuild, p))
	}
	b.Do(root)
}
