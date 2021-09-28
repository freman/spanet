package wifly

import "strings"

// DefaultReplacementCharacter is the character that the WiFLY device uses
// by default to replace spaces in SSID and password
const DefaultReplacementCharacter = "$"

// AlternativeReplacementCharacters is a list of characters that may possibly
// be used if someone has put a $ in their SSID or password
var AlternativeReplacementCharacters = []string{
	"!", "\"", "'", "#", "$", "%", "&", "(",
	")", "*", "+", ",", "-", ".", "/", ":",
	";", "<", ">", "=", "?", "@", "[", "]",
	"\\", "^", "_", "`", "{", "|", "}", "~",
}

// ReplaceSpaces will replace the spaces in a SSID and password, returning the replacement character used
func ReplaceSpaces(ssid, password string) (replacement, replacedssid, replacedpassword string) {
	replacement = DefaultReplacementCharacter
	if strings.Contains(ssid, replacement) || strings.Contains(password, replacement) {
		tmp := ssid + password
		for _, ch := range AlternativeReplacementCharacters {
			if !strings.Contains(tmp, ch) {
				replacement = ch
				break
			}
		}
		if replacement == DefaultReplacementCharacter {
			panic("Unable to find viable replacement character")
		}
	}

	return replacement,
		strings.ReplaceAll(ssid, " ", replacement),
		strings.ReplaceAll(password, " ", replacement)
}
