package messaging

import (
	"fmt"
	"strings"

	"github.com/pocketbase/pocketbase"
	"github.com/pocketbase/pocketbase/forms"
	"github.com/pocketbase/pocketbase/models"
)

// handleMessage is the top-level dispatcher.
// It routes based on the session's current state, falling back to command parsing.
func handleMessage(app *pocketbase.PocketBase, msg incomingMessage) *BotReply {
	s := getOrCreate(msg.From)

	switch s.State {
	case stateSignupName:
		return onSignupName(app, msg.From, msg.Body, s)
	case stateSignupEmail:
		return onSignupEmail(msg.From, msg.Body, s)
	case stateSignupPassword:
		return onSignupPassword(app, msg.From, msg.Body, s)
	case stateSignupTower:
		return onSignupTower(app, msg.From, msg.Body, s)
	case stateSignupApartment:
		return onSignupApartment(app, msg.From, msg.Body, s)
	case stateSignupCreateApt:
		return onSignupCreateApt(app, msg.From, msg.Body, s)
	case stateSignupPasskey:
		return onSignupPasskey(app, msg.From, msg.Body, s)
	case stateLoginEmail:
		return onLoginEmail(msg.From, msg.Body, s)
	case stateLoginPassword:
		return onLoginPassword(app, msg.From, msg.Body, s)
	case statePostType:
		return onPostType(msg.From, msg.Body, s)
	case statePostTitle:
		return onPostTitle(app, msg.From, msg.Body, s)
	}

	// Idle: parse as a command.
	return dispatchCommand(app, msg, s)
}

// ── Idle command dispatch ──────────────────────────────────────────────────

func dispatchCommand(app *pocketbase.PocketBase, msg incomingMessage, s *session) *BotReply {
	cmd, _ := parseCommand(msg.Body)

	switch cmd {
	case "", "start", "hello", "hi":
		return welcome(app, s)
	case "signup", "register", "create account":
		return startSignup(msg.From, s)
	case "login", "sign in":
		return startLogin(msg.From, s)
	case "logout":
		return doLogout(msg.From)
	case "whoami":
		return doWhoami(s)
	case "post", "sell", "share":
		return startPost(msg.From, s)
	case "listings", "browse":
		return doListListings(app, "")
	case "towers":
		return doListTowers(app)
	case "ping":
		return reply("pong 🏓")
	case "help":
		return doHelp(s)
	default:
		return replyWith(
			fmt.Sprintf("I didn't understand %q.", msg.Body),
			act("Show help", "help"),
			act("Browse listings", "listings"),
		)
	}
}

func parseCommand(body string) (cmd, args string) {
	parts := strings.SplitN(strings.TrimSpace(body), " ", 2)
	cmd = strings.ToLower(parts[0])
	if len(parts) > 1 {
		args = strings.TrimSpace(parts[1])
	}
	return
}

// ── Welcome ────────────────────────────────────────────────────────────────

func welcome(app *pocketbase.PocketBase, s *session) *BotReply {
	if s.User != nil {
		name := s.User.GetString("display_name")
		if name == "" {
			name = s.User.Email()
		}
		return replyWith(
			fmt.Sprintf("Welcome back, %s! 👋\nWhat would you like to do?", name),
			act("Browse listings 📋", "listings"),
			act("Post something ✏️", "post"),
			act("My towers 🏢", "towers"),
			act("Help ❓", "help"),
		)
	}
	return replyWith(
		"Welcome to TowerShare 🏢\n\nShare resources with your neighbours.\nWhat would you like to do?",
		act("Create account", "signup"),
		act("Sign in", "login"),
	)
}

func doHelp(s *session) *BotReply {
	if s.User == nil {
		return replyWith(
			"TowerShare helps you share things with your building.\n\nGet started:",
			act("Create account", "signup"),
			act("Sign in", "login"),
		)
	}
	return replyWith(
		"TowerShare 🏢\n\nWhat can I help with?",
		act("Browse listings 📋", "listings"),
		act("Post something ✏️", "post"),
		act("My towers 🏢", "towers"),
		act("Sign out", "logout"),
	)
}

