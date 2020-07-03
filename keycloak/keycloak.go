package keycloak

import (
	"context"

	"github.com/caspr-io/mu-kit/util"

	"github.com/Nerzal/gocloak/v5"
	"github.com/rs/zerolog/log"
)

type KeycloakConfig struct {
	URL           string `split_words:"true" required:"true"`
	AdminRealm    string `split_words:"true" required:"true" default:"master"`
	AdminUsername string `split_words:"true" required:"true" default:"keycloak"`
	AdminPassword string `split_words:"true" required:"true"`
}

type Keycloak struct {
	config *KeycloakConfig
	Client gocloak.GoCloak
	Jwt    *gocloak.JWT
}

func ConnectToKeycloak(ctx context.Context, config *KeycloakConfig) (*Keycloak, error) {
	log.Ctx(ctx).Debug().
		Str("url", config.URL).
		Str("admin-realm", config.AdminRealm).
		Str("admin-username", config.AdminUsername).
		Msg("Logging in to Keycloak")

	client := gocloak.NewClient(config.URL)

	token, err := client.LoginAdmin(config.AdminUsername, config.AdminPassword, config.AdminRealm)
	if err != nil {
		log.Ctx(ctx).Error().Err(err).Msg("Cannot connect to Keycloak")
		return nil, err
	}

	return &Keycloak{
		config: config,
		Client: client,
		Jwt:    token,
	}, nil
}

func (kc *Keycloak) NewRealm(realmName string) (string, error) {
	return kc.Client.CreateRealm(kc.Token(), gocloak.RealmRepresentation{
		Realm: util.StringP(realmName),
	})
}

func (kc *Keycloak) ExistsRealm(realmName string) bool {
	realm, err := kc.Client.GetRealm(kc.Token(), realmName)
	if err != nil {
		return false
	}

	return realm != nil
}

func (kc *Keycloak) Token() string {
	return kc.Jwt.AccessToken
}
