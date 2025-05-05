package middleware

import (
	"net/http"

	"github.com/hurricanerix/http-helper/config"
	"github.com/rs/cors"
)

const defaultCORSAllowedOrigins = "*"
const defaultCORSAllowedMethods = "HEAD,GET"
const defaultCORSAllowCredentials = true

func CORS(next http.Handler) http.Handler {
	corsMiddleware := cors.New(cors.Options{
		AllowedOrigins:   config.StringSliceEnv("HH_CORS_ALLOWED_ORIGINS", defaultCORSAllowedOrigins),
		AllowedMethods:   config.StringSliceEnv("HH_CORS_ALLOWED_METHOD", defaultCORSAllowedMethods),
		AllowCredentials: config.BoolEnv("HH_CORS_ALLOW_CREDENTIALS", defaultCORSAllowCredentials),
	})

	return corsMiddleware.Handler(next)
}
