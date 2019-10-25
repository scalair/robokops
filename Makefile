build:
	go build

install:
	mkdir -p /etc/robokops
	sudo cp bom.yaml /etc/robokops/bom.yaml
	sudo cp robokops /usr/local/bin/robokops

clean:
	rm robokops