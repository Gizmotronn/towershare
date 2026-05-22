import { useState, type FormEvent } from 'react'
import { Link, useNavigate } from 'react-router-dom'
import { AlertTriangle, ArrowLeft, Info } from 'lucide-react'
import { useAppData } from '../context/AppDataContext'
import { useAuth } from '../context/AuthContext'
import { pb } from '../lib/pb'
import { CATEGORY_LABEL, type ItemCategory } from '../types'
import { Button } from '../components/ui/Button'
import { Input } from '../components/ui/Input'
import { Label } from '../components/ui/Label'
import { Textarea } from '../components/ui/Textarea'
import { Select } from '../components/ui/Select'
import { ImageUpload } from '../components/ImageUpload'
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
  const { isPbAuth, pbUser } = useAuth()
  const [title, setTitle] = useState('')
  const [description, setDescription] = useState('')
  const [imageFiles, setImageFiles] = useState<File[]>([])
  const [category, setCategory] = useState<ItemCategory>('furniture')
  const [estimatedM3, setEstimatedM3] = useState('0.3')
  const [pickupBy, setPickupBy] = useState(
    new Date(Date.now() + 14 * 86400000).toISOString().slice(0, 10),
  )
  const [submitting, setSubmitting] = useState(false)
  const [uploadError, setUploadError] = useState<string | null>(null)

  const handleSubmit = async (e: FormEvent) => {
    e.preventDefault()
    setUploadError(null)
    setSubmitting(true)

    try {
      let photoUrl = photoFor(title)
      let images: string[] = []

      if (isPbAuth && pbUser && imageFiles.length > 0) {
        // Upload to PocketBase and get back a real listing record
        const formData = new FormData()
        formData.append('title', title.trim())
        formData.append('description', description.trim())
        formData.append('category', category)
        formData.append('status', 'available')
        formData.append('owner', pbUser.id)
        formData.append('estimated_m3', estimatedM3)
        formData.append('pickup_by', new Date(pickupBy).toISOString())
        imageFiles.forEach((f) => formData.append('images', f))

        try {
          const record = await pb.collection('listings').create(formData)
          // Build PocketBase file URLs
          images = (record.images as string[]).map(
            (name) =>
              `${pb.baseUrl}/api/files/${record.collectionId}/${record.id}/${name}`,
          )
          photoUrl = images[0] ?? photoUrl
        } catch (err) {
          console.warn('PocketBase listing creation failed, using local preview', err)
          // Fall through to object URL fallback
        }
      }

      // Object URL fallback for demo or if PB upload failed
      if (images.length === 0 && imageFiles.length > 0) {
        images = imageFiles.map((f) => URL.createObjectURL(f))
        photoUrl = images[0]
      }

      const id = createListing({
        title: title.trim(),
        description: description.trim(),
        photoUrl,
        category,
        estimatedM3: parseFloat(estimatedM3) || 0,
        pickupBy: new Date(pickupBy).toISOString(),
        images,
      })
      if (id) navigate(`/listings/${id}`)
    } catch {
      setUploadError('Something went wrong. Please try again.')
    } finally {
      setSubmitting(false)
    }
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
              <Label>Photos</Label>
              <ImageUpload files={imageFiles} onChange={setImageFiles} maxFiles={5} />
              {isPbAuth && (
                <p className="text-xs text-brand-700">
                  Photos will be saved to your account.
                </p>
              )}
              {!isPbAuth && imageFiles.length > 0 && (
                <p className="text-xs text-slate-400">
                  Demo mode — photos are previewed locally and not persisted.
                </p>
              )}
            </div>

            {uploadError && (
              <p className="text-sm text-rose-600 bg-rose-50 border border-rose-100 rounded-md px-3 py-2">
                {uploadError}
              </p>
            )}

            <div className="flex gap-2 pt-2">
              <Button type="submit" size="md" disabled={submitting}>
                {submitting ? 'Posting…' : 'Post item'}
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
