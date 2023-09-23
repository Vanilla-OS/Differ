all: differ
rebuild: clean differ

differ:
	go build -o differ .

.PHONY: clean

clean:
	rm -f differ