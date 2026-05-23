package messaging

import (
	"fmt"
	"strings"

	"github.com/pocketbase/pocketbase"
	"github.com/pocketbase/pocketbase/forms"
	"github.com/pocketbase/pocketbase/models"
)

func handleMessage(app *pocketbase.PocketBase, msg incomingMessage) string {
	cmd, args := parseCommand(msg.Body)

	switch cmd {
	case "help", "":
		return helpText()
	case "signup":
		return cmdSignup(app, msg.From, args)
	case "login":
		return cmdLogin(app, msg.From, args)
	case "logout":
		return cmdLogout(msg.From)
	case "whoami":
		return cmdWhoami(msg.From)
	case "towers":
		return listTowers(app)
	case "join":
		if args == "" {
			return "Usage: join <tower name>"
		}
		return joinTower(app, msg.From, args)
	case "post":
		return cmdPost(app, msg.From, args)
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

Account
  signup <name> <email> <password>
  login <email> <password>
  logout
  whoami

Listings
  post <lend|give|request> <title>
  listings [tower name]

Building
  towers
  join <tower name>

  ping  •  help`
}

func parseCommand(body string) (cmd, args string) {
	parts := strings.SplitN(strings.ToLower(strings.TrimSpace(body)), " ", 2)
	cmd = strings.ToLower(parts[0])
	if len(parts) > 1 {
		args = strings.TrimSpace(parts[1])
	}
	return
}

// ── Auth commands ──────────────────────────────────────────────────────────

func cmdSignup(app *pocketbase.PocketBase, from, args string) string {
	// signup <name> <email> <password>
	parts := strings.SplitN(args, " ", 3)
	if len(parts) < 3 {
		return "Usage: signup <name> <email> <password>"
	}
	name, email, password := parts[0], parts[1], parts[2]

	usersCol, err := app.Dao().FindCollectionByNameOrId("users")
	if err != nil {
		return "Error finding users collection: " + err.Error()
	}

	record := models.NewRecord(usersCol)
	form := forms.NewRecordUpsert(app, record)
	form.SetFullManageAccess(true)

	if err := form.LoadData(map[string]any{
		"email":           email,
		"password":        password,
		"passwordConfirm": password,
		"display_name":    name,
	}); err != nil {
		return "Invalid data: " + err.Error()
	}
	if err := form.Submit(); err != nil {
		return "Sign up failed: " + err.Error()
	}

	// Auto-verify so they can log in immediately.
	record.SetVerified(true)
	if err := app.Dao().SaveRecord(record); err != nil {
		return "Verification failed: " + err.Error()
	}

	setSession(from, &session{User: record})
	return fmt.Sprintf("Account created! Welcome, %s 👋\n\nYou're now logged in. Send 'towers' to get started.", name)
}

func cmdLogin(app *pocketbase.PocketBase, from, args string) string {
	// login <email> <password>
	parts := strings.SplitN(args, " ", 2)
	if len(parts) < 2 {
		return "Usage: login <email> <password>"
	}
	email, password := parts[0], parts[1]

	record, err := app.Dao().FindAuthRecordByEmail("users", email)
	if err != nil {
		return "No account found for that email."
	}
	if !record.ValidatePassword(password) {
		return "Wrong password."
	}

	setSession(from, &session{User: record})
	name := record.GetString("display_name")
	if name == "" {
		name = record.Email()
	}
	return fmt.Sprintf("Welcome back, %s! 👋\n\nSend 'help' to see what you can do.", name)
}

func cmdLogout(from string) string {
	clearSession(from)
	return "Logged out. Send 'login <email> <password>' to sign back in."
}

func cmdWhoami(from string) string {
	s := getSession(from)
	if s == nil {
		return "Not logged in. Send 'login <email> <password>'."
	}
	name := s.User.GetString("display_name")
	if name == "" {
		name = s.User.Email()
	}
	return fmt.Sprintf("Logged in as %s (%s)", name, s.User.Email())
}

// ── Listing commands ───────────────────────────────────────────────────────

func cmdPost(app *pocketbase.PocketBase, from, args string) string {
	s := getSession(from)
	if s == nil {
		return "You need to be logged in to post.\n\nSend: login <email> <password>"
	}

	// post <lend|give|request> <title>
	parts := strings.SplitN(args, " ", 2)
	if len(parts) < 2 {
		return "Usage: post <lend|give|request> <title>\n\nExample: post give Moving boxes"
	}
	category, title := strings.ToLower(parts[0]), parts[1]

	validCategories := map[string]bool{"lend": true, "give": true, "request": true}
	if !validCategories[category] {
		return fmt.Sprintf("Category must be lend, give, or request.\nGot: %q", category)
	}

	// Pick the first available tower.
	towers, err := app.Dao().FindRecordsByFilter("towers", "id != ''", "name", 1, 0, nil)
	if err != nil || len(towers) == 0 {
		return "No towers registered yet. Ask your building manager to set one up."
	}
	tower := towers[0]

	listingsCol, err := app.Dao().FindCollectionByNameOrId("listings")
	if err != nil {
		return "Error: " + err.Error()
	}

	record := models.NewRecord(listingsCol)
	record.Set("title", title)
	record.Set("category", category)
	record.Set("status", "active")
	record.Set("owner", s.User.Id)
	record.Set("tower", tower.Id)

	if err := app.Dao().SaveRecord(record); err != nil {
		return "Failed to post listing: " + err.Error()
	}

	return fmt.Sprintf("Posted! ✅\n\n[%s] %s\nTower: %s\n\nSend 'listings' to see all active listings.",
		strings.ToUpper(category), title, tower.GetString("name"))
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
		return fmt.Sprintf("No tower found matching %q.\nSend 'towers' to see all.", towerName)
	}
	return fmt.Sprintf(
		"Request noted for %q.\nA manager will be in touch shortly.",
		records[0].GetString("name"),
	)
}

func listListings(app *pocketbase.PocketBase, towerFilter string) string {
	filter := "status = 'active'"
	if towerFilter != "" {
		filter += fmt.Sprintf(" && tower.name ~ '%s'",
			strings.ReplaceAll(towerFilter, "'", "''"))
	}

	records, err := app.Dao().FindRecordsByFilter("listings", filter, "-created", 10, 0, nil)
	if err != nil || len(records) == 0 {
		return "No active listings found."
	}

	var sb strings.Builder
	sb.WriteString("Active listings:\n")
	for _, r := range records {
		sb.WriteString(fmt.Sprintf("  [%s] %s\n",
			strings.ToUpper(r.GetString("category")), r.GetString("title")))
	}
	return strings.TrimRight(sb.String(), "\n")
}