// ── Sign up flow ───────────────────────────────────────────────────────────

func startSignup(from string, s *session) *BotReply {
	s.State = stateSignupName
	s.Pending = make(map[string]string)
	return reply("Let's get you set up 🚀\n\nWhat's your name?")
}

func onSignupName(app *pocketbase.PocketBase, from, body string, s *session) *BotReply {
	name := strings.TrimSpace(body)
	if name == "" {
		return reply("Please enter your name.")
	}
	s.Pending["name"] = name
	s.State = stateSignupEmail
	return reply(fmt.Sprintf("Nice to meet you, %s! 👋\n\nWhat's your email address?", name))
}

func onSignupEmail(from, body string, s *session) *BotReply {
	email := strings.TrimSpace(strings.ToLower(body))
	if !strings.Contains(email, "@") {
		return reply("That doesn't look like a valid email. Try again:")
	}
	s.Pending["email"] = email
	s.State = stateSignupPassword
	return reply("Choose a password\n(minimum 8 characters)")
}

func onSignupPassword(app *pocketbase.PocketBase, from, body string, s *session) *BotReply {
	password := strings.TrimSpace(body)
	if len(password) < 8 {
		return reply("Password must be at least 8 characters. Try again:")
	}
	s.Pending["password"] = password
	s.State = stateSignupTower
	return towerPickerReply(app, "Which tower do you live in?")
}

func onSignupTower(app *pocketbase.PocketBase, from, body string, s *session) *BotReply {
	// Value is either "tower:<id>:<name>" (from action chip) or free text.
	if strings.HasPrefix(body, "tower:") {
		parts := strings.SplitN(body, ":", 3)
		if len(parts) == 3 {
			s.Pending["tower_id"] = parts[1]
			s.Pending["tower_name"] = parts[2]
		}
	} else {
		// Free-text tower lookup.
		records, err := app.Dao().FindRecordsByFilter(
			"towers",
			fmt.Sprintf("name ~ '%s'", strings.ReplaceAll(body, "'", "''")),
			"name", 5, 0, nil,
		)
		if err != nil || len(records) == 0 {
			return towerPickerReply(app, fmt.Sprintf("No tower found for %q. Please choose one:", body))
		}
		s.Pending["tower_id"] = records[0].Id
		s.Pending["tower_name"] = records[0].GetString("name")
	}

	s.State = stateSignupApartment
	return replyWith(
		fmt.Sprintf("Great — %s! 🏢\n\nWhat's your apartment number?", s.Pending["tower_name"]),
		act("Skip →", "skip"),
	)
}

func onSignupApartment(app *pocketbase.PocketBase, from, body string, s *session) *BotReply {
	if strings.ToLower(body) == "skip" {
		return finaliseSignup(app, from, s)
	}

	unitNumber := strings.TrimSpace(body)
	towerID := s.Pending["tower_id"]

	// Look up the apartment.
	apartments, err := app.Dao().FindRecordsByFilter(
		"apartments",
		fmt.Sprintf("unit_number = '%s' && tower = '%s'",
			strings.ReplaceAll(unitNumber, "'", "''"), towerID),
		"unit_number", 1, 0, nil,
	)
	if err != nil || len(apartments) == 0 {
		code := demoJoinCode(unitNumber)
		s.Pending["pending_unit"] = unitNumber
		s.Pending["pending_code"] = code
		s.State = stateSignupCreateApt
		return replyWith(
			fmt.Sprintf("Apartment %s isn't registered yet in %s.\n\nCreate it now? Your join code will be:\n\n  %s\n\nShare this code with your flatmates so they can join the same apartment.",
				unitNumber, s.Pending["tower_name"], code),
			act(fmt.Sprintf("Create apt %s", unitNumber), "create"),
			act("Try a different number", "retry"),
			act("Skip →", "skip"),
		)
	}

	apt := apartments[0]
	s.Pending["apartment_id"] = apt.Id
	joinCode := apt.GetString("join_code")

	if joinCode != "" {
		s.State = stateSignupPasskey
		return reply("Enter the apartment passkey to verify you live here:")
	}

	// No passkey set — proceed directly.
	return finaliseSignup(app, from, s)
}

