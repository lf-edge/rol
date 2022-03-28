package utils

// CutIndexingString cut the string for indexing in db
func CutIndexingString(s string) string {
	if len(s) > 191 {
		return s[0:191]
	}
	return s
}
