package log

import (
	"fmt"
	"os"

	"github.com/gookit/color"
)

type prettyPrinter struct{}

const (
	BLOOD = "#c41f3e"
	RED   = "#e63758"
)

var Pretty = prettyPrinter{}

func (p prettyPrinter) Error(message string) {
	fmt.Fprint(os.Stderr, color.HEX(RED).Sprintf("Error: %s\n", message))
}

func (p prettyPrinter) Errorf(template string, more ...any) {
	p.Error(fmt.Sprintf(template, more...))
}

func (p prettyPrinter) Error1(message string) {
	p.Error(message)
	os.Exit(1)
}

func (p prettyPrinter) Errorf1(template string, more ...any) {
	p.Errorf(template, more...)
	os.Exit(1)
}

func (p prettyPrinter) Fatal(message string) {
	fmt.Fprint(os.Stderr, color.HEX(BLOOD).Sprintf("boom 💥, %s\n", message))
	os.Exit(1)
}

func (p prettyPrinter) Fatalf(template string, more ...any) {
	p.Fatal(fmt.Sprintf(template, more...))
}
