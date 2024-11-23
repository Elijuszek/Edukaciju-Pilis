package color

import (
	"fmt"

	"github.com/fatih/color"
)

const (
	NONE = iota
	RED
	GREEN
	YELLOW
	BLUE
	PURPLE
)

func Format(c int, text string) string {
	switch c {
	case RED:
		return color.New(color.FgRed).SprintFunc()(text)
	case GREEN:
		return color.New(color.FgGreen).SprintFunc()(text)
	case YELLOW:
		return color.New(color.FgYellow).SprintFunc()(text)
	case BLUE:
		return color.New(color.FgBlue).SprintFunc()(text)
	case PURPLE:
		return color.New(color.FgMagenta).SprintFunc()(text)
	default:
		return text
	}
}

func main() {
	fmt.Println(Format(RED, "This is red text"))
	fmt.Println(Format(GREEN, "This is green text"))
	fmt.Println(Format(YELLOW, "This is yellow text"))
	fmt.Println(Format(BLUE, "This is blue text"))
	fmt.Println(Format(PURPLE, "This is purple text"))
	fmt.Println(Format(NONE, "This is normal text"))
}
