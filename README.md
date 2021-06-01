# Backstage Software Catalog Performance testing service

This service sole purpose is to answer the catalog with `catalog-info.yaml` files and locations.

`curl http://localhost:9191/foo/bar1337/blob/master/catalog-info.yaml` responds with an entity that has the same path as the "repo" in the url.
The description is the current timestamp so that it later can be used in the catalog to detect if and how quickly new changes are picked up.

```yaml
apiVersion: backstage.io/v1alpha1
kind: Component
metadata:
  description: "2021-05-26T13:43:50+02:00"
  name: bar1337
spec:
  lifecycle: experimental
  owner: wAcbaAolwo
  type: website
```

`curl http://localhost:9191/locations/42/catalog-info.yaml` produces 42 locations.

### Adding 1000 locations to the catalog

Adding a whole lot of entities can either be done by adding each location manually using bash.

```bash
for i in {1..1000}; do curl -X POST -H 'Content-Type: application/json' localhost:7000/api/catalog/locations -d '{"type": "url", "target": "http://localhost:9191/foo/'bar$i'/blob/master/catalog-info.yaml"}'; done
```

Or by adding a location that references to 1000 other locations.

```bash
curl -X POST -H 'Content-Type: application/json' localhost:7000/api/catalog/locations -d '{"type": "url", "target": "http://localhost:9191/locations/1000/catalog-info.yaml"}'
```

### Running/Building/Deploying

```bash
# requirements
brew install golang
# run locally
go run cmd/catman/main.go
# build for current architecture.
go build ./cmd/catman
# build for linux
GOOS=linux go build ./cmd/catman

# running on a specific port
./catman -port 9191
# setting base url, needed if this is suppose to be deployed on something else than localhost
# to produce valid urls for generated locations
./catman -baseurl http://mydomain.example.com
```

### Hand crafted code to benchmark individual methods in the catalog

See `utils.ts` for producing lightweight metrics.

```javascript
const timer = createTimer("some work");
await work();
timer.end();
```

### Permit catalog to read entities from catman

The catalog does not permit reading from unknown locations so the hostname of catman needs to be appended to the allow list.

```yaml
backend:
  # ... other config...
  reading:
    allow:
      - host: example.com
      - host: "*.mozilla.org"
      - host: localhost:9191
```
