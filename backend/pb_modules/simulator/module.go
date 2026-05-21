// Package simulator serves a browser chat UI at /simulator.
// All markup, styles, and scripts are Go string constants — no .html files.
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

const page = head + body + script

// ── head ──────────────────────────────────────────────────────────────────

const head = `<!DOCTYPE html>
<html lang="en">
<head>
<meta charset="UTF-8">
<meta name="viewport" content="width=device-width,initial-scale=1,viewport-fit=cover">
<title>TowerShare</title>
<style>
*,*::before,*::after{box-sizing:border-box;margin:0;padding:0}

:root {
  --bg:       #0a0a0a;
  --surface:  #161616;
  --surface2: #1e1e1e;
  --border:   #262626;
  --text:     #f2f2f7;
  --muted:    #636366;
  --accent:   #00b4d8;
  --adim:     rgba(0,180,216,.1);
  --sent:     #00b4d8;
  --recv:     #1c1c1e;
  --green:    #34c759;
}

html, body {
  height: 100%;
  background: var(--bg);
  color: var(--text);
  font-family: -apple-system, BlinkMacSystemFont, "Segoe UI", sans-serif;
  font-size: 16px;
  -webkit-font-smoothing: antialiased;
  overflow: hidden;
}

/* ── Shell ── */
#shell {
  display: flex;
  flex-direction: column;
  height: 100dvh;
  max-width: 480px;
  margin: 0 auto;
  border-inline: 1px solid var(--border);
}

/* ── Header ── */
#hdr {
  flex-shrink: 0;
  display: flex;
  align-items: center;
  gap: 10px;
  padding: 10px 14px;
  padding-top: max(10px, env(safe-area-inset-top));
  background: rgba(10,10,10,.9);
  backdrop-filter: blur(24px);
  border-bottom: 1px solid var(--border);
}
.avatar {
  width: 38px; height: 38px;
  border-radius: 50%;
  background: linear-gradient(135deg, #00b4d8, #0077b6);
  display: flex; align-items: center; justify-content: center;
  font-size: 18px; flex-shrink: 0;
}
.hdr-text { flex: 1; min-width: 0; }
.hdr-name { font-weight: 600; font-size: 15px; }
.hdr-sub  { font-size: 12px; color: var(--green); display: flex; align-items: center; gap: 4px; }
.hdr-sub::before {
  content: ''; width: 6px; height: 6px;
  border-radius: 50%; background: var(--green);
}
.channels { display: flex; gap: 5px; }
.ch {
  padding: 3px 9px; border-radius: 12px; font-size: 11px; font-weight: 500;
  border: 1px solid var(--border); color: var(--muted);
  background: transparent; cursor: pointer; transition: all .15s;
}
.ch.on { background: var(--accent); border-color: var(--accent); color: #fff; }

/* ── Messages ── */
#msgs {
  flex: 1;
  overflow-y: auto;
  padding: 14px 14px 6px;
  display: flex;
  flex-direction: column;
  scroll-behavior: smooth;
}
#msgs::-webkit-scrollbar { width: 0; }

/* ── Rows ── */
.row {
  display: flex;
  flex-direction: column;
  margin-bottom: 3px;
}
.row.sent  { align-items: flex-end;   margin-bottom: 8px; }
.row.recv  { align-items: flex-start; }

.meta { font-size: 11px; color: var(--muted); padding: 0 4px; margin-bottom: 3px; }
.time { font-size: 10px; color: var(--muted); padding: 2px 4px 0; }

/* ── Bubbles ── */
.bub {
  max-width: 78%;
  padding: 9px 13px;
  font-size: 15px;
  line-height: 1.5;
  white-space: pre-wrap;
  word-break: break-word;
}
.sent .bub {
  background: var(--sent);
  color: #fff;
  border-radius: 18px 18px 4px 18px;
}
.recv .bub {
  background: var(--recv);
  color: var(--text);
  border-radius: 18px 18px 18px 4px;
}

/* ── Option bubbles ── */
/* Options are grouped below a bot message. Each is its own tappable bubble. */
.opts {
  display: flex;
  flex-direction: column;
  align-items: flex-start;
  gap: 6px;
  margin: 4px 0 10px;
}
.opt {
  max-width: 82%;
  padding: 9px 13px;
  border-radius: 18px 18px 18px 4px;
  border: 1.5px solid var(--accent);
  color: var(--accent);
  background: var(--adim);
  font-size: 15px;
  line-height: 1.4;
  cursor: pointer;
  display: flex;
  align-items: center;
  justify-content: space-between;
  gap: 10px;
  transition: background .12s, transform .08s, opacity .2s, border-color .15s;
  user-select: none;
}
.opt:hover  { background: rgba(0,180,216,.18); }
.opt:active { transform: scale(.98); }
.opt .arr   { opacity: .6; flex-shrink: 0; font-size: 14px; }

/* States after tapping */
.opt.chosen {
  background: var(--accent);
  color: #fff;
  border-color: var(--accent);
}
.opt.chosen .arr { opacity: 1; }
.opt.fade {
  opacity: 0;
  transform: translateX(-6px) scale(.97);
  pointer-events: none;
}

/* ── Typing indicator ── */
#typing {
  display: none;
  align-items: flex-start;
  margin-bottom: 10px;
}
#typing .bub {
  padding: 10px 14px;
  display: flex;
  gap: 5px;
  align-items: center;
}
.dot {
  width: 7px; height: 7px;
  border-radius: 50%;
  background: var(--muted);
  animation: bob 1.2s ease-in-out infinite;
}
.dot:nth-child(2) { animation-delay: .2s; }
.dot:nth-child(3) { animation-delay: .4s; }
@keyframes bob {
  0%,60%,100% { transform: translateY(0); }
  30%          { transform: translateY(-6px); }
}

/* ── Composer ── */
#composer {
  flex-shrink: 0;
  display: flex;
  align-items: flex-end;
  gap: 8px;
  padding: 10px 12px;
  padding-bottom: max(10px, env(safe-area-inset-bottom));
  background: rgba(10,10,10,.9);
  backdrop-filter: blur(24px);
  border-top: 1px solid var(--border);
}
#inp {
  flex: 1;
  background: var(--surface2);
  border: 1px solid var(--border);
  border-radius: 20px;
  padding: 9px 15px;
  color: var(--text);
  font-family: inherit;
  font-size: 15px;
  outline: none;
  resize: none;
  max-height: 100px;
  overflow-y: auto;
  line-height: 1.4;
  transition: border-color .15s;
}
#inp:focus        { border-color: var(--accent); }
#inp::placeholder { color: var(--muted); }
#sendbtn {
  width: 36px; height: 36px;
  border-radius: 50%;
  background: var(--accent);
  border: none; cursor: pointer; flex-shrink: 0;
  display: flex; align-items: center; justify-content: center;
  transition: opacity .15s;
}
#sendbtn:disabled { opacity: .35; cursor: default; }
#sendbtn svg { fill: #fff; width: 16px; height: 16px; }
</style>
</head>
`

