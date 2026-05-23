import type { FurnitureListing, ItemCategory, ListingStatus } from '../types'

const PB_URL = (import.meta.env.VITE_PB_URL as string | undefined) ?? 'http://localhost:8090'

export type PbUser = {
  id: string
  email: string
  display_name: string
}

export type PbTower = {
  id: string
  name: string
  address: string
}

export type PbListing = {
  id: string
  title: string
  description: string
  category: string
  status: string
  estimated_m3?: number
  pickup_by?: string
  owner: string
  tower: string
  images: string[]
  created: string
  collectionId: string
}

const TOKEN_KEY = 'pb_token'
const USER_KEY = 'pb_user'

let _token: string | null = localStorage.getItem(TOKEN_KEY)
let _user: PbUser | null = JSON.parse(localStorage.getItem(USER_KEY) ?? 'null')

export function getPbUser(): PbUser | null {
  return _user
}

export function setPbSession(token: string, user: PbUser): void {
  _token = token
  _user = user
  localStorage.setItem(TOKEN_KEY, token)
  localStorage.setItem(USER_KEY, JSON.stringify(user))
}

export function clearPbSession(): void {
  _token = null
  _user = null
  localStorage.removeItem(TOKEN_KEY)
  localStorage.removeItem(USER_KEY)
}

function authHeader(): Record<string, string> {
  return _token ? { Authorization: _token } : {}
}

export async function pbAuth(
  email: string,
  password: string,
): Promise<{ token: string; record: PbUser }> {
  const res = await fetch(`${PB_URL}/api/collections/users/auth-with-password`, {
    method: 'POST',
    headers: { 'Content-Type': 'application/json' },
    body: JSON.stringify({ identity: email, password }),
  })
  if (!res.ok) {
    const err = (await res.json().catch(() => ({}))) as { message?: string }
    throw new Error(err.message ?? 'Login failed')
  }
  return res.json() as Promise<{ token: string; record: PbUser }>
}

export async function pbGetTowers(): Promise<PbTower[]> {
  const res = await fetch(
    `${PB_URL}/api/collections/towers/records?sort=name&perPage=50`,
    { headers: authHeader() },
  )
  if (!res.ok) throw new Error('Failed to fetch towers')
  const data = (await res.json()) as { items: PbTower[] }
  return data.items ?? []
}

export type PbListingExpanded = PbListing & {
  expand?: { owner?: PbUser }
}

export function pbFileUrl(collectionId: string, recordId: string, filename: string): string {
  return `${PB_URL}/api/files/${collectionId}/${recordId}/${filename}`
}

export async function pbGetListings(): Promise<PbListingExpanded[]> {
  const res = await fetch(
    `${PB_URL}/api/collections/listings/records?sort=-created&perPage=100&expand=owner`,
    { headers: authHeader() },
  )
  if (!res.ok) throw new Error('Failed to fetch listings')
  const data = (await res.json()) as { items: PbListingExpanded[] }
  return data.items ?? []
}

const STATUS_MAP: Record<string, ListingStatus> = {
  active: 'available',
  claimed: 'claimed',
  closed: 'collected',
  available: 'available',
  collected: 'collected',
}

const VALID_CATEGORIES: ItemCategory[] = ['furniture', 'whitegoods', 'ewaste', 'mattress']

export function pbToFurnitureListing(pb: PbListingExpanded): FurnitureListing {
  const imageUrl = pb.images?.length
    ? pbFileUrl(pb.collectionId, pb.id, pb.images[0])
    : `https://picsum.photos/seed/${encodeURIComponent(pb.title.toLowerCase().replace(/\s+/g, '-'))}/600/400`

  return {
    id: pb.id,
    postedById: pb.owner,
    title: pb.title,
    description: pb.description ?? '',
    photoUrl: imageUrl,
    category: VALID_CATEGORIES.includes(pb.category as ItemCategory)
      ? (pb.category as ItemCategory)
      : 'furniture',
    estimatedM3: pb.estimated_m3 ?? 0,
    pickupBy: pb.pickup_by ?? '',
    status: STATUS_MAP[pb.status] ?? 'available',
    createdAt: pb.created,
  }
}

export async function pbCreateListing(body: FormData): Promise<PbListing> {
  const res = await fetch(`${PB_URL}/api/collections/listings/records`, {
    method: 'POST',
    // No Content-Type header — browser sets the multipart boundary automatically.
    headers: authHeader(),
    body,
  })
  if (!res.ok) {
    const err = (await res.json().catch(() => ({}))) as { message?: string }
    throw new Error(err.message ?? 'Failed to create listing')
  }
  return res.json() as Promise<PbListing>
}
