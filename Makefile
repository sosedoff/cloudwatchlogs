build:
	go-assets-builder static -o assets.go
	go build

docker: build
	docker build -t sosedoff/cloudwatchlogs .

docker-release: docker
	docker push sosedoff/cloudwatchlogs