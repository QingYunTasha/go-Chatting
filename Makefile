web-build:
	docker build -t web:latest -f Dockerfile/web.dockerfile .

web-push:
	docker push web:latest

web-benchmark:
	k6 run benchmark/web/index.js

web-all: web-build web-push