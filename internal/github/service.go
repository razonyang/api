package github

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"os"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/eko/gocache/lib/v4/cache"
	"github.com/eko/gocache/lib/v4/store"
	log "github.com/sirupsen/logrus"
)

var regRepos = regexp.MustCompile(`([\d,]+)\s+Repositories`)
var regPkgs = regexp.MustCompile(`([\d,]+)\s+Packages`)

type Service struct {
	cache cache.CacheInterface[string]
}

func NewService(cache cache.CacheInterface[string]) *Service {
	return &Service{
		cache: cache,
	}
}

func (s *Service) Dependents(ctx context.Context, owner, repo string) (*Dependents, error) {
	d := &Dependents{}
	cacheKey := fmt.Sprintf("github-dependents -%s-%s", owner, repo)
	cachedVal, err := s.cache.Get(ctx, cacheKey)
	if err == nil {
		err = json.Unmarshal([]byte(cachedVal), d)
	}

	if err != nil {
		log.Debugf("failed to fetch dependents from cache: %s", err)

		api := fmt.Sprintf("%s/%s/%s/network/dependents", os.Getenv("GITHUB_URL"), owner, repo)
		client := http.Client{
			Timeout: 60 * time.Second,
		}
		resp, err := client.Get(api)
		if err != nil {
			return nil, err
		}

		if resp.StatusCode != 200 {
			return nil, errors.New(resp.Status)
		}

		defer resp.Body.Close()
		body, err := io.ReadAll(resp.Body)
		if err != nil {
			return nil, err
		}

		repoMatches := regRepos.FindAllSubmatch(body, 1)
		if len(repoMatches) > 0 {
			d.Repositories, _ = strconv.Atoi(strings.ReplaceAll(string(repoMatches[0][1]), ",", ""))
		}

		pkgMatches := regPkgs.FindAllSubmatch(body, 1)
		if len(pkgMatches) > 0 {
			d.Packages, _ = strconv.Atoi(strings.ReplaceAll(string(pkgMatches[0][1]), ",", ""))
		}

		data, err := json.Marshal(d)
		if err != nil {
			log.Infof("failed to marshal dependents: %s", err)
		} else {
			if err = s.cache.Set(ctx, cacheKey, string(data), store.WithExpiration(6*time.Hour)); err != nil {
				log.Infof("failed to cache: %s", err)
			}
		}
	}

	return d, nil
}
