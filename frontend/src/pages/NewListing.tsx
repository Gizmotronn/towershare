import { useState, type FormEvent } from 'react'
import { Link, useNavigate } from 'react-router-dom'
import { ArrowLeft } from 'lucide-react'
import { useAppData } from '../context/AppDataContext'
import { useAuth } from '../context/AuthContext'
import { Button } from '../components/ui/Button'
import { Input } from '../components/ui/Input'
import { Label } from '../components/ui/Label'
import { Textarea } from '../components/ui/Textarea'
import { Select } from '../components/ui/Select'
import { pbCreateListing } from '../lib/pb'
import {
  Card,
  CardContent,
  CardDescription,
  CardHeader,
  CardTitle,
} from '../components/ui/Card'

type PbCategory = 'give' | 'lend' | 'request'

const CATEGORY_LABEL: Record<PbCategory, string> = {
  give:    'Give away (free)',
  lend:    'Lend (borrow & return)',
  request: 'Request (looking for)',
}

export function NewListing() {
  const navigate = useNavigate()
  const { createListing, building } = useAppData()
  const { pbUser } = useAuth()
  const [title, setTitle] = useState('')
  const [description, setDescription] = useState('')
  const [category, setCategory] = useState<PbCategory>('give')
  const [imageFile, setImageFile] = useState<File | null>(null)
  const [submitting, setSubmitting] = useState(false)
  const [error, setError] = useState<string | null>(null)

  const handleSubmit = async (e: FormEvent) => {
    e.preventDefault()
    setError(null)

    if (pbUser) {
      if (!building.id) {
        setError('No tower found. Is the backend running?')
        return
      }
      setSubmitting(true)
      try {
        const fd = new FormData()
        fd.append('title', title.trim())
        fd.append('description', description.trim())
        fd.append('category', category)
        fd.append('status', 'active')
        fd.append('owner', pbUser.id)
        fd.append('tower', building.id)
        if (imageFile) fd.append('images', imageFile)
        await pbCreateListing(fd)
        navigate('/marketplace')
      } catch (err) {
        setError(err instanceof Error ? err.message : 'Failed to create listing')
      } finally {
        setSubmitting(false)
      }
      return
    }

    // No PocketBase session — write to local (demo) state only.
    const photoUrl = imageFile
      ? URL.createObjectURL(imageFile)
      : `https://picsum.photos/seed/${encodeURIComponent(title.toLowerCase().replace(/\s+/g, '-') || 'item')}/600/400`

    const id = createListing({
      title: title.trim(),
      description: description.trim(),
      photoUrl,
      category: 'furniture', // local seed uses item-type categories; default to furniture
      estimatedM3: 0.3,
      pickupBy: new Date(Date.now() + 14 * 86400000).toISOString(),
    })
    if (id) navigate(`/listings/${id}`)
  }

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
              <Label htmlFor="category">Listing type</Label>
              <Select
                id="category"
                value={category}
                onChange={(e) => setCategory(e.target.value as PbCategory)}
              >
                {(Object.keys(CATEGORY_LABEL) as PbCategory[]).map((k) => (
                  <option key={k} value={k}>
                    {CATEGORY_LABEL[k]}
                  </option>
                ))}
              </Select>
            </div>

            <div className="space-y-1.5">
              <Label htmlFor="photo">Photo (optional)</Label>
              <input
                id="photo"
                type="file"
                accept="image/jpeg,image/png,image/webp,image/gif"
                className="block w-full text-sm text-slate-700 file:mr-3 file:py-1.5 file:px-3 file:rounded-md file:border-0 file:text-sm file:font-medium file:bg-slate-100 file:text-slate-700 hover:file:bg-slate-200 cursor-pointer"
                onChange={(e) => setImageFile(e.target.files?.[0] ?? null)}
              />
              <p className="text-xs text-slate-500">JPEG, PNG, WebP or GIF · max 5 MB</p>
            </div>

            {error && (
              <p className="text-sm text-rose-600 bg-rose-50 border border-rose-100 rounded-md px-3 py-2">
                {error}
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
