package main

import (
	"fmt"
	"regexp"
	"strconv"
	"time"

	domain "github.com/tushar7789/E_commerce_ms/internal/domain"
	env "github.com/tushar7789/E_commerce_ms/internal/env"
)

var durationRe = regexp.MustCompile(`(?i)^([0-9]+(?:\.[0-9]+)?)(ns|us|µs|ms|s|m|h|d|w)$`)

func Load() domain.Config {
	return domain.Config{
		AccessSecret:  env.getEnv("ACCESS_SECRET", "dev-access-secret-change-me"),
		RefreshSecret: env.getEnv("REFRESH_SECRET", "dev-refresh-secret-change-me"),
		AccessTTL:     mustDuration(env.getEnv("ACCESS_TTL", "15m")),
		RefreshTTL:    mustDuration(env.getEnv("REFRESH_TTL", "2h")),
	}
}

func mustDuration(v string) time.Duration {
	d, err := parseDurationDW(v)
	if err != nil {
		return 15 * time.Minute
	}
	return d
}

func parseDurationDW(s string) (time.Duration, error) {
	m := durationRe.FindStringSubmatch(s)
	if m == nil {
		return 0, fmt.Errorf("invalid duration: %q", s)
	}

	val, err := strconv.ParseFloat(m[1], 64)
	if err != nil {
		return 0, err
	}

	switch m[2] {
	case "ns":
		return time.Duration(val), nil
	case "us", "µs":
		return time.Duration(val * float64(time.Microsecond)), nil
	case "ms":
		return time.Duration(val * float64(time.Millisecond)), nil
	case "s":
		return time.Duration(val * float64(time.Second)), nil
	case "m":
		return time.Duration(val * float64(time.Minute)), nil
	case "h":
		return time.Duration(val * float64(time.Hour)), nil
	case "d":
		return time.Duration(val * float64(24*time.Hour)), nil
	case "w":
		return time.Duration(val * float64(7*24*time.Hour)), nil
	default:
		return 0, fmt.Errorf("invalid duration unit: %q", m[2])
	}
}