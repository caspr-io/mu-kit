package realm

import (
	"context"

	"github.com/caspr-io/mu-kit/keycloak/api"

	"github.com/Nerzal/gocloak/v5"
	"github.com/caspr-io/mu-kit/util"
)

//nolint:golint
type RealmAPI interface {
	Create(ctx context.Context, realmName string) (string, error)
	Exists(ctx context.Context, realmName string) bool
}

type realmAPI api.Service

func GetRealmService(s *api.Service) RealmAPI {
	return (*realmAPI)(s)
}

func (r *realmAPI) Create(ctx context.Context, realmName string) (string, error) {
	return r.Client.CreateRealm(r.Jwt.AccessToken, gocloak.RealmRepresentation{
		Realm:   util.StringP(realmName),
		Enabled: util.BoolP(true),
	})
}

func (r *realmAPI) Exists(ctx context.Context, realmName string) bool {
	realm, err := r.Client.GetRealm(r.Jwt.AccessToken, realmName)
	if err != nil {
		return false
	}

	return realm != nil
}
