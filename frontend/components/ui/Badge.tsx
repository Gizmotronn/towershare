import type { HTMLAttributes } from 'react'
import { cn } from '../../lib/utils'

type Variant = 'default' | 'success' | 'warning' | 'neutral' | 'outline'

const variantClasses: Record<Variant, string> = {
  default: 'bg-brand-50 text-brand-700 border-brand-100',
  success: 'bg-emerald-50 text-emerald-700 border-emerald-100',
  warning: 'bg-amber-50 text-amber-700 border-amber-100',
  neutral: 'bg-slate-100 text-slate-700 border-slate-200',
  outline: 'bg-white text-slate-700 border-slate-300',
}

export function Badge({
  className,
  variant = 'default',
  ...props
}: HTMLAttributes<HTMLSpanElement> & { variant?: Variant }) {
  return (
    <span
      className={cn(
        'inline-flex items-center gap-1 rounded-full border px-2.5 py-0.5 text-xs font-medium',
        variantClasses[variant],
        className,
      )}
      {...props}
    />
  )
}
