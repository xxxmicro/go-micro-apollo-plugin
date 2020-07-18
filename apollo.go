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
	namespaceName string
	opts          source.Options
}

func (a *apolloSource) String() string {
	return "apollo"
}

func (a *apolloSource) Read() (*source.ChangeSet, error) {	
	c := agollo.GetConfig(a.namespaceName).GetCache()

	kv := map[string]string{}
	
	c.Range(func(key interface{}, value interface{}) bool {
		log.Info(fmt.Sprintf("%s=%s", key, string(value.([]byte)))
		
		kv[key.(string)] = string(value.([]byte))
		return true
	})
	data, err := makeMap(a.opts.Encoder, kv)
	if err != nil {
		return nil, fmt.Errorf("error reading data: %v", err)
	}

	b, err := a.opts.Encoder.Encode(data)
	if err != nil {
		return nil, fmt.Errorf("error reading source: %v", err)
	}

	cs := &source.ChangeSet{
		Timestamp: time.Now(),
		Format: a.opts.Encoder.String(),
		Source: a.String(),
		Data:   b,
	}
	cs.Checksum = cs.Sum()
	return cs, nil
}

func (a *apolloSource) Watch() (source.Watcher, error) {
	watcher, err := newWatcher(a.String(), a.opts.Encoder)
	storage.AddChangeListener(watcher)
	return watcher, err
}

func (a *apolloSource) Write(cs *source.ChangeSet) error {
	return nil
}

func NewSource(opts ...source.Option) source.Source {
	options := source.NewOptions(opts...)
	namespace, ok := options.Context.Value(namespaceName{}).(string)
	
	address, ok := options.Context.Value(addressKey{}).(string)	
	
	backupConfigPath, ok := options.Context.Value(backupConfigPathKey{}).(string)
	if !ok {
		backupConfigPath = "./config"
	}

	cluster, ok := options.Context.Value(clusterKey{}).(string)
	if !ok {
		cluster = "dev"
	}

	appId, ok := options.Context.Value(appIdKey{}).(string)

	log.Logf(fmt.Sprintf("address: %s, namespace: %s", address, namespace))

	readyConfig := &config.AppConfig{
		IsBackupConfig:   true,
		BackupConfigPath: backupConfigPath,
		AppID:            appId,
		Cluster:          cluster,
		NamespaceName:    namespace,
		IP:               address,
	}
	agollo.InitCustomConfig(func() (*config.AppConfig, error) {
		return readyConfig, nil
	})

	if err := agollo.Start(); err != nil {
		log.Error(err)
	}

	return &apolloSource{
		opts: options, 
		namespaceName: namespace,
	}
}
