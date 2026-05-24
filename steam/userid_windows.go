package steam

import (
	"fmt"

	"golang.org/x/sys/windows/registry"
)

func ActiveUserID() (uint32, error) {

	k, err := registry.OpenKey(registry.CURRENT_USER, `Software\Valve\Steam\ActiveProcess`, registry.QUERY_VALUE)
	if err != nil {
		return 0, fmt.Errorf("failed to open Steam registry key: %w", err)
	}
	defer k.Close() // nolint:errcheck

	p, _, err := k.GetIntegerValue("ActiveUser")
	if err != nil {
		return 0, fmt.Errorf("failed to read ActiveUser from registry: %w", err)
	}

	return uint32(p), nil
}
