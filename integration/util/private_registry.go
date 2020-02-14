package util

import (
	"os"

	. "github.com/onsi/ginkgo"
)

func GetEiriniDockerHubPassword() string {
	password := os.Getenv("EIRINIUSER_PASSWORD")
	if password == "" {
		Skip("eiriniuser password not provided. Please expoert EIRINIUSER_PASSWORD")
	}
	return password
}
