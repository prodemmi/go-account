package geoip

import (
	"context"
	"errors"
	"net"
	"os"
	"strings"

	"aman/internal/ports/outbound"

	"github.com/oschwald/geoip2-golang"
	"go.uber.org/fx"
)

type resolver struct {
	db *geoip2.Reader
}

type unavailable struct{}

func (unavailable) Lookup(string) (outbound.IPLocation, error) {
	return outbound.IPLocation{}, errors.New("GEOIP_DATABASE_PATH is not configured")
}

// New opens the explicitly configured local MaxMind City database. No client
// IPs are sent to an external geolocation service.
func New(lifecycle fx.Lifecycle) (outbound.GeoIPResolver, error) {
	path := strings.TrimSpace(os.Getenv("GEOIP_DATABASE_PATH"))
	if path == "" {
		return unavailable{}, nil
	}
	db, err := geoip2.Open(path)
	if err != nil {
		return nil, err
	}
	lifecycle.Append(fx.Hook{OnStop: func(context.Context) error {
		db.Close()
		return nil
	}})
	return resolver{db: db}, nil
}

func (r resolver) Lookup(rawIP string) (outbound.IPLocation, error) {
	ip := net.ParseIP(strings.TrimSpace(rawIP))
	if ip == nil {
		return outbound.IPLocation{}, errors.New("invalid client IP")
	}
	record, err := r.db.City(ip)
	if err != nil {
		return outbound.IPLocation{}, err
	}
	location := outbound.IPLocation{
		Country: record.Country.IsoCode,
		City:    record.City.Names["en"],
	}
	if len(record.Subdivisions) > 0 {
		location.Region = record.Subdivisions[0].Names["en"]
	}
	return location, nil
}

func Module() fx.Option {
	return fx.Provide(New)
}
