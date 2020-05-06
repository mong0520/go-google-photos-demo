APP=go-google-photos-demo
NAMESPACE := mong0520

build:
	docker build -t ${NAMESPACE}/${APP} .

dev:
	docker-compose up api

down:
	docker-compose down

push:
	@docker tag ${APP} ${NAMESPACE}/${APP}
	@docker push ${NAMESPACE}/${APP}
	@heroku container:push web

release:
	@heroku container:release web