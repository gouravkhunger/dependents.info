package env

type Environment string

const (
	EnvProduction  = Environment("production")
	EnvDevelopment = Environment("development")
)

func EnvFromString(str string) Environment {
	switch str {
	case "production":
		return EnvProduction
	default:
		return EnvDevelopment
	}
}
