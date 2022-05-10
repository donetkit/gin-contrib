package uuid

import (
	"github.com/google/uuid"
	"strings"
)

func GoogleId() string {
	u, _ := uuid.NewRandom()
	return strings.ReplaceAll(u.String(), "-", "")
}
