# PRD — Shared Hard Waste

**Hackathon theme:** Community & citizen initiatives — resilience built from the ground up.

## 1. One-liner

A building-scale platform that lets Melbourne apartment residents **pool their unused hard waste collection entitlements** and **give away unwanted furniture** to neighbours — turning a wasted council allowance into a circular, community-owned resource.

## 2. Problem

The City of Melbourne gives every household **one hard waste collection per financial year, up to 1m³**. In apartment buildings this creates three failures:

1. **Wasted entitlements.** Most residents never book their collection. The allowance expires silently each July.
2. **Stranded waste.** Residents with a couch, mattress, or bed frame they can't fit into 1m³ have no legal option — they wait months, or they dump.
3. **Illegal dumping.** Bulky items left on nature strips, in laneways, and beside bins are one of the council's most expensive and visible amenity problems.

Meanwhile, much of the "waste" is still usable. A bookshelf one resident is throwing out is a bookshelf the new tenant on level 4 was about to buy.

## 3. Why this fits the theme

The solution is **resident-driven, building-scale, and circular**:

- Residents self-organise around a shared resource (their building's pooled entitlements).
- Reuse happens **inside the building first** — items never need to leave before finding a new owner.
- What can't be reused gets disposed of legally and efficiently, using entitlements that would otherwise have been wasted.
- The council is a partner, not a gatekeeper.

## 4. Users

**Primary persona — Apartment Resident (Renter or Owner).**
Lives in a multi-unit building. Has a couch to get rid of, or wants the neighbour's couch, or has an entitlement they're not using. Doesn't want to deal with the council website.

**Secondary persona — Building Coordinator.**
A resident, owners' corp member, or building manager who oversees their building's pool. Confirms residency, books the actual collection, schedules the curbside drop window.

**Tertiary persona — City of Melbourne Waste Team.**
The council partner. Receives consolidated, verified bookings instead of fragmented individual ones. Sees data on diversion and dumping reduction.

### 4.1 User research findings

Three findings from talking to apartment residents shape the product:

- **R1 — Residents resist creating yet another account.** Asked about login, the response was consistently "do I have to?" Implication: onboarding must be passwordless and one-step. Building-issued code + unit confirmation + (optionally) a magic link. No password screen.
- **R2 — Owners' corps and building managers already ship resident apps.** Wumbo, Residential, and Home App are commonly installed in Melbourne apartment estates and cover maintenance requests, bills, shared-space booking, and message boards. Residents will not install a second app for an occasional behaviour like hard waste. Implication: the long-term destination is **embedded inside / SSO'd from** one of those incumbent apps, not a standalone install. For the hackathon we ship standalone, but the data model and auth assume "delegated identity from the building's app" is the v1 path.
- **R3 — WhatsApp is the informal incumbent, and it's failing.** Many buildings have a resident WhatsApp group already coordinating give-aways. Two consistent complaints: notification noise (every message pings everyone) and findability (listings scroll out of view within hours and can't be searched). Implication: the product must beat WhatsApp on the two axes residents actually complained about — **structured listings with status states** (reserved / claimed / gone, so dead posts disappear) and **quiet-by-default notifications** (opt-in by category, not broadcast).

### 4.2 Existing solutions & positioning

| Tool | What it does | Where we fit |
|---|---|---|
| Wumbo / Residential / Home App | OC-managed: maintenance, bills, shared-space booking, building message board | **Complementary.** One feature deep on hard waste + reuse. Future: embed as a module or SSO in. Do not compete on bills or maintenance. |
| Resident WhatsApp groups | Informal give-away coordination | **Direct displacement.** Same job-to-be-done, fixes the noise + findability failures residents already named. |
| Council hard waste portal | Per-household booking | **Sits behind us.** We aggregate bookings into the council portal; resident never visits it directly. |

The pitch sentence: *"What your building's WhatsApp group is already trying to do — but searchable, statusful, and plugged into council's entitlement system."*

## 5. Goals

- **G1.** Make it trivially easy for a resident to donate an unused entitlement to their building's pool.
- **G2.** Make it trivially easy for a resident to list a piece of furniture as "free to a good home" — visible to their building first.
- **G3.** Let the building coordinator book a single consolidated hard waste collection that uses the pooled entitlements for whatever's left over.
- **G4.** Give the council a clean, auditable record of which entitlement was used by which household.

## 6. Non-goals (for the hackathon MVP)

- Payments, money, or escrow.
- Inter-building or city-wide marketplaces.
- Houses / non-apartment dwellings.
- Logistics beyond the building's existing curbside collection point.
- A native mobile app (web responsive is enough for the demo).
- Traditional username + password login (see R1 — passwordless only).
- Competing with Wumbo / Residential / Home App on bills, maintenance, or shared-space booking (see R2 — stay one-feature-deep).

## 7. Core user flows

### Flow A — Donate an entitlement
1. Resident logs in, sees "You have 1m³ of hard waste collection available this financial year."
2. Taps **Donate to building pool**.
3. Confirms — entitlement now sits in the building's shared pool, visible to all residents.

