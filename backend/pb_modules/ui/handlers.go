package ui

import (
	"fmt"
	"html/template"
	"io"
	"net/http"
	"strings"
	"time"

	"github.com/labstack/echo/v5"
	"github.com/pocketbase/pocketbase"
	"github.com/pocketbase/pocketbase/forms"
	"github.com/pocketbase/pocketbase/models"
	"github.com/pocketbase/pocketbase/tools/filesystem"
)

const cookieName = "ts_uid"
const cookieTTL = 7 * 24 * time.Hour

// ── Auth middleware ────────────────────────────────────────────────────────

func requireAuth(app *pocketbase.PocketBase, c echo.Context) (*models.Record, bool) {
	cookie, err := c.Request().Cookie(cookieName)
	if err != nil {
		return nil, false
	}
	user, err := app.Dao().FindRecordById("users", cookie.Value)
	if err != nil {
		return nil, false
	}
	return user, true
}

func setAuthCookie(c echo.Context, userID string) {
	http.SetCookie(c.Response().Writer, &http.Cookie{
		Name:     cookieName,
		Value:    userID,
		Path:     "/",
		MaxAge:   int(cookieTTL.Seconds()),
		HttpOnly: true,
		SameSite: http.SameSiteLaxMode,
	})
}

func clearAuthCookie(c echo.Context) {
	http.SetCookie(c.Response().Writer, &http.Cookie{
		Name:   cookieName,
		Value:  "",
		Path:   "/",
		MaxAge: -1,
	})
}

func userName(u *models.Record) string {
	n := u.GetString("display_name")
	if n == "" {
		return u.Email()
	}
	return n
}

// ── Login ──────────────────────────────────────────────────────────────────

func handleLogin(app *pocketbase.PocketBase) echo.HandlerFunc {
	return func(c echo.Context) error {
		if _, ok := requireAuth(app, c); ok {
			return c.Redirect(http.StatusFound, "/app/feed")
		}
		if c.Request().Method == http.MethodGet {
			return render(c, loginTmpl, map[string]any{"Title": "Sign In"})
		}
		email := strings.TrimSpace(c.FormValue("email"))
		password := c.FormValue("password")

		record, err := app.Dao().FindAuthRecordByEmail("users", email)
		if err != nil || !record.ValidatePassword(password) {
			return render(c, loginTmpl, map[string]any{
				"Title": "Sign In", "Error": "Invalid email or password.", "Email": email,
			})
		}
		setAuthCookie(c, record.Id)
		return c.Redirect(http.StatusFound, "/app/feed")
	}
}

// ── Signup ─────────────────────────────────────────────────────────────────

func handleSignup(app *pocketbase.PocketBase) echo.HandlerFunc {
	return func(c echo.Context) error {
		if _, ok := requireAuth(app, c); ok {
			return c.Redirect(http.StatusFound, "/app/feed")
		}
		if c.Request().Method == http.MethodGet {
			return render(c, signupTmpl, map[string]any{"Title": "Sign Up"})
		}

		email := strings.TrimSpace(strings.ToLower(c.FormValue("email")))
		password := c.FormValue("password")
		displayName := strings.TrimSpace(c.FormValue("display_name"))

		usersCol, _ := app.Dao().FindCollectionByNameOrId("users")
		record := models.NewRecord(usersCol)
		form := forms.NewRecordUpsert(app, record)
		form.SetFullManageAccess(true)

		if err := form.LoadData(map[string]any{
			"email":           email,
			"password":        password,
			"passwordConfirm": password,
			"display_name":    displayName,
		}); err != nil {
			return render(c, signupTmpl, map[string]any{
				"Title": "Sign Up", "Error": err.Error(),
				"Email": email, "DisplayName": displayName,
			})
		}
		if err := form.Submit(); err != nil {
			return render(c, signupTmpl, map[string]any{
				"Title": "Sign Up", "Error": err.Error(),
				"Email": email, "DisplayName": displayName,
			})
		}
		record.SetVerified(true)
		_ = app.Dao().SaveRecord(record)

		setAuthCookie(c, record.Id)
		return c.Redirect(http.StatusFound, "/app/feed")
	}
}

// ── Logout ─────────────────────────────────────────────────────────────────

func handleLogout(c echo.Context) error {
	clearAuthCookie(c)
	return c.Redirect(http.StatusFound, "/app/login")
}

// ── Feed ───────────────────────────────────────────────────────────────────

type listingView struct {
	Title      string
	Category   string
	TowerName  string
	ImageURL   string
	TimeAgo    string
}

