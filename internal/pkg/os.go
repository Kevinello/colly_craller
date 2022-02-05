package pkg

import "os"

func GetEnv(key string, defaultVal string) (val string) {
	val = os.Getenv(key)
	if val == "" {
		val = defaultVal
	}
	return
}
