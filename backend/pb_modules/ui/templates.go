package ui

import "html/template"

// All templates are defined as Go string constants — no .html files.

const css = `
*,*::before,*::after{box-sizing:border-box;margin:0;padding:0}
:root{
  --bg:#0a0a0a;--surface:#141414;--surface2:#1e1e1e;--border:#2a2a2a;
  --text:#f0f0f0;--muted:#707070;--accent:#00b4d8;
  --green:#34c759;--orange:#ff9500;--red:#ff3b30;--radius:12px;
}
html,body{height:100%;background:var(--bg);color:var(--text);
  font-family:-apple-system,BlinkMacSystemFont,"Segoe UI",sans-serif;font-size:15px}
a{color:var(--accent)}
button,input[type=submit]{cursor:pointer;font-family:inherit;font-size:14px;
  border:none;border-radius:8px;padding:10px 18px;transition:opacity .15s}
button.primary,input[type=submit]{background:var(--accent);color:#fff;font-weight:600;width:100%}
button.primary:hover,input[type=submit]:hover{opacity:.85}
button.ghost{background:transparent;color:var(--accent);border:1px solid var(--accent)}
input[type=text],input[type=email],input[type=password],textarea,select{
  width:100%;background:var(--surface2);border:1px solid var(--border);border-radius:8px;
  padding:10px 14px;color:var(--text);font-family:inherit;font-size:15px;outline:none}
input:focus,textarea:focus,select:focus{border-color:var(--accent)}
textarea{resize:vertical;min-height:80px}
label{display:block;font-size:13px;color:var(--muted);margin-bottom:6px;font-weight:500}
.field{margin-bottom:18px}
.err{color:var(--red);font-size:13px;margin-top:8px}
.ok{color:var(--green);font-size:13px;margin-top:8px}
nav{background:var(--surface);border-bottom:1px solid var(--border);
  padding:14px 20px;display:flex;align-items:center;gap:12px}
nav .logo{font-size:18px;font-weight:700;color:var(--accent);flex:1}
nav .user{font-size:13px;color:var(--muted)}
.page{max-width:480px;margin:0 auto;padding:24px 20px}
.card{background:var(--surface);border:1px solid var(--border);
  border-radius:var(--radius);overflow:hidden;margin-bottom:12px}
.card-img{width:100%;height:180px;background:var(--surface2);
  display:flex;align-items:center;justify-content:center;font-size:40px}
.card-img img{width:100%;height:100%;object-fit:cover;display:block}
.card-body{padding:14px}
.badge{font-size:11px;font-weight:600;padding:3px 8px;border-radius:6px;text-transform:uppercase;display:inline-block;margin-bottom:8px}
.badge-lend{background:rgba(0,180,216,.15);color:var(--accent)}
.badge-give{background:rgba(52,199,89,.15);color:var(--green)}
.badge-request{background:rgba(255,149,0,.15);color:var(--orange)}
.card-title{font-weight:600;font-size:15px;margin-bottom:4px}
.card-meta{font-size:13px;color:var(--muted)}
.tabs{display:flex;gap:16px;margin-bottom:24px;border-bottom:1px solid var(--border);padding-bottom:0}
.tab{padding:8px 4px 10px;color:var(--muted);text-decoration:none;border-bottom:2px solid transparent;margin-bottom:-1px}
.tab.active{color:var(--accent);border-bottom-color:var(--accent)}
.filter-row{display:flex;gap:8px;margin-bottom:16px;flex-wrap:wrap}
.chip{padding:5px 12px;border-radius:20px;background:var(--surface);
  border:1px solid var(--border);color:var(--muted);text-decoration:none;font-size:13px}
.chip.active{background:var(--accent);border-color:var(--accent);color:#fff}
.fab{position:fixed;bottom:24px;right:calc(50% - 220px);width:52px;height:52px;
  border-radius:50%;background:var(--accent);color:#fff;font-size:24px;
  display:flex;align-items:center;justify-content:center;
  text-decoration:none;box-shadow:0 4px 20px rgba(0,180,216,.4)}
.preview-grid{display:grid;grid-template-columns:repeat(3,1fr);gap:8px;margin-bottom:18px}
.preview-grid img{width:100%;aspect-ratio:1;object-fit:cover;border-radius:8px}
h2{font-size:18px;font-weight:700;margin-bottom:16px}
.auth-wrap{display:flex;flex-direction:column;justify-content:center;min-height:100vh;
  max-width:380px;margin:0 auto;padding:32px 20px}
.auth-logo{text-align:center;margin-bottom:32px}
.auth-logo h1{font-size:28px;font-weight:800;color:var(--accent)}
.auth-logo p{color:var(--muted);margin-top:6px;font-size:14px}
.auth-switch{text-align:center;margin-top:16px;font-size:14px;color:var(--muted)}
.empty{text-align:center;padding:48px 0;color:var(--muted)}
.empty .icon{font-size:40px;margin-bottom:12px}
`

const baseOpen = `<!DOCTYPE html>
<html lang="en">
<head>
<meta charset="UTF-8">
<meta name="viewport" content="width=device-width,initial-scale=1">
<title>TowerShare{{if .Title}} — {{.Title}}{{end}}</title>
<style>` + css + `</style>
</head>
<body>`

const baseClose = `</body></html>`

// ── Page templates ─────────────────────────────────────────────────────────

