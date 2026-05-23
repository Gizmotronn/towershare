import { Link } from 'react-router-dom'
import {
  ArrowRight,
  CalendarDays,
  PackagePlus,
  Recycle,
  Share2,
  Sofa,
} from 'lucide-react'
import { useAppData } from '../context/AppDataContext'
import {
  Card,
  CardContent,
  CardDescription,
  CardHeader,
  CardTitle,
} from '../components/ui/Card'
import { EntitlementBar } from '../components/EntitlementBar'
import { Badge } from '../components/ui/Badge'
import { formatDate, formatM3, daysUntil } from '../lib/utils'

export function Dashboard() {
  const {
    currentUser,
    listings,
    shares,
    collectionDay,
    collectionItems,
    building,
    getResident,
  } = useAppData()

  if (!currentUser) return null

  const itemsGivenAway = listings.filter(
    (l) => l.status === 'claimed' || l.status === 'collected',
  ).length
  const m3Shared = shares
    .filter((s) => s.status === 'claimed')
    .reduce((sum, s) => sum + s.m3Amount, 0)
  const m2Pooled = collectionItems.reduce((sum, ci) => sum + ci.estimatedM3, 0)

  const daysToCollection = daysUntil(collectionDay.scheduledDate)
  const organizer = getResident(collectionDay.organizerId)

  return (
    <div className="space-y-8">
      <div>
        <p className="text-sm text-slate-500">
          Welcome back, {currentUser.name.split(' ')[0]} 👋
        </p>
        <h1 className="text-2xl sm:text-3xl font-semibold text-slate-900 mt-1">
          Your building, working together.
        </h1>
      </div>

      <div className="grid grid-cols-1 lg:grid-cols-3 gap-5">
        <Card className="lg:col-span-2">
          <CardHeader>
            <CardTitle>Your hard-waste entitlement</CardTitle>
            <CardDescription>
              {currentUser.entitlementYear} · Resets each financial year
            </CardDescription>
          </CardHeader>
          <CardContent className="space-y-5">
            <EntitlementBar
              remaining={currentUser.entitlementRemainingM3}
              capacity={1}
            />
            <div className="flex flex-wrap gap-2 text-sm">
              <Link
                to="/entitlements"
                className="inline-flex items-center gap-1.5 h-9 px-3 rounded-md border border-slate-300 bg-white text-slate-900 text-sm font-medium hover:bg-slate-50"
              >
                <Share2 className="h-4 w-4" /> Offer or claim m³
              </Link>
              <Link
                to="/collection-day"
                className="inline-flex items-center gap-1.5 h-9 px-3 rounded-md border border-slate-300 bg-white text-slate-900 text-sm font-medium hover:bg-slate-50"
              >
                <CalendarDays className="h-4 w-4" /> Add to pool
              </Link>
            </div>
          </CardContent>
        </Card>

        <Card className="bg-gradient-to-br from-brand-600 to-brand-700 text-white border-brand-700">
          <CardHeader className="border-b-white/10">
            <p className="text-xs uppercase tracking-wider text-brand-50/80 font-medium">
              Next collection day
            </p>
            <p className="text-2xl font-semibold mt-1">
              {formatDate(collectionDay.scheduledDate)}
            </p>
            <p className="text-sm text-brand-50/80 mt-1">
              {daysToCollection > 0
                ? `${daysToCollection} days away`
                : daysToCollection === 0
                  ? 'Today'
                  : 'Has passed'}
            </p>
          </CardHeader>
          <CardContent className="text-sm text-brand-50/90 space-y-3">
            <p>
              Organised by {organizer?.name ?? 'the building coordinator'}. Add
              your items to pool against the building's combined entitlement.
            </p>
            <Link
              to="/collection-day"
              className="inline-flex items-center gap-1.5 text-white font-medium hover:underline"
            >
              View pool <ArrowRight className="h-4 w-4" />
            </Link>
          </CardContent>
        </Card>
      </div>

      <div className="grid grid-cols-2 sm:grid-cols-4 gap-4">
        <StatTile
          icon={<Sofa className="h-4 w-4" />}
          label="Items given away"
          value={itemsGivenAway}
        />
        <StatTile
          icon={<Share2 className="h-4 w-4" />}
          label="m³ shared between neighbours"
          value={formatM3(m3Shared)}
        />
        <StatTile
          icon={<PackagePlus className="h-4 w-4" />}
          label="m³ in pickup pool"
          value={formatM3(m2Pooled)}
        />
        <StatTile
          icon={<Recycle className="h-4 w-4" />}
          label="Active residents"
          value={`${building.totalUnits}`}
        />
      </div>

      <div className="grid grid-cols-1 lg:grid-cols-2 gap-5">
        <Card>
          <CardHeader className="flex items-start justify-between">
            <div>
              <CardTitle>Recent give-aways</CardTitle>
              <CardDescription>From neighbours in the block</CardDescription>
            </div>
            <Link
              to="/marketplace"
              className="text-sm text-brand-700 hover:underline font-medium"
            >
              See all
            </Link>
          </CardHeader>
          <CardContent className="space-y-3">
            {listings.slice(0, 4).map((listing) => {
              const poster = getResident(listing.postedById)
              return (
                <Link
                  key={listing.id}
                  to={`/listings/${listing.id}`}
                  className="flex items-center gap-3 rounded-lg hover:bg-slate-50 -mx-2 px-2 py-2"
                >
                  <img
                    src={listing.photoUrl}
                    alt=""
                    className="h-12 w-16 object-cover rounded-md bg-slate-100"
                  />
                  <div className="flex-1 min-w-0">
                    <p className="font-medium text-sm text-slate-900 truncate">
                      {listing.title}
                    </p>
                    <p className="text-xs text-slate-500 truncate">
                      {poster?.name} · {poster?.aptNumber}
                    </p>
                  </div>
                  {listing.status === 'available' ? (
                    <Badge variant="success">Available</Badge>
                  ) : listing.status === 'claimed' ? (
                    <Badge variant="warning">Claimed</Badge>
                  ) : (
                    <Badge variant="neutral">Collected</Badge>
                  )}
                </Link>
              )
            })}
          </CardContent>
        </Card>

        <Card>
          <CardHeader className="flex items-start justify-between">
            <div>
              <CardTitle>m³ currently on offer</CardTitle>
              <CardDescription>
                Neighbours with entitlement to spare
              </CardDescription>
            </div>
            <Link
              to="/entitlements"
              className="text-sm text-brand-700 hover:underline font-medium"
            >
              See all
            </Link>
          </CardHeader>
          <CardContent className="space-y-3">
            {shares.filter((s) => s.status === 'offered').length === 0 && (
              <p className="text-sm text-slate-500">
                No m³ on offer right now.
              </p>
            )}
            {shares
              .filter((s) => s.status === 'offered')
              .slice(0, 4)
              .map((share) => {
                const offerer = getResident(share.offeredById)
                return (
                  <div
                    key={share.id}
                    className="flex items-start gap-3 rounded-lg -mx-2 px-2 py-2"
                  >
                    <div className="h-9 w-9 rounded-full bg-brand-100 text-brand-700 grid place-items-center text-xs font-semibold shrink-0">
                      {offerer?.name
                        .split(' ')
                        .map((p) => p[0])
                        .slice(0, 2)
                        .join('') ?? '?'}
                    </div>
                    <div className="flex-1 min-w-0">
                      <p className="font-medium text-sm text-slate-900">
                        {offerer?.name}{' '}
                        <span className="text-slate-400 font-normal">
                          · {offerer?.aptNumber}
                        </span>
                      </p>
                      <p className="text-xs text-slate-600 line-clamp-1">
                        {share.note ?? 'Open to whoever needs it.'}
                      </p>
                    </div>
                    <Badge variant="default">{formatM3(share.m3Amount)}</Badge>
                  </div>
                )
              })}
          </CardContent>
        </Card>
      </div>
    </div>
  )
}

function StatTile({
  icon,
  label,
  value,
}: {
  icon: React.ReactNode
  label: string
  value: number | string
}) {
  return (
    <Card>
      <CardContent className="space-y-1.5">
        <div className="inline-flex items-center justify-center h-8 w-8 rounded-md bg-brand-50 text-brand-700">
          {icon}
        </div>
        <p className="text-2xl font-semibold text-slate-900 tabular-nums">
          {value}
        </p>
        <p className="text-xs text-slate-500 leading-snug">{label}</p>
      </CardContent>
    </Card>
  )
}
