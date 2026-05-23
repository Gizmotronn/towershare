export type Resident = {
  id: string
  name: string
  aptNumber: string
  email: string
  entitlementRemainingM3: number
  entitlementYear: string
}

export type ItemCategory = 'furniture' | 'whitegoods' | 'ewaste' | 'mattress'

export type ListingStatus = 'available' | 'claimed' | 'collected'

export type FurnitureListing = {
  id: string
  postedById: string
  title: string
  description: string
  photoUrl: string
  category: ItemCategory
  estimatedM3: number
  pickupBy: string
  status: ListingStatus
  claimedById?: string
  createdAt: string
}

export type EntitlementShareStatus = 'offered' | 'claimed'

export type EntitlementShare = {
  id: string
  offeredById: string
  m3Amount: number
  note?: string
  status: EntitlementShareStatus
  claimedById?: string
  createdAt: string
}

export type CollectionDayStatus = 'planned' | 'locked' | 'completed'

export type CollectionDay = {
  id: string
  scheduledDate: string
  organizerId: string
  status: CollectionDayStatus
}

export type CollectionItem = {
  id: string
  collectionDayId: string
  residentId: string
  title: string
  category: ItemCategory
  estimatedM3: number
}

export type Building = {
  id: string
  name: string
  address: string
  totalUnits: number
}

export const CATEGORY_LABEL: Record<ItemCategory, string> = {
  furniture: 'Household furniture',
  whitegoods: 'White goods',
  ewaste: 'E-waste',
  mattress: 'Mattress',
}

export const CATEGORY_SHORT: Record<ItemCategory, string> = {
  furniture: 'Furniture',
  whitegoods: 'White goods',
  ewaste: 'E-waste',
  mattress: 'Mattress',
}
