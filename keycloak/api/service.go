package api

import "github.com/Nerzal/gocloak/v5"

type Service struct {
	Client gocloak.GoCloak
	Jwt    *gocloak.JWT
}

func (s *Service) Token() string {
	return s.Jwt.AccessToken
}
