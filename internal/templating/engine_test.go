package templating

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestEngineAlternateSeparator(t *testing.T) {
	defaultTemplate, _ := NewDefaultTemplateWithPattern("topic*")
	engine, err := NewEngine("_", defaultTemplate, []string{
		"/ /*/*/* /measurement/origin/measurement*",
	})
	require.NoError(t, err)
	name, tags, field, err := engine.Apply("/opsagent/host01/cpu")
	require.NoError(t, err)
	require.Equal(t, "opsagent_cpu", name)
	require.Equal(t, map[string]string{
		"origin": "host01",
	}, tags)
	require.Equal(t, "", field)
}
