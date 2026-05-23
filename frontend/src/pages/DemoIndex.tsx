import { Link, useNavigate } from 'react-router-dom'
import {
  ArrowRight,
  Building2,
  CalendarDays,
  Compass,
  LayoutDashboard,
  LogIn,
  PackagePlus,
  QrCode,
  Share2,
  Sofa,
  Sparkles,
  User,
  UserCircle2,
} from 'lucide-react'
import { useAuth } from '../context/AuthContext'
import { useAppData } from '../context/AppDataContext'
import { Button } from '../components/ui/Button'
import { Select } from '../components/ui/Select'
import { Label } from '../components/ui/Label'
import { formatM3 } from '../lib/utils'

type RouteCard = {
  to: string
  title: string
  blurb: string
  icon: React.ReactNode
  badge?: 'public' | 'protected'
}

const entryFlow: RouteCard[] = [
  {
    to: '/qr-poster',
    title: 'QR poster',
    blurb:
      'The printable A4 lobby sign. Point a phone camera at it to enter the flow as a resident would.',
    icon: <QrCode className="h-4 w-4" />,
    badge: 'public',
  },
  {
    to: '/welcome?b=bldg_brunswick_heights',
    title: 'Welcome',
    blurb:
      'Where the QR scan lands. Building branding, live stats, "How it works" primer, and the CTA into login.',
    icon: <Sparkles className="h-4 w-4" />,
    badge: 'public',
  },
  {
    to: '/login?b=bldg_brunswick_heights',
    title: 'Login',
    blurb:
      'Email/password form plus a persona quick-switch. The QR banner shows when arriving via ?b=.',
    icon: <LogIn className="h-4 w-4" />,
    badge: 'public',
  },
]

const residentApp: RouteCard[] = [
  {
    to: '/dashboard',
    title: 'Dashboard',
    blurb:
      "The resident's home. Their entitlement bar, the next collection day, recent give-aways, m³ on offer.",
    icon: <LayoutDashboard className="h-4 w-4" />,
    badge: 'protected',
  },
  {
    to: '/marketplace',
    title: 'Give-aways',
    blurb:
      'The building feed of furniture. Filter, claim, and see who posted what.',
    icon: <Sofa className="h-4 w-4" />,
    badge: 'protected',
  },
  {
    to: '/listings/new',
    title: 'New listing',
    blurb:
      'The form a resident fills out to post a piece of furniture before it becomes waste.',
    icon: <PackagePlus className="h-4 w-4" />,
    badge: 'protected',
  },
  {
    to: '/entitlements',
    title: 'Share m³',
    blurb:
      'Donate unused m³ to the pool or claim from neighbours who have spare entitlement.',
    icon: <Share2 className="h-4 w-4" />,
    badge: 'protected',
  },
  {
    to: '/collection-day',
    title: 'Collection day',
    blurb:
      'The consolidated booking. Add items to the pool, see whose m³ is being drawn down.',
    icon: <CalendarDays className="h-4 w-4" />,
    badge: 'protected',
  },
  {
    to: '/profile',
    title: 'Profile',
    blurb: "Resident's own details and a record of what they've shared.",
    icon: <UserCircle2 className="h-4 w-4" />,
    badge: 'protected',
  },
]

