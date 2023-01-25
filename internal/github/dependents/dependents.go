package dependents

import (
	"errors"
	"fmt"
	"io"
	"net/http"
	"os"
	"regexp"
	"strconv"
	"strings"
	"time"
)

var regRepos = regexp.MustCompile(`([\d,]+)\s+Repositories`)
var regPkgs = regexp.MustCompile(`([\d,]+)\s+Packages`)

type Dependents struct {
	Repositories int
	Packages     int
}

func DependentsByOwnerAndRepo(owner, repo string) (*Dependents, error) {
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

	d := &Dependents{0, 0}

	repoMatches := regRepos.FindAllSubmatch(body, 1)
	if len(repoMatches) > 0 {
		d.Repositories, _ = strconv.Atoi(strings.ReplaceAll(string(repoMatches[0][1]), ",", ""))
	}

	pkgMatches := regPkgs.FindAllSubmatch(body, 1)
	if len(pkgMatches) > 0 {
		d.Packages, _ = strconv.Atoi(strings.ReplaceAll(string(pkgMatches[0][1]), ",", ""))
	}

	return d, nil
}
