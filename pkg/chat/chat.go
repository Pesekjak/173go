package chat

import "strings"

const ColorSymbol rune = '§'

const (
	Black       string = "§0"
	DarkBlue           = "§1"
	DarkGreen          = "§2"
	DarkAqua           = "§3"
	DarkRed            = "§4"
	DarkPurple         = "§5"
	Gold               = "§6"
	Gray               = "§7"
	DarkGray           = "§8"
	Blue               = "§9"
	Green              = "§a"
	Aqua               = "§b"
	Red                = "§c"
	LightPurple        = "§d"
	Yellow             = "§e"
	White              = "§f"

	Magic         = "§k"
	Bold          = "§l"
	Strikethrough = "§m"
	Underline     = "§n"
	Italic        = "§o"

	Reset = "§r"
)

func TranslateColorCodes(c rune, msg string) string {
	cs := string(c)
	replacer := strings.NewReplacer(
		cs+"0", Black,
		cs+"1", DarkBlue,
		cs+"2", DarkGreen,
		cs+"3", DarkAqua,
		cs+"4", DarkRed,
		cs+"5", DarkPurple,
		cs+"6", Gold,
		cs+"7", Gray,
		cs+"8", DarkGray,
		cs+"9", Blue,
		cs+"a", Green,
		cs+"b", Aqua,
		cs+"c", Red,
		cs+"d", LightPurple,
		cs+"e", Yellow,
		cs+"f", White,
		cs+"k", Magic,
		cs+"l", Bold,
		cs+"m", Strikethrough,
		cs+"n", Underline,
		cs+"o", Italic,
		cs+"r", Reset,
	)
	return replacer.Replace(msg)
}

var stripReplacer = strings.NewReplacer(
	Black, "",
	DarkBlue, "",
	DarkGreen, "",
	DarkAqua, "",
	DarkRed, "",
	DarkPurple, "",
	Gold, "",
	Gray, "",
	DarkGray, "",
	Blue, "",
	Green, "",
	Aqua, "",
	Red, "",
	LightPurple, "",
	Yellow, "",
	White, "",
	Magic, "",
	Bold, "",
	Strikethrough, "",
	Underline, "",
	Italic, "",
	Reset, "",
)

func StripColorCodes(msg string) string {
	return stripReplacer.Replace(msg)
}

// IsColor checks if the given string constant is a color code.
func IsColor(code string) bool {
	switch code {
	case Black, DarkBlue, DarkGreen, DarkAqua, DarkRed, DarkPurple, Gold, Gray, DarkGray, Blue, Green, Aqua, Red,
		LightPurple, Yellow, White:
		return true
	default:
		return false
	}
}

// IsStyle checks if the given string constant is a style code.
func IsStyle(code string) bool {
	switch code {
	case Magic, Bold, Strikethrough, Underline, Italic:
		return true
	default:
		return false
	}
}
