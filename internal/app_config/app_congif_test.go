package appconfig

import (
	"fmt"
	"testing"
)

func TestBuildJWTString(t *testing.T) {
	a, b, f, v, s, c, t1 := ParseFlags()
	fmt.Println(a, b, f, v, s, c, t1)
}
