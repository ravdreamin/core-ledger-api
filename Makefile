setup:
	@echo "🚀 Starting Infrastructure..."
	docker-compose up -d
	@echo "⏳ Waiting for Database to be ready..."
	@sleep 5
	@echo "🏗️ Running Migrations..."
	migrate -path ./migrations -database "postgres://admin:gk123456@localhost:5432/core-ledger_db?sslmode=disable" up
	@echo "👤 Seeding Admin User..."
	docker exec -it $$(docker ps -qf "name=db") psql -U admin -d core-ledger_db -c "INSERT INTO users (username, password_hash, role) VALUES ('admin', '\$$2a\$$10\$$N9qo8uLOickgx2ZMRZoMyeIjZAgcfl7p92ldGxad68LJZdL17lhWy', 'admin') ON CONFLICT DO NOTHING;"
	@echo "✅ Setup Complete. Type 'make run' to start the API."

run:
	go run main.go