package consul

import (
	"testing"
)

func TestRegistrarDiscovery(t *testing.T) {
	rd := RegistrarDiscovery(_address)

	_ = rd
}
