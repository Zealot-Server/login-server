package api

import (
	"strconv"
	"strings"
)

// clientUpdateRequired applies the optional application-version gate.
// An invalid minimum disables the gate so a bad setting cannot lock out players.
func clientUpdateRequired(minimum, current string) bool {
	if strings.TrimSpace(minimum) == "" {
		return false
	}

	minParts, valid := parseAppVersion(minimum)
	if !valid {
		return false
	}

	currentParts, valid := parseAppVersion(current)
	if !valid {
		return true
	}

	for index := range minParts {
		if currentParts[index] != minParts[index] {
			return currentParts[index] < minParts[index]
		}
	}
	return false
}

func parseAppVersion(version string) ([3]uint64, bool) {
	var result [3]uint64
	parts := strings.Split(strings.TrimSpace(version), ".")
	if len(parts) != len(result) {
		return result, false
	}
	for index, part := range parts {
		if part == "" {
			return result, false
		}
		value, err := strconv.ParseUint(part, 10, 64)
		if err != nil {
			return result, false
		}
		result[index] = value
	}
	return result, true
}
