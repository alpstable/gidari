package morningstar

// * This is a generated file, do not edit

// Token is an OAuth bearer token authorized using a base64-encoded user id and password.
type Token struct {
	AccessToken string `bson:"access_token" json:"access_token" sql:"access_token"`
	ExpiresIn   int    `bson:"expires_in" json:"expires_in" sql:"expires_in"`
	TokenType   string `bson:"token_type" json:"token_type" sql:"token_type"`
}
