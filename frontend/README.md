# Brunswick Heights — Hard Waste Hub

A community-built prototype for an apartment block in Melbourne to share
hard-waste collection entitlements and keep usable furniture out of landfill.

Built for the **Community & Citizen Initiatives** hackathon theme — resilience
from the ground up, residents driving their own circular-economy solutions.

## The problem

Melbourne residents each receive **one 1m³ hard-waste collection** per
financial year. Two failure modes follow:

1. Residents waste their entitlement when they have nothing bulky to throw out.
2. Residents with too many items exceed their entitlement and contribute to
   the city's illegal-dumping problem.

Apartment blocks concentrate both: many neighbours, isolated entitlements, no
coordination.

## The solution

A web app for residents of a single block that lets them:

- **Give furniture away** to neighbours before it becomes waste.
- **Share unused m³** of their annual entitlement with a neighbour who needs
  more.
- **Pool a single collection day** for the building, drawing against everyone's
  combined entitlement.

## What's in the prototype

- Public **welcome** + **printable QR poster** for the building lobby
- Fake **login** with a one-click persona switcher across 12 seeded residents
- **Dashboard** with entitlement bar, next collection day, and quick stats
- **Furniture marketplace** with post / claim / withdraw flow
- **Entitlement sharing** marketplace with offer + claim
- **Building collection day** with a pooled-m³ meter and per-resident grouping
- **Profile** showing your listings, balance, and pool contributions

All data is in-memory and resets on refresh — perfect for hackathon demoing
two sides of a swap by switching personas.

## Tech

- Vite + React 19 + TypeScript
- React Router v6
- Tailwind CSS v3 with a small in-repo shadcn-style component kit
- `lucide-react` icons
- QR codes via the `api.qrserver.com` public endpoint

## Running locally

```bash
npm install
npm run dev
```

Open <http://localhost:5173>. The dashboard is at `/`; the public lobby
landing is at `/welcome` and the printable poster is at `/qr-poster`.

## Demo path (2-minute walkthrough)

1. From the login screen, pick **Amelia Tran** in the demo dropdown.
2. Go to **Share m³** → claim Kira's 0.5m³ offer. Your balance jumps 1m³ → 1.5m³.
3. Go to **Give-aways** → post an item. It appears in the grid.
4. Switch persona (log out → log back in as **Ben Hartley**) → claim your
   posted item.
5. Either resident: **Collection day** → add an item to the pool. Pool meter
   updates and the per-resident entitlement bar decrements.

## What's intentionally out of scope

- Real authentication or backend
- Photo uploads (use URL field, defaults to a placeholder)
- Multi-building support
- Notifications, chat, real maps
- Tests
