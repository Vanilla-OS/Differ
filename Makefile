all: clean differ

differ:
	go build -o differ -tags="sonic avx" .

test:
	go test -v ./.../...

.PHONY: clean

clean:
	rm -f differ