### Flow B — Give away furniture
1. Resident taps **Give away an item**.
2. Snaps a photo, picks a category (couch, mattress, table, …), adds rough dimensions.
3. Item posts to the building's internal board for 48 hours before, optionally, going to a wider feed.
4. Another resident taps **Claim** → they coordinate pickup in-app (chat thread, no addresses needed beyond unit number).
5. If unclaimed, item rolls into the next scheduled hard waste collection.

### Flow C — Coordinator books a collection
1. Coordinator sees the building's queue: claimed items (resolved), unclaimed items (need disposal), available entitlements in the pool.
2. Taps **Book collection** — the app computes which residents' entitlements are being drawn down and submits to council.
3. Each resident whose entitlement was used gets a notification and an audit record.

## 8. MVP feature list

| # | Feature | Priority |
|---|---|---|
| F1 | Passwordless onboarding — building code + unit confirmation, no password (see R1) | P0 |
| F2 | Entitlement balance display + "Donate to pool" action | P0 |
| F3 | Building pool view (entitlements available, contributors) | P0 |
| F4 | Furniture listing — photo, category, dimensions, status | P0 |
| F5 | Building-internal feed + claim flow, with status states (available / reserved / claimed / gone) so dead listings disappear (see R3) | P0 |
| F6 | Coordinator dashboard: queue, pool, book collection | P0 |
| F7 | Council submission stub (mocked API, real-looking payload) | P0 |
| F8 | Eco Footprint view — kg diverted from landfill, estimated CO₂ saved, monthly chart by category | P0 |
| F9 | Audit trail per entitlement (who donated, who used it, when) | P1 |
| F10 | In-app chat for pickup coordination | P1 |
| F11 | Quiet-by-default notifications — opt-in by category, not broadcast (see R3) | P1 |
| F12 | Building leaderboard — residents ranked by items diverted + entitlements donated | P1 |
| F13 | Badges — first listing, first claim, first donor, streaks | P1 |

## 9. Council integration (P0 for the pitch, mocked for the build)

The product assumes a partnership with the **City of Melbourne**. For the demo we mock the council-side API but design the contract realistically:

- `GET /entitlements/{household_id}` → balance + expiry.
- `POST /collections` → submit a booking referencing one or more `entitlement_id`s.
- `POST /entitlements/{id}/transfer` → reassign an entitlement to the building pool (still attributed to the originating household for compliance).

The pitch position: this gives council **better data, fewer fragmented bookings, and a measurable reduction in illegal dumping** in exchange for the API surface.

## 10. Success metrics

For the demo:
- ≥ 80% of seeded residents in the demo building can donate an entitlement in under 30 seconds.
- A coordinator can go from "see queue" to "submit booking" in under 60 seconds.
- The audit trail correctly attributes every pooled entitlement back to its original household.

For a real pilot (the slide-ware metric):
- % of issued entitlements that get used (baseline vs. pilot).
- Items claimed within-building / total items listed.
- Reported illegal dumping incidents in the pilot building's block.

## 11. Out of scope / future

- City-wide furniture marketplace once a building has saturated demand.
- Integration with op-shops and charity pickups (Sacred Heart Mission, Vinnies) for items worth refurbishing.
- Extending to other councils (Yarra, Port Phillip, Stonnington) — same model, different API.
- **Embedding inside incumbent building apps (Wumbo, Residential, Home App)** via SSO / module — the eventual distribution play per R2.

## 12. Tech notes

- Stack: existing Vite + React + TypeScript + Tailwind scaffold in this repo.
- Persistence: in-memory / localStorage for the demo; document the data model so council integration is plausible.
- Auth: passwordless by design (see R1). For the hackathon: tap a seeded resident from a building-code splash — no password screen, no email roundtrip. For v1 real: building code + unit + magic link. Data model assumes a future identity provider (council, OC, or incumbent building app) can vouch for the resident.
- Council API: TypeScript interface + fake implementation that returns deterministic responses for the demo.

## 13. Demo script (3 minutes)

1. **30s — the problem.** Show council's hard waste page, then a screenshot of a noisy resident WhatsApp group with give-away posts buried under chat. "Every household gets 1m³. Most apartments never use it. The ones that need more dump on the street. The ones already trying to share — do it in WhatsApp, and it doesn't work."
2. **60s — Flow A + B.** Tap-to-enter as a resident (no login). Donate an entitlement. Another resident lists a couch. A third resident claims it — listing flips to "claimed" and disappears from the active feed. The couch never enters landfill.
3. **60s — Flow C.** Coordinator books a consolidated collection for the unclaimed items, drawing from the pool. Show the council payload. Cut to the Eco Footprint view — kg diverted ticks up, building leaderboard reshuffles.
4. **30s — the pitch.** "Same allowance, no new infrastructure, less dumping. Built by residents, partnered with council. Designed to slot into the building apps they already have."
