all: build sankey

build:
	go build -o sankey

sankey: build
	cat access.log | ./sankey && python plot.py

clean:
	rm sankey data.csv
