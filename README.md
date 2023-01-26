# My Public APIs

| Endpoint | Version |
|---|:---:|
| https://api.razonyang.com | `v1`

The instance is launched on my small k3s (k8s) cluster, you can also [host it yourself](#self-hosted).

## APIs

### GitHub

#### GitHub Dependents

Show how many projects are using your repository.

`{endpoint}/{version}/github/dependents/{owner}/{repo}`

```json
{
  "packages": 520,
  "packages_humanize": "520",
  "repositories": 468751,
  "repositories_humanize": "468.8k"
}
```

**Integrate with Shields.io**

[![GitHub Used By](https://img.shields.io/badge/dynamic/json?color=success&label=used%20by&query=repositories_humanize&logo=github&url=https://api.razonyang.com/v1/github/dependents/twbs/bootstrap)](https://github.com/twbs/bootstrap/network/dependents)
[![Hugo Used By](https://img.shields.io/badge/dynamic/json?color=success&label=used%20by&query=repositories_humanize&logo=hugo&url=https://api.razonyang.com/v1/github/dependents/razonyang/hugo-theme-bootstrap)](https://github.com/razonyang/hugo-theme-bootstrap/network/dependents)

`https://img.shields.io/badge/dynamic/json?color=success&label=used-by&query=repositories_humanize&logo=github&url={endpoint}/{version}/github/dependents/{owner}/{repo}`

### Hugo

#### Hugo Module Info

Returns Hugo module/theme info, i.e. ![Hugo Requirements](https://img.shields.io/badge/dynamic/json?color=important&label=requirements&query=requirements&logo=hugo&url=https://api.razonyang.com/v1/hugo/modules/github.com/razonyang/hugo-mod-search).

`{endpoint}/{version}/hugo/modules/github.com/{owner}/{repo}`.

```json
{
   "hugoVersion":{
      "extended":true,
      "min":"0.99.0",
      "max":"0.111.1"
   },
   "requirements":"\u003e=0.99.0 \u003c=0.111.1 extended"
}
```

## Self-Hosted

### Requirements

- Redis: for caching.

### Environment Variables

| Env | Required | Default | Description
|:-:|:-:|:-:|---
| `PORT` | N | `8080` | The HTTP server port.
| `REDIS_ADDR` | N | `127.0.0.1:6379` | Redis address in form `host:port`.
| `REDIS_PASSWORD` | N | - | Redis password.
| `GITHUB_TOKEN` | N | - | GitHub API access token.

The `.env` file will be loaded if presents.

### Deployments

#### Build from Source

```sh
$ git clone https://github.com/razonyang/api
$ cd api
$ go build
$ ./api
```

Or via `go install`.

```sh
$ go install github.com/razonyang/api@latest
$ api
```

#### Deploy via Docker

```sh
$ docker run -p 8080:8080 razonyang/api
```
