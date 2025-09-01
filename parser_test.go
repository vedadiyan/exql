package exql

import (
	"fmt"
	"testing"
)

func TestParser(t *testing.T) {
	ctx := NewDefaultContext(WithBuiltInLibrary())
	fmt.Println(Eval("parse('https://www.abc.com')", ctx))
}
