.PHONY: up down logs rebuild migrate shell pb-admin dev dev-backend dev-frontend status clean help

# ── config ────────────────────────────────────────────────────────────────────
PB_ADMIN_EMAIL ?= admin@towershare.local
PB_ADMIN_PASS  ?= towershare-dev

# ── primary targets ───────────────────────────────────────────────────────────

# Start PocketBase in Docker + Vite dev server locally
up: .env
	docker compose up -d
	@echo ""
	@echo "  Frontend         → http://localhost:5173"
	@echo "  PocketBase admin → http://localhost:8090/_/"
	@echo "  API              → http://localhost:8090/api/"
	@echo ""
	yarn --cwd frontend dev

# Tear down containers (keeps volumes)
down:
	docker compose down

# Rebuild image (frontend + backend) after any changes, then restart
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

# Run both PocketBase and Vite dev server concurrently (requires Go 1.22+ and Node)
dev:
	make -j2 dev-backend dev-frontend

# Run PocketBase directly on the host (faster iteration; requires Go 1.22+)
dev-backend:
	cd backend && go run . serve --http=0.0.0.0:8090

# Run Vite dev server (proxies /api and /_ to PocketBase on :8090)
dev-frontend:
	cd frontend && yarn dev

# Open the browser chat simulator (requires PocketBase running)
chat:
	open http://localhost:8090/simulator 2>/dev/null || xdg-open http://localhost:8090/simulator

# Terminal (TUI) alternative
chat-tui:
	cd tools/simulator && go run . --url=http://localhost:8090

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
	@echo "  make up             Boot PocketBase + Caddy (frontend served from built assets)"
	@echo "  make up-dev         Boot all services incl. Vite dev server (hot-reload)"
	@echo "  make down           Stop services"
	@echo "  make rebuild        Rebuild full image (frontend + backend) after any changes"
	@echo "  make logs           Tail all logs  (svc=pocketbase for one service)"
	@echo "  make migrate        Run pending DB migrations"
	@echo "  make shell          Shell into pocketbase container"
	@echo "  make pb-admin       Upsert demo superadmin"
	@echo "  make dev            Run PocketBase + Vite locally (no Docker)"
	@echo "  make dev-backend    Run PocketBase directly (no Docker)"
	@echo "  make dev-frontend   Run Vite dev server only"
	@echo "  make chat           Open local chat simulator in browser"
	@echo "  make send MSG=...   Send a test bot message via curl"
	@echo "  make status         Show container status"
	@echo "  make clean          Stop + delete volumes (destructive!)"
	@echo ""
