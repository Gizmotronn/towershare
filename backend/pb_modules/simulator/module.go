// Package simulator serves a browser-based chat UI at /simulator.
// The entire page — HTML, CSS and JS — is defined as Go string constants.
// Only active when APP_ENV != "production".
package simulator

import (
	"github.com/labstack/echo/v5"
	"github.com/pocketbase/pocketbase"
	"github.com/pocketbase/pocketbase/core"
)

func Register(app *pocketbase.PocketBase) {
	app.OnBeforeServe().Add(func(e *core.ServeEvent) error {
		e.Router.GET("/simulator", func(c echo.Context) error {
			c.Response().Header().Set("Content-Type", "text/html; charset=utf-8")
			_, err := c.Response().Write([]byte(page))
			return err
		})
		return nil
	})
}

// page is the complete chat simulator — HTML structure, CSS styles, and JS
// logic — defined as a Go string constant. No .html files anywhere.
const page = doctype + styles + body + scripts

const doctype = `<!DOCTYPE html>
<html lang="en">
<head>
<meta charset="UTF-8">
<meta name="viewport" content="width=device-width,initial-scale=1,viewport-fit=cover">
<title>TowerShare — Chat</title>
`

const styles = `<style>
*,*::before,*::after{box-sizing:border-box;margin:0;padding:0}
:root{
  --bg:#0a0a0a;--surface:#141414;--surface2:#1e1e1e;--border:#252525;
  --accent:#00b4d8;--accent-dim:rgba(0,180,216,.12);
  --text:#f2f2f7;--muted:#636366;--sent-bg:#00b4d8;--recv-bg:#1c1c1e;
  --green:#34c759;--red:#ff3b30;--orange:#ff9500;
  --safe-bottom:env(safe-area-inset-bottom,0px);
}
html,body{height:100%;overflow:hidden;background:var(--bg);color:var(--text);
  font-family:-apple-system,BlinkMacSystemFont,"Segoe UI",Helvetica,sans-serif;
  font-size:16px;-webkit-font-smoothing:antialiased}

/* Shell */
#shell{display:flex;flex-direction:column;height:100vh;max-width:480px;
  margin:0 auto;background:var(--bg);position:relative;
  border-left:1px solid var(--border);border-right:1px solid var(--border)}

/* Header */
#header{
  display:flex;align-items:center;gap:10px;
  padding:12px 16px 10px;
  padding-top:calc(12px + env(safe-area-inset-top,0px));
  background:rgba(10,10,10,.85);
  backdrop-filter:blur(20px);-webkit-backdrop-filter:blur(20px);
  border-bottom:1px solid var(--border);flex-shrink:0;
}
.avatar{width:36px;height:36px;border-radius:50%;flex-shrink:0;
  background:linear-gradient(135deg,var(--accent),#0077b6);
  display:flex;align-items:center;justify-content:center;font-size:18px}
.hdr-info{flex:1;min-width:0}
.hdr-name{font-weight:600;font-size:15px;color:var(--text)}
.hdr-status{font-size:12px;color:var(--green);display:flex;align-items:center;gap:4px}
.hdr-status::before{content:'';width:6px;height:6px;border-radius:50%;background:var(--green);flex-shrink:0}
.channel-pills{display:flex;gap:6px;flex-shrink:0}
.pill{padding:3px 10px;border-radius:12px;font-size:11px;font-weight:500;
  border:1px solid var(--border);color:var(--muted);cursor:pointer;background:transparent;
  transition:all .15s;white-space:nowrap}
.pill.active{background:var(--accent);border-color:var(--accent);color:#fff}

/* Messages */
#messages{flex:1;overflow-y:auto;overflow-x:hidden;padding:12px 16px;
  display:flex;flex-direction:column;gap:2px;scroll-behavior:smooth}
#messages::-webkit-scrollbar{width:0}

/* Bubble rows */
.row{display:flex;flex-direction:column;margin-bottom:8px}
.row.sent{align-items:flex-end}
.row.recv{align-items:flex-start}
.sender{font-size:11px;color:var(--muted);margin-bottom:3px;padding:0 4px}

/* Bubbles */
.bubble{
  max-width:78%;padding:9px 13px;
  font-size:15px;line-height:1.45;
  white-space:pre-wrap;word-break:break-word;
}
.sent .bubble{
  background:var(--sent-bg);color:#fff;
  border-radius:18px 18px 4px 18px;
}
.recv .bubble{
  background:var(--recv-bg);color:var(--text);
  border-radius:18px 18px 18px 4px;
}
.bubble-time{font-size:10px;color:var(--muted);margin-top:3px;padding:0 4px}

/* Action chips */
.actions{display:flex;flex-wrap:wrap;gap:7px;margin-top:6px;padding:0 2px}
.chip{
  padding:7px 14px;border-radius:18px;
  border:1.5px solid var(--accent);background:var(--accent-dim);
  color:var(--accent);font-size:13px;font-weight:500;cursor:pointer;
  transition:background .12s,transform .1s,opacity .2s;
  white-space:nowrap;
}
.chip:hover{background:rgba(0,180,216,.22)}
.chip:active{transform:scale(.96)}
.chip.used{opacity:.35;pointer-events:none;border-color:var(--border);color:var(--muted);background:transparent}

/* Typing indicator */
#typing{display:none;align-items:flex-start;gap:0;margin-bottom:8px;padding:0 2px}
#typing .bubble{padding:10px 14px;display:flex;gap:5px;align-items:center}
.dot{width:7px;height:7px;border-radius:50%;background:var(--muted);
  animation:bounce 1.2s infinite ease-in-out}
.dot:nth-child(2){animation-delay:.2s}
.dot:nth-child(3){animation-delay:.4s}
@keyframes bounce{0%,60%,100%{transform:translateY(0)}30%{transform:translateY(-6px)}}

/* Input area */
#composer{
  padding:10px 12px;
  padding-bottom:calc(10px + var(--safe-bottom));
  background:rgba(10,10,10,.9);
  backdrop-filter:blur(20px);-webkit-backdrop-filter:blur(20px);
  border-top:1px solid var(--border);
  display:flex;align-items:flex-end;gap:8px;flex-shrink:0;
}
#input{
  flex:1;background:var(--surface2);border:1px solid var(--border);
  border-radius:20px;padding:9px 15px;color:var(--text);
  font-family:inherit;font-size:15px;outline:none;
  max-height:120px;overflow-y:auto;resize:none;
  line-height:1.4;
  transition:border-color .15s;
}
#input:focus{border-color:var(--accent)}
#input::placeholder{color:var(--muted)}
#send{
  width:36px;height:36px;border-radius:50%;
  background:var(--accent);border:none;cursor:pointer;
  display:flex;align-items:center;justify-content:center;
  flex-shrink:0;transition:opacity .15s;
}
#send:disabled{opacity:.35;cursor:default}
#send svg{fill:#fff;width:16px;height:16px}
</style>
`

