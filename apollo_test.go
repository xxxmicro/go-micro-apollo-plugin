package apollo

import(
	"os"
	"testing"
	"github.com/micro/go-micro/v2/util/log"
	"github.com/micro/go-micro/v2/config"
	"github.com/micro/go-micro/v2/config/source"
	"github.com/micro/go-micro/v2/config/encoder/yaml"
)

type statsConfig struct {
	Name string
}


func TestApollo(t *testing.T) {
	t.Logf("listen 1")

	e := yaml.NewEncoder()
	if err := config.Load(NewSource(
		WithIp("http://apollo-dev.dev.lucfish.com:8080"),	
		WithNamespaceName("application"),
		source.WithEncoder(e),
	)); err != nil {
    	log.Error(err)
	}

	t.Logf("listen 2")

	StatsConfig := statsConfig{}

	if err := config.Scan(&StatsConfig); err != nil {
	    log.Error(err)
	}

	t.Logf("listen config change")

	go func() {
		for {
			w, err := config.Watch("mongo.port")
			if err != nil {
				log.Error(err)
			}
			// wait for next value
			v, err := w.Next()
			if err != nil {
				log.Error(err)
			}
			if err := v.Scan(&StatsConfig); err != nil {
				log.Error(err)
			}
			// TODO
			log.Info(StatsConfig.Name)
		}
	}()

	c := make(chan os.Signal)
	_ = <-c
	t.Logf("退出")
}
