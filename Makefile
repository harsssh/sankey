.PHONY: all
all: build sankey

.PHONY: build
build:
	go build -o sankey

.PHONY: sankey
sankey: build
	cat log/access.log | ./sankey && python plot.py

.PHONY: log
log:
	scp isu-1:/var/log/nginx/access.log log/

.PHONY: clean
clean:
	rm sankey data.csv
