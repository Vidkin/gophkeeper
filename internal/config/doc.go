// Package config provides configuration management for the agent, including
// parsing command-line flags and environment variables.
//
// server_address.go defines the `ServerAddress` struct, which holds the host and port
// information for a server. It provides functionality to initialize a server address
// with default values, convert the address to a string representation, and set the
// address from a formatted string.
//
// interval.go defines the Interval type, which represents a time interval
// in seconds. It includes a custom JSON unmarshalling method to parse
// interval strings that are expected to have a suffix of "s" (for seconds).
package config
