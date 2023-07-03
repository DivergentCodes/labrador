// Package record is a canonical, intermediate representation of a value from a remote system.
package record

type Record struct {
	// The key for the variable.
	Key string

	// The value of the variable.
	Value string

	// Remote service that the key/value pair came from.
	Source string

	// Additional attributes about this record that might be useful.
	Metadata map[string]string
}
