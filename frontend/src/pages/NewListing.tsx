import { useState, type FormEvent } from 'react'
import { Link, useNavigate } from 'react-router-dom'
import { AlertTriangle, ArrowLeft, Info } from 'lucide-react'
import { useAppData } from '../context/AppDataContext'
import { CATEGORY_LABEL, type ItemCategory } from '../types'
import { Button } from '../components/ui/Button'
import { Input } from '../components/ui/Input'
import { Label } from '../components/ui/Label'
import { Textarea } from '../components/ui/Textarea'
import { Select } from '../components/ui/Select'
import {
  Card,
  CardContent,
  CardDescription,
  CardHeader,
  CardTitle,
} from '../components/ui/Card'

const photoFor = (title: string) =>
  `https://picsum.photos/seed/${encodeURIComponent(title.toLowerCase().replace(/\s+/g, '-') || 'furniture')}/600/400`

const RESTRICTIONS: Partial<Record<ItemCategory, { tone: 'warn' | 'info'; text: string }[]>> = {
  ewaste: [
    {
      tone: 'warn',
      text: 'Council rule: remove lithium batteries from e-waste before collection.',
    },
  ],
  whitegoods: [
    {
      tone: 'warn',
      text: 'Council rule: remove washing-machine doors before kerbside placement.',
    },
    {
      tone: 'info',
      text: 'White goods only accepted whole — no parts (e.g. drums, motors).',
    },
  ],
  mattress: [
    {
      tone: 'info',
      text: 'Tip: keep mattresses dry — place out the day before, not earlier.',
    },
  ],
  furniture: [],
}

export function NewListing() {
  const navigate = useNavigate()
  const { createListing } = useAppData()
  const [title, setTitle] = useState('')
  const [description, setDescription] = useState('')
  const [photoUrl, setPhotoUrl] = useState('')
  const [category, setCategory] = useState<ItemCategory>('furniture')
  const [estimatedM3, setEstimatedM3] = useState('0.3')
  const [pickupBy, setPickupBy] = useState(
    new Date(Date.now() + 14 * 86400000).toISOString().slice(0, 10),
  )

  const handleSubmit = (e: FormEvent) => {
    e.preventDefault()
    const id = createListing({
      title: title.trim(),
      description: description.trim(),
      photoUrl: photoUrl.trim() || photoFor(title),
      category,
      estimatedM3: parseFloat(estimatedM3) || 0,
      pickupBy: new Date(pickupBy).toISOString(),
    })
    if (id) navigate(`/listings/${id}`)
  }

  const reminders = RESTRICTIONS[category] ?? []

  return (
    <div className="max-w-2xl mx-auto space-y-6">
      <Link
        to="/marketplace"
        className="inline-flex items-center gap-1.5 text-sm text-slate-600 hover:text-slate-900"
      >
        <ArrowLeft className="h-4 w-4" /> Back to give-aways
      </Link>

      <Card>
        <CardHeader>
          <CardTitle>Post a furniture item</CardTitle>
          <CardDescription>
            Give a piece to a neighbour before it becomes hard waste.
          </CardDescription>
        </CardHeader>
        <CardContent>
          <form onSubmit={handleSubmit} className="space-y-4">
            <div className="space-y-1.5">
              <Label htmlFor="title">Title</Label>
              <Input
                id="title"
                placeholder="e.g. Two-seater couch, charcoal"
                required
                value={title}
                onChange={(e) => setTitle(e.target.value)}
              />
            </div>
            <div className="space-y-1.5">
              <Label htmlFor="description">Description</Label>
              <Textarea
                id="description"
                placeholder="Condition, dimensions, why you're giving it away..."
                required
                value={description}
                onChange={(e) => setDescription(e.target.value)}
              />
            </div>
            <div className="space-y-1.5">
              <Label htmlFor="category">Category</Label>
              <Select
                id="category"
                value={category}
                onChange={(e) => setCategory(e.target.value as ItemCategory)}
              >
                <option value="furniture">{CATEGORY_LABEL.furniture}</option>
                <option value="whitegoods">{CATEGORY_LABEL.whitegoods}</option>
                <option value="ewaste">{CATEGORY_LABEL.ewaste}</option>
                <option value="mattress">{CATEGORY_LABEL.mattress}</option>
              </Select>
              <p className="text-xs text-slate-500">
                Matches the four categories the City of Melbourne accepts at
                kerbside.
              </p>
            </div>

            {reminders.length > 0 && (
              <div className="space-y-2">
                {reminders.map((r, i) => (
                  <div
                    key={i}
                    className={
                      r.tone === 'warn'
                        ? 'flex items-start gap-2 rounded-md bg-amber-50 border border-amber-200 px-3 py-2 text-sm text-amber-900'
                        : 'flex items-start gap-2 rounded-md bg-slate-50 border border-slate-200 px-3 py-2 text-sm text-slate-700'
                    }
                  >
                    {r.tone === 'warn' ? (
                      <AlertTriangle className="h-4 w-4 mt-0.5 shrink-0" />
                    ) : (
                      <Info className="h-4 w-4 mt-0.5 shrink-0" />
                    )}
                    <span>{r.text}</span>
                  </div>
                ))}
              </div>
            )}

            <div className="grid grid-cols-1 sm:grid-cols-2 gap-4">
              <div className="space-y-1.5">
                <Label htmlFor="m3">Estimated size (m³)</Label>
                <Input
                  id="m3"
                  type="number"
                  step="0.1"
                  min="0.1"
                  max="1"
                  required
                  value={estimatedM3}
                  onChange={(e) => setEstimatedM3(e.target.value)}
                />
                <p className="text-xs text-slate-500">
                  Approximate volume if it ended up at the kerbside.
                </p>
              </div>
              <div className="space-y-1.5">
                <Label htmlFor="pickup">Available until</Label>
                <Input
                  id="pickup"
                  type="date"
                  required
                  value={pickupBy}
                  onChange={(e) => setPickupBy(e.target.value)}
                />
              </div>
            </div>
            <div className="space-y-1.5">
              <Label htmlFor="photo">Photo URL (optional)</Label>
              <Input
                id="photo"
                placeholder="Leave blank to use a placeholder image"
                value={photoUrl}
                onChange={(e) => setPhotoUrl(e.target.value)}
              />
            </div>
            <div className="flex gap-2 pt-2">
              <Button type="submit" size="md">
                Post item
              </Button>
              <Link
                to="/marketplace"
                className="inline-flex items-center justify-center h-10 px-4 rounded-lg text-sm font-medium text-slate-700 hover:bg-slate-100"
              >
                Cancel
              </Link>
            </div>
          </form>
        </CardContent>
      </Card>
    </div>
  )
}
