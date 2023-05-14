package conf

import (
	"fmt"
	"testing"
)

func TestClient(t *testing.T) {
	Init("./")
	fmt.Println(GetClientDiscoveryAddr())
}
