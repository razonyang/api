package github

import (
	"context"
	"errors"
	"fmt"
	"io"
	"net/http"
	"os"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/eko/gocache/lib/v4/store"
	"github.com/google/go-github/v49/github"
	"github.com/razonyang/api/internal/app"
	"github.com/razonyang/api/internal/helper"
	log "github.com/sirupsen/logrus"
	"golang.org/x/mod/semver"
)

var regRepos = regexp.MustCompile(`([\d,]+)\s+Repositories`)
var regPkgs = regexp.MustCompile(`([\d,]+)\s+Packages`)

type Service struct {
	cache *app.CacheService
}

func NewService(cache *app.CacheService) *Service {
	return &Service{
		cache: cache,
	}
}

func (s *Service) Dependents(ctx context.Context, owner, repo string) (*Dependents, error) {
	d := &Dependents{}
	cacheKey := fmt.Sprintf("github-dependents -%s-%s", owner, repo)
	_, err := s.cache.Get(ctx, cacheKey, d)
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

		if err = s.cache.Set(ctx, cacheKey, d, store.WithExpiration(time.Hour)); err != nil {
			log.Infof("failed to cache: %s", err)
		}
	}

	return d, nil
}

func (s *Service) Tag(ctx context.Context, owner, repo, prefix string) (*Tag, error) {
	t := &Tag{}
	cacheKey := fmt.Sprintf("github-tag-%s-%s-%s", owner, repo, prefix)
	_, err := s.cache.Get(ctx, cacheKey, t)
	if err != nil {
		log.Debugf("failed to fetch tag from cache: %s", err)
		tags, err := s.allTags(ctx, owner, repo)
		if err != nil {
			return nil, err
		}
		if len(tags) == 0 {
			return nil, fmt.Errorf("no tags found")
		}

		if prefix != "" {
			for i := len(tags) - 1; i >= 0; i-- {
				if strings.HasPrefix(tags[i], prefix) {
					t.Name = tags[i][len(prefix):]
					break
				}
			}
		} else {
			t.Name = tags[len(tags)-1]
		}
		if err = s.cache.Set(ctx, cacheKey, t, store.WithExpiration(15*time.Minute)); err != nil {
			log.Errorf("failed to cache tag: %s", err)
		}
	}

	return t, nil
}

func (s *Service) allTags(ctx context.Context, owner, repo string) (tags []string, err error) {
	cacheKey := fmt.Sprintf("github-tags-%s-%s", owner, repo)
	_, err = s.cache.Get(ctx, cacheKey, &tags)
	if err != nil {
		log.Debugf("failed to fetch tags from cache: %s", err)

		client := helper.NewGitHubClent(ctx)
		for page := 1; true; page++ {
			items, _, err := client.Repositories.ListTags(ctx, owner, repo, &github.ListOptions{
				Page:    page,
				PerPage: 100,
			})
			if err != nil {
				log.Warn(err)
				break
			}

			if len(items) == 0 {
				break
			}

			for _, tag := range items {
				tags = append(tags, *tag.Name)
			}
		}

		semver.Sort(tags)
		if err = s.cache.Set(ctx, cacheKey, tags, store.WithExpiration(15*time.Minute)); err != nil {
			log.Errorf("failed to cache tags: %s", err)
		}
		fmt.Println(tags)
	}

	return
}
