package resources

import (
	"strings"
	"testing"
)

// TestHelloName calls greetings.Hello with a name, checking
// for a valid return value.
func TestIndexTemplate(t *testing.T) {
	want := "index.tmpl"
	msg := Templates().DefinedTemplates()
	if !strings.Contains(msg, want) {
		t.Fatalf(`Templates %v, want %v`, msg, want)
	}
}
