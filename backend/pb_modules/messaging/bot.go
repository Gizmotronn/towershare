package messaging

import (
	"fmt"
	"strings"

	"github.com/pocketbase/pocketbase"
	"github.com/pocketbase/pocketbase/forms"
	"github.com/pocketbase/pocketbase/models"
)

// handleMessage is the top-level dispatcher.
func handleMessage(app *pocketbase.PocketBase, msg incomingMessage) *BotReply {
	s := getOrCreate(msg.From)

	switch s.State {
	case stateSignupName:
		return onSignupName(msg.From, msg.Body, s)
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

	return dispatchIdle(app, msg, s)
}

// ── Idle ───────────────────────────────────────────────────────────────────

func dispatchIdle(app *pocketbase.PocketBase, msg incomingMessage, s *session) *BotReply {
	body := strings.ToLower(strings.TrimSpace(msg.Body))

	switch body {
	case "", "start", "hello", "hi", "hey":
		return welcome(s)
	case "signup", "create account", "register":
		return startSignup(msg.From, s)
	case "login", "sign in", "signin":
		return startLogin(msg.From, s)
	case "logout", "sign out", "signout":
		return doLogout(msg.From)
	case "whoami", "me", "account":
		return doWhoami(s)
	case "post", "sell", "share", "new listing":
		return startPost(msg.From, s)
	case "listings", "browse", "feed":
		return doListings(app)
	case "towers", "buildings":
		return doTowers(app)
	case "ping":
		return reply("pong 🏓")
	case "help":
		return doHelp(s)
	}

	if s.User != nil {
		return replyWith(
			fmt.Sprintf("I didn't catch that. What would you like to do?"),
			act("Browse listings 📋", "listings"),
			act("Post something ✏️", "post"),
			act("Help", "help"),
		)
	}
	return replyWith(
		"I didn't catch that.",
		act("Create an account", "signup"),
		act("Sign in", "login"),
	)
}

// ── Welcome / help ─────────────────────────────────────────────────────────

func welcome(s *session) *BotReply {
	if s.User != nil {
		name := displayName(s.User)
		return replyWith(
			fmt.Sprintf("Hey %s! 👋\nWhat would you like to do?", name),
			act("Browse listings 📋", "listings"),
			act("Post something ✏️", "post"),
			act("My building 🏢", "towers"),
			act("Sign out", "logout"),
		)
	}
	return replyWith(
		"Hey! 👋 I'm TowerShare.\nI help neighbours share things in their building.",
		act("Create an account", "signup"),
		act("Sign in", "login"),
	)
}

func doHelp(s *session) *BotReply {
	if s.User != nil {
		return replyWith(
			"What can I help with?",
			act("Browse listings 📋", "listings"),
			act("Post something ✏️", "post"),
			act("My building 🏢", "towers"),
			act("Sign out", "logout"),
		)
	}
	return replyWith(
		"What would you like to do?",
		act("Create an account", "signup"),
		act("Sign in", "login"),
	)
}

// ── Sign-up flow ───────────────────────────────────────────────────────────

func startSignup(from string, s *session) *BotReply {
	s.State = stateSignupName
	s.Pending = make(map[string]string)
	return reply("Let's get you set up.\n\nWhat's your name?")
}

func onSignupName(from, body string, s *session) *BotReply {
	name := strings.TrimSpace(body)
	if name == "" {
		return reply("What's your name?")
	}
	s.Pending["name"] = name
	s.State = stateSignupEmail
	return reply(fmt.Sprintf("Hi %s! What's your email address?", name))
}

func onSignupEmail(from, body string, s *session) *BotReply {
	email := strings.TrimSpace(strings.ToLower(body))
	if !strings.Contains(email, "@") || !strings.Contains(email, ".") {
		return reply("That doesn't look right. Enter a valid email:")
	}
	s.Pending["email"] = email
	s.State = stateSignupPassword
	return reply("And a password — at least 8 characters.")
}

func onSignupPassword(app *pocketbase.PocketBase, from, body string, s *session) *BotReply {
	pw := strings.TrimSpace(body)
	if len(pw) < 8 {
		return reply("Too short — needs at least 8 characters. Try again:")
	}
	s.Pending["password"] = pw
	s.State = stateSignupTower
	return towerPickerReply(app, "Which tower do you live in?")
}

func onSignupTower(app *pocketbase.PocketBase, from, body string, s *session) *BotReply {
	if strings.HasPrefix(body, "tower:") {
		parts := strings.SplitN(body, ":", 3)
		if len(parts) == 3 {
			s.Pending["tower_id"] = parts[1]
			s.Pending["tower_name"] = parts[2]
		}
	} else {
		records, err := app.Dao().FindRecordsByFilter(
			"towers",
			fmt.Sprintf("name ~ '%s'", strings.ReplaceAll(body, "'", "''")),
			"name", 5, 0, nil,
		)
		if err != nil || len(records) == 0 {
			return towerPickerReply(app, fmt.Sprintf("No match for %q. Pick one:", body))
		}
		s.Pending["tower_id"] = records[0].Id
		s.Pending["tower_name"] = records[0].GetString("name")
	}

	s.State = stateSignupApartment
	return replyWith(
		"Got it. What's your apartment number?",
		act("Skip for now →", "skip"),
	)
}

func onSignupApartment(app *pocketbase.PocketBase, from, body string, s *session) *BotReply {
	if strings.ToLower(strings.TrimSpace(body)) == "skip" {
		return finaliseSignup(app, from, s)
	}

	unitNumber := strings.TrimSpace(body)
	towerID := s.Pending["tower_id"]
	towerName := s.Pending["tower_name"]

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
			fmt.Sprintf(
				"Apartment %s isn't in %s yet.\n\nCreate it? Your invite code will be %s — share it with flatmates so they can join.",
				unitNumber, towerName, code,
			),
			act(fmt.Sprintf("Create apartment %s", unitNumber), "create"),
			act("Try a different number", "retry"),
			act("Skip for now →", "skip"),
		)
	}

	apt := apartments[0]
	s.Pending["apartment_id"] = apt.Id

	if code := apt.GetString("join_code"); code != "" {
		s.State = stateSignupPasskey
		return reply(fmt.Sprintf("Apartment %s found. Enter the invite code to verify:", unitNumber))
	}

	return finaliseSignup(app, from, s)
}

