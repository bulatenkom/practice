package main

import "fmt"

// ANSI escape codes
const (
	c_reset = "\033[0m"
	// 3-4 bit
	c_red    = "\033[31m"
	c_yellow = "\033[33m"
	c_bold   = "\033[1m"
	// 24 bit TrueColor
	fg_black       = "\033[38;2;0;0;0m"
	bg_med_brown   = "\033[48;2;209;139;71m"
	bg_light_peach = "\033[48;2;255;206;158m"
)

func fgTrueColor(r, g, b int) string {
	return fmt.Sprintf("\033[38;2;%v;%v;%vm", r, g, b)
}

func bgTrueColor(r, g, b int) string {
	return fmt.Sprintf("\033[48;2;%v;%v;%vm", r, g, b)
}

func ansiYellow(s string) string {
	return c_yellow + s + c_reset
}

func ansiRed(s string) string {
	return c_red + s + c_reset
}
