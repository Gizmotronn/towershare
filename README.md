# TowerShare

PocketBase backend + Swift frontend, managed with Docker Compose.

## Quick start

```bash
make up
```

That's it. On first run it creates `.env` with a generated encryption key and boots PocketBase + Caddy.

---

## Local / demo credentials

| What | Value |
|---|---|
| PocketBase admin UI | http://localhost:8090/_/ |
| Admin email | `admin@towershare.local` |
| Admin password | `towershare-dev` |
| API root | http://localhost:8090/api/ |
| Health check | http://localhost:8090/api/health |

Create the superadmin account after first boot:

```bash
make pb-admin
```

> These creds are for local development only. Change them before any shared or production deployment.

---

## All make targets

```
make up             Boot all services (creates .env if missing)
make down           Stop services
make rebuild        Rebuild backend image after Go changes
make logs           Tail all logs  (svc=pocketbase for one service)
make migrate        Run pending DB migrations
make shell          Shell into pocketbase container
make pb-admin       Upsert demo superadmin (admin@towershare.local / towershare-dev)
make dev-backend    Run PocketBase directly on host — no Docker, faster iteration
make status         Show container status
make clean          Stop + delete volumes (destructive)
```

---

## Project layout

```
TowerShare/
├── backend/
│   ├── main.go                    ← PocketBase entry point
│   ├── Dockerfile
│   ├── go.mod / go.sum
│   ├── pb_migrations/             ← versioned schema migrations (Go)
│   ├── pb_hooks/                  ← JS hooks (hot-reload, no rebuild needed)
│   └── pb_modules/                ← Go extension modules
│       └── example_module/
├── frontend/TowerShare/           ← Swift package (SPM)
│   └── Sources/TowerShare/
│       ├── Core/                  ← shared PocketBase client
│       ├── Auth/                  ← AuthService + LoginView
│       └── Features/              ← app screens
├── docker/
│   └── Caddyfile                  ← reverse proxy (HTTPS on towershare.localhost)
├── docs/
│   ├── getting-started.md
│   └── extending-pocketbase.md    ← how to add backend modules
├── Makefile
├── knowns.sh                      ← project knowledge + task CLI
└── docker-compose.yml
```

---

## Adding a backend module

```bash
./knowns.sh new-module my_feature
# wire it into backend/main.go, then:
make rebuild
```

See [`docs/extending-pocketbase.md`](docs/extending-pocketbase.md) for patterns (routes, hooks, cron, email).

---

## Swift frontend

```bash
open frontend/TowerShare/Package.swift
```

Xcode resolves the [PocketBase Swift SDK](https://github.com/pocketbase/swift-sdk) automatically. The app reads `PB_URL` from the environment (defaults to `http://localhost:8090`).

For on-device testing with local HTTPS, trust Caddy's certificate once:

```bash
docker compose exec caddy caddy trust   # macOS only
```

Then connect to `https://towershare.localhost`.
