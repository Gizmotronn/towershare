# Extending PocketBase with Modules

TowerShare uses Go modules to extend PocketBase. Each module lives in `backend/pb_modules/<name>/` and registers hooks, routes, or scheduled jobs against the `*pocketbase.PocketBase` app.

## Quick start

```bash
./knowns.sh new-module my_feature
```

This scaffolds `backend/pb_modules/my_feature/module.go` with a `Register(app)` function.

Then wire it into `backend/main.go`:

```go
import (
    _ "github.com/gizmotronn/towershare/pb_modules/my_feature"
)
```

If your module uses an `init()` self-registration pattern (advanced), the blank import is enough. Otherwise call `my_feature.Register(app)` before `app.Start()`.

---

## Module anatomy

```
pb_modules/
  my_feature/
    module.go        ← Register(app) entry point
    routes.go        ← optional: HTTP route handlers
    hooks.go         ← optional: record / model event hooks
    schema.go        ← optional: migration helpers specific to this module
```

Split files are optional — one file is fine for small modules.

---

## Common patterns

### Custom REST route

```go
app.OnServe().BindFunc(func(se *core.ServeEvent) error {
    se.Router.GET("/api/towershare/my-endpoint", func(re *core.RequestEvent) error {
        info, _ := re.RequestInfo()
        if info.Auth == nil {
            return re.JSON(http.StatusUnauthorized, map[string]string{"error": "unauthenticated"})
        }
        return re.JSON(http.StatusOK, map[string]any{"user": info.Auth.Id})
    })
    return se.Next()
})
```

### Record lifecycle hook

```go
// Fires after a "posts" record is successfully created.
app.OnRecordAfterCreateSuccess("posts").BindFunc(func(e *core.RecordEvent) error {
    app.Logger().Info("post created", "id", e.Record.Id)
    return e.Next()  // always call Next() to continue the chain
})
```

### Cron job

```go
app.Cron().MustAdd("nightly_cleanup", "0 2 * * *", func() {
    app.Logger().Info("running nightly cleanup")
    // your logic here
})
```

### Sending email

```go
message := &mailer.Message{
    From:    mail.Address{Address: app.Settings().Meta.SenderAddress},
    To:      []mail.Address{{Address: "user@example.com"}},
    Subject: "Welcome to TowerShare",
    HTML:    "<p>Thanks for joining!</p>",
}
if err := app.NewMailClient().Send(message); err != nil {
    app.Logger().Error("mail send failed", "err", err)
}
```

---

## JS hooks vs Go modules

| | JS hooks (`pb_hooks/*.pb.js`) | Go modules (`pb_modules/`) |
|---|---|---|
| Hot reload | Yes | No (requires rebuild) |
| Performance | Interpreted | Compiled |
| Access to Go stdlib | No | Yes |
| Best for | Simple event logging, quick prototypes | Complex business logic, 3rd-party SDKs |

---

## Adding a migration in a module

Create a file in `backend/pb_migrations/` and use the `m.Register` pattern shown in `1_init_schema.go`. Migrations run automatically in `APP_ENV=development` and must be triggered manually in production:

```bash
./knowns.sh migrate
```

---

## Auth extension

PocketBase's built-in `users` collection handles email/password and OAuth2. To add custom fields:

1. Add them in a migration (`pb_migrations/`).
2. Pass them when creating a user from Swift (`AuthService.signUp`).
3. Read them via `pb.collection("users").getOne(id)`.

OAuth2 providers are configured in the PocketBase admin UI at `/_/` → Settings → Auth providers. No code changes required.
