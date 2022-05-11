package uuid

import (
	"strings"
)

func NewUUID() string {
	uuid, err := GenerateGoogleUUID()
	if err == nil {
		return strings.ReplaceAll(uuid, "-", "")
	}
	uuid, _ = GenerateUUID()
	return uuid
}
