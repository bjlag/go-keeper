package master_key

import (
	"context"
	"fmt"

	model "github.com/bjlag/go-keeper/internal/domain/client"
	"github.com/bjlag/go-keeper/internal/infrastructure/crypt/master_key"
)

type Usecase struct {
	tokens   tokens
	options  options
	salter   salter
	keymaker keymaker
}

func NewUsecase(tokens tokens, options options, salter salter, keymaker keymaker) *Usecase {
	return &Usecase{
		tokens:   tokens,
		options:  options,
		salter:   salter,
		keymaker: keymaker,
	}
}

func (u *Usecase) Do(ctx context.Context, data Data) error {
	const op = "usecase.master_key.Do"

	option, err := u.options.OptionBySlug(ctx, model.OptSaltKey)
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	var salt *master_key.Salt

	if option == nil {
		salt, err = u.salter.GenerateSalt()
		if err != nil {
			return fmt.Errorf("%s: %w", op, err)
		}

		err = u.options.SaveOption(ctx, model.Option{
			Slug:  model.OptSaltKey,
			Value: salt.ToString(),
		})
		if err != nil {
			return fmt.Errorf("%s: %w", op, err)
		}
	} else {
		salt, err = master_key.ParseString(option.Value)
		if err != nil {
			return fmt.Errorf("%s: %w", op, err)
		}
	}

	masterKey := u.keymaker.GenerateMasterKey([]byte(data.Password), salt.Value())

	u.tokens.SaveMasterKey(masterKey)

	return nil
}
