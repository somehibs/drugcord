package drugcord

import (
	"fmt"
	"strings"
)

// Define an interface for formattable types.
// This is good for translating from JSON into lazy maps.
type Formattable interface {
	Fields() *map[string]string
	TableFields() *map[string]map[string]map[string]string
	MultipleFields() *map[string][]string
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
			if s != "" {
				s += "\n"
			}
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
	complexFields := df.FormatComplexFields(f)
	ret = append(ret, checkFirst(&fields, []string{"summary", "duration"}))
	tf := df.FormatTableFields(f)
	for _, v := range tf {
		ret = append(ret, v)
	}
	for k, v := range fields {
		ret = append(ret, fmt.Sprintf(kvFmt, k, v))
	}
	ret = append(ret, complexFields...)
	return
}

type Column struct {
	name string
	size int
}

var allTypes = []string{"Threshold", "Light", "Common", "Strong", "Heavy"}
var dxmTypes = []string{"First-Plateau", "Second-Plateau", "Third-Plateau", "Fourth"}

func FormatOneTable(name string, tables *map[string]map[string]string) (ret []string) {
	columns := []Column{Column{name, 0}}
	largestType := 0
	types := map[string]bool{}
	for k, v := range *tables {
		// k is title
		largestItem := 0
		for dataType, dataItem := range v {
			types[dataType] = true
			typeLen := len(dataType)
			itemLen := len(dataItem)
			if typeLen > largestType {
				largestType = typeLen
			}
			if itemLen > largestItem {
				largestItem = itemLen
			}
		}
		if len(k) > largestItem {
			columns = append(columns, Column{name: k, size: len(k) + 2})
		} else {
			columns = append(columns, Column{name: k, size: largestItem + 2})
		}
	}
	columns[0].size = largestType + 2
	fmt.Printf("%s\n", largestType)
	fmt.Printf("%s\n", columns)
	oneLine := "```"
	// Now, output all the column headers
	for _, column := range columns {
		oneLine += getColumnC(column)
	}
	oneLine += "|"
	ret = append(ret, oneLine)
	// Sweet, output a breaker line
	oneLine = ""
	for _, column := range columns {
		oneLine += getLinePiece(column.size)
	}
	oneLine += "|"
	ret = append(ret, oneLine)
	// Iterate over all the maps of values and output them
	for _, t := range allTypes {
		oneLine = ""
		dataFound := false
		for i, column := range columns {
			m := (*tables)[column.name]
			if i == 0 {
				oneLine += getColumn(t, columns[i].size)
			} else if types[t] == true {
				// Fetch rows
				oneLine += getColumn(m[t], columns[i].size)
				dataFound = true
			}
		}
		if dataFound {
			oneLine += "|"
			ret = append(ret, oneLine)
		}
	}
	oneLine = "```"
	ret = append(ret, oneLine)
	return
}

func getLinePiece(size int) (r string) {
	r = "|"
	for i := 0; i < size; i += 1 {
		r += "-"
	}
	return
}

func getColumnC(column Column) string {
	return getColumn(column.name, column.size)
}
func getColumn(name string, size int) string {
	if name == "" {
		name = "-"
	}
	if len(name) < size {
		sizeLess := size - len(name)
		prefix := sizeLess / 2
		for i := 0; i < prefix; i += 1 {
			name = " " + name
		}
		postfix := sizeLess - prefix
		for i := 0; i < postfix; i += 1 {
			name = name + " "
		}
	}
	return fmt.Sprintf("|%s", name)
}

func (df DiscordFormatter) FormatTableFields(f Formattable) (ret []string) {
	for k, v := range *f.TableFields() {
		ret = append(ret, FormatOneTable(k, &v)...)
	}
	return
}

func trim(s string) string {
	return s[:len(s)-2]
}

func (df DiscordFormatter) FormatComplexFields(f Formattable) (ret []string) {
	for k, v := range *f.MultipleFields() {
		sep := ", "
		vt := trim(strings.Join(v, sep))
		//if k == "categories" {
		//	// Format like that
		//	sep = "http://drugs.tripsit.me/category/"
		//	vt = ""
		//	for _, s := range v {
		//		vt += sep + s + ", "
		//	}
		//	vt = trim(vt)
		//}
		// Format like this
		ret = append(ret, fmt.Sprintf("`%s` %s ", k, vt))
	}
	return
}
