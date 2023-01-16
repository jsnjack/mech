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
  - [x] http
  - [ ] icmp
  - [ ] dns
  - [ ] tcp
  - [ ] ssl cert

## DNS
 - [ ] Domain records
   - [x] A
   - [x] AAAA
   - [x] ANAME
   - [ ] CAA
   - [ ] CERT
   - [x] CNAME
   - [ ] HINFO
   - [ ] HTTP
   - [x] MX
   - [ ] NAPTR
   - [ ] NS
   - [ ] PTR
   - [ ] RP
   - [ ] SPF
   - [ ] SRV
   - [x] TXT
   - [ ] pools?

 - [x] GeoProximity
   - [ ] Renaming

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
  dns:
    surfly.gratis:
      - file4.yaml
```

> Use `mech sonar discover static -t http` command to print existing configuration

## Resource naming

Some of the resource (e.g. Sonar HTTP check ID in failover configuration) can be specified in 2 different ways:
 - ID of the resource, int
 - dynamically discovered value (e.g. `@sonar,http:test-online`). When parsing the configuration `mech` will call Constellix
   Sonar REST API and retrieve all available http checks. If one of the http checks has name `test-online`, it's ID will be
   used as `sonarCheckId`

# Resources
 - [Constellix DNS REST API v4](https://api.dns.constellix.com/v4/docs#tag/Domains)
 - [Constellix Sonar Rest API](https://api-docs.constellix.com/)
 - [Constellix go client](https://github.com/Constellix/constellix-go-client) (just for reference)
