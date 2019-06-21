docker:
	docker build -t sosedoff/cloudwatchlogs .

docker-release: docker
	docker push sosedoff/cloudwatchlogs