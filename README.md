# Energomera Exporter for Prometheus written on Go

This is a simple server written in Go that queries energomera smart power meters and exports data via HTTP for Prometheus consumption.

# Getting Started
To run it in foreground:

```bash
./energomera_exporter [flags]
```

Help on flags:
```bash
./energomera_exporter --help
```

To run it as service (work in progress... )

For more information check the [source code documentation](https://pkg.go.dev/github.com/peak-load/energomera_exporter). All of the core developers are accessible via the [Prometheus Developers mailinglist](https://groups.google.com/forum/?fromgroups#!forum/prometheus-developers).

# License
MIT License, see [LICENSE](https://github.com/peak-load/energomera_exporter/blob/main/LICENSE)
