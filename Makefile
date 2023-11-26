all: clean differ

differ:
	go build -o differ -tags="sonic avx" .

.PHONY: clean

clean:
	rm -f differ