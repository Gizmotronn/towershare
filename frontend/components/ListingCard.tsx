import { Link } from 'react-router-dom'
import { Calendar, Ruler } from 'lucide-react'
import { CATEGORY_SHORT, type FurnitureListing } from '../types'
import { useAppData } from '../context/AppDataContext'
import { formatDateShort, formatM3 } from '../lib/utils'
import { Badge } from './ui/Badge'

export function ListingCard({ listing }: { listing: FurnitureListing }) {
  const { getResident } = useAppData()
  const poster = getResident(listing.postedById)

  return (
    <Link
      to={`/listings/${listing.id}`}
      className="group block rounded-xl border border-slate-200 bg-white shadow-card overflow-hidden transition-shadow hover:shadow-md focus-visible:outline-none focus-visible:ring-2 focus-visible:ring-brand-600"
    >
      <div className="aspect-[4/3] overflow-hidden bg-slate-100 relative">
        <img
          src={listing.photoUrl}
          alt={listing.title}
          className="h-full w-full object-cover transition-transform group-hover:scale-[1.02]"
          loading="lazy"
        />
        <div className="absolute top-3 right-3">
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
        <div className="absolute top-3 left-3">
          <Badge variant="outline" className="bg-white/90 backdrop-blur-sm">
            {CATEGORY_SHORT[listing.category]}
          </Badge>
        </div>
      </div>
      <div className="p-4 space-y-2">
        <h3 className="font-semibold text-slate-900 leading-snug">
          {listing.title}
        </h3>
        <p className="text-xs text-slate-500">
          From {poster?.name ?? 'a neighbour'} · {poster?.aptNumber ?? ''}
        </p>
        <div className="flex items-center gap-3 text-xs text-slate-600 pt-1">
          <span className="inline-flex items-center gap-1">
            <Ruler className="h-3.5 w-3.5" />
            {formatM3(listing.estimatedM3)}
          </span>
          <span className="inline-flex items-center gap-1">
            <Calendar className="h-3.5 w-3.5" />
            By {formatDateShort(listing.pickupBy)}
          </span>
        </div>
      </div>
    </Link>
  )
}
