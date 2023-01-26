package hugo

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"path"
	"time"

	"github.com/eko/gocache/lib/v4/store"
	"github.com/google/go-github/v49/github"
	"github.com/pelletier/go-toml/v2"
	"github.com/razonyang/api/internal/app"
	log "github.com/sirupsen/logrus"
	"golang.org/x/oauth2"
	"gopkg.in/yaml.v3"
)

var configFilenames = []string{"hugo.toml", "hugo.yaml", "hugo.json", "config.toml", "config.yaml", "config.json"}

type Service struct {
	cache *app.CacheService
}

func NewService(cache *app.CacheService) *Service {
	return &Service{
		cache: cache,
	}
}

func (service *Service) newClient(ctx context.Context) *github.Client {
	token := os.Getenv("GITHUB_TOKEN")
	if token == "" {
		return github.NewClient(nil)
	}

	ts := oauth2.StaticTokenSource(
		&oauth2.Token{AccessToken: token},
	)
	tc := oauth2.NewClient(ctx, ts)

	return github.NewClient(tc)
}

func (service *Service) Config(ctx context.Context, vendor, owner, repo string) (*Config, error) {
	cfg := &Config{}
	cacheKey := fmt.Sprintf("hugo-module-%s-%s-%s-config", vendor, owner, repo)
	_, err := service.cache.Get(ctx, cacheKey, cfg)
	if err != nil {
		log.Debugf("failed to fetch hugo module config from cache: %s", err)

		client := service.newClient(ctx)

		for _, filename := range configFilenames {
			f, _, _, err := client.Repositories.GetContents(ctx, owner, repo, filename, nil)
			if err != nil {
				log.Debug(err)
				continue
			}
			log.Debugf("found hugo module config file: %s", filename)
			content, err := f.GetContent()
			if err != nil {
				log.Panic("failed to get content: %s", err)
			}
			data := []byte(content)
			switch path.Ext(filename) {
			case ".toml":
				err = toml.Unmarshal(data, cfg)
			case ".json":
				err = json.Unmarshal(data, cfg)
			case ".yaml":
				err = yaml.Unmarshal(data, cfg)
			}
			log.Debug(cfg)

			if err != nil {
				log.Panic("failed to parse config: %s", err)
			}

			if err = service.cache.Set(ctx, cacheKey, cfg, store.WithExpiration(time.Hour)); err != nil {
				log.Infof("failed to cache hugo module config: %s", err)
			}

			break
		}
	}

	return cfg, nil
}
