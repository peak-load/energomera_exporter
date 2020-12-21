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

# Credits 
## Original Python code: 
* https://support.wirenboard.com/t/schityvanie-pokazanij-i-programmirovanie-elektroschetchika-energomera-se102m-po-rs-485/212                                                                                                                                               
* https://github.com/sj-asm/energomera

## Documentation / resources
* GOST-R MEK 61107-2001 (RU) https://standartgost.ru/g/ГОСТ_Р_МЭК_61107-2001
* Manufacturer website (RU) http://www.energomera.ru
* Power meter users manual (RU) http://www.energomera.ru/documentations/ce102m_full_re.pdf
* Power meter basic setup guide (RU) https://shop.energomera.kharkov.ua/DOC/ASKUE-485/meter_settings_network_RS485.pdf

# License
MIT License, see [LICENSE](https://github.com/peak-load/energomera_exporter/blob/main/LICENSE)
