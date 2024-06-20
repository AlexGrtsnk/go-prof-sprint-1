package appconfig

import (
	"fmt"
	"testing"
)

func TestBuildJWTString(t *testing.T) {
	a, b, f, v := ParseFlags()
	fmt.Println(a, b, f, v)
}
