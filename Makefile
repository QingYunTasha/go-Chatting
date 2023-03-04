web-build:
	docker build -t web:latest -f Dockerfile/web.dockerfile .

web-push:
	docker push web:latest

web-all: web-build web-push