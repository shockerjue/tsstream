default: build

build:
	go build -x -o a.out main.go

	rm -rf bin
	mkdir bin
	mkdir bin/conf
	cp a.out bin/

	rm a.out

normal:
	cp conf/normal.conf conf/app.conf
	cp conf/normal.conf bin/conf/app.conf

	./bin/a.out

extra:
	cp conf/extra.conf conf/app.conf
	cp conf/extra.conf bin/conf/app.conf

	./bin/a.out

moni:
	cp conf/monitor.conf conf/app.conf
	cp conf/monitor.conf bin/conf/app.conf

	./bin/a.out

install:
	cp bin/a.out demo/server/normal/bin/normal
	cp bin/a.out demo/server/extra1/bin/extra1
	cp bin/a.out demo/server/extra2/bin/extra2
	cp bin/a.out demo/monitor/bin/monitor
	cp monitor/monitor.html demo/monitor/

	cp -rf demo/server /opt/
	cp -rf demo/client /opt/
	cp -rf demo/push /opt/
	cp -rf demo/monitor /opt/
	
uninstall:
	rm demo/server/normal/bin/normal
	rm demo/server/extra1/bin/extra1
	rm demo/server/extra2/bin/extra2
	rm demo/monitor/bin/monitor

clean:
	rm -rf bin
	rm demo/server/normal/bin/normal
	rm demo/server/extra1/bin/extra1
	rm demo/server/extra2/bin/extra2
