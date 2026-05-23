import { Link, Navigate, useNavigate, useParams } from 'react-router-dom'
import { ArrowLeft, Calendar, Ruler, User } from 'lucide-react'
import { useAppData } from '../context/AppDataContext'
import { CATEGORY_LABEL } from '../types'
import { Button } from '../components/ui/Button'
import { Badge } from '../components/ui/Badge'
import {
  Card,
  CardContent,
  CardHeader,
  CardTitle,
} from '../components/ui/Card'
import { formatDate, formatM3 } from '../lib/utils'

export function ListingDetail() {
  const { id } = useParams()
  const navigate = useNavigate()
  const { listings, currentUser, getResident, claimListing, withdrawListing } =
    useAppData()

  const listing = listings.find((l) => l.id === id)
  if (!listing) return <Navigate to="/marketplace" replace />

  const poster = getResident(listing.postedById)
  const claimer = getResident(listing.claimedById)
  const isMine = currentUser?.id === listing.postedById
  const canClaim = !isMine && listing.status === 'available'

  return (
    <div className="max-w-3xl mx-auto space-y-6">
      <Link
        to="/marketplace"
        className="inline-flex items-center gap-1.5 text-sm text-slate-600 hover:text-slate-900"
      >
        <ArrowLeft className="h-4 w-4" /> Back to give-aways
      </Link>

      <Card className="overflow-hidden">
        <div className="aspect-[16/9] bg-slate-100 overflow-hidden">
          <img
            src={listing.photoUrl}
            alt={listing.title}
            className="h-full w-full object-cover"
          />
        </div>
        <CardHeader>
          <div className="flex items-start justify-between gap-4">
            <CardTitle className="text-xl">{listing.title}</CardTitle>
            {listing.status === 'available' && (
              <Badge variant="success">Available</Badge>
            )}
            {listing.status === 'claimed' && (
              <Badge variant="warning">Claimed</Badge>
            )}
            {listing.status === 'collected' && (
              <Badge variant="neutral">Collected</Badge>
            )}
          </div>
        </CardHeader>
        <CardContent className="space-y-5">
          <div className="flex flex-wrap gap-4 text-sm text-slate-600">
            <span className="inline-flex items-center gap-1.5">
              <User className="h-4 w-4 text-slate-400" />
              From {poster?.name} ({poster?.aptNumber})
            </span>
            <Badge variant="outline">{CATEGORY_LABEL[listing.category]}</Badge>
            <span className="inline-flex items-center gap-1.5">
              <Ruler className="h-4 w-4 text-slate-400" />
              {formatM3(listing.estimatedM3)} footprint
            </span>
            <span className="inline-flex items-center gap-1.5">
              <Calendar className="h-4 w-4 text-slate-400" />
              Available until {formatDate(listing.pickupBy)}
            </span>
          </div>

          <p className="text-slate-700 leading-relaxed whitespace-pre-wrap">
            {listing.description}
          </p>

          {listing.status === 'claimed' && claimer && (
            <div className="rounded-lg bg-amber-50 border border-amber-100 px-4 py-3 text-sm text-amber-900">
              Claimed by {claimer.name} ({claimer.aptNumber}). They'll be in
              touch about pickup.
            </div>
          )}

          <div className="flex flex-wrap gap-2 pt-2">
            {canClaim && (
              <Button onClick={() => claimListing(listing.id)}>
                Claim this item
              </Button>
            )}
            {isMine && listing.status === 'available' && (
              <Button
                variant="outline"
                onClick={() => {
                  withdrawListing(listing.id)
                  navigate('/marketplace')
                }}
              >
                Withdraw listing
              </Button>
            )}
            {!canClaim && !isMine && listing.status !== 'available' && (
              <p className="text-sm text-slate-500">
                This item is no longer available.
              </p>
            )}
          </div>
        </CardContent>
      </Card>
    </div>
  )
}