const body = `</head>
<body>
<div id="shell">

  <div id="header">
    <div class="avatar">🏢</div>
    <div class="hdr-info">
      <div class="hdr-name">TowerShare Bot</div>
      <div class="hdr-status" id="status-text">online</div>
    </div>
    <div class="channel-pills">
      <button class="pill active" data-ch="whatsapp" onclick="setChannel(this)">WhatsApp</button>
      <button class="pill" data-ch="imessage" onclick="setChannel(this)">iMessage</button>
      <button class="pill" data-ch="sms" onclick="setChannel(this)">SMS</button>
    </div>
  </div>

  <div id="messages">
    <div id="typing" class="row recv">
      <div class="bubble recv">
        <span class="dot"></span><span class="dot"></span><span class="dot"></span>
      </div>
    </div>
  </div>

  <div id="composer">
    <textarea id="input" rows="1" placeholder="Message"
      oninput="autoResize(this)"
      onkeydown="onKey(event)"></textarea>
    <button id="send" onclick="submit()" disabled>
      <svg viewBox="0 0 24 24"><path d="M2 21l21-9L2 3v7l15 2-15 2z"/></svg>
    </button>
  </div>

</div>
`

const scripts = `<script>
'use strict';

let channel = 'whatsapp';
let activeChipRow = null;
let busy = false;
const from = '+61400000001';

const $ = id => document.getElementById(id);
const msgs = $('messages');
const inp = $('input');

// ── Utilities ──────────────────────────────────────────────────────────────

function now() {
  return new Date().toLocaleTimeString([], {hour:'2-digit', minute:'2-digit'});
}

function scrollBottom() {
  msgs.scrollTo({top: msgs.scrollHeight, behavior: 'smooth'});
}

function setStatus(text, ok = true) {
  const el = $('status-text');
  el.textContent = text;
  el.style.color = ok ? 'var(--green)' : 'var(--muted)';
}

function autoResize(el) {
  el.style.height = 'auto';
  el.style.height = Math.min(el.scrollHeight, 120) + 'px';
  $('send').disabled = !el.value.trim();
}

function setChannel(btn) {
  document.querySelectorAll('.pill').forEach(p => p.classList.remove('active'));
  btn.classList.add('active');
  channel = btn.dataset.ch;
}

// ── Chip management ────────────────────────────────────────────────────────

function retireChips() {
  if (activeChipRow) {
    activeChipRow.querySelectorAll('.chip').forEach(c => c.classList.add('used'));
    activeChipRow = null;
  }
}

// ── Bubble rendering ───────────────────────────────────────────────────────

function appendSent(text) {
  retireChips();
  const row = document.createElement('div');
  row.className = 'row sent';
  row.innerHTML =
    '<div class="bubble">' + esc(text) + '</div>' +
    '<div class="bubble-time">' + now() + '</div>';
  msgs.insertBefore(row, $('typing'));
  scrollBottom();
}

function appendRecv(text, actions) {
  const row = document.createElement('div');
  row.className = 'row recv';
  row.innerHTML =
    '<div class="sender">Bot</div>' +
    '<div class="bubble">' + esc(text) + '</div>' +
    '<div class="bubble-time">' + now() + '</div>';
  msgs.insertBefore(row, $('typing'));

  if (actions && actions.length) {
    const chips = document.createElement('div');
    chips.className = 'actions';
    actions.forEach(a => {
      const btn = document.createElement('button');
      btn.className = 'chip';
      btn.textContent = a.label;
      btn.onclick = () => chipClick(a.label, a.value);
      chips.appendChild(btn);
    });
    msgs.insertBefore(chips, $('typing'));
    activeChipRow = chips;
  }
  scrollBottom();
}

function esc(str) {
  return str
    .replace(/&/g,'&amp;').replace(/</g,'&lt;').replace(/>/g,'&gt;')
    .replace(/"/g,'&quot;').replace(/\n/g,'<br>');
}

// ── Typing indicator ───────────────────────────────────────────────────────

function showTyping() {
  $('typing').style.display = 'flex';
  scrollBottom();
}
function hideTyping() {
  $('typing').style.display = 'none';
}

// ── API call ───────────────────────────────────────────────────────────────

async function callBot(body) {
  if (busy) return;
  busy = true;
  $('send').disabled = true;
  setStatus('typing…', false);
  showTyping();

  try {
    const res = await fetch('/api/towershare/message', {
      method: 'POST',
      headers: {'Content-Type': 'application/json'},
      body: JSON.stringify({from, body, channel}),
    });
    hideTyping();
    const data = await res.json();
    if (data.reply)   appendRecv(data.reply, data.actions);
    else if (data.error) appendRecv('⚠️  ' + data.error, []);
    setStatus('online');
  } catch (e) {
    hideTyping();
    appendRecv('⚠️  Could not reach PocketBase. Is it running?', []);
    setStatus('offline', false);
  } finally {
    busy = false;
    $('send').disabled = !inp.value.trim();
  }
}

// ── User actions ───────────────────────────────────────────────────────────

async function submit() {
  const text = inp.value.trim();
  if (!text || busy) return;
  inp.value = '';
  inp.style.height = 'auto';
  $('send').disabled = true;
  appendSent(text);
  await callBot(text);
  inp.focus();
}

async function chipClick(label, value) {
  retireChips();
  appendSent(label);
  await callBot(value);
  inp.focus();
}

function onKey(e) {
  if (e.key === 'Enter' && !e.shiftKey) {
    e.preventDefault();
    submit();
  }
}

// ── Boot ───────────────────────────────────────────────────────────────────

window.addEventListener('load', async () => {
  inp.focus();
  await callBot('');   // empty body → welcome message
});
</script>
</body></html>`
