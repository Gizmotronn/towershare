# Getting Started — TowerShare

## Prerequisites

- Docker + Docker Compose
- Xcode 15+ (for the Swift frontend)
- Go 1.22+ (only needed if editing the backend outside Docker)

---

## 1. Bootstrap

```bash
cp .env.example .env
# Generate an encryption key:
echo "PB_ENCRYPTION_KEY=$(openssl rand -hex 32)" >> .env
```

## 2. Start the backend

```bash
./knowns.sh up
```

PocketBase admin UI: **http://localhost:8090/_/**

On first run, PocketBase will prompt you to create a superadmin. You can also do it via:

```bash
./knowns.sh pb-admin admin@towershare.local changeme
```

## 3. Open the Swift project

```
open frontend/TowerShare/Package.swift
```

Xcode resolves the PocketBase Swift SDK automatically. Run the app on a simulator — it will hit `http://localhost:8090`.

For local HTTPS (required on device), Caddy runs at `https://towershare.localhost` after you trust its certificate:

```bash
# macOS only — add Caddy's root CA to your keychain
docker compose exec caddy caddy trust
```

## 4. knowns.sh cheatsheet

```bash
./knowns.sh up               # start Docker services
./knowns.sh down             # stop
./knowns.sh logs             # tail pocketbase logs
./knowns.sh rebuild          # rebuild after Go changes
./knowns.sh new-module foo   # scaffold a new backend module
./knowns.sh migrate          # run pending DB migrations
./knowns.sh pb-admin <e> <p> # upsert superadmin
./knowns.sh task add "..."   # track a task
./knowns.sh note "..."       # quick project note
./knowns.sh decision "..."   # record an architecture decision
./knowns.sh status           # containers + open tasks
```

---

## Project layout

```
TowerShare/
├── backend/
│   ├── main.go                  ← PocketBase entry point
│   ├── Dockerfile
│   ├── pb_migrations/           ← Go migration files
│   ├── pb_hooks/                ← JS hot-reload hooks
│   ├── pb_modules/              ← Go extension modules
│   │   └── example_module/
│   └── pb_data/                 ← runtime data (gitignored)
├── frontend/
│   └── TowerShare/
│       ├── Package.swift
│       └── Sources/TowerShare/
│           ├── Core/            ← shared client, config
│           ├── Auth/            ← AuthService + LoginView
│           ├── Features/        ← feature screens
│           └── Models/          ← data models
├── docker/
│   └── Caddyfile
├── docs/
│   ├── getting-started.md       ← you are here
│   └── extending-pocketbase.md
├── .knowns/                     ← local project notes (gitignored)
├── docker-compose.yml
├── knowns.sh                    ← project management CLI
└── .env.example
```
