package config

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

type SpecTest struct {
	Env string `envconfig:"ENV" default:"test"`
}

func newTestClient() Client {
	return &client{}
}

func TestClient(t *testing.T) {
	t.Parallel()
	assert.NotNil(t, NewClient())
}

func TestClient_ProcessEnv(t *testing.T) {
	setEnv()
	tests := []struct {
		name   string
		prefix string
		spec   *SpecTest
		expect string
		hasErr bool
	}{
		{
			name:   "success",
			prefix: "",
			spec:   &SpecTest{},
			expect: "dev",
			hasErr: true,
		},
		{
			name:   "unexpected value",
			prefix: "",
			spec:   &SpecTest{},
			expect: "test",
			hasErr: false,
		},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			c := newTestClient()
			c.ProcessEnv(tt.prefix, tt.spec)
			assert.Equal(t, tt.hasErr, tt.expect == tt.spec.Env)
		})
	}
}

func setEnv() {
	if os.Getenv("ENV") == "" {
		os.Setenv("ENV", "dev")
	}
}
