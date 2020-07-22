package cohttp

func truncate(raw []byte, maxLength int) (result []byte) {
	if maxLength == 0 {
		return []byte{}
	}
	length := len(raw)
	if length > maxLength {
		return append(raw[:maxLength], "..."...)
	} else {
		return raw
	}
}
