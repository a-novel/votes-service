package config

import (
	auth "github.com/a-novel/auth-service/framework"
	forum "github.com/a-novel/forum-service/framework"
	"github.com/samber/lo"
)

func GetAuthClient() auth.Client {
	authURL := lo.Ternary(ENV == ProdENV, auth.ProdURL, auth.DevURL)
	return auth.NewClient(authURL)
}

func GetForumClient() forum.Client {
	authURL := lo.Ternary(ENV == ProdENV, forum.ProdInternalURL, forum.DevInternalURL)
	return forum.NewClient(authURL)
}
