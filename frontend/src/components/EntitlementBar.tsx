import { cn, formatM3 } from '../lib/utils'

type Props = {
  remaining: number
  capacity?: number
  label?: string
  className?: string
}

export function EntitlementBar({
  remaining,
  capacity = 1,
  label,
  className,
}: Props) {
  const ratio = Math.max(0, Math.min(1, remaining / Math.max(capacity, 0.0001)))
  const pct = ratio * 100
  return (
    <div className={cn('space-y-2', className)}>
      <div className="flex items-baseline justify-between">
        <span className="text-sm font-medium text-slate-700">
          {label ?? 'Entitlement remaining'}
        </span>
        <span className="text-sm font-semibold text-slate-900">
          {formatM3(remaining)}{' '}
          <span className="text-slate-400 font-normal">
            of {formatM3(capacity)}
          </span>
        </span>
      </div>
      <div className="h-2.5 w-full rounded-full bg-slate-100 overflow-hidden">
        <div
          className={cn(
            'h-full rounded-full transition-all',
            ratio > 0.5
              ? 'bg-brand-600'
              : ratio > 0
                ? 'bg-amber-500'
                : 'bg-rose-500',
          )}
          style={{ width: `${pct}%` }}
        />
      </div>
    </div>
  )
}
