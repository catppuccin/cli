package schema_test

import (
	"testing"

	"github.com/catppuccin/cli/schema"
	"github.com/catppuccin/cli/testdata"

	"github.com/matryer/is"
)

func TestValidSchema(t *testing.T) {
	assert := is.New(t)

	fi, err := testdata.Schemas.Open("schema/helix.catppuccin.yaml")
	assert.NoErr(err) // Should open file

	err = schema.Lint(fi)
	assert.NoErr(err) // Helix schema should be valid
}

func TestInvalidSchema(t *testing.T) {
	assert := is.New(t)

	fi, err := testdata.Schemas.Open("schema/vscode.catppuccin.yaml")
	assert.NoErr(err) // Should open file

	err = schema.Lint(fi)
	assert.True(err != nil) // VSCode schema should be invalid; has an extra property
}
