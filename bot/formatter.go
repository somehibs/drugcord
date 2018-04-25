package drugcord

import (
	"fmt"
	"strings"
)

// Define an interface for formattable types.
// This is good for translating from JSON into lazy maps.
type Formattable interface {
	Fields() *map[string]string
	TableFields() *map[string]map[string]string
	ComplexFields() *map[string]map[string]string
}

// Define interfaces for formatters.
type FieldFormatter interface {
	FormatFields(f Formattable) string
}

type Formatter interface {
	FormatAll(f Formattable) string
	FormatOne(f Formattable) string
	FieldFormatter(f Formattable) []string
	FormatTableFields(f Formattable) []string
	FormatComplexFields(f Formattable) []string
}

// Define some common formatters. Not all formatters have to support all interfaces.
type DiscordFormatter struct {
	Formatter
}

func missingItem() []string {
	return []string{"Could not find any fields."}
}

func (df DiscordFormatter) FormatAll(f Formattable) (ret string) {
	return strings.Join(df.FormatFields(f), "\n")
}

const kvFmt = "`%s` %s"

func checkFirst(fields *map[string]string, names []string) (s string) {
	for _, name := range names {
		if (*fields)[name] != "" {
			s += fmt.Sprintf(kvFmt, name, (*fields)[name])
			delete(*fields, name)
		}
	}
	return
}

func (df DiscordFormatter) FormatFields(f Formattable) (ret []string) {
	if f.Fields() == nil {
		return missingItem()
	}
	fields := *f.Fields()
	ret = append(ret, checkFirst(&fields, []string{"summary", "duration"}))
	tf := df.FormatTableFields(f)
	for _, v := range tf {
		ret = append(ret, v)
	}
	for k, v := range fields {
		ret = append(ret, fmt.Sprintf(kvFmt, k, v))
	}
	return
}

func (df DiscordFormatter) FormatTableFields(f Formattable) (ret []string) {
	for k, v := range *f.TableFields() {
		ret = append(ret, fmt.Sprintf("`dose %s` %s", k, v))
	}
	return
}

func (df DiscordFormatter) FormatComplexFields(f Formattable) (ret []string) {
	for k, v := range *f.ComplexFields() {
		ret = append(ret, fmt.Sprintf("%s %s ", k, v))
	}
	return
}
