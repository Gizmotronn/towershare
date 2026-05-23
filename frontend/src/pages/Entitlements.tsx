import { useState, type FormEvent } from 'react'
import { Gift, HandHeart, Sparkles } from 'lucide-react'
import { useAppData } from '../context/AppDataContext'
import {
  Card,
  CardContent,
  CardDescription,
  CardHeader,
  CardTitle,
} from '../components/ui/Card'
import { Button } from '../components/ui/Button'
import { Input } from '../components/ui/Input'
import { Label } from '../components/ui/Label'
import { Textarea } from '../components/ui/Textarea'
import { Badge } from '../components/ui/Badge'
import { EntitlementBar } from '../components/EntitlementBar'
import { formatDateShort, formatM3 } from '../lib/utils'

export function Entitlements() {
  const { shares, currentUser, getResident, offerShare, claimShare } =
    useAppData()
  const [m2, setM2] = useState('0.25')
  const [note, setNote] = useState('')
  const [error, setError] = useState<string | null>(null)
  const [flash, setFlash] = useState<string | null>(null)

  if (!currentUser) return null

  const availableShares = shares.filter(
    (s) => s.status === 'offered' && s.offeredById !== currentUser.id,
  )
  const mySharesOffered = shares.filter((s) => s.offeredById === currentUser.id)
  const sharesIReceived = shares.filter(
    (s) => s.status === 'claimed' && s.claimedById === currentUser.id,
  )

  const handleOffer = (e: FormEvent) => {
    e.preventDefault()
    setError(null)
    setFlash(null)
    const amount = parseFloat(m2)
    if (isNaN(amount) || amount <= 0) {
      setError('Pick an amount greater than 0.')
      return
    }
    if (amount > currentUser.entitlementRemainingM3) {
      setError(
        `You only have ${formatM3(currentUser.entitlementRemainingM3)} left to share.`,
      )
      return
    }
    const ok = offerShare(amount, note.trim() || undefined)
    if (ok) {
      setFlash(`Offered ${formatM3(amount)} to the block.`)
      setM2('0.25')
      setNote('')
    }
  }

  const handleClaim = (id: string) => {
    claimShare(id)
    setFlash('Share claimed — your balance has been topped up.')
  }

  return (
    <div className="space-y-6">
      <div>
        <h1 className="text-2xl font-semibold text-slate-900">
          Share your m³
        </h1>
        <p className="text-sm text-slate-500 mt-1">
          Lend or claim unused hard-waste entitlement. Each resident starts
          with 1m³ per financial year.
        </p>
      </div>

      <Card>
        <CardHeader>
          <CardTitle>Your balance</CardTitle>
        </CardHeader>
        <CardContent>
          <EntitlementBar
            remaining={currentUser.entitlementRemainingM3}
            capacity={1}
          />
        </CardContent>
      </Card>

      {flash && (
        <div className="rounded-lg bg-brand-50 border border-brand-100 px-4 py-3 text-sm text-brand-900 flex items-center gap-2">
          <Sparkles className="h-4 w-4" />
          {flash}
        </div>
      )}

      <div className="grid grid-cols-1 lg:grid-cols-2 gap-5">
        <Card>
          <CardHeader>
            <div className="flex items-center gap-2">
              <Gift className="h-5 w-5 text-brand-700" />
              <CardTitle>Available from neighbours</CardTitle>
            </div>
            <CardDescription>
              Free m³ being offered right now.
            </CardDescription>
          </CardHeader>
          <CardContent className="space-y-3">
            {availableShares.length === 0 && (
              <p className="text-sm text-slate-500">
                Nothing on offer at the moment.
              </p>
            )}
            {availableShares.map((share) => {
              const offerer = getResident(share.offeredById)
              return (
                <div
                  key={share.id}
                  className="border border-slate-200 rounded-lg p-4 flex items-start gap-3"
                >
                  <div className="h-10 w-10 rounded-full bg-brand-100 text-brand-700 grid place-items-center text-sm font-semibold shrink-0">
                    {offerer?.name
                      .split(' ')
                      .map((p) => p[0])
                      .slice(0, 2)
                      .join('')}
                  </div>
                  <div className="flex-1 min-w-0">
                    <div className="flex items-center gap-2 flex-wrap">
                      <p className="font-medium text-sm text-slate-900">
                        {offerer?.name}
                      </p>
                      <span className="text-xs text-slate-400">
                        · {offerer?.aptNumber}
                      </span>
                      <Badge variant="default">
                        {formatM3(share.m3Amount)}
                      </Badge>
                    </div>
                    {share.note && (
                      <p className="text-sm text-slate-600 mt-1.5">
                        "{share.note}"
                      </p>
                    )}
                    <p className="text-xs text-slate-400 mt-1">
                      Offered {formatDateShort(share.createdAt)}
                    </p>
                  </div>
                  <Button size="sm" onClick={() => handleClaim(share.id)}>
                    Claim
                  </Button>
                </div>
              )
            })}
          </CardContent>
        </Card>

        <Card>
          <CardHeader>
            <div className="flex items-center gap-2">
              <HandHeart className="h-5 w-5 text-brand-700" />
              <CardTitle>Offer some of yours</CardTitle>
            </div>
            <CardDescription>
              Won't need it this year? Make it available to a neighbour.
            </CardDescription>
          </CardHeader>
          <CardContent>
            <form onSubmit={handleOffer} className="space-y-4">
              <div className="space-y-1.5">
                <Label htmlFor="m2">Amount to offer (m³)</Label>
                <Input
                  id="m2"
                  type="number"
                  step="0.05"
                  min="0.05"
                  max={currentUser.entitlementRemainingM3}
                  value={m2}
                  onChange={(e) => setM2(e.target.value)}
                  required
                />
              </div>
              <div className="space-y-1.5">
                <Label htmlFor="note">Note (optional)</Label>
                <Textarea
                  id="note"
                  placeholder="Anything you want neighbours to know about this offer..."
                  value={note}
                  onChange={(e) => setNote(e.target.value)}
                />
              </div>
              {error && (
                <p className="text-sm text-rose-600 bg-rose-50 border border-rose-100 rounded-md px-3 py-2">
                  {error}
                </p>
              )}
              <Button
                type="submit"
                disabled={currentUser.entitlementRemainingM3 <= 0}
              >
                Offer to the block
              </Button>
              {currentUser.entitlementRemainingM3 <= 0 && (
                <p className="text-xs text-slate-500">
                  You don't have any entitlement left to offer right now.
                </p>
              )}
            </form>
          </CardContent>
        </Card>
      </div>

      <div className="grid grid-cols-1 lg:grid-cols-2 gap-5">
        <Card>
          <CardHeader>
            <CardTitle className="text-base">
              Shares you've offered
            </CardTitle>
          </CardHeader>
          <CardContent className="space-y-3">
            {mySharesOffered.length === 0 ? (
              <p className="text-sm text-slate-500">
                You haven't offered any m³ yet.
              </p>
            ) : (
              mySharesOffered.map((share) => {
                const claimer = getResident(share.claimedById)
                return (
                  <div
                    key={share.id}
                    className="flex items-center justify-between gap-3 py-2 border-b border-slate-100 last:border-0"
                  >
                    <div>
                      <p className="text-sm font-medium text-slate-900">
                        {formatM3(share.m3Amount)}
                        {share.note && (
                          <span className="text-slate-500 font-normal">
                            {' '}
                            — "{share.note.slice(0, 40)}
                            {share.note.length > 40 ? '…' : ''}"
                          </span>
                        )}
                      </p>
                      <p className="text-xs text-slate-500">
                        {share.status === 'claimed'
                          ? `Claimed by ${claimer?.name ?? 'a neighbour'}`
                          : 'Still on offer'}
                      </p>
                    </div>
                    <Badge
                      variant={
                        share.status === 'claimed' ? 'neutral' : 'success'
                      }
                    >
                      {share.status === 'claimed' ? 'Claimed' : 'Offered'}
                    </Badge>
                  </div>
                )
              })
            )}
          </CardContent>
        </Card>

        <Card>
          <CardHeader>
            <CardTitle className="text-base">
              Shares you've received
            </CardTitle>
          </CardHeader>
          <CardContent className="space-y-3">
            {sharesIReceived.length === 0 ? (
              <p className="text-sm text-slate-500">
                You haven't claimed any shared m³ yet.
              </p>
            ) : (
              sharesIReceived.map((share) => {
                const offerer = getResident(share.offeredById)
                return (
                  <div
                    key={share.id}
                    className="flex items-center justify-between gap-3 py-2 border-b border-slate-100 last:border-0"
                  >
                    <div>
                      <p className="text-sm font-medium text-slate-900">
                        {formatM3(share.m3Amount)} from {offerer?.name}
                      </p>
                      <p className="text-xs text-slate-500">
                        {offerer?.aptNumber}
                      </p>
                    </div>
                    <Badge variant="default">Received</Badge>
                  </div>
                )
              })
            )}
          </CardContent>
        </Card>
      </div>
    </div>
  )
}
