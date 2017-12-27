# safepath

* License: MIT (see LICENSE file)
* Current state: stable and maintained (Dec 2017)
* Godoc: https://godoc.org/github.com/fxkr/safepath

A library for safely handling absolute and relative paths. This helps prevents bugs (that could for example lead to path traversal attacks) at compile time.

The idea is that if you get a path as a string, you wrap it in a safepath.Path or safepath.AbsolutePath immediately, and convert it back to a string only at the last possible moment.

Note: This library is UNIX specific. It does not handle Windows paths and the Windows-specific special cases at all.
