package validation

import (
	"context"

	"github.com/compose-spec/compose-go/v2/loader"
	"github.com/compose-spec/compose-go/v2/types"
	u "github.com/quollix/common/utils"
)

const composeConfigConsistencyCheckFailed = "compose config consistency check failed"

type ComposeSyntaxValidator interface {
	ValidateComposeSyntax(composeFileBytes []byte) error
}

type ComposeSyntaxValidatorImpl struct{}

func NewComposeSyntaxValidator() ComposeSyntaxValidator {
	return &ComposeSyntaxValidatorImpl{}
}

func (v *ComposeSyntaxValidatorImpl) ValidateComposeSyntax(composeFileBytes []byte) error {
	configDetails := types.ConfigDetails{
		ConfigFiles: []types.ConfigFile{
			{
				Filename: "docker-compose.yml",
				Content:  composeFileBytes,
			},
		},
	}

	_, err := loader.LoadModelWithContext(context.Background(), configDetails, func(opts *loader.Options) {
		opts.SkipInterpolation = true
		opts.SkipInclude = true
		opts.SkipExtends = true
		opts.SkipResolveEnvironment = true
		opts.ResolvePaths = false
		opts.SetProjectName("compose_check", true)
	})
	if err != nil {
		return u.Logger.NewError(composeConfigConsistencyCheckFailed, "compose_error_message", err.Error())
	}
	return nil
}