export function DemoIndex() {
  const navigate = useNavigate()
  const { currentUserId, loginAs, logout } = useAuth()
  const { residents, building, listings } = useAppData()

  const currentResident = residents.find((r) => r.id === currentUserId)
  const sampleListing = listings[0]

  const handlePersona = (residentId: string) => {
    if (residentId === '__none') {
      logout()
      return
    }
    loginAs(residentId)
  }

  return (
    <div className="min-h-screen bg-slate-50">
      <div className="max-w-4xl mx-auto px-4 sm:px-6 py-10">
        <div className="flex items-center gap-3">
          <div className="h-10 w-10 rounded-xl bg-slate-900 grid place-items-center text-white">
            <Compass className="h-5 w-5" />
          </div>
          <div className="leading-tight">
            <p className="text-xs uppercase tracking-wider text-slate-500 font-semibold">
              Demo navigator
            </p>
            <p className="text-sm text-slate-700">
              {building.name} · Hard Waste Hub prototype
            </p>
          </div>
        </div>

        <h1 className="mt-6 text-2xl sm:text-3xl font-semibold text-slate-900">
          Jump to any page in the demo.
        </h1>
        <p className="mt-2 text-slate-600 max-w-2xl">
          The full resident flow starts from the QR poster. The protected pages
          need a persona — set one below and you can navigate freely without
          logging in each time.
        </p>

        <div className="mt-8 rounded-2xl border border-slate-200 bg-white p-5 sm:p-6 shadow-sm">
          <div className="flex items-start gap-3">
            <div className="h-9 w-9 rounded-lg bg-brand-50 text-brand-700 grid place-items-center shrink-0">
              <User className="h-4 w-4" />
            </div>
            <div className="flex-1">
              <p className="text-sm font-semibold text-slate-900">
                Active persona
              </p>
              <p className="text-sm text-slate-600 mt-0.5">
                {currentResident ? (
                  <>
                    Signed in as{' '}
                    <span className="font-medium text-slate-900">
                      {currentResident.name}
                    </span>{' '}
                    ({currentResident.aptNumber}) ·{' '}
                    {formatM3(currentResident.entitlementRemainingM3)} left
                  </>
                ) : (
                  'No one signed in. Protected pages will auto-pick Amelia Tran (Apt 1A) — pick someone else above to switch.'
                )}
              </p>
            </div>
          </div>

          <div className="mt-4 grid sm:grid-cols-[1fr_auto] gap-3">
            <div className="space-y-1.5">
              <Label htmlFor="persona-switch">Switch persona</Label>
              <Select
                id="persona-switch"
                value={currentUserId ?? '__none'}
                onChange={(e) => handlePersona(e.target.value)}
              >
                <option value="__none">— Signed out —</option>
                {residents.map((r) => (
                  <option key={r.id} value={r.id}>
                    {r.name} ({r.aptNumber}) —{' '}
                    {formatM3(r.entitlementRemainingM3)} left
                  </option>
                ))}
              </Select>
            </div>
            <div className="flex items-end">
              <Button
                variant="outline"
                size="md"
                className="bg-white w-full sm:w-auto"
                onClick={() => navigate('/dashboard')}
                disabled={!currentUserId}
              >
                Open dashboard <ArrowRight className="h-4 w-4" />
              </Button>
            </div>
          </div>
        </div>

        <Section
          title="Entry flow"
          description="What happens between the lobby QR and the dashboard."
          icon={<Building2 className="h-4 w-4" />}
        >
          {entryFlow.map((r) => (
            <RouteTile key={r.to} {...r} />
          ))}
        </Section>

        <Section
          title="Resident app"
          description="Pages behind the auth gate. Pick a persona above first."
          icon={<LayoutDashboard className="h-4 w-4" />}
        >
          {residentApp.map((r) => (
            <RouteTile key={r.to} {...r} />
          ))}
        </Section>

        {sampleListing && (
          <Section
            title="Sample deep link"
            description="A real listing detail page from the seed data."
            icon={<Sofa className="h-4 w-4" />}
          >
            <RouteTile
              to={`/listings/${sampleListing.id}`}
              title={sampleListing.title}
              blurb={`Listing detail for "${sampleListing.title}" — posted by a seeded resident, ${formatM3(sampleListing.estimatedM3)}.`}
              icon={<Sofa className="h-4 w-4" />}
              badge="protected"
            />
          </Section>
        )}

        <p className="mt-12 text-center text-xs text-slate-400">
          This page is unlisted and intended for hackathon demos only.
        </p>
      </div>
    </div>
  )
}

function Section({
  title,
  description,
  icon,
  children,
}: {
  title: string
  description: string
  icon: React.ReactNode
  children: React.ReactNode
}) {
  return (
    <section className="mt-10">
      <div className="flex items-center gap-2">
        <div className="h-7 w-7 rounded-md bg-slate-900 text-white grid place-items-center">
          {icon}
        </div>
        <h2 className="text-lg font-semibold text-slate-900">{title}</h2>
      </div>
      <p className="mt-1 text-sm text-slate-500">{description}</p>
      <div className="mt-4 grid sm:grid-cols-2 gap-3">{children}</div>
    </section>
  )
}

function RouteTile({ to, title, blurb, icon, badge }: RouteCard) {
  return (
    <Link
      to={to}
      className="group rounded-xl border border-slate-200 bg-white p-4 hover:border-brand-500 hover:shadow-md transition-all"
    >
      <div className="flex items-start gap-3">
        <div className="h-8 w-8 rounded-md bg-brand-50 text-brand-700 grid place-items-center shrink-0">
          {icon}
        </div>
        <div className="flex-1 min-w-0">
          <div className="flex items-center gap-2">
            <p className="font-medium text-sm text-slate-900">{title}</p>
            {badge === 'public' && (
              <span className="text-[10px] uppercase tracking-wider font-semibold text-brand-700 bg-brand-50 px-1.5 py-0.5 rounded">
                Public
              </span>
            )}
            {badge === 'protected' && (
              <span className="text-[10px] uppercase tracking-wider font-semibold text-amber-700 bg-amber-50 px-1.5 py-0.5 rounded">
                Protected
              </span>
            )}
          </div>
          <p className="mt-1 text-xs text-slate-600 leading-snug line-clamp-3">
            {blurb}
          </p>
          <p className="mt-2 text-xs font-mono text-slate-400 truncate">
            {to}
          </p>
        </div>
        <ArrowRight className="h-4 w-4 text-slate-300 group-hover:text-brand-600 mt-1 shrink-0" />
      </div>
    </Link>
  )
}
