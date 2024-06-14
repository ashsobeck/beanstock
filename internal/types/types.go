package types

type Website struct {
	Id           string                 `bson:"id,omitempty"`
	Url          string                 `bson:"url,omitempty"`
	ShopProvider string                 `bson:"shop_provider,omitempty"`
	Json         map[string]interface{} `bson:"site_json,omitempty"`
	LastHash     int                    `bson:"last_hash,omitempty"`
}
