import { useState } from 'react'
import { Link } from 'react-router-dom'
import { Plus } from 'lucide-react'
import { useAppData } from '../context/AppDataContext'
import { ListingCard } from '../components/ListingCard'
import { cn } from '../lib/utils'
import { CATEGORY_SHORT, type ItemCategory, type ListingStatus } from '../types'

type OwnerFilter = 'all' | 'mine'

export function Marketplace() {
  const { listings, currentUser } = useAppData()
  const [filter, setFilter] = useState<OwnerFilter>('all')
  const [statusFilter, setStatusFilter] = useState<ListingStatus | 'all'>('all')
  const [categoryFilter, setCategoryFilter] = useState<ItemCategory | 'all'>(
    'all',
  )

  let filtered = listings
  if (filter === 'mine' && currentUser) {
    filtered = filtered.filter((l) => l.postedById === currentUser.id)
  }
  if (statusFilter !== 'all') {
    filtered = filtered.filter((l) => l.status === statusFilter)
  }
  if (categoryFilter !== 'all') {
    filtered = filtered.filter((l) => l.category === categoryFilter)
  }

  return (
    <div className="space-y-6">
      <div className="flex items-end justify-between gap-4 flex-wrap">
        <div>
          <h1 className="text-2xl font-semibold text-slate-900">
            Furniture give-aways
          </h1>
          <p className="text-sm text-slate-500 mt-1">
            Free pieces from your neighbours — keep usable items out of
            landfill.
          </p>
        </div>
        <Link
          to="/listings/new"
          className="inline-flex items-center gap-1.5 h-10 px-4 rounded-lg bg-brand-600 text-white text-sm font-medium hover:bg-brand-700 shadow-sm"
        >
          <Plus className="h-4 w-4" />
          Post an item
        </Link>
      </div>

      <div className="flex flex-wrap gap-2 items-center">
        <div className="flex gap-1 p-1 bg-slate-100 rounded-lg">
          <FilterChip
            active={filter === 'all'}
            onClick={() => setFilter('all')}
          >
            All
          </FilterChip>
          <FilterChip
            active={filter === 'mine'}
            onClick={() => setFilter('mine')}
          >
            My posts
          </FilterChip>
        </div>
        <span className="text-slate-300">·</span>
        <div className="flex gap-1 p-1 bg-slate-100 rounded-lg">
          <FilterChip
            active={statusFilter === 'all'}
            onClick={() => setStatusFilter('all')}
          >
            Any status
          </FilterChip>
          <FilterChip
            active={statusFilter === 'available'}
            onClick={() => setStatusFilter('available')}
          >
            Available
          </FilterChip>
          <FilterChip
            active={statusFilter === 'claimed'}
            onClick={() => setStatusFilter('claimed')}
          >
            Claimed
          </FilterChip>
          <FilterChip
            active={statusFilter === 'collected'}
            onClick={() => setStatusFilter('collected')}
          >
            Collected
          </FilterChip>
        </div>
        <span className="text-slate-300">·</span>
        <div className="flex gap-1 p-1 bg-slate-100 rounded-lg overflow-x-auto">
          <FilterChip
            active={categoryFilter === 'all'}
            onClick={() => setCategoryFilter('all')}
          >
            All categories
          </FilterChip>
          <FilterChip
            active={categoryFilter === 'furniture'}
            onClick={() => setCategoryFilter('furniture')}
          >
            {CATEGORY_SHORT.furniture}
          </FilterChip>
          <FilterChip
            active={categoryFilter === 'whitegoods'}
            onClick={() => setCategoryFilter('whitegoods')}
          >
            {CATEGORY_SHORT.whitegoods}
          </FilterChip>
          <FilterChip
            active={categoryFilter === 'ewaste'}
            onClick={() => setCategoryFilter('ewaste')}
          >
            {CATEGORY_SHORT.ewaste}
          </FilterChip>
          <FilterChip
            active={categoryFilter === 'mattress'}
            onClick={() => setCategoryFilter('mattress')}
          >
            {CATEGORY_SHORT.mattress}
          </FilterChip>
        </div>
      </div>

      {filtered.length === 0 ? (
        <div className="text-center py-16 border-2 border-dashed border-slate-200 rounded-xl">
          <p className="text-slate-500">No listings match this filter.</p>
          <Link
            to="/listings/new"
            className="inline-flex items-center gap-1.5 mt-3 text-brand-700 font-medium hover:underline"
          >
            <Plus className="h-4 w-4" /> Post the first one
          </Link>
        </div>
      ) : (
        <div className="grid grid-cols-1 sm:grid-cols-2 lg:grid-cols-3 gap-5">
          {filtered.map((listing) => (
            <ListingCard key={listing.id} listing={listing} />
          ))}
        </div>
      )}
    </div>
  )
}

function FilterChip({
  active,
  onClick,
  children,
}: {
  active: boolean
  onClick: () => void
  children: React.ReactNode
}) {
  return (
    <button
      onClick={onClick}
      className={cn(
        'px-3 py-1.5 text-sm font-medium rounded-md transition-colors whitespace-nowrap',
        active
          ? 'bg-white text-slate-900 shadow-sm'
          : 'text-slate-600 hover:text-slate-900',
      )}
    >
      {children}
    </button>
  )
}
