package api

import (
	"github.com/gin-gonic/gin"
	"opentf-server/internal/api/iam"
)

func RegisterRoutes(r *gin.Engine) {

	// OfferProperty endpoints
	r.GET("/api/offers/:offer_id/properties", GetOfferProperties)
	r.POST("/api/offers/:offer_id/properties", CreateOfferProperty)
	r.GET("/api/offers/:offer_id/properties/:id", GetOfferProperty)
	r.PUT("/api/offers/:offer_id/properties/:id", UpdateOfferProperty)
	r.DELETE("/api/offers/:offer_id/properties/:id", DeleteOfferProperty)

	// OfferCategory endpoints
	r.GET("/api/offer_categories", JWTAuthMiddleware(), GetOfferCategories)
	r.POST("/api/offer_categories", JWTAuthMiddleware(), CreateOfferCategory)
	r.PUT("/api/offer_categories/:id", JWTAuthMiddleware(), UpdateOfferCategory)
	r.DELETE("/api/offer_categories/:id", JWTAuthMiddleware(), DeleteOfferCategory)

	// Offer endpoints
	r.GET("/api/offers", JWTAuthMiddleware(), GetOffers)
	r.POST("/api/offers", JWTAuthMiddleware(), CreateOffer)
	r.GET("/api/offers/:offer_id", JWTAuthMiddleware(), GetOffer)
	r.PUT("/api/offers/:offer_id", JWTAuthMiddleware(), UpdateOffer)
	r.DELETE("/api/offers/:offer_id", JWTAuthMiddleware(), DeleteOffer)

	// SuggestedValue endpoints imbriquées par domaine
	r.GET("/api/domains/:domain_id/suggested_values", JWTAuthMiddleware(), GetSuggestedValuesByDomain)
	r.POST("/api/domains/:domain_id/suggested_values", JWTAuthMiddleware(), CreateSuggestedValueByDomain)
	r.PUT("/api/domains/:domain_id/suggested_values/:id", JWTAuthMiddleware(), UpdateSuggestedValueByDomain)
	r.DELETE("/api/domains/:domain_id/suggested_values/:id", JWTAuthMiddleware(), DeleteSuggestedValueByDomain)
	// Domain endpoints
	r.GET("/api/domains", JWTAuthMiddleware(), GetDomains)
	r.POST("/api/domains", JWTAuthMiddleware(), CreateDomain)
	r.GET("/api/domains/:domain_id", JWTAuthMiddleware(), GetDomain)
	r.PUT("/api/domains/:domain_id", JWTAuthMiddleware(), UpdateDomain)
	r.DELETE("/api/domains/:domain_id", JWTAuthMiddleware(), DeleteDomain)

	// SuggestedValue endpoints
	r.GET("/api/suggested_values", JWTAuthMiddleware(), GetSuggestedValues)
	r.POST("/api/suggested_values", JWTAuthMiddleware(), CreateSuggestedValue)
	r.GET("/api/suggested_values/:id", JWTAuthMiddleware(), GetSuggestedValue)
	r.PUT("/api/suggested_values/:id", JWTAuthMiddleware(), UpdateSuggestedValue)
	r.DELETE("/api/suggested_values/:id", JWTAuthMiddleware(), DeleteSuggestedValue)

	// User endpoints
	r.GET("/api/users/me", JWTAuthMiddleware(), GetCurrentUser)
	r.GET("/api/users/:id", JWTAuthMiddleware(), GetUser)
	r.PUT("/api/users/:id", JWTAuthMiddleware(), UpdateUser)
	r.DELETE("/api/users/:id", JWTAuthMiddleware(), DeleteUser)
	r.GET("/api/users", JWTAuthMiddleware(), GetUsers)
	r.POST("/api/users", JWTAuthMiddleware(), CreateUser)

	// Authentication endpoints
	r.POST("/api/users/login", LoginUser)            // locale
	r.POST("/api/users/login_oidc", LoginOIDCUser)   // OIDC
	r.POST("/api/users/login_token", LoginTokenUser) // token API

	// Module endpoints
	r.GET("/modules", JWTAuthMiddleware(), ListModules)
	r.POST("/modules", JWTAuthMiddleware(), CreateModule)
	r.GET("/modules/:module_id/properties", JWTAuthMiddleware(), ListProperties)
	r.POST("/modules/:module_id/properties", JWTAuthMiddleware(), CreateProperty)
	r.GET("/modules/:module_id/properties/:property_id", JWTAuthMiddleware(), GetProperty)
	r.PUT("/modules/:module_id/properties/:property_id", JWTAuthMiddleware(), UpdateProperty)
	r.DELETE("/modules/:module_id/properties/:property_id", JWTAuthMiddleware(), DeleteProperty)
	r.GET("/api/instances", JWTAuthMiddleware(), GetInstances)
	r.POST("/api/instances", JWTAuthMiddleware(), CreateInstance)
	r.GET("/api/instances/:id", JWTAuthMiddleware(), GetInstance)
	r.PUT("/api/instances/:id", JWTAuthMiddleware(), UpdateInstance)
	r.DELETE("/api/instances/:id", JWTAuthMiddleware(), DeleteInstance)
	r.GET("/jobs", JWTAuthMiddleware(), ListJobs)
	r.POST("/jobs", JWTAuthMiddleware(), CreateJob)

	// IAM OIDC config endpoints
	r.GET("/api/iam/auth/oidc", JWTAuthMiddleware(), iam.GetOIDCConfig)
	r.POST("/api/iam/auth/oidc", JWTAuthMiddleware(), iam.SetOIDCConfig)

	// IAM Auth methods endpoints
	r.GET("/api/iam/auth_methods", iam.GetAuthMethods)
	r.POST("/api/iam/auth_methods", JWTAuthMiddleware(), iam.SetAuthMethod)

	// Group management endpoints
	r.GET("/api/groups", JWTAuthMiddleware(), GetGroups)
	r.POST("/api/groups", JWTAuthMiddleware(), CreateGroup)
	r.PUT("/api/groups/:id", JWTAuthMiddleware(), UpdateGroup)
	r.DELETE("/api/groups/:id", JWTAuthMiddleware(), DeleteGroup)

	r.GET("/api/modules/:id/archive", GetModuleArchive)
	r.POST("/api/modules/:id/update", UpdateModuleArchive)
	// ...existing code...
	r.GET("/api/modules", GetModules)
	r.POST("/api/modules", CreateModule)
	r.GET("/api/providers", GetProviders)
	r.POST("/api/providers", CreateProvider)
	// ...existing code...
	backend := r.Group("/backendapi")
	backend.POST("/statelocks", AcquireStateLock)
	backend.DELETE("/statelocks/:id", ReleaseStateLock)
	backend.GET("/statelocks/:instance_id", GetStateLock)
	// ...existing code...
	{
		backend.GET("/statefiles/:id", GetStateFile)
		backend.POST("/statefiles", CreateOrUpdateStateFile)
	}
}
