mech
=======

`mech` automate Constellix DNS configuration

# Supported features

> [Sonar REST API](https://api-docs.constellix.com/)

## Sonar
- [ ] static configuration
  - [x] http
  - [x] tcp
  - [ ] icmp
  - [ ] dns
  - [ ] ssl cert
- [ ] runtime data
  - [ ] http
  - [ ] icmp
  - [ ] dns
  - [ ] tcp
  - [ ] ssl cert

## DNS
 - [ ] Domains
 - [ ] Domain records
 - [ ] GeoProximity

# Configuration format
```
constellix:
  sonar:
    http_checks:
      - file1.yaml
      - file2.yaml
      - ...
    tcp_checks:
      - file3.yaml
```

> Use `mech sonar discover static -t http` command to print existing configuration

# Resources
 - [Constellix DNS REST API v4](https://api.dns.constellix.com/v4/docs#tag/Domains)
 - [Constellix Sonar Rest API](https://api-docs.constellix.com/)
 - [Constellix go client](https://github.com/Constellix/constellix-go-client) (just for reference)
