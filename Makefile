start:
	docker-compose up -d && docker image prune -f
stop:
	docker-compose down

.PHONY: start stop