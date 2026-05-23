#!/usr/bin/env bash
# knowns.sh — TowerShare project management
# Usage: ./knowns.sh <command> [args]
set -euo pipefail

KNOWNS_DIR=".knowns"
BLUE='\033[0;34m'; GREEN='\033[0;32m'; YELLOW='\033[1;33m'; RED='\033[0;31m'; NC='\033[0m'

mkdir -p "$KNOWNS_DIR"

cmd="${1:-help}"
shift || true

# ── helpers ───────────────────────────────────────────────────────────────────

_ts()  { date '+%Y-%m-%d %H:%M'; }
_slug() { echo "$*" | tr '[:upper:]' '[:lower:]' | tr -cs 'a-z0-9' '-' | sed 's/-\+/-/g; s/^-//; s/-$//'; }

# ── commands ──────────────────────────────────────────────────────────────────

cmd_up() {
    echo -e "${GREEN}Starting TowerShare...${NC}"
    docker compose up -d "$@"
    echo -e "${GREEN}PocketBase admin:${NC} http://localhost:8090/_/"
}

cmd_down() {
    echo -e "${YELLOW}Stopping TowerShare...${NC}"
    docker compose down "$@"
}

cmd_logs() {
    docker compose logs -f "${1:-pocketbase}"
}

cmd_rebuild() {
    echo -e "${YELLOW}Rebuilding backend image...${NC}"
    docker compose build --no-cache pocketbase
    docker compose up -d pocketbase
}

cmd_migrate() {
    echo -e "${BLUE}Running migrations...${NC}"
    docker compose exec pocketbase /pb/pocketbase migrate up \
        --dir=/pb/pb_data --migrationsDir=/pb/pb_migrations
}

cmd_shell() {
    docker compose exec pocketbase sh
}

cmd_pb_admin() {
    email="${1:?Usage: knowns.sh pb-admin <email> <password>}"
    pass="${2:?Usage: knowns.sh pb-admin <email> <password>}"
    # Try create first; fall back to update if the account already exists.
    docker compose exec pocketbase /pb/pocketbase admin create "$email" "$pass" --dir=/pb/pb_data 2>/dev/null || \
    docker compose exec pocketbase /pb/pocketbase admin update "$email" "$pass" --dir=/pb/pb_data
}

cmd_note() {
    # Add a project note: ./knowns.sh note "Something we learned today"
    text="${*:?Usage: knowns.sh note <text>}"
    file="$KNOWNS_DIR/notes.md"
    echo "- [$(_ts)] $text" >> "$file"
    echo -e "${GREEN}Noted.${NC}"
}

cmd_notes() {
    file="$KNOWNS_DIR/notes.md"
    [[ -f "$file" ]] && cat "$file" || echo "(no notes yet)"
}

cmd_decision() {
    # Record an architecture decision: ./knowns.sh decision "use X because Y"
    text="${*:?Usage: knowns.sh decision <text>}"
    slug="$(_slug "$text")"
    file="$KNOWNS_DIR/decisions/${slug}.md"
    mkdir -p "$KNOWNS_DIR/decisions"
    cat > "$file" <<EOF
# Decision: $text

**Date:** $(_ts)
**Status:** accepted

## Context
<!-- Why did we need to decide this? -->

## Decision
$text

## Consequences
<!-- What changes as a result? -->
EOF
    echo -e "${GREEN}Decision recorded:${NC} $file"
    ${EDITOR:-vim} "$file"
}

