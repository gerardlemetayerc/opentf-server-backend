package iam

type OIDCClaims struct {
	Nom         string `json:"nom"`
	Prenom      string `json:"prenom"`
	Mail        string `json:"mail"`
	Groupes     string `json:"groupes"`
	DisplayName string `json:"displayName"`
}

type OIDCConfig struct {
	ID                    uint       `gorm:"primary_key" json:"id"`
	LocalEnabled          bool       `json:"localEnabled"`
	OIDCEnabled           bool       `json:"oidcEnabled"`
	Issuer                string     `json:"issuer"`
	ClientID              string     `json:"clientId"`
	ClientSecret          string     `json:"clientSecret"`
	Scopes                string     `json:"scopes"`
	TokenEndpoint         string     `json:"tokenEndpoint"`
	AuthorizationEndpoint string     `json:"authorizationEndpoint"`
	Strict                bool       `json:"strict"`
	Claims                OIDCClaims `gorm:"embedded" json:"claims"`
}
