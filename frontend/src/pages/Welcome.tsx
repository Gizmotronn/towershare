import { useEffect, useMemo } from 'react'
import { Link, Navigate, useSearchParams } from 'react-router-dom'
import {
  ArrowRight,
  Building2,
  CheckCircle2,
  Recycle,
  Share2,
  Sofa,
  Sparkles,
} from 'lucide-react'
import { useAuth } from '../context/AuthContext'
import { useAppData } from '../context/AppDataContext'
import { Button } from '../components/ui/Button'
import { formatM3 } from '../lib/utils'

export function Welcome() {
  const [params] = useSearchParams()
  const buildingParam = params.get('b')
  const { currentUserId } = useAuth()
  const { building, listings, shares, residents } = useAppData()

  const buildingMatches = !buildingParam || buildingParam === building.id

  const stats = useMemo(() => {
    const itemsRehomed = listings.filter(
      (l) => l.status === 'claimed' || l.status === 'collected',
    ).length
    const m3Shared = shares
      .filter((s) => s.status === 'claimed')
      .reduce((sum, s) => sum + s.m3Amount, 0)
    return {
      itemsRehomed,
      m3Shared,
      neighbours: residents.length,
    }
  }, [listings, shares, residents])

  useEffect(() => {
    document.title = `${building.name} · Hard Waste Hub`
  }, [building.name])

  if (currentUserId) return <Navigate to="/dashboard" replace />

  return (
    <div className="min-h-screen bg-gradient-to-b from-brand-50 via-white to-slate-50">
      <div className="max-w-md mx-auto px-5 pt-10 pb-16">
        <div className="flex items-center gap-3">
          <div className="h-11 w-11 rounded-xl bg-brand-600 grid place-items-center text-white shadow-lg shadow-brand-600/30">
            <Building2 className="h-5 w-5" />
          </div>
          <div className="leading-tight">
            <p className="text-xs uppercase tracking-wider text-brand-700 font-semibold">
              Hard Waste Hub
            </p>
            <p className="text-sm text-slate-600">{building.address}</p>
          </div>
        </div>

        <div className="mt-8">
          <div className="inline-flex items-center gap-1.5 rounded-full bg-white border border-brand-100 px-3 py-1 text-xs font-medium text-brand-700 shadow-sm">
            <Sparkles className="h-3.5 w-3.5" />
            You scanned the lobby code
          </div>
          <h1 className="mt-4 text-3xl sm:text-4xl font-semibold text-slate-900 leading-tight">
            Welcome to{' '}
            <span className="text-brand-700">{building.name}</span>'s
            give-away & share board.
          </h1>
          <p className="mt-3 text-slate-600 leading-relaxed">
            Pool your council hard-waste entitlement with your neighbours, list
            furniture you no longer need, and grab the things you do — all
            without a single thing hitting the kerb.
          </p>
        </div>

        {!buildingMatches && (
          <div className="mt-6 rounded-lg border border-amber-200 bg-amber-50 px-4 py-3 text-sm text-amber-900">
            We didn't recognise the building code{' '}
            <code className="font-mono text-xs bg-white/60 px-1 rounded">
              {buildingParam}
            </code>
            . Showing the {building.name} board for now.
          </div>
        )}

        <div className="mt-8 grid grid-cols-3 gap-3">
          <StatTile
            icon={<Sofa className="h-4 w-4" />}
            value={stats.itemsRehomed}
            label="items rehomed"
          />
          <StatTile
            icon={<Share2 className="h-4 w-4" />}
            value={formatM3(stats.m3Shared)}
            label="m³ shared"
          />
          <StatTile
            icon={<Recycle className="h-4 w-4" />}
            value={stats.neighbours}
            label="neighbours"
          />
        </div>

        <div className="mt-10 space-y-3">
          <Link
            to={`/login?b=${building.id}`}
            className="block"
            aria-label="Continue as a resident"
          >
            <Button size="lg" className="w-full justify-center">
              I'm a resident <ArrowRight className="h-4 w-4" />
            </Button>
          </Link>
          <p className="text-center text-xs text-slate-500">
            12 apartments. Free to use. Council-partnered.
          </p>
        </div>

        <div className="mt-10 rounded-xl border border-brand-100 bg-brand-50/60 p-4">
          <p className="text-[11px] uppercase tracking-wider text-brand-700 font-semibold">
            From the City of Melbourne
          </p>
          <p className="mt-1.5 text-sm text-brand-900 leading-relaxed italic">
            "If you live in a high rise building, you can ask your building
            manager to book a collection for you."
          </p>
          <p className="mt-2 text-xs text-slate-600">
            One free Hard Waste collection — for one cubic metre — per
            address, per financial year.
          </p>
        </div>

        <div className="mt-12">
          <p className="text-xs uppercase tracking-wider text-slate-500 font-semibold">
            How it works
          </p>
          <ul className="mt-3 space-y-3">
            <Step
              n={1}
              title="Donate or claim m³"
              body="Got an unused entitlement? Drop it in the pool. Need more for a couch you're tossing? Take some from a neighbour."
            />
            <Step
              n={2}
              title="List furniture before it's waste"
              body="A bookshelf you don't need is one your new neighbour was about to buy. Snap a photo, post it to the building feed."
            />
            <Step
              n={3}
              title="One collection. Less dumping."
              body="Whatever's left gets booked as a single hard-waste pickup, drawing from the pooled entitlement."
            />
          </ul>
        </div>

        <p className="mt-10 text-center text-[11px] text-slate-400">
          A prototype built with the City of Melbourne in mind.
        </p>
      </div>
    </div>
  )
}

function StatTile({
  icon,
  value,
  label,
}: {
  icon: React.ReactNode
  value: number | string
  label: string
}) {
  return (
    <div className="rounded-xl bg-white border border-slate-200 px-3 py-3 shadow-sm">
      <div className="inline-flex h-7 w-7 items-center justify-center rounded-md bg-brand-50 text-brand-700">
        {icon}
      </div>
      <p className="mt-2 text-xl font-semibold text-slate-900 tabular-nums">
        {value}
      </p>
      <p className="text-[11px] text-slate-500 leading-tight mt-0.5">{label}</p>
    </div>
  )
}

function Step({
  n,
  title,
  body,
}: {
  n: number
  title: string
  body: string
}) {
  return (
    <li className="flex gap-3">
      <div className="h-7 w-7 rounded-full bg-brand-600 text-white grid place-items-center text-xs font-semibold shrink-0">
        {n}
      </div>
      <div>
        <p className="font-medium text-sm text-slate-900 flex items-center gap-1.5">
          {title}
          <CheckCircle2 className="h-3.5 w-3.5 text-brand-600" />
        </p>
        <p className="text-sm text-slate-600 leading-snug mt-0.5">{body}</p>
      </div>
    </li>
  )
}
