package handlers

import (
	"api-barebone/handlers/bbAuth"
	"api-barebone/types"
)

// Routes the application's routes
var Routes = types.Route{
	SubRoutes: map[string]types.Route{
		"auth": types.Route{
			Post: bbAuth.AuthorizationGrantTypePassword,
			SubRoutes: map[string]types.Route{
				"logout": types.Route{
					Get:  bbAuth.Logout,
					Post: bbAuth.Logout,
				},
			},
		},
	},
}
