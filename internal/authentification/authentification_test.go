package authentification

import "testing"

func TestBuildJWTString(t *testing.T) {
	_, err := BuildJWTString()
	if err != nil {
		t.Errorf("IntMin(2, -2) = %d; want -2", err)
	}
}
