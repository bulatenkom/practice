package main

// ANSI escape codes
const (
	c_reset = "\033[0m"
	// 3-4 bit
	c_red    = "\033[31m"
	c_yellow = "\033[33m"
	c_bold   = "\033[1m"
	// 24 bit TrueColor
	fg_black               = "\033[38;2;0;0;0m"
	bg_med_brown           = "\033[48;2;209;139;71m"
	bg_light_peach         = "\033[48;2;255;206;158m"
	bg_cell_selected_light = "\033[48;2;168;168;240m"
	bg_cell_selected_dark  = "\033[48;2;148;112;192m"
)

func AnsiYellow(s string) string {
	return c_yellow + s + c_reset
}

func AnsiRed(s string) string {
	return c_red + s + c_reset
}
