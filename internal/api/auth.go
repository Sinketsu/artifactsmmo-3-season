package api

import (
	"context"

	"github.com/Sinketsu/artifactsmmo-3-season/gen/oas"
)

type Auth struct {
	Token string
}

func (a *Auth) HTTPBasic(_ context.Context, _ string) (oas.HTTPBasic, error) {
	return oas.HTTPBasic{}, nil
}

func (a *Auth) JWTBearer(_ context.Context, _ string) (oas.JWTBearer, error) {
	return oas.JWTBearer{Token: a.Token}, nil
}
