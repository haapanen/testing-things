package main

import (
	"github.com/gin-gonic/gin"
	"github.com/kelseyhightower/envconfig"
)

type AppConfig struct {
	CaCertPath string `envconfig:"CA_CERT_PATH" required:"true"`
	CaKeyPath  string `envconfig:"CA_KEY_PATH" required:"true"`
}

var appConfig AppConfig
var keyGenerator *KeyGenerator

func main() {
	envconfig.MustProcess("", &appConfig)

	keyGenerator = NewKeyGenerator(appConfig.CaCertPath, appConfig.CaKeyPath)

	r := gin.Default()

	r.POST("/create-keys-and-certificates", createKeysAndCertificatesHandler)

	r.Run(":8080")
}

func createKeysAndCertificatesHandler(c *gin.Context) {
	kp, err := keyGenerator.GenerateKeysAndCertificates()
	if err != nil {
		c.JSON(500, gin.H{
			"error": err.Error(),
		})
		return
	}

	c.JSON(200, kp)
}