// ── body ──────────────────────────────────────────────────────────────────

const body = `<body>
<div id="shell">

  <div id="hdr">
    <div class="avatar">🏢</div>
    <div class="hdr-text">
      <div class="hdr-name">TowerShare</div>
      <div class="hdr-sub" id="sub">online</div>
    </div>
    <div class="channels">
      <button class="ch on"  data-ch="whatsapp" onclick="setCh(this)">WhatsApp</button>
      <button class="ch"     data-ch="imessage" onclick="setCh(this)">iMessage</button>
      <button class="ch"     data-ch="sms"      onclick="setCh(this)">SMS</button>
    </div>
  </div>

  <div id="msgs">
    <!-- typing indicator always lives at the bottom of the list -->
    <div id="typing" class="row recv">
      <div class="bub recv">
        <span class="dot"></span>
        <span class="dot"></span>
        <span class="dot"></span>
      </div>
    </div>
  </div>

  <div id="composer">
    <textarea id="inp" rows="1" placeholder="Message"
              oninput="grow(this)" onkeydown="onkey(event)"></textarea>
    <button id="sendbtn" disabled onclick="submit()">
      <svg viewBox="0 0 24 24"><path d="M2 21l21-9L2 3v7l15 2-15 2z"/></svg>
    </button>
  </div>

</div>
`

// ── script ────────────────────────────────────────────────────────────────

