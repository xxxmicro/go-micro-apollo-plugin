package apollo

import(
	"os"
	"fmt"
	"testing"
	"github.com/micro/go-micro/v2/util/log"
	"github.com/micro/go-micro/v2/config"
	"github.com/micro/go-micro/v2/config/source"
	"github.com/micro/go-micro/v2/config/encoder/yaml"
)

type MongoConfig struct {
	Host string `json:"host"`
	Port int `json:"port"`
}

func TestApollo(t *testing.T) {
	e := yaml.NewEncoder()
	if err := config.Load(NewSource(
		WithAddress("http://apollo.dev.com:8080"),	
		WithNamespace("application"),
		WithAppId("xpay-api"),
		WithCluster("dev"),
		WithBackupConfigPath("./"),
		source.WithEncoder(e),
	)); err != nil {
    	log.Error(err)
	}


	var mc MongoConfig
	if err := config.Get("mongo").Scan(&mc); err != nil {
		log.Error(err)
	}

	fmt.Printf("host: %s\n", mc.Host)
	fmt.Printf("port: %d\n", mc.Port)

	go func() {
		for {
			w, err := config.Watch()
			if err != nil {
				log.Error(err)
			}
			// wait for next value
			v, err := w.Next()
			if err != nil {
				log.Error(err)
			}
			log.Info(v)

			var mc MongoConfig
			if err := config.Get("mongo").Scan(&mc); err != nil {
				log.Error(err)
			}

			fmt.Printf("host: %s\n", mc.Host)
			fmt.Printf("port: %d\n", mc.Port)
		}
	}()

	c := make(chan os.Signal)
	_ = <-c
	log.Info("退出")
}
