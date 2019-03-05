default: build 

build:
	go build -x -o localApp main.go

	rm -rf bin
	mkdir bin
	mkdir bin/conf
	cp localApp bin/

	rm localApp

normal:
	cp conf/normal.conf conf/app.conf
	cp conf/normal.conf bin/conf/app.conf

	./bin/localApp

extra:
	cp conf/extra.conf conf/app.conf
	cp conf/extra.conf bin/conf/app.conf

	./bin/localApp

clean:
	rm -rf bin

