package validation

import (
	"testing"

	"github.com/quollix/common/assert"
	u "github.com/quollix/common/utils"
)

func TestValidateComposeSyntax_Success(t *testing.T) {
	composeSyntaxValidator := NewComposeSyntaxValidator()
	err := composeSyntaxValidator.ValidateComposeSyntax([]byte(`
services:
  app:
    image: nginx:1.27
`))
	assert.Nil(t, err)
}

func TestValidateComposeSyntax_Fail(t *testing.T) {
	composeSyntaxValidator := NewComposeSyntaxValidator()
	err := composeSyntaxValidator.ValidateComposeSyntax([]byte("hello"))
	assert.Equal(t, composeConfigConsistencyCheckFailed, u.ExtractError(err))
}

// We want validation to work completely in-memory, so we assert that the system does not try to access these non-existing files.
func TestValidateComposeSyntax_DoesNotLoadReferencedFiles(t *testing.T) {
	composeSyntaxValidator := NewComposeSyntaxValidator()
	err := composeSyntaxValidator.ValidateComposeSyntax([]byte(`
include:
  - missing-compose.yml
services:
  app:
    extends:
      file: missing-base.yml
      service: base
    image: nginx:1.27
`))
	assert.Nil(t, err)
}
