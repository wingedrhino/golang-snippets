package main

import "fmt"
import "strings"

func main() {
	var hello M
	hello.Set("Hello")
	var world M
	world.Set("World")
	messages := []M{hello, world}
	SayHello(messages)
}

// SayHello says Hello
func SayHello(messages []M) {
	fmt.Println(joinMessages(messages))
}

// joinMessages joins messages
func joinMessages(messages []M) string {
	var words []string
	for _, m := range messages {
		words = append(words, m.Get())
	}
	return strings.Join(words, ", ")
}

// M is a message
type M struct {
	message string
}

// Set Sets the message
func (m *M) Set(s string) {
	m.message = s
}

// Get reads the message
func (m M) Get() string {
	return m.message
}