func onSignupPasskey(app *pocketbase.PocketBase, from, body string, s *session) *BotReply {
	aptID := s.Pending["apartment_id"]
	apt, err := app.Dao().FindRecordById("apartments", aptID)
	if err != nil {
		return reply("Something went wrong. Please try again.")
	}
	if apt.GetString("join_code") != strings.TrimSpace(body) {
		return replyWith(
			"Wrong passkey. Try again, or contact your manager.",
			act("Try again", ""),
			act("Start over", "signup"),
		)
	}
	return finaliseSignup(app, from, s)
}

func finaliseSignup(app *pocketbase.PocketBase, from string, s *session) *BotReply {
	usersCol, err := app.Dao().FindCollectionByNameOrId("users")
	if err != nil {
		return reply("Error creating account: " + err.Error())
	}

	record := models.NewRecord(usersCol)
	form := forms.NewRecordUpsert(app, record)
	form.SetFullManageAccess(true)

	data := map[string]any{
		"email":           s.Pending["email"],
		"password":        s.Pending["password"],
		"passwordConfirm": s.Pending["password"],
		"display_name":    s.Pending["name"],
	}
	if aptID := s.Pending["apartment_id"]; aptID != "" {
		data["apartment"] = aptID
	}

	if err := form.LoadData(data); err != nil {
		return reply("Invalid data: " + err.Error())
	}
	if err := form.Submit(); err != nil {
		return reply("Sign up failed: " + err.Error())
	}

	record.SetVerified(true)
	if err := app.Dao().SaveRecord(record); err != nil {
		return reply("Could not verify account: " + err.Error())
	}

	s.User = record
	s.State = stateIdle
	s.Pending = make(map[string]string)

	apt := ""
	if s.Pending["apartment_id"] != "" {
		apt = "\nApartment linked ✓"
	}

	return replyWith(
		fmt.Sprintf("Welcome to TowerShare, %s! 🎉%s\n\nYou're all set.", s.Pending["name"], apt),
		act("Browse listings 📋", "listings"),
		act("Post something ✏️", "post"),
	)
}

// ── Login flow ─────────────────────────────────────────────────────────────

func startLogin(from string, s *session) *BotReply {
	s.State = stateLoginEmail
	s.Pending = make(map[string]string)
	return reply("Enter your email:")
}

func onLoginEmail(from, body string, s *session) *BotReply {
	email := strings.TrimSpace(strings.ToLower(body))
	if !strings.Contains(email, "@") {
		return reply("That doesn't look like a valid email. Try again:")
	}
	s.Pending["email"] = email
	s.State = stateLoginPassword
	return reply("Enter your password:")
}

func onLoginPassword(app *pocketbase.PocketBase, from, body string, s *session) *BotReply {
	record, err := app.Dao().FindAuthRecordByEmail("users", s.Pending["email"])
	if err != nil {
		return replyWith(
			"No account found for that email.",
			act("Try again", "login"),
			act("Create account", "signup"),
		)
	}
	if !record.ValidatePassword(strings.TrimSpace(body)) {
		return replyWith(
			"Wrong password. Try again:",
			act("Forgot? Start over", "login"),
		)
	}

	s.User = record
	s.State = stateIdle
	s.Pending = make(map[string]string)

	name := record.GetString("display_name")
	if name == "" {
		name = record.Email()
	}
	return replyWith(
		fmt.Sprintf("Welcome back, %s! 👋", name),
		act("Browse listings 📋", "listings"),
		act("Post something ✏️", "post"),
		act("My towers 🏢", "towers"),
	)
}

