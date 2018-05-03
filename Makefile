.PHONY: docker-build docker-run

docker-build:
	docker build -t historislack:latest .

docker-run:
	docker run -d \
		-v `pwd`/data:/var/lib/postgresql/data \
		-p 5432:5432 \
		historislack:latest
