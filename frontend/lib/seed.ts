import type {
  Building,
  CollectionDay,
  CollectionItem,
  EntitlementShare,
  FurnitureListing,
  Resident,
} from '../types'

export const BUILDING: Building = {
  id: 'bldg_brunswick_heights',
  name: 'Brunswick Heights',
  address: '42 Sydney Rd, Brunswick VIC 3056',
  totalUnits: 12,
}

export const FY = 'FY2025-26'

const today = new Date()
const inDays = (d: number) => {
  const t = new Date(today)
  t.setDate(t.getDate() + d)
  return t.toISOString()
}
const daysAgo = (d: number) => inDays(-d)

export const RESIDENTS: Resident[] = [
  { id: 'r_amelia', name: 'Amelia Tran', aptNumber: 'Apt 1A', email: 'amelia@brunswickheights.local', entitlementRemainingM3: 1.0, entitlementYear: FY },
  { id: 'r_ben', name: 'Ben Hartley', aptNumber: 'Apt 1B', email: 'ben@brunswickheights.local', entitlementRemainingM3: 0.5, entitlementYear: FY },
  { id: 'r_chen', name: 'Chen Liu', aptNumber: 'Apt 2A', email: 'chen@brunswickheights.local', entitlementRemainingM3: 1.0, entitlementYear: FY },
  { id: 'r_dani', name: 'Dani Okafor', aptNumber: 'Apt 2B', email: 'dani@brunswickheights.local', entitlementRemainingM3: 0.25, entitlementYear: FY },
  { id: 'r_elena', name: 'Elena Rossi', aptNumber: 'Apt 3A', email: 'elena@brunswickheights.local', entitlementRemainingM3: 0, entitlementYear: FY },
  { id: 'r_finn', name: 'Finn Whittaker', aptNumber: 'Apt 3B', email: 'finn@brunswickheights.local', entitlementRemainingM3: 1.0, entitlementYear: FY },
  { id: 'r_gita', name: 'Gita Iyer', aptNumber: 'Apt 4A', email: 'gita@brunswickheights.local', entitlementRemainingM3: 0.75, entitlementYear: FY },
  { id: 'r_hugo', name: 'Hugo Martinez', aptNumber: 'Apt 4B', email: 'hugo@brunswickheights.local', entitlementRemainingM3: 1.0, entitlementYear: FY },
  { id: 'r_isla', name: 'Isla Nguyen', aptNumber: 'Apt 5A', email: 'isla@brunswickheights.local', entitlementRemainingM3: 0.5, entitlementYear: FY },
  { id: 'r_jamal', name: 'Jamal Yusuf', aptNumber: 'Apt 5B', email: 'jamal@brunswickheights.local', entitlementRemainingM3: 1.0, entitlementYear: FY },
  { id: 'r_kira', name: 'Kira Petrov', aptNumber: 'Apt 6A', email: 'kira@brunswickheights.local', entitlementRemainingM3: 1.0, entitlementYear: FY },
  { id: 'r_leo', name: 'Leo Park', aptNumber: 'Apt 6B', email: 'leo@brunswickheights.local', entitlementRemainingM3: 0.5, entitlementYear: FY },
]

const photo = (seed: string) =>
  `https://picsum.photos/seed/${seed}/600/400`