func doLogout(from string) *BotReply {
	clearSession(from)
	return replyWith(
		"You've been signed out.",
		act("Sign back in", "login"),
	)
}

func doWhoami(s *session) *BotReply {
	if s.User == nil {
		return replyWith("Not signed in.",
			act("Sign in", "login"),
			act("Create account", "signup"),
		)
	}
	name := s.User.GetString("display_name")
	if name == "" {
		name = s.User.Email()
	}
	return reply(fmt.Sprintf("Signed in as %s\n%s", name, s.User.Email()))
}

// ── Post listing flow ──────────────────────────────────────────────────────

func startPost(from string, s *session) *BotReply {
	if s.User == nil {
		return replyWith(
			"You need to be signed in to post.",
			act("Sign in", "login"),
			act("Create account", "signup"),
		)
	}
	s.State = statePostType
	s.Pending = make(map[string]string)
	return replyWith(
		"What type of listing?",
		act("Lend 🔄 — borrow & return", "lend"),
		act("Give 🎁 — free to keep", "give"),
		act("Request 🙏 — looking for something", "request"),
	)
}

func onPostType(from, body string, s *session) *BotReply {
	cat := strings.ToLower(strings.TrimSpace(body))
	if cat != "lend" && cat != "give" && cat != "request" {
		return replyWith(
			"Choose a category:",
			act("Lend 🔄", "lend"),
			act("Give 🎁", "give"),
			act("Request 🙏", "request"),
		)
	}
	s.Pending["category"] = cat
	s.State = statePostTitle

	prompts := map[string]string{
		"lend":    "What are you lending? Give it a short title:",
		"give":    "What are you giving away?",
		"request": "What are you looking for?",
	}
	return replyWith(prompts[cat],
		act("← Cancel", "cancel"),
	)
}

func onPostTitle(app *pocketbase.PocketBase, from, body string, s *session) *BotReply {
	if strings.ToLower(strings.TrimSpace(body)) == "cancel" {
		s.State = stateIdle
		return replyWith("Cancelled.",
			act("Browse listings 📋", "listings"),
		)
	}

	title := strings.TrimSpace(body)
	if title == "" {
		return reply("Please enter a title:")
	}
	if len(title) > 200 {
		return reply("Title too long (max 200 chars). Try again:")
	}

	category := s.Pending["category"]

	towers, err := app.Dao().FindRecordsByFilter("towers", "id != ''", "name", 1, 0, nil)
	if err != nil || len(towers) == 0 {
		s.State = stateIdle
		return reply("No towers found. Ask your manager to set one up first.")
	}
	tower := towers[0]

	listingsCol, err := app.Dao().FindCollectionByNameOrId("listings")
	if err != nil {
		s.State = stateIdle
		return reply("Error: " + err.Error())
	}

	record := models.NewRecord(listingsCol)
	record.Set("title", title)
	record.Set("category", category)
	record.Set("status", "active")
	record.Set("owner", s.User.Id)
	record.Set("tower", tower.Id)

	if err := app.Dao().SaveRecord(record); err != nil {
		s.State = stateIdle
		return reply("Failed to post: " + err.Error())
	}

	s.State = stateIdle
	s.Pending = make(map[string]string)

	return replyWith(
		fmt.Sprintf("Posted! ✅\n\n[%s] %s\nTower: %s", strings.ToUpper(category), title, tower.GetString("name")),
		act("Browse listings 📋", "listings"),
		act("Post another ✏️", "post"),
	)
}

// ── Listings / towers ──────────────────────────────────────────────────────