func onSignupCreateApt(app *pocketbase.PocketBase, from, body string, s *session) *BotReply {
	switch strings.ToLower(strings.TrimSpace(body)) {
	case "retry":
		s.State = stateSignupApartment
		return reply("What's your apartment number?")

	case "skip":
		return finaliseSignup(app, from, s)

	case "create":
		unitNumber := s.Pending["pending_unit"]
		towerID := s.Pending["tower_id"]
		code := s.Pending["pending_code"]

		aptCol, err := app.Dao().FindCollectionByNameOrId("apartments")
		if err != nil {
			return reply("Something went wrong: " + err.Error())
		}

		record := models.NewRecord(aptCol)
		record.Set("unit_number", unitNumber)
		record.Set("tower", towerID)
		record.Set("join_code", code)

		if err := app.Dao().SaveRecord(record); err != nil {
			return reply("Couldn't create apartment: " + err.Error())
		}

		s.Pending["apartment_id"] = record.Id
		return finaliseSignup(app, from, s)

	default:
		unitNumber := s.Pending["pending_unit"]
		code := s.Pending["pending_code"]
		return replyWith(
			fmt.Sprintf("Create apartment %s with invite code %s?", unitNumber, code),
			act(fmt.Sprintf("Create apartment %s", unitNumber), "create"),
			act("Try a different number", "retry"),
			act("Skip for now →", "skip"),
		)
	}
}

func onSignupPasskey(app *pocketbase.PocketBase, from, body string, s *session) *BotReply {
	apt, err := app.Dao().FindRecordById("apartments", s.Pending["apartment_id"])
	if err != nil {
		return reply("Something went wrong. Try again.")
	}
	if apt.GetString("join_code") != strings.TrimSpace(body) {
		return replyWith(
			"That code doesn't match. Try again:",
			act("Start over", "signup"),
		)
	}
	return finaliseSignup(app, from, s)
}

