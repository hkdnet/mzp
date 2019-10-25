bin/mzp: main.go
	mkdir -p $$(basename $@)
	GO111MODULE=on go build -o $@ $^
