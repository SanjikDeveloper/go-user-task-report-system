package config

type AuthConfig struct {
	JWTSigningKey string `env:"JWT_SIGNING_KEY"`
	PasswordSalt  string `env:"PASSWORD_SALT"`
}
