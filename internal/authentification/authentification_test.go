package authentification

import "testing"

func TestBuildJWTString(t *testing.T) {
	_, err := BuildJWTString()
	if err != nil {
		t.Errorf("this is err = %d", err)
	}
}

func TestGetUserID(t *testing.T) {
	err := GetUserID("")
	if err != -1 {
		t.Errorf("this is err = %d", err)
	}
}

func BenchmarkBuildJWTString(b *testing.B) {
	for i := 0; i < b.N; i++ {
		_, _ = BuildJWTString()
	}
}
