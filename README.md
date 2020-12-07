# kgen

## Install

- Install from source

```sh
make install
```

## Usage

### Annotate the pod that needs to be scraped

- `cm.example.com/scrape` (required)
- `cm.example.com/port` (optional, Prometheus scrapes container serving port by default)
- `cm.example.com/path` (optional, Prometheus scrapes `/metrics` by default)

Example:

```sh
POD=some_pod
# create
kubectl annotate --overwrite pods $POD 'cm.example.com/scrape'='true' 'cm.example.com/port'='9990'
# remove
kubectl annotate --overwrite pods $POD 'cm.example.com/scrape-' 'cm.example.com/port-'
```

### Annotate the node that needs to be scrapes

- `cm.example.com/scrape` (required)

Example:

```sh
NODE=some_node
# create
kubectl annotate --overwrite nodes $NODE 'cm.example.com/scrape'='true'
# remove
kubectl annotate --overwrite nodes $NODE 'cm.example.com/scrape-'
```
