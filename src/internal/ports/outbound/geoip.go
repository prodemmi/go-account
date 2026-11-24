package outbound

// IPLocation contains the fields used to enrich an authentication session.
type IPLocation struct {
	Country string
	City    string
	Region  string
}

// GeoIPResolver looks up an IP in the configured local geolocation database.
type GeoIPResolver interface {
	Lookup(ip string) (IPLocation, error)
}
