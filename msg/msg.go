package msg

import (
	"fmt"

	"github.com/seriiserii825/go-libs/color"
)

func Info(message string) {
	fmt.Println(color.Wrap(color.Blue, "[INFO]: "+message))
}

func Success(message string) {
	fmt.Println(color.Wrap(color.Green, "[SUCCESS]: "+message))
}
func Error(message string) {
	fmt.Println(color.Wrap(color.Red, "[ERROR]: "+message))
}
func Warning(message string) {
	fmt.Println(color.Wrap(color.Yellow, "[WARNING]: "+message))
}
