package version

import (
	"fmt"

	"github.com/fatih/color"
)

const Major string = "0"
const Minor string = "0"
const Patch string = "1"

var Version string = fmt.Sprintf("%s.%s.%s", Major, Minor, Patch)

func PrintVersion() {
	fmt.Printf("%s - version %s\n", color.BlueString("wc"), color.GreenString(Version))
	fmt.Printf(" %s: %s\n", color.YellowString("API Versions"), color.GreenString("v1"))
}
