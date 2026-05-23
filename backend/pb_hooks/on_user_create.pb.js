// PocketBase JS hooks run in the embedded JS runtime (goja).
// This file is hot-reloaded in dev — no server restart needed.
// For complex logic, prefer Go modules in pb_modules/.

onRecordAfterCreateSuccess((e) => {
    console.log("JS hook: user created", e.record.id);
    e.next();
}, "users");
