package uuid

import (
	"fmt"

	"github.com/google/uuid"
)

func generateUUID() string {
	id, err := uuid.NewV7()
	if err != nil {
		fmt.Printf("Erro ao gerar UUID: %v\n", err)
		return ""
	}
	return id.String()
}
