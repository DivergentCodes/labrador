// Package build contains compile-time variables and defaults, overridden by LDFLAGS during build.
package core

// Build-time variables set by LDFLAGS.
var (
	AuthorName  = "Jack Sullivan"
	AuthorEmail = "jack@divergent.codes"
	Version     = "dev-snapshot"
	Commit      = "unspecified"
	Date        = "unspecified"
	BuiltBy     = "unspecified"
)