func handleFeed(app *pocketbase.PocketBase) echo.HandlerFunc {
	return func(c echo.Context) error {
		user, ok := requireAuth(app, c)
		if !ok {
			return c.Redirect(http.StatusFound, "/app/login")
		}

		cat := c.QueryParam("cat")
		filter := "status = 'active'"
		if cat != "" {
			filter += fmt.Sprintf(" && category = '%s'", cat)
		}

		records, err := app.Dao().FindRecordsByFilter("listings", filter, "-created", 50, 0, nil)
		if err != nil {
			records = nil
		}

		var listings []listingView
		for _, r := range records {
			lv := listingView{
				Title:    r.GetString("title"),
				Category: r.GetString("category"),
				TimeAgo:  timeAgo(r.GetDateTime("created").Time()),
			}
			// Expand tower name.
			if towerID := r.GetString("tower"); towerID != "" {
				if t, err := app.Dao().FindRecordById("towers", towerID); err == nil {
					lv.TowerName = t.GetString("name")
				}
			}
			// First image thumb URL.
			images := r.GetStringSlice("images")
			if len(images) > 0 {
				lv.ImageURL = fmt.Sprintf("/api/files/%s/%s/%s?thumb=600x0",
					r.Collection().Id, r.Id, images[0])
			}
			listings = append(listings, lv)
		}

		return render(c, feedTmpl, map[string]any{
			"Title":    "Feed",
			"UserName": userName(user),
			"Filter":   cat,
			"Listings": listings,
		})
	}
}

// ── New listing ────────────────────────────────────────────────────────────

type towerOption struct {
	Id   string
	Name string
}

func handleNewListing(app *pocketbase.PocketBase) echo.HandlerFunc {
	return func(c echo.Context) error {
		user, ok := requireAuth(app, c)
		if !ok {
			return c.Redirect(http.StatusFound, "/app/login")
		}

		towerRecords, _ := app.Dao().FindRecordsByFilter("towers", "id != ''", "name", 50, 0, nil)
		towers := make([]towerOption, len(towerRecords))
		for i, t := range towerRecords {
			towers[i] = towerOption{Id: t.Id, Name: t.GetString("name")}
		}

		if c.Request().Method == http.MethodGet {
			return render(c, newListingTmpl, map[string]any{
				"Title": "New Listing", "Towers": towers,
			})
		}

		// Parse multipart form (max 32 MB for images).
		if err := c.Request().ParseMultipartForm(32 << 20); err != nil {
			return render(c, newListingTmpl, map[string]any{
				"Title": "New Listing", "Towers": towers,
				"Error": "Could not parse form: " + err.Error(),
			})
		}

		title := strings.TrimSpace(c.FormValue("title"))
		description := strings.TrimSpace(c.FormValue("description"))
		category := c.FormValue("category")
		towerID := c.FormValue("tower")

		listingsCol, _ := app.Dao().FindCollectionByNameOrId("listings")
		record := models.NewRecord(listingsCol)
		form := forms.NewRecordUpsert(app, record)
		form.SetFullManageAccess(true)

		if err := form.LoadData(map[string]any{
			"title":       title,
			"description": description,
			"category":    category,
			"status":      "active",
			"owner":       user.Id,
			"tower":       towerID,
		}); err != nil {
			return render(c, newListingTmpl, map[string]any{
				"Title": "New Listing", "Towers": towers, "Error": err.Error(),
				"ListingTitle": title, "Description": description,
				"Category": category, "TowerID": towerID,
			})
		}

		// Attach uploaded images.
		if files := c.Request().MultipartForm.File["images"]; len(files) > 0 {
			max := 5
			if len(files) < max {
				max = len(files)
			}
			for _, fh := range files[:max] {
				pbFile, err := filesystem.NewFileFromMultipart(fh)
				if err == nil {
					_ = form.AddFiles("images", pbFile)
				}
			}
		}

		if err := form.Submit(); err != nil {
			return render(c, newListingTmpl, map[string]any{
				"Title": "New Listing", "Towers": towers, "Error": err.Error(),
			})
		}

		return c.Redirect(http.StatusFound, "/app/feed")
	}
}

// ── Helpers ────────────────────────────────────────────────────────────────

func render(c echo.Context, tmpl *template.Template, data map[string]any) error {
	c.Response().Header().Set("Content-Type", "text/html; charset=utf-8")
	c.Response().WriteHeader(http.StatusOK)
	return tmpl.Execute(io.Writer(c.Response().Writer), data)
}

func timeAgo(t time.Time) string {
	d := time.Since(t)
	switch {
	case d < time.Minute:
		return "just now"
	case d < time.Hour:
		return fmt.Sprintf("%dm ago", int(d.Minutes()))
	case d < 24*time.Hour:
		return fmt.Sprintf("%dh ago", int(d.Hours()))
	default:
		return fmt.Sprintf("%dd ago", int(d.Hours()/24))
	}
}
