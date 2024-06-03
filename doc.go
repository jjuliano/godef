// Copyright 2009 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

/*
Godef prints the source location of definitions in Go programs.

Usage:

	godef [-t] [-a] [-A] [-o offset] [-i] [-f file] [-acme] [expr]

Options:

	-f file
		Specifies the source file in which to evaluate expr.
	-o offset
		If expr is not given, specifies a location within file, which should be within or adjacent to an identifier or field selector.
	-t
		Prints the type of the expression.
	-a
		Prints all the public members (fields and methods) of the expression, along with their locations.
	-A
		Prints private members as well as public members.
	-i
		Reads the source from standard input; file must still be specified to locate other files in the same source package.
	-acme
		Reads the offset, file name, and contents from the current acme window.

Expr must be an identifier or a Go expression terminated with a field selector.

Example:

	$ cd $GOROOT
	$ godef -f src/pkg/xml/read.go 'NewParser().Skip'
	src/pkg/xml/read.go:384:18
	$
*/
package main