func finaliseSignup(app *pocketbase.PocketBase, from string, s *session) *BotReply {
	// Capture before clearing pending.
	name  := s.Pending["name"]
	email := s.Pending["email"]
	pw    := s.Pending["password"]
	aptID := s.Pending["apartment_id"]

	usersCol, err := app.Dao().FindCollectionByNameOrId("users")
	if err != nil {
		return reply("Error creating account: " + err.Error())
	}

	record := models.NewRecord(usersCol)
	form := forms.NewRecordUpsert(app, record)
	form.SetFullManageAccess(true)

	data := map[string]any{
		"email":           email,
		"password":        pw,
		"passwordConfirm": pw,
		"display_name":    name,
	}
	if aptID != "" {
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
		return reply("Couldn't verify account: " + err.Error())
	}

	s.User = record
	s.State = stateIdle
	s.Pending = make(map[string]string)

	suffix := ""
	if aptID != "" {
		suffix = "\nApartment linked ✓"
	}

	return replyWith(
		fmt.Sprintf("🎉 You're in, %s!\nWelcome to TowerShare.%s", name, suffix),
		act("Browse listings 📋", "listings"),
		act("Post something ✏️", "post"),
	)
}

// ── Login flow ─────────────────────────────────────────────────────────────

func startLogin(from string, s *session) *BotReply {
	s.State = stateLoginEmail
	s.Pending = make(map[string]string)
	return reply("What's your email?")
}

func onLoginEmail(from, body string, s *session) *BotReply {
	email := strings.TrimSpace(strings.ToLower(body))
	if !strings.Contains(email, "@") {
		return reply("That doesn't look right. Enter your email:")
	}
	s.Pending["email"] = email
	s.State = stateLoginPassword
	return reply("And your password?")
}

func onLoginPassword(app *pocketbase.PocketBase, from, body string, s *session) *BotReply {
	record, err := app.Dao().FindAuthRecordByEmail("users", s.Pending["email"])
	if err != nil || !record.ValidatePassword(strings.TrimSpace(body)) {
		return replyWith(
			"Email or password didn't match.",
			act("Try again", "login"),
			act("Create an account", "signup"),
		)
	}

	s.User = record
	s.State = stateIdle
	s.Pending = make(map[string]string)

	return replyWith(
		fmt.Sprintf("Welcome back, %s! 👋", displayName(record)),
		act("Browse listings 📋", "listings"),
		act("Post something ✏️", "post"),
		act("My building 🏢", "towers"),
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
		return replyWith("You're not signed in.",
			act("Sign in", "login"),
			act("Create an account", "signup"),
		)
	}
	return reply(fmt.Sprintf("Signed in as %s\n%s", displayName(s.User), s.User.Email()))
}

// ── Post listing flow ──────────────────────────────────────────────────────

func startPost(from string, s *session) *BotReply {
	if s.User == nil {
		return replyWith(
			"You need to be signed in to post.",
			act("Sign in", "login"),
			act("Create an account", "signup"),
		)
	}
	s.State = statePostType
	s.Pending = make(map[string]string)
	return replyWith(
		"What type of listing?",
		act("Lend 🔄 — they return it", "lend"),
		act("Give 🎁 — free to keep", "give"),
		act("Request 🙏 — looking for something", "request"),
	)
}

func onPostType(from, body string, s *session) *BotReply {
	cat := strings.ToLower(strings.TrimSpace(body))
	if cat != "lend" && cat != "give" && cat != "request" {
		return replyWith("Choose a type:",
			act("Lend 🔄", "lend"),
			act("Give 🎁", "give"),
			act("Request 🙏", "request"),
		)
	}
	s.Pending["category"] = cat
	s.State = statePostTitle
	prompts := map[string]string{
		"lend":    "What are you lending?",
		"give":    "What are you giving away?",
		"request": "What are you looking for?",
	}
	return replyWith(prompts[cat], act("← Cancel", "cancel"))
}

