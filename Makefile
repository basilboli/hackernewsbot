GCLOUD_PROJECT:=$(shell gcloud config list project --format="value(core.project)")

.PHONY: all
all: push

.PHONY: build
build:
	godep save
	docker build -t basilboli/hackernewsbot .

.PHONY: run
run: build
	docker-compose up

.PHONY: push
push: build
	docker push basilboli/hackernewsbot

.PHONY: clean
clean:
	docker stop `docker ps -f name=hackernewsbot --no-trunc -aq`
	docker stop `docker ps -f name=redis --no-trunc -aq`
	docker rm `docker ps -f name=hackernewsbot --no-trunc -aq`
	docker rm `docker ps -f name=redis --no-trunc -aq`

	