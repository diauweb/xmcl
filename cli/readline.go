package cli

import (
	"fmt"

	"github.com/gookit/color"
)

func Ask(prompt string) string {
	color.Style{color.Blue, color.OpBold}.Printf("%s: ", prompt)
	var input string
	fmt.Scanln(&input)
	return input
}