const script = `<script>
'use strict';

const msgs   = document.getElementById('msgs');
const inp    = document.getElementById('inp');
const typing = document.getElementById('typing');

let channel = 'whatsapp';
let busy    = false;
const from  = '+61400000001';

// ── utils ─────────────────────────────────────────────────────────────────

function esc(s) {
  return s.replace(/&/g,'&amp;').replace(/</g,'&lt;').replace(/>/g,'&gt;')
          .replace(/\n/g,'<br>');
}
function now() {
  return new Date().toLocaleTimeString([],{hour:'2-digit',minute:'2-digit'});
}
function scroll() {
  msgs.scrollTo({top: msgs.scrollHeight, behavior:'smooth'});
}
function setStatus(t, ok) {
  const el = document.getElementById('sub');
  el.textContent = t;
  el.style.color = ok === false ? 'var(--muted)' : 'var(--green)';
}
function grow(el) {
  el.style.height = 'auto';
  el.style.height = Math.min(el.scrollHeight, 100) + 'px';
  document.getElementById('sendbtn').disabled = !el.value.trim();
}
function setCh(btn) {
  document.querySelectorAll('.ch').forEach(b => b.classList.remove('on'));
  btn.classList.add('on');
  channel = btn.dataset.ch;
}

// ── option management ─────────────────────────────────────────────────────

let liveOpts = null; // the current .opts container

// Retire (fade out) the current live options group without selecting anything.
function retireOpts() {
  if (!liveOpts) return;
  liveOpts.querySelectorAll('.opt').forEach(o => o.classList.add('fade'));
  const old = liveOpts;
  setTimeout(() => old.remove(), 220);
  liveOpts = null;
}

// Called when user taps an option bubble.
function pickOpt(label, value, optEl, group) {
  // Highlight chosen, fade the rest.
  optEl.classList.add('chosen');
  group.querySelectorAll('.opt').forEach(o => {
    if (o !== optEl) o.classList.add('fade');
  });
  liveOpts = null;

  setTimeout(() => {
    group.remove();
    addSent(label);
    call(value);
  }, 180);
}

// Render a list of options as tappable bubbles below the last bot message.
function addOpts(actions) {
  retireOpts();

  const group = document.createElement('div');
  group.className = 'opts';

  actions.forEach(a => {
    const o = document.createElement('div');
    o.className = 'opt';
    o.innerHTML = '<span>' + esc(a.label) + '</span><span class="arr">›</span>';
    o.onclick = () => pickOpt(a.label, a.value, o, group);
    group.appendChild(o);
  });

  msgs.insertBefore(group, typing);
  liveOpts = group;
  scroll();
}

// ── bubble rendering ──────────────────────────────────────────────────────

function addSent(text) {
  retireOpts();
  const row = document.createElement('div');
  row.className = 'row sent';
  row.innerHTML =
    '<div class="bub">' + esc(text) + '</div>' +
    '<div class="time">' + now() + '</div>';
  msgs.insertBefore(row, typing);
  scroll();
}

function addRecv(text, actions) {
  const row = document.createElement('div');
  row.className = 'row recv';
  row.innerHTML =
    '<div class="meta">Bot · ' + now() + '</div>' +
    '<div class="bub">' + esc(text) + '</div>';
  msgs.insertBefore(row, typing);

  if (actions && actions.length) {
    addOpts(actions);
  }
  scroll();
}

// ── typing indicator ──────────────────────────────────────────────────────

function showTyping()  { typing.style.display = 'flex'; scroll(); }
function hideTyping()  { typing.style.display = 'none'; }

// ── API ───────────────────────────────────────────────────────────────────

async function call(body) {
  if (busy) return;
  busy = true;
  document.getElementById('sendbtn').disabled = true;
  setStatus('typing…', false);
  showTyping();

  try {
    const res  = await fetch('/api/towershare/message', {
      method:  'POST',
      headers: {'Content-Type':'application/json'},
      body:    JSON.stringify({from, body, channel}),
    });
    const data = await res.json();
    hideTyping();
    if (data.reply) addRecv(data.reply, data.actions);
    else if (data.error) addRecv('⚠️  ' + data.error, []);
    setStatus('online', true);
  } catch (e) {
    hideTyping();
    addRecv('⚠️  Could not reach PocketBase. Is it running?\n\nmake up  or  make dev-backend', []);
    setStatus('offline', false);
  }

  busy = false;
  document.getElementById('sendbtn').disabled = !inp.value.trim();
}

// ── user actions ──────────────────────────────────────────────────────────

async function submit() {
  const text = inp.value.trim();
  if (!text || busy) return;
  inp.value = '';
  inp.style.height = 'auto';
  document.getElementById('sendbtn').disabled = true;
  addSent(text);
  await call(text);
  inp.focus();
}

function onkey(e) {
  if (e.key === 'Enter' && !e.shiftKey) { e.preventDefault(); submit(); }
}

// ── boot ──────────────────────────────────────────────────────────────────

window.addEventListener('load', () => {
  inp.focus();
  call(''); // empty body → welcome message + first options
});
</script>
</body></html>`