const loginTmplSrc = baseOpen + `
<div class="auth-wrap">
  <div class="auth-logo"><h1>🏢 TowerShare</h1><p>Share resources with your building</p></div>
  <div class="tabs">
    <a class="tab active" href="/app/login">Sign In</a>
    <a class="tab" href="/app/signup">Sign Up</a>
  </div>
  {{if .Error}}<div class="err">{{.Error}}</div>{{end}}
  <form method="POST" action="/app/login">
    <div class="field"><label>Email</label>
      <input type="email" name="email" placeholder="you@example.com" required autocomplete="email"
             value="{{.Email}}"></div>
    <div class="field"><label>Password</label>
      <input type="password" name="password" placeholder="Password" required autocomplete="current-password"></div>
    <input type="submit" class="primary" value="Sign In">
  </form>
  <div class="auth-switch">No account? <a href="/app/signup">Sign up</a></div>
</div>
` + baseClose

const signupTmplSrc = baseOpen + `
<div class="auth-wrap">
  <div class="auth-logo"><h1>🏢 TowerShare</h1><p>Share resources with your building</p></div>
  <div class="tabs">
    <a class="tab" href="/app/login">Sign In</a>
    <a class="tab active" href="/app/signup">Sign Up</a>
  </div>
  {{if .Error}}<div class="err">{{.Error}}</div>{{end}}
  <form method="POST" action="/app/signup">
    <div class="field"><label>Display name</label>
      <input type="text" name="display_name" placeholder="Your name" required value="{{.DisplayName}}"></div>
    <div class="field"><label>Email</label>
      <input type="email" name="email" placeholder="you@example.com" required autocomplete="email"
             value="{{.Email}}"></div>
    <div class="field"><label>Password</label>
      <input type="password" name="password" placeholder="Min 8 characters" required minlength="8"
             autocomplete="new-password"></div>
    <input type="submit" class="primary" value="Create Account">
  </form>
  <div class="auth-switch">Already have an account? <a href="/app/login">Sign in</a></div>
</div>
` + baseClose

const feedTmplSrc = baseOpen + `
<nav>
  <span class="logo">TowerShare</span>
  <span class="user">{{.UserName}}</span>
  <a href="/app/logout" style="font-size:13px;color:var(--muted)">Sign out</a>
</nav>
<div class="page">
  <div style="display:flex;align-items:center;margin-bottom:16px">
    <h2 style="margin:0;flex:1">Listings</h2>
  </div>
  <div class="filter-row">
    <a class="chip{{if eq .Filter ""}} active{{end}}" href="/app/feed">All</a>
    <a class="chip{{if eq .Filter "lend"}} active{{end}}" href="/app/feed?cat=lend">Lend</a>
    <a class="chip{{if eq .Filter "give"}} active{{end}}" href="/app/feed?cat=give">Give</a>
    <a class="chip{{if eq .Filter "request"}} active{{end}}" href="/app/feed?cat=request">Request</a>
  </div>
  {{if .Listings}}
    {{range .Listings}}
    <div class="card">
      {{if .ImageURL}}
        <div class="card-img"><img src="{{.ImageURL}}" loading="lazy"></div>
      {{else}}
        <div class="card-img">📦</div>
      {{end}}
      <div class="card-body">
        <span class="badge badge-{{.Category}}">{{.Category}}</span>
        <div class="card-title">{{.Title}}</div>
        <div class="card-meta">🏢 {{.TowerName}} · {{.TimeAgo}}</div>
      </div>
    </div>
    {{end}}
  {{else}}
    <div class="empty"><div class="icon">📭</div><p>No listings yet.</p></div>
  {{end}}
</div>
<a class="fab" href="/app/listings/new" title="New listing">+</a>
` + baseClose

const newListingTmplSrc = baseOpen + `
<nav>
  <span class="logo">TowerShare</span>
  <a href="/app/feed" style="font-size:13px;color:var(--muted)">← Back</a>
</nav>
<div class="page">
  <h2>New Listing</h2>
  {{if .Error}}<div class="err">{{.Error}}</div>{{end}}
  {{if .Success}}<div class="ok">{{.Success}}</div>{{end}}
  <form method="POST" action="/app/listings/new" enctype="multipart/form-data">
    <div class="field"><label>Photos (up to 5)</label>
      <input type="file" name="images" accept="image/*" multiple style="color:var(--muted)"></div>
    <div class="field"><label>Title</label>
      <input type="text" name="title" placeholder="What are you sharing?" required maxlength="200"
             value="{{.Title}}"></div>
    <div class="field"><label>Description</label>
      <textarea name="description" placeholder="Details, condition...">{{.Description}}</textarea></div>
    <div class="field"><label>Category</label>
      <select name="category" required>
        <option value="">Select...</option>
        <option value="lend"{{if eq .Category "lend"}} selected{{end}}>Lend — borrow and return</option>
        <option value="give"{{if eq .Category "give"}} selected{{end}}>Give — free to keep</option>
        <option value="request"{{if eq .Category "request"}} selected{{end}}>Request — looking for something</option>
      </select></div>
    <div class="field"><label>Tower</label>
      <select name="tower" required>
        <option value="">Select tower...</option>
        {{range .Towers}}
          <option value="{{.Id}}"{{if eq $.TowerID .Id}} selected{{end}}>{{.Name}}</option>
        {{end}}
      </select></div>
    <input type="submit" class="primary" value="Post Listing">
  </form>
</div>
` + baseClose

// ── Compiled templates ─────────────────────────────────────────────────────

var (
	loginTmpl      = template.Must(template.New("login").Parse(loginTmplSrc))
	signupTmpl     = template.Must(template.New("signup").Parse(signupTmplSrc))
	feedTmpl       = template.Must(template.New("feed").Parse(feedTmplSrc))
	newListingTmpl = template.Must(template.New("new").Parse(newListingTmplSrc))
)
