import {
  createContext,
  useCallback,
  useContext,
  useEffect,
  useMemo,
  useState,
  type ReactNode,
} from 'react'
import type {
  Building,
  CollectionDay,
  CollectionItem,
  EntitlementShare,
  FurnitureListing,
  ItemCategory,
  Resident,
} from '../types'
import {
  BUILDING,
  COLLECTION_DAY,
  COLLECTION_ITEMS,
  LISTINGS,
  RESIDENTS,
  SHARES,
} from '../lib/seed'
import { uid } from '../lib/utils'
import { pbGetTowers } from '../lib/pb'
import { useAuth } from './AuthContext'

type NewListingInput = {
  title: string
  description: string
  photoUrl: string
  category: ItemCategory
  estimatedM3: number
  pickupBy: string
}

type NewCollectionItemInput = {
  title: string
  category: ItemCategory
  estimatedM3: number
}

type AppData = {
  building: Building
  residents: Resident[]
  listings: FurnitureListing[]
  shares: EntitlementShare[]
  collectionDay: CollectionDay
  collectionItems: CollectionItem[]
  currentUser: Resident | null

  getResident: (id: string | undefined) => Resident | undefined

  createListing: (input: NewListingInput) => string | null
  claimListing: (id: string) => void
  withdrawListing: (id: string) => void

  offerShare: (m3: number, note?: string) => boolean
  claimShare: (id: string) => void

  addCollectionItem: (input: NewCollectionItemInput) => boolean
  removeCollectionItem: (id: string) => void
}

const AppDataContext = createContext<AppData | null>(null)

export function AppDataProvider({ children }: { children: ReactNode }) {
  const { currentUserId } = useAuth()
  const [building, setBuilding] = useState<Building>(BUILDING)
  const [residents, setResidents] = useState<Resident[]>(RESIDENTS)

  useEffect(() => {
    pbGetTowers()
      .then(([first]) => {
        if (!first) return
        setBuilding((prev) => ({
          ...prev,
          id: first.id,
          name: first.name,
          address: first.address,
        }))
      })
      .catch(() => {/* backend not running — keep seed building */})
  }, [])
  const [listings, setListings] = useState<FurnitureListing[]>(LISTINGS)
  const [shares, setShares] = useState<EntitlementShare[]>(SHARES)
  const [collectionDay] = useState<CollectionDay>(COLLECTION_DAY)
  const [collectionItems, setCollectionItems] =
    useState<CollectionItem[]>(COLLECTION_ITEMS)

  const currentUser = useMemo(
    () => residents.find((r) => r.id === currentUserId) ?? null,
    [residents, currentUserId],
  )

  const getResident = useCallback(
    (id: string | undefined) => residents.find((r) => r.id === id),
    [residents],
  )

  const adjustEntitlement = useCallback(
    (residentId: string, delta: number) => {
      setResidents((prev) =>
        prev.map((r) =>
          r.id === residentId
            ? {
                ...r,
                entitlementRemainingM3: Math.max(
                  0,
                  Math.round((r.entitlementRemainingM3 + delta) * 100) / 100,
                ),
              }
            : r,
        ),
      )
    },
    [],
  )

  const createListing = useCallback(
    (input: NewListingInput): string | null => {
      if (!currentUserId) return null
      const id = uid('l')
      const listing: FurnitureListing = {
        id,
        postedById: currentUserId,
        status: 'available',
        createdAt: new Date().toISOString(),
        ...input,
      }
      setListings((prev) => [listing, ...prev])
      return id
    },
    [currentUserId],
  )

  const claimListing = useCallback(
    (id: string) => {
      if (!currentUserId) return
      setListings((prev) =>
        prev.map((l) =>
          l.id === id && l.status === 'available' && l.postedById !== currentUserId
            ? { ...l, status: 'claimed', claimedById: currentUserId }
            : l,
        ),
      )
    },
    [currentUserId],
  )

  const withdrawListing = useCallback(
    (id: string) => {
      if (!currentUserId) return
      setListings((prev) =>
        prev.filter((l) => !(l.id === id && l.postedById === currentUserId)),
      )
    },
    [currentUserId],
  )

  const offerShare = useCallback(
    (m3: number, note?: string): boolean => {
      if (!currentUserId) return false
      if (!currentUser || currentUser.entitlementRemainingM3 < m3) return false
      const share: EntitlementShare = {
        id: uid('s'),
        offeredById: currentUserId,
        m3Amount: m3,
        note,
        status: 'offered',
        createdAt: new Date().toISOString(),
      }
      setShares((prev) => [share, ...prev])
      adjustEntitlement(currentUserId, -m3)
      return true
    },
    [currentUserId, currentUser, adjustEntitlement],
  )

  const claimShare = useCallback(
    (id: string) => {
      if (!currentUserId) return
      const share = shares.find((s) => s.id === id)
      if (!share || share.status !== 'offered') return
      if (share.offeredById === currentUserId) return
      setShares((prev) =>
        prev.map((s) =>
          s.id === id ? { ...s, status: 'claimed', claimedById: currentUserId } : s,
        ),
      )
      adjustEntitlement(currentUserId, share.m3Amount)
    },
    [currentUserId, shares, adjustEntitlement],
  )

  const addCollectionItem = useCallback(
    (input: NewCollectionItemInput): boolean => {
      if (!currentUserId) return false
      if (!currentUser || currentUser.entitlementRemainingM3 < input.estimatedM3)
        return false
      const item: CollectionItem = {
        id: uid('ci'),
        collectionDayId: collectionDay.id,
        residentId: currentUserId,
        title: input.title,
        category: input.category,
        estimatedM3: input.estimatedM3,
      }
      setCollectionItems((prev) => [...prev, item])
      adjustEntitlement(currentUserId, -input.estimatedM3)
      return true
    },
    [currentUserId, currentUser, collectionDay.id, adjustEntitlement],
  )

  const removeCollectionItem = useCallback(
    (id: string) => {
      if (!currentUserId) return
      const item = collectionItems.find((c) => c.id === id)
      if (!item || item.residentId !== currentUserId) return
      setCollectionItems((prev) => prev.filter((c) => c.id !== id))
      adjustEntitlement(currentUserId, item.estimatedM3)
    },
    [currentUserId, collectionItems, adjustEntitlement],
  )

  const value: AppData = {
    building,
    residents,
    listings,
    shares,
    collectionDay,
    collectionItems,
    currentUser,
    getResident,
    createListing,
    claimListing,
    withdrawListing,
    offerShare,
    claimShare,
    addCollectionItem,
    removeCollectionItem,
  }

  return (
    <AppDataContext.Provider value={value}>{children}</AppDataContext.Provider>
  )
}

export function useAppData() {
  const ctx = useContext(AppDataContext)
  if (!ctx) throw new Error('useAppData must be used within AppDataProvider')
  return ctx
}
