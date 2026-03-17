.PHONY: generate clean

generate:
	buf dep update
	buf generate

clean:
	rm -rf gen/
