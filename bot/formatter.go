package drugcord

// Define an interface for formattable types.
type Formattable interface {
	Fields() map[string]string
	TableFields() []map[string]string
	ComplexFields() []map[string]map[string]string
}

// Define an interface for formatters.
type FieldFormatter interface {
	FormatFields(f *Formattable) string
}

type Formatter interface {
	FieldFormatter
	FormatTableFields(f *Formattable) string
	FormatComplexFields(f *Formattable) string
}

// Define some common formatters. Not all formatters have to support all interfaces.
