package validation

import (
	"fmt"
	"strings"
)

type Error struct {
	Field  string
	Value  any
	Reason string
}

var _ error = (*Error)(nil)

func (e Error) Error() string {
	var builder strings.Builder
	fmt.Fprintf(&builder, "field %q", e.Field)
	if e.Value != nil {
		fmt.Fprintf(&builder, " value %v", e.Value)
	}
	fmt.Fprintf(&builder, ": %s", e.Reason)
	return builder.String()
}
