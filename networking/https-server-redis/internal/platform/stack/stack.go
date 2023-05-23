// Package stack contains a stack implementation for use in application
package stack

// Stack represents a Stack that can only either Push data or read all data
type Stack interface {
	Push(m string) error
	Read() ([]string, error)
}