cmd_decisions() {
    dir="$KNOWNS_DIR/decisions"
    [[ -d "$dir" ]] || { echo "(no decisions yet)"; return; }
    echo -e "${BLUE}Architecture decisions:${NC}"
    for f in "$dir"/*.md; do
        title=$(head -1 "$f" | sed 's/# //')
        echo "  $title  ($f)"
    done
}

cmd_task() {
    action="${1:-list}"; shift || true
    case "$action" in
        add)
            text="${*:?Usage: knowns.sh task add <text>}"
            echo "- [ ] $text" >> "$KNOWNS_DIR/tasks.md"
            echo -e "${GREEN}Task added.${NC}"
            ;;
        done)
            text="${*:?Usage: knowns.sh task done <text>}"
            sed -i.bak "s|- \[ \] $text|- [x] $text|" "$KNOWNS_DIR/tasks.md" && rm -f "$KNOWNS_DIR/tasks.md.bak"
            echo -e "${GREEN}Marked done.${NC}"
            ;;
        list|*)
            file="$KNOWNS_DIR/tasks.md"
            [[ -f "$file" ]] && cat "$file" || echo "(no tasks yet)"
            ;;
    esac
}

cmd_status() {
    echo -e "${BLUE}=== TowerShare Status ===${NC}"
    docker compose ps
    echo ""
    echo -e "${BLUE}Open tasks:${NC}"
    grep '\- \[ \]' "$KNOWNS_DIR/tasks.md" 2>/dev/null || echo "  (none)"
}

cmd_env_check() {
    echo -e "${BLUE}Checking required env vars...${NC}"
    warn=0
    for var in PB_ENCRYPTION_KEY; do
        if [[ -z "${!var:-}" ]]; then
            echo -e "  ${YELLOW}WARN${NC}: $var is not set"
            warn=1
        else
            echo -e "  ${GREEN}OK${NC}: $var"
        fi
    done
    [[ $warn -eq 0 ]] && echo -e "${GREEN}All env vars present.${NC}"
}

cmd_new_module() {
    name="${1:?Usage: knowns.sh new-module <module_name>}"
    slug="$(_slug "$name")"
    dir="backend/pb_modules/$slug"
    [[ -d "$dir" ]] && { echo -e "${RED}Module already exists: $dir${NC}"; exit 1; }
    mkdir -p "$dir"
    cat > "$dir/module.go" <<GOEOF
package ${slug//-/_}

import (
	"github.com/pocketbase/pocketbase"
	"github.com/pocketbase/pocketbase/core"
)

func Register(app *pocketbase.PocketBase) {
	app.OnServe().BindFunc(func(se *core.ServeEvent) error {
		// TODO: register routes / hooks for ${name}
		return se.Next()
	})
}
GOEOF
    echo -e "${GREEN}Module scaffolded:${NC} $dir/module.go"
    echo -e "${YELLOW}Remember to import it in backend/main.go:${NC}"
    echo "  _ \"github.com/gizmotronn/towershare/pb_modules/${slug//-/_}\""
}

cmd_help() {
    cat <<EOF
${BLUE}knowns.sh${NC} — TowerShare project management

Docker
  up [args]              Start all services
  down [args]            Stop all services
  logs [service]         Tail logs (default: pocketbase)
  rebuild                Rebuild and restart pocketbase image
  shell                  Shell into pocketbase container
  pb-admin <email> <pw>  Upsert superadmin account

Database
  migrate                Run pending migrations

Development
  new-module <name>      Scaffold a new pb_module
  env-check              Verify required env vars

Project knowledge
  note <text>            Append a quick note
  notes                  List all notes
  decision <text>        Record an architecture decision (opens \$EDITOR)
  decisions              List recorded decisions
  task add <text>        Add a task
  task done <text>       Mark a task done
  task list              List tasks

  status                 Container status + open tasks
  help                   This message
EOF
}

# ── dispatch ──────────────────────────────────────────────────────────────────

case "$cmd" in
    up)            cmd_up "$@" ;;
    down)          cmd_down "$@" ;;
    logs)          cmd_logs "$@" ;;
    rebuild)       cmd_rebuild ;;
    migrate)       cmd_migrate ;;
    shell)         cmd_shell ;;
    pb-admin)      cmd_pb_admin "$@" ;;
    note)          cmd_note "$@" ;;
    notes)         cmd_notes ;;
    decision)      cmd_decision "$@" ;;
    decisions)     cmd_decisions ;;
    task)          cmd_task "$@" ;;
    status)        cmd_status ;;
    new-module)    cmd_new_module "$@" ;;
    env-check)     cmd_env_check ;;
    help|--help|-h) cmd_help ;;
    *)             echo -e "${RED}Unknown command: $cmd${NC}"; cmd_help; exit 1 ;;
esac
