package printer

import (
	"fmt"
	"strings"

	"github.com/spf13/cobra"
)

// JSONPrinterFactory holds all configuration options for the JSONPrinter.
// It is an implementation of CustomPrinterFactory and acts as a factory for JSONPrinter
type JSONPrinterFactory struct {
	Compact bool
}

// NewJSONPrinterFactory creates new JSONPrinterFactory with the injected default values
func NewJSONPrinterFactory(compact bool) *JSONPrinterFactory {
	return &JSONPrinterFactory{Compact: compact}
}

// AddFlags will add all JSONPrinter specific flags to the cobra.Command
func (j *JSONPrinterFactory) AddFlags(cmd *cobra.Command) {
	cmd.PersistentFlags().BoolVar(&j.Compact, "compact-output", j.Compact, "Print JSON output compact")
}

// SupportedFormats returns the supported printer format that can be created by JSONPrinterFactory
func (j *JSONPrinterFactory) SupportedFormats() []string {
	return []string{"json"}
}

// CreatePrinter creates a JSONPrinter from the options set. If the format is unsupported, or it is not possible
// to create an ObjectPrinter with the current configuration it will return an error
func (j *JSONPrinterFactory) CreatePrinter(format string) (ObjectPrinter, error) {
	if err := j.validate(); err != nil {
		return nil, err
	}
	switch strings.ToLower(format) {
	case "json":
		panic("json printer implementation missing")
	default:
		return nil, fmt.Errorf("invalid output format used for JSON Printer: %q", format)
	}
}

// Validate verifies whether the current configuration can be used to create an ObjectPrinter. It will return an error
// if it is not possible
func (j *JSONPrinterFactory) validate() error {
	return nil
}