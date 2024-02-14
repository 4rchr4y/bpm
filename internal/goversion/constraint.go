//go:build !go1.18
// +build !go1.18

package goversion

const MinGoVersionMessage = "This program requires Go version " + MinGoVersion + " or higher for compilation"

func init() {
	MinGoVersionMessage
}
