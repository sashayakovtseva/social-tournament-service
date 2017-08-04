bin:
	go build -x -o st-service -ldflags "-extldflags '-static'" .
build:
	sudo docker build -t sashayakovtseva/sts:latest .
push:build
	sudo docker push sashayakovtseva/sts:latest