TARGET=energomera-exporter

all: energomera_exporter.go
	go build -o $(TARGET)

clean:
	go clean
	rm -f $(TARGET)
