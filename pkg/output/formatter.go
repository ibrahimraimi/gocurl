package output

import (
	"encoding/json"
	"fmt"
	"io"
)

// Formatter defines an interface for formatting HTTP responses
type Formatter interface {
	Format(data interface{}, w io.Writer) error
}

// JSONFormatter formats data as JSON
type JSONFormatter struct {
	Pretty bool
}

// NewJSONFormatter creates a new JSONFormatter
func NewJSONFormatter(pretty bool) *JSONFormatter {
	return &JSONFormatter{
		Pretty: pretty,
	}
}

// Format formats the data as JSON and writes it to the writer
func (f *JSONFormatter) Format(data interface{}, w io.Writer) error {
	var bytes []byte
	var err error

	if f.Pretty {
		bytes, err = json.MarshalIndent(data, "", "  ")
	} else {
		bytes, err = json.Marshal(data)
	}

	if err != nil {
		return err
	}

	_, err = fmt.Fprintln(w, string(bytes))
	return err
}

// TextFormatter formats data as plain text
type TextFormatter struct{}

// NewTextFormatter creates a new TextFormatter
func NewTextFormatter() *TextFormatter {
	return &TextFormatter{}
}

// Format formats the data as text and writes it to the writer
func (f *TextFormatter) Format(data interface{}, w io.Writer) error {
	_, err := fmt.Fprintln(w, data)
	return err
}