func onPostTitle(app *pocketbase.PocketBase, from, body string, s *session) *BotReply {
	if strings.ToLower(strings.TrimSpace(body)) == "cancel" {
		s.State = stateIdle
		return replyWith("Cancelled.", act("Browse listings 📋", "listings"))
	}

	title := strings.TrimSpace(body)
	if title == "" {
		return reply("Enter a title:")
	}
	if len(title) > 200 {
		return reply("Too long — keep it under 200 characters.")
	}

	towers, err := app.Dao().FindRecordsByFilter("towers", "id != ''", "name", 1, 0, nil)
	if err != nil || len(towers) == 0 {
		s.State = stateIdle
		return reply("No towers found. Ask your manager to set one up first.")
	}

	listingsCol, err := app.Dao().FindCollectionByNameOrId("listings")
	if err != nil {
		s.State = stateIdle
		return reply("Error: " + err.Error())
	}

	record := models.NewRecord(listingsCol)
	record.Set("title", title)
	record.Set("category", s.Pending["category"])
	record.Set("status", "active")
	record.Set("owner", s.User.Id)
	record.Set("tower", towers[0].Id)

	if err := app.Dao().SaveRecord(record); err != nil {
		s.State = stateIdle
		return reply("Couldn't post: " + err.Error())
	}

	cat := s.Pending["category"]
	s.State = stateIdle
	s.Pending = make(map[string]string)

	return replyWith(
		fmt.Sprintf("Posted! ✅\n\n%s  %s\n%s",
			categoryEmoji(cat), title, towers[0].GetString("name")),
		act("Browse listings 📋", "listings"),
		act("Post another ✏️", "post"),
	)
}

// ── Listings / towers ──────────────────────────────────────────────────────

func doTowers(app *pocketbase.PocketBase) *BotReply {
	records, err := app.Dao().FindRecordsByFilter("towers", "id != ''", "name", 20, 0, nil)
	if err != nil || len(records) == 0 {
		return reply("No towers registered yet.")
	}
	var sb strings.Builder
	sb.WriteString("Towers:\n")
	for _, r := range records {
		sb.WriteString("  • " + r.GetString("name") + "\n")
	}
	return reply(strings.TrimRight(sb.String(), "\n"))
}

func doListings(app *pocketbase.PocketBase) *BotReply {
	records, err := app.Dao().FindRecordsByFilter(
		"listings", "status = 'active'", "-created", 10, 0, nil)
	if err != nil || len(records) == 0 {
		return replyWith("No active listings right now.", act("Post something ✏️", "post"))
	}
	var sb strings.Builder
	sb.WriteString("Active listings:\n\n")
	for _, r := range records {
		sb.WriteString(categoryEmoji(r.GetString("category")) + "  " + r.GetString("title") + "\n")
	}
	return replyWith(strings.TrimRight(sb.String(), "\n"), act("Post something ✏️", "post"))
}

// ── Helpers ────────────────────────────────────────────────────────────────

func displayName(r *models.Record) string {
	if n := r.GetString("display_name"); n != "" {
		return n
	}
	return r.Email()
}

func categoryEmoji(cat string) string {
	return map[string]string{"lend": "🔄", "give": "🎁", "request": "🙏"}[cat]
}

// demoJoinCode repeats unitNumber until ≥ 6 chars, then takes the first 6.
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

func towerPickerReply(app *pocketbase.PocketBase, text string) *BotReply {
	towers, err := app.Dao().FindRecordsByFilter("towers", "id != ''", "name", 10, 0, nil)
	if err != nil || len(towers) == 0 {
		return reply("No towers registered yet. Ask your manager to set one up.")
	}
	actions := make([]BotAction, len(towers))
	for i, t := range towers {
		actions[i] = act(t.GetString("name"),
			fmt.Sprintf("tower:%s:%s", t.Id, t.GetString("name")))
	}
	return replyWith(text, actions...)
}
