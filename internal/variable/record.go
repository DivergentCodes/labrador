// Package variable is a canonical, intermediate representation of a value from a remote system.
package variable

type Variable struct {
	// The key for the variable.
	Key string

	// The value of the variable.
	Value string

	// Remote service that the key/value pair came from.
	Source string

	// Additional attributes about this variable that might be useful.
	Metadata map[string]string
}
