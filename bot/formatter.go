package drugcord

import (
	"fmt"
)

// Define an interface for formattable types.
// This is good for translating from JSON into lazy maps.
type Formattable interface {
	Fields() *map[string]string
	TableFields() map[string]map[string]string
	ComplexFields() *map[string]map[string]string
}

// Define interfaces for formatters.
type FieldFormatter interface {
	FormatFields(f Formattable) string
}

type Formatter interface {
	FormatAll(f Formattable) string
	FieldFormatter(f Formattable) string
	FormatTableFields(f Formattable) string
	FormatComplexFields(f Formattable) string
}

// Define some common formatters. Not all formatters have to support all interfaces.
type DiscordFormatter struct {
}

func (df *DiscordFormatter) FormatFields(f Formattable) (ret string) {
	for k, v := range *f.Fields() {
		ret += fmt.Sprintf("%s %s", k, v)
	}
	return
}
