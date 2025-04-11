docker-up:
	@echo "Stopping and removing existing Docker containers..."
	docker-compose down
	@echo "Starting Docker containers..."
	docker-compose up -d --build
	@echo "Tailing logs..."
	docker-compose logs -f tinyliar

docker-down:
	docker-compose down

docker-logs:
	docker-compose logs -f tinyliar

docker-rebuild:
	docker-compose down
	docker-compose build --no-cache
	docker-compose up -d

docker-ps:
	docker-compose ps
