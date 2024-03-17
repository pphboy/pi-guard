package tool

import (
	"strings"

	"github.com/google/uuid"
)

func GetUUIDUpper() string {
	return strings.ToUpper(strings.ReplaceAll(uuid.NewString(), "-", ""))
}
