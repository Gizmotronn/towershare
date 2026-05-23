package messaging

import (
	"fmt"
	"strings"

	"github.com/pocketbase/pocketbase"
)

// handleMessage is the brain of the bot. Add commands here.
func handleMessage(app *pocketbase.PocketBase, msg incomingMessage) string {
	cmd, args := parseCommand(msg.Body)

	switch cmd {
	case "help", "":
		return helpText()

	case "towers":
		return listTowers(app)

	case "join":
		if args == "" {
			return "Usage: join <tower name>\nExample: join Sunset Tower"
		}
		return joinTower(app, msg.From, args)

	case "listings":
		return listListings(app, args)

	case "ping":
		return "pong 🏓"

	default:
		return fmt.Sprintf("Unknown command: %q\n\n%s", cmd, helpText())
	}
}

func helpText() string {
	return `TowerShare Bot 🏢

Commands:
  towers           — list all towers
  join <name>      — request access to a tower
  listings [tower] — see active listings
  ping             — health check
  help             — this message`
}

func parseCommand(body string) (cmd, args string) {
	parts := strings.SplitN(strings.ToLower(strings.TrimSpace(body)), " ", 2)
	cmd = parts[0]
	if len(parts) > 1 {
		args = strings.TrimSpace(parts[1])
	}
	return
}

func listTowers(app *pocketbase.PocketBase) string {
	records, err := app.Dao().FindRecordsByFilter("towers", "id != ''", "name", 20, 0, nil)
	if err != nil || len(records) == 0 {
		return "No towers registered yet."
	}
	var sb strings.Builder
	sb.WriteString("Towers:\n")
	for _, r := range records {
		sb.WriteString(fmt.Sprintf("  • %s\n", r.GetString("name")))
	}
	return strings.TrimRight(sb.String(), "\n")
}

func joinTower(app *pocketbase.PocketBase, from, towerName string) string {
	records, err := app.Dao().FindRecordsByFilter(
		"towers",
		fmt.Sprintf("name ~ '%s'", strings.ReplaceAll(towerName, "'", "''")),
		"name", 5, 0, nil,
	)
	if err != nil || len(records) == 0 {
		return fmt.Sprintf("No tower found matching %q.\nSend 'towers' to see all towers.", towerName)
	}
	tower := records[0]
	return fmt.Sprintf(
		"Access request noted for %q (from %s).\nA tower manager will be in touch shortly.",
		tower.GetString("name"), from,
	)
	// TODO: create a notifications/requests record and alert managers.
}

func listListings(app *pocketbase.PocketBase, towerFilter string) string {
	filter := "status = 'active'"
	if towerFilter != "" {
		filter += fmt.Sprintf(" && tower.name ~ '%s'", strings.ReplaceAll(towerFilter, "'", "''"))
	}

	records, err := app.Dao().FindRecordsByFilter("listings", filter, "-created", 10, 0, nil)
	if err != nil || len(records) == 0 {
		return "No active listings found."
	}

	var sb strings.Builder
	sb.WriteString("Active listings:\n")
	for _, r := range records {
		sb.WriteString(fmt.Sprintf("  [%s] %s\n", r.GetString("category"), r.GetString("title")))
	}
	return strings.TrimRight(sb.String(), "\n")
}
