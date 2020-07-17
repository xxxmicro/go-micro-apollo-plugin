package apollo

import (
	"github.com/micro/go-micro/v2/config/source"
	"github.com/micro/go-micro/v2/util/log"
	"github.com/zouyx/agollo/v3"
	"github.com/zouyx/agollo/v3/storage"
	"github.com/zouyx/agollo/v3/env/config"
	"time"
	"fmt"
)

type apolloSource struct {
	ip string
	serviceName   string
	namespaceName string
	opts          source.Options
}

func (a *apolloSource) String() string {
	return "apollo"
}

func (a *apolloSource) Read() (*source.ChangeSet, error) {
	log.Logf(fmt.Sprintf("ip: %s, namespace: %s", a.ip, a.namespaceName))
	
	readyConfig := &config.AppConfig{
		IsBackupConfig:   true,
		BackupConfigPath: "./",
		AppID:            "xpay-api",
		Cluster:          "dev",
		NamespaceName:    a.namespaceName,
		IP:               a.ip,
	}
	agollo.InitCustomConfig(func() (*config.AppConfig, error) {
		return readyConfig, nil
	})

	if err := agollo.Start(); err != nil {
		log.Error(err)
	}
	c := agollo.GetConfig(a.namespaceName)
	content := []byte(c.GetValue("content"))
	//b, err := a.opts.Encoder.Encode(content)
	//if err != nil {
	//	return nil, fmt.Errorf("error reading source: %v", err)
	//}

	cs := &source.ChangeSet{
		Timestamp: time.Now(),
		// TODO 根据 namespaceName 适配
		Format: "yaml",
		Source: a.String(),
		Data:   content,
	}
	cs.Checksum = cs.Sum()
	return cs, nil
}

func (a *apolloSource) Watch() (source.Watcher, error) {
	watcher, err := newWatcher(a.String())
	storage.AddChangeListener(watcher)
	return watcher, err
}

func (a *apolloSource) Write(cs *source.ChangeSet) error {
	return nil
}

func NewSource(opts ...source.Option) source.Source {
	options := source.NewOptions(opts...)
	var nName string
	namespaceName, ok := options.Context.Value(namespaceName{}).(string)
	if ok {
		nName = namespaceName
	}

	ip, ok := options.Context.Value("ip").(string)
	
	return &apolloSource{
		ip: ip,
		opts: options, 
		namespaceName: nName,
	}
}
