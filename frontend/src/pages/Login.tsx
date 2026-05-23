import { useEffect, useState, type FormEvent } from 'react'
import { Link, useNavigate, useSearchParams } from 'react-router-dom'
import { ArrowLeft, Building2, Sparkles } from 'lucide-react'
import { useAuth } from '../context/AuthContext'
import { useAppData } from '../context/AppDataContext'
import { Button } from '../components/ui/Button'
import { Input } from '../components/ui/Input'
import { Label } from '../components/ui/Label'
import { Select } from '../components/ui/Select'
import { formatM3 } from '../lib/utils'

export function Login() {
  const navigate = useNavigate()
  const [params] = useSearchParams()
  const fromQr = params.has('b')
  const { login, loginAs, currentUserId } = useAuth()
  const { residents, building } = useAppData()
  const [email, setEmail] = useState('')
  const [password, setPassword] = useState('')
  const [error, setError] = useState<string | null>(null)
  const [personaId, setPersonaId] = useState<string>('r_amelia')

  useEffect(() => {
    if (currentUserId) navigate('/dashboard', { replace: true })
  }, [currentUserId, navigate])

  const handleSubmit = (e: FormEvent) => {
    e.preventDefault()
    setError(null)
    const ok = login(email)
    if (!ok) {
      setError(
        'No resident found with that email. Try the demo login below to switch personas.',
      )
      return
    }
    navigate('/dashboard', { replace: true })
  }

  const handlePersona = () => {
    loginAs(personaId)
    navigate('/dashboard', { replace: true })
  }

  return (
    <div className="min-h-screen bg-gradient-to-br from-brand-50 via-white to-slate-100 flex items-center justify-center p-4">
      <div className="w-full max-w-md">
        <Link
          to="/welcome"
          className="inline-flex items-center gap-1.5 text-sm text-slate-500 hover:text-slate-900 mb-4"
        >
          <ArrowLeft className="h-4 w-4" /> Back
        </Link>
        <div className="flex items-center gap-3 justify-center mb-6">
          <div className="h-12 w-12 rounded-xl bg-brand-600 grid place-items-center text-white shadow-lg shadow-brand-600/20">
            <Building2 className="h-6 w-6" />
          </div>
          <div className="leading-tight">
            <p className="font-semibold text-slate-900 text-lg">
              {building.name}
            </p>
            <p className="text-sm text-slate-500">Hard Waste Hub</p>
          </div>
        </div>

        {fromQr && (
          <div className="mb-4 rounded-lg bg-brand-50 border border-brand-100 px-3 py-2 text-xs text-brand-900 flex items-center gap-2">
            <Sparkles className="h-3.5 w-3.5 shrink-0" />
            <span>
              Scanned from the lobby — sign in as one of the 12 residents to
              continue.
            </span>
          </div>
        )}

        <div className="bg-white rounded-2xl shadow-xl border border-slate-200 p-6 sm:p-8 space-y-6">
          <div>
            <h1 className="text-2xl font-semibold text-slate-900">
              Welcome back
            </h1>
            <p className="text-sm text-slate-500 mt-1">
              Sign in to share furniture and pool your hard-waste entitlement
              with your neighbours.
            </p>
          </div>

          <form onSubmit={handleSubmit} className="space-y-4">
            <div className="space-y-1.5">
              <Label htmlFor="email">Email</Label>
              <Input
                id="email"
                type="email"
                placeholder="you@brunswickheights.local"
                value={email}
                onChange={(e) => setEmail(e.target.value)}
                autoComplete="email"
              />
            </div>
            <div className="space-y-1.5">
              <Label htmlFor="password">Password</Label>
              <Input
                id="password"
                type="password"
                placeholder="••••••••"
                value={password}
                onChange={(e) => setPassword(e.target.value)}
                autoComplete="current-password"
              />
            </div>
            {error && (
              <p className="text-sm text-rose-600 bg-rose-50 border border-rose-100 rounded-md px-3 py-2">
                {error}
              </p>
            )}
            <Button type="submit" size="lg" className="w-full">
              Sign in
            </Button>
          </form>

          <div className="relative">
            <div className="absolute inset-0 flex items-center">
              <div className="w-full border-t border-slate-200" />
            </div>
            <div className="relative flex justify-center">
              <span className="bg-white px-3 text-xs uppercase tracking-wider text-slate-400 font-medium">
                Quick demo login
              </span>
            </div>
          </div>

          <div className="rounded-lg bg-brand-50/60 border border-brand-100 p-4 space-y-3">
            <div className="flex items-start gap-2">
              <Sparkles className="h-4 w-4 text-brand-700 mt-0.5 shrink-0" />
              <p className="text-xs text-brand-900 leading-relaxed">
                This prototype seeds 12 residents. Pick any of them to demo both
                sides of a furniture swap or m³ share.
              </p>
            </div>
            <div className="space-y-1.5">
              <Label htmlFor="persona">Sign in as</Label>
              <Select
                id="persona"
                value={personaId}
                onChange={(e) => setPersonaId(e.target.value)}
              >
                {residents.map((r) => (
                  <option key={r.id} value={r.id}>
                    {r.name} ({r.aptNumber}) —{' '}
                    {formatM3(r.entitlementRemainingM3)} left
                  </option>
                ))}
              </Select>
            </div>
            <Button
              type="button"
              variant="outline"
              size="md"
              className="w-full bg-white"
              onClick={handlePersona}
            >
              Continue as this resident
            </Button>
          </div>
        </div>

        <p className="text-center text-xs text-slate-400 mt-6">
          {building.address}
        </p>
      </div>
    </div>
  )
}
