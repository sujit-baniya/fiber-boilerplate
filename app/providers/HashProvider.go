package providers

import hashing "github.com/thomasvvugt/fiber-hashing"

var hash hashing.Driver

func HashProvider() hashing.Driver {
	return hash
}

func SetHashProvider(config ...hashing.Config) {
	if len(config) > 0 {
		hash = hashing.New(config[0])
	} else {
		hash = hashing.New()
	}
}