export const LISTINGS: FurnitureListing[] = [
  {
    id: 'l_couch',
    postedById: 'r_amelia',
    title: 'Two-seater couch, charcoal',
    description: 'Comfy IKEA Karlstad two-seater. Some wear on one armrest but very sturdy. Great for a first apartment or study nook. Pickup from foyer.',
    photoUrl: photo('couch-charcoal'),
    category: 'furniture',
    estimatedM3: 0.6,
    pickupBy: inDays(12),
    status: 'available',
    createdAt: daysAgo(3),
  },
  {
    id: 'l_bookshelf',
    postedById: 'r_chen',
    title: 'Tall pine bookshelf (5 shelf)',
    description: '180cm tall, solid pine, all hardware intact. Moving interstate next month — would love it to stay in the building.',
    photoUrl: photo('bookshelf-pine'),
    category: 'furniture',
    estimatedM3: 0.3,
    pickupBy: inDays(18),
    status: 'available',
    createdAt: daysAgo(1),
  },
  {
    id: 'l_desk',
    postedById: 'r_finn',
    title: 'IKEA Linnmon desk + Alex drawers',
    description: 'White desk (120x60) plus 3-drawer Alex unit. Used for WFH the last two years. Small scratch on the desktop but otherwise excellent.',
    photoUrl: photo('desk-linnmon'),
    category: 'furniture',
    estimatedM3: 0.5,
    pickupBy: inDays(9),
    status: 'claimed',
    claimedById: 'r_kira',
    createdAt: daysAgo(5),
  },
  {
    id: 'l_mattress',
    postedById: 'r_jamal',
    title: 'Queen mattress (~3 years old)',
    description: 'Sealy Posturepedic queen. Always used with a mattress protector. Free to whoever can carry it down two flights.',
    photoUrl: photo('mattress-queen'),
    category: 'mattress',
    estimatedM3: 0.8,
    pickupBy: inDays(20),
    status: 'available',
    createdAt: daysAgo(2),
  },
  {
    id: 'l_fridge',
    postedById: 'r_isla',
    title: 'Bar fridge, runs cold',
    description: 'Husky 90L bar fridge, perfect for a studio. Working order — defrosted and cleaned. Lithium battery in the temp display has been removed.',
    photoUrl: photo('bar-fridge-husky'),
    category: 'whitegoods',
    estimatedM3: 0.4,
    pickupBy: inDays(15),
    status: 'available',
    createdAt: daysAgo(2),
  },
  {
    id: 'l_chairs',
    postedById: 'r_gita',
    title: 'Set of 4 dining chairs, oak',
    description: 'Solid oak chairs, matching set of 4. One has a loose joint that needs a screw. Otherwise solid.',
    photoUrl: photo('chairs-oak'),
    category: 'furniture',
    estimatedM3: 0.4,
    pickupBy: inDays(14),
    status: 'available',
    createdAt: daysAgo(6),
  },
  {
    id: 'l_lamp',
    postedById: 'r_hugo',
    title: 'Floor lamp, brass + linen shade',
    description: 'Mid-century style brass floor lamp, takes a standard E27 bulb. Bulb not included.',
    photoUrl: photo('lamp-brass'),
    category: 'furniture',
    estimatedM3: 0.1,
    pickupBy: inDays(7),
    status: 'collected',
    claimedById: 'r_isla',
    createdAt: daysAgo(10),
  },
]

export const SHARES: EntitlementShare[] = [
  {
    id: 's_kira_offer',
    offeredById: 'r_kira',
    m3Amount: 0.5,
    note: "Won't be using mine this year — happy for it to help a neighbour.",
    status: 'offered',
    createdAt: daysAgo(4),
  },
  {
    id: 's_hugo_offer',
    offeredById: 'r_hugo',
    m3Amount: 0.25,
    note: 'Have a bit spare after the lamp went to a good home.',
    status: 'offered',
    createdAt: daysAgo(2),
  },
  {
    id: 's_finn_claimed',
    offeredById: 'r_finn',
    m3Amount: 0.5,
    note: 'Take it, just give credit at the door if I ever need a favour.',
    status: 'claimed',
    claimedById: 'r_ben',
    createdAt: daysAgo(8),
  },
  {
    id: 's_jamal_claimed',
    offeredById: 'r_jamal',
    m3Amount: 0.5,
    status: 'claimed',
    claimedById: 'r_dani',
    createdAt: daysAgo(11),
  },
]

export const COLLECTION_DAY: CollectionDay = {
  id: 'cd_2026_q1',
  scheduledDate: inDays(21),
  organizerId: 'r_amelia',
  status: 'planned',
}

export const COLLECTION_ITEMS: CollectionItem[] = [
  { id: 'ci_1', collectionDayId: 'cd_2026_q1', residentId: 'r_elena', title: 'Old microwave + broken vacuum', category: 'ewaste', estimatedM3: 0.4 },
  { id: 'ci_2', collectionDayId: 'cd_2026_q1', residentId: 'r_dani', title: 'Cracked TV stand', category: 'furniture', estimatedM3: 0.5 },
  { id: 'ci_3', collectionDayId: 'cd_2026_q1', residentId: 'r_ben', title: 'Office chair (broken hydraulic)', category: 'furniture', estimatedM3: 0.3 },
  { id: 'ci_4', collectionDayId: 'cd_2026_q1', residentId: 'r_leo', title: 'Two bedside tables', category: 'furniture', estimatedM3: 0.5 },
  { id: 'ci_5', collectionDayId: 'cd_2026_q1', residentId: 'r_dani', title: 'Rolled-up rug', category: 'furniture', estimatedM3: 0.3 },
  { id: 'ci_6', collectionDayId: 'cd_2026_q1', residentId: 'r_finn', title: 'Old dishwasher (door removed)', category: 'whitegoods', estimatedM3: 0.4 },
]
