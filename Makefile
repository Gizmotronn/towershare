.PHONY: up down logs rebuild migrate shell pb-admin dev-backend dev status clean help

# ── config ────────────────────────────────────────────────────────────────────
PB_ADMIN_EMAIL ?= admin@towershare.local
PB_ADMIN_PASS  ?= towershare-dev

# ── primary targets ───────────────────────────────────────────────────────────

# Start everything (first-run safe: generates .env if missing, then boots)
up: .env
	docker compose up -d
	@echo ""
	@echo "  PocketBase admin → http://localhost:8090/_/"
	@echo "  API              → http://localhost:8090/api/"
	@echo ""

# Tear down containers (keeps volumes)
down:
	docker compose down

# Rebuild the pocketbase image after Go changes, then restart
rebuild:
	docker compose build --no-cache pocketbase
	docker compose up -d pocketbase

# Tail logs (default: all; pass svc=pocketbase etc.)
logs:
	docker compose logs -f $(svc)

# Run DB migrations inside the container
migrate:
	docker compose exec pocketbase /pb/pocketbase migrate up \
		--dir=/pb/pb_data --migrationsDir=/pb/pb_migrations

# Shell into pocketbase container
shell:
	docker compose exec pocketbase sh

# Upsert superadmin (uses PB_ADMIN_EMAIL / PB_ADMIN_PASS vars above)
pb-admin:
	docker compose exec pocketbase /pb/pocketbase admin create \
		$(PB_ADMIN_EMAIL) $(PB_ADMIN_PASS) --dir=/pb/pb_data || \
	docker compose exec pocketbase /pb/pocketbase admin update \
		$(PB_ADMIN_EMAIL) $(PB_ADMIN_PASS) --dir=/pb/pb_data

# ── local dev (no Docker) ─────────────────────────────────────────────────────

# Run PocketBase directly on the host (faster iteration; requires Go 1.22+)
dev-backend:
	cd backend && go run . serve --http=0.0.0.0:8090

# Open the chat simulator (served by PocketBase at /simulator)
chat:
	open http://localhost:8090/simulator 2>/dev/null || xdg-open http://localhost:8090/simulator

# Quick curl to send a bot message from the terminal
# Usage: make send MSG="towers"
send:
	@curl -s -X POST http://localhost:8090/api/towershare/message \
		-H "Content-Type: application/json" \
		-d "{\"from\":\"+61400000001\",\"body\":\"$(MSG)\"}" | python3 -m json.tool

# ── project state ─────────────────────────────────────────────────────────────

status:
	@docker compose ps

clean:
	docker compose down -v --remove-orphans
	@echo "Volumes removed. pb_data is gone — you'll need to re-initialise."

# ── .env bootstrap ────────────────────────────────────────────────────────────
.env:
	@echo "Creating .env from .env.example..."
	cp .env.example .env
	@printf "PB_ENCRYPTION_KEY=" >> .env
	@openssl rand -hex 32 >> .env
	@echo ".env created. Review it before going to production."

# ── help ──────────────────────────────────────────────────────────────────────
help:
	@echo ""
	@echo "  make up             Boot all services (creates .env if missing)"
	@echo "  make down           Stop services"
	@echo "  make rebuild        Rebuild backend image after Go changes"
	@echo "  make logs           Tail all logs  (svc=pocketbase for one service)"
	@echo "  make migrate        Run pending DB migrations"
	@echo "  make shell          Shell into pocketbase container"
	@echo "  make pb-admin       Upsert demo superadmin"
	@echo "  make dev-backend    Run PocketBase directly (no Docker)"
	@echo "  make chat           Open local chat simulator in browser"
	@echo "  make send MSG=...   Send a test bot message via curl"
	@echo "  make status         Show container status"
	@echo "  make clean          Stop + delete volumes (destructive!)"
	@echo ""
