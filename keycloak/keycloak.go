package keycloak

import (
	"context"

	"github.com/caspr-io/mu-kit/keycloak/api"
	"github.com/caspr-io/mu-kit/keycloak/realm"

	"github.com/Nerzal/gocloak/v5"
	"github.com/rs/zerolog/log"
)

//nolint:golint
type KeycloakConfig struct {
	URL           string `split_words:"true" required:"true"`
	AdminRealm    string `split_words:"true" required:"true" default:"master"`
	AdminUsername string `split_words:"true" required:"true" default:"keycloak"`
	AdminPassword string `split_words:"true" required:"true"`
}

type Keycloak struct {
	config  *KeycloakConfig
	service api.Service
	realm   realm.RealmAPI
}

func ConnectToKeycloak(ctx context.Context, config *KeycloakConfig) (*Keycloak, error) {
	log.Ctx(ctx).Debug().
		Str("url", config.URL).
		Str("admin-realm", config.AdminRealm).
		Str("admin-username", config.AdminUsername).
		Msg("Logging in to Keycloak")

	kc := &Keycloak{config: config}
	kc.service.Client = gocloak.NewClient(config.URL)

	token, err := kc.service.Client.LoginAdmin(config.AdminUsername, config.AdminPassword, config.AdminRealm)
	if err != nil {
		log.Ctx(ctx).Error().Err(err).Msg("Cannot connect to Keycloak")
		return nil, err
	}

	kc.service.Jwt = token
	kc.realm = realm.GetRealmService(&kc.service)

	return kc, nil
}

func (kc *Keycloak) Realm() realm.RealmAPI {
	return kc.realm
}

func (kc *Keycloak) Client() gocloak.GoCloak {
	return kc.service.Client
}

func (kc *Keycloak) Jwt() *gocloak.JWT {
	return kc.service.Jwt
}
