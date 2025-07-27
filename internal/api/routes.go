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
	r.GET("/api/offer_categories", GetOfferCategories)
	r.POST("/api/offer_categories", CreateOfferCategory)
	r.PUT("/api/offer_categories/:id", UpdateOfferCategory)
	r.DELETE("/api/offer_categories/:id", DeleteOfferCategory)

	// Offer endpoints
	r.GET("/api/offers", GetOffers)
	r.POST("/api/offers", CreateOffer)
	r.GET("/api/offers/:offer_id", GetOffer)
	r.PUT("/api/offers/:offer_id", UpdateOffer)
	r.DELETE("/api/offers/:offer_id", DeleteOffer)
	// SuggestedValue endpoints imbriquées par domaine
	r.GET("/api/domains/:domain_id/suggested_values", GetSuggestedValuesByDomain)
	r.POST("/api/domains/:domain_id/suggested_values", CreateSuggestedValueByDomain)
	r.PUT("/api/domains/:domain_id/suggested_values/:id", UpdateSuggestedValueByDomain)
	r.DELETE("/api/domains/:domain_id/suggested_values/:id", DeleteSuggestedValueByDomain)
	// Domain endpoints
	r.GET("/api/domains", GetDomains)
	r.POST("/api/domains", CreateDomain)
	r.GET("/api/domains/:domain_id", GetDomain)
	r.PUT("/api/domains/:domain_id", UpdateDomain)
	r.DELETE("/api/domains/:domain_id", DeleteDomain)

	// SuggestedValue endpoints
	r.GET("/api/suggested_values", GetSuggestedValues)
	r.POST("/api/suggested_values", CreateSuggestedValue)
	r.GET("/api/suggested_values/:id", GetSuggestedValue)
	r.PUT("/api/suggested_values/:id", UpdateSuggestedValue)
	r.DELETE("/api/suggested_values/:id", DeleteSuggestedValue)
	r.GET("/api/users/:id", GetUser)
	r.PUT("/api/users/:id", UpdateUser)
	r.DELETE("/api/users/:id", DeleteUser)
	// User management endpoints
	r.GET("/api/users", GetUsers)
	r.POST("/api/users", CreateUser)
	r.POST("/api/users/login", LoginUser)
	r.GET("/modules", ListModules)
	r.POST("/modules", CreateModule)
	r.GET("/modules/:module_id/properties", ListProperties)
	r.POST("/modules/:module_id/properties", CreateProperty)
	r.GET("/modules/:module_id/properties/:property_id", GetProperty)
	r.PUT("/modules/:module_id/properties/:property_id", UpdateProperty)
	r.DELETE("/modules/:module_id/properties/:property_id", DeleteProperty)
	r.GET("/instances", ListInstances)
	r.POST("/instances", CreateInstance)
	r.GET("/jobs", ListJobs)
	r.POST("/jobs", CreateJob)

	// IAM OIDC config endpoints
	r.GET("/api/iam/auth/oidc", iam.GetOIDCConfig)
	r.POST("/api/iam/auth/oidc", iam.SetOIDCConfig)

	// IAM Auth methods endpoints
	r.GET("/api/iam/auth_methods", iam.GetAuthMethods)
	r.POST("/api/iam/auth_methods", iam.SetAuthMethod)
	// Group management endpoints
	r.GET("/api/groups", GetGroups)
	r.POST("/api/groups", CreateGroup)
	r.PUT("/api/groups/:id", UpdateGroup)
	r.DELETE("/api/groups/:id", DeleteGroup)
}
