bin/mzp: main.go
	mkdir -p $$(basename $@)
	go build -o $@ $^
