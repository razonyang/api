# My Public APIs

| Endpoint | Version |
|---|:---:|
| https://api.razonyang.com | `v1`

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

[![GitHub Used By](https://img.shields.io/badge/dynamic/json?color=success&label=used%20by&query=repositories_humanize&logo=github&url=https%3A%2F%2Fapi.razonyang.com%2Fv1%2Fgithub%2Fdependents%2Ftwbs%2Fbootstrap)](https://github.com/twbs/bootstrap/network/dependents)
[![Hugo Used By](https://img.shields.io/badge/dynamic/json?color=success&label=used%20by&query=repositories_humanize&logo=hugo&url=https%3A%2F%2Fapi.razonyang.com%2Fv1%2Fgithub%2Fdependents%2Frazonyang%2Fhugo-theme-bootstrap)](https://github.com/razonyang/hugo-theme-bootstrap/network/dependents)

`https://img.shields.io/badge/dynamic/json?color=success&label=used-by&query=repositories_humanize&logo=github&url={endpoint}/{version}/github/dependents/{owner}/{repo}`
