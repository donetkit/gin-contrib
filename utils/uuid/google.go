package uuid

import (
	"github.com/google/uuid"
)

func GenerateGoogleUUID() (string, error) {
	u, err := uuid.NewRandom()
	if err != nil {
		return "", err
	}
	return u.String(), nil
}
