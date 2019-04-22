Check your [N26](https://n26.com) balance.

## Build
```
cd cmd/n26balance
go get
go build
```

## Usage:
Populate a configuration file `cfg.yaml`

```
user: "n26 username"
pass: "n26 password"
```

```
cd cmd/n26balance
./n26balance --config cfg.yaml
```