func doListTowers(app *pocketbase.PocketBase) *BotReply {
	records, err := app.Dao().FindRecordsByFilter("towers", "id != ''", "name", 20, 0, nil)
	if err != nil || len(records) == 0 {
		return reply("No towers registered yet.")
	}
	var sb strings.Builder
	sb.WriteString("Towers:\n")
	for _, r := range records {
		sb.WriteString(fmt.Sprintf("  • %s\n", r.GetString("name")))
	}
	return reply(strings.TrimRight(sb.String(), "\n"))
}

func doListListings(app *pocketbase.PocketBase, towerFilter string) *BotReply {
	filter := "status = 'active'"
	if towerFilter != "" {
		filter += fmt.Sprintf(" && tower.name ~ '%s'",
			strings.ReplaceAll(towerFilter, "'", "''"))
	}

	records, err := app.Dao().FindRecordsByFilter("listings", filter, "-created", 10, 0, nil)
	if err != nil || len(records) == 0 {
		return replyWith(
			"No active listings right now.",
			act("Post something ✏️", "post"),
		)
	}

	var sb strings.Builder
	sb.WriteString("Active listings:\n\n")
	for _, r := range records {
		sb.WriteString(fmt.Sprintf("[%s] %s\n",
			strings.ToUpper(r.GetString("category")), r.GetString("title")))
	}
	return replyWith(strings.TrimRight(sb.String(), "\n"),
		act("Post something ✏️", "post"),
	)
}

// ── Helpers ────────────────────────────────────────────────────────────────

func onSignupCreateApt(app *pocketbase.PocketBase, from, body string, s *session) *BotReply {
	switch strings.ToLower(strings.TrimSpace(body)) {
	case "retry":
		s.State = stateSignupApartment
		return reply("Enter your apartment number:")

	case "skip":
		s.State = stateIdle
		return finaliseSignup(app, from, s)

	case "create":
		unitNumber := s.Pending["pending_unit"]
		towerID := s.Pending["tower_id"]
		code := s.Pending["pending_code"]

		aptCol, err := app.Dao().FindCollectionByNameOrId("apartments")
		if err != nil {
			return reply("Error finding apartments collection: " + err.Error())
		}

		record := models.NewRecord(aptCol)
		record.Set("unit_number", unitNumber)
		record.Set("tower", towerID)
		record.Set("join_code", code)

		if err := app.Dao().SaveRecord(record); err != nil {
			return reply("Could not create apartment: " + err.Error())
		}

		s.Pending["apartment_id"] = record.Id
		return finaliseSignup(app, from, s)

	default:
		// They typed something unexpected — re-show the options.
		unitNumber := s.Pending["pending_unit"]
		code := s.Pending["pending_code"]
		return replyWith(
			fmt.Sprintf("Create apartment %s with join code %s?", unitNumber, code),
			act(fmt.Sprintf("Create apt %s", unitNumber), "create"),
			act("Try a different number", "retry"),
			act("Skip →", "skip"),
		)
	}
}

// demoJoinCode repeats the unit number string until it reaches 6 characters.
// e.g. "12" → "121212", "123" → "123123", "1" → "111111"
func demoJoinCode(unitNumber string) string {
	if unitNumber == "" {
		return "000000"
	}
	s := unitNumber
	for len(s) < 6 {
		s += unitNumber
	}
	return s[:6]
}

// towerPickerReply builds a reply with one action chip per tower.
func towerPickerReply(app *pocketbase.PocketBase, text string) *BotReply {
	towers, err := app.Dao().FindRecordsByFilter("towers", "id != ''", "name", 10, 0, nil)
	if err != nil || len(towers) == 0 {
		return reply("No towers registered yet. Ask your manager to set one up first.")
	}
	actions := make([]BotAction, len(towers))
	for i, t := range towers {
		actions[i] = act(t.GetString("name"), fmt.Sprintf("tower:%s:%s", t.Id, t.GetString("name")))
	}
	return replyWith(text, actions...)
}
