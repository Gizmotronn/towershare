import { useRef, useState, useCallback, type DragEvent, type ChangeEvent } from 'react'
import { ImagePlus, X } from 'lucide-react'

type Preview = {
  file: File
  url: string
}

type Props = {
  files: File[]
  onChange: (files: File[]) => void
  maxFiles?: number
}

const ACCEPT = 'image/jpeg,image/png,image/webp,image/gif'
const MAX_SIZE_MB = 5

function validFiles(incoming: FileList | File[]): File[] {
  return Array.from(incoming).filter((f) => {
    if (!ACCEPT.split(',').includes(f.type)) return false
    if (f.size > MAX_SIZE_MB * 1024 * 1024) return false
    return true
  })
}

export function ImageUpload({ files, onChange, maxFiles = 5 }: Props) {
  const inputRef = useRef<HTMLInputElement>(null)
  const [dragging, setDragging] = useState(false)
  const [previews] = useState<Preview[]>([])

  // Keep preview object URLs in sync — derive from files prop
  const previewsForFiles: Preview[] = files.map((f, i) => ({
    file: f,
    url: previews[i]?.file === f ? previews[i].url : URL.createObjectURL(f),
  }))

  const addFiles = useCallback(
    (incoming: FileList | File[]) => {
      const valid = validFiles(incoming)
      const merged = [...files, ...valid].slice(0, maxFiles)
      onChange(merged)
    },
    [files, maxFiles, onChange],
  )

  const remove = (index: number) => {
    URL.revokeObjectURL(previewsForFiles[index].url)
    onChange(files.filter((_, i) => i !== index))
  }

  const onDragOver = (e: DragEvent) => {
    e.preventDefault()
    setDragging(true)
  }
  const onDragLeave = () => setDragging(false)
  const onDrop = (e: DragEvent) => {
    e.preventDefault()
    setDragging(false)
    if (e.dataTransfer.files.length) addFiles(e.dataTransfer.files)
  }
  const onInputChange = (e: ChangeEvent<HTMLInputElement>) => {
    if (e.target.files?.length) addFiles(e.target.files)
    e.target.value = ''
  }

  const canAdd = files.length < maxFiles

  return (
    <div className="space-y-3">
      {canAdd && (
        <button
          type="button"
          onClick={() => inputRef.current?.click()}
          onDragOver={onDragOver}
          onDragLeave={onDragLeave}
          onDrop={onDrop}
          className={[
            'w-full flex flex-col items-center justify-center gap-2 rounded-xl border-2 border-dashed px-4 py-8 text-sm transition-colors cursor-pointer',
            dragging
              ? 'border-brand-500 bg-brand-50 text-brand-700'
              : 'border-slate-200 text-slate-400 hover:border-brand-400 hover:bg-brand-50/40 hover:text-brand-600',
          ].join(' ')}
        >
          <ImagePlus className="h-7 w-7" />
          <span>
            Drop photos here or <span className="font-medium underline underline-offset-2">browse</span>
          </span>
          <span className="text-xs">
            JPEG, PNG or WebP · up to {MAX_SIZE_MB} MB each · {maxFiles - files.length} slot
            {maxFiles - files.length !== 1 ? 's' : ''} left
          </span>
        </button>
      )}

      <input
        ref={inputRef}
        type="file"
        accept={ACCEPT}
        multiple
        className="sr-only"
        onChange={onInputChange}
      />

      {previewsForFiles.length > 0 && (
        <div className="flex flex-wrap gap-2">
          {previewsForFiles.map((p, i) => (
            <div key={p.url} className="relative group h-24 w-24 shrink-0">
              <img
                src={p.url}
                alt={`Photo ${i + 1}`}
                className="h-full w-full rounded-lg object-cover border border-slate-200"
              />
              <button
                type="button"
                onClick={() => remove(i)}
                className="absolute -top-1.5 -right-1.5 h-5 w-5 rounded-full bg-slate-900 text-white flex items-center justify-center opacity-0 group-hover:opacity-100 transition-opacity"
                aria-label="Remove photo"
              >
                <X className="h-3 w-3" />
              </button>
              {i === 0 && (
                <span className="absolute bottom-1 left-1 text-[10px] font-medium bg-black/60 text-white rounded px-1 leading-5">
                  Cover
                </span>
              )}
            </div>
          ))}
        </div>
      )}
    </div>
  )
}
