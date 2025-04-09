# Learning Bosh

Followed the [Bosh docs](https://bosh.io/docs/create-release/) to create a simple bosh release.

I vendored in the package for golang using this: https://github.com/cloudfoundry/bosh-package-golang-release

I added it as an instance group to an existing CloudFoundry Bosh deployment, like this:
```
instance_groups:
- azs:
  - z1
  instances: 1
  jobs:
  - name: server
    properties:
      port: 8080
    release: clay
  - name: route_registrar
    properties:
      nats:
        tls:
          client_cert: ((nats_client_cert.certificate))
          client_key: ((nats_client_cert.private_key))
          enabled: true
      route_registrar:
        routes:
        - name: server-route
          port: 8080
          registration_interval: 10s
          uris:
          - bosh-proxy.SYSTEM_DOMAIN
    release: routing
  - name: cf-cli-8-linux
    release: cf-cli
  name: clay-server
  networks:
  - name: default
  stemcell: default
  vm_type: minimal
```

The route registrar stuff was optional, I just did that so I could hit my server through Gorouter.
