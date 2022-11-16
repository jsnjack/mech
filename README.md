mech
=======

`mech` automate Constellix DNS configuration

# Supported features

> [Sonar REST API](https://api-docs.constellix.com/)

## Sonar
- [ ] static configuration
  - [x] http checks
  - [x] tcp checks

## Configuration format
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

> Use `mech sonar static` command to print existing configuration
