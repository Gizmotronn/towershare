import { useEffect, useState } from 'react'
import { Link } from 'react-router-dom'
import { ArrowLeft, Building2, Printer, QrCode } from 'lucide-react'
import { useAppData } from '../context/AppDataContext'
import { Button } from '../components/ui/Button'

export function QrPoster() {
  const { building } = useAppData()
  const [origin, setOrigin] = useState('')

  useEffect(() => {
    setOrigin(window.location.origin)
  }, [])

  const target = `${origin}/welcome?b=${building.id}`
  const qrSrc = origin
    ? `https://api.qrserver.com/v1/create-qr-code/?size=600x600&margin=2&data=${encodeURIComponent(
        target,
      )}`
    : ''

  return (
    <div className="min-h-screen bg-slate-100">
      <div className="print:hidden bg-white border-b border-slate-200">
        <div className="max-w-4xl mx-auto px-4 sm:px-6 py-3 flex items-center justify-between">
          <Link
            to="/welcome"
            className="inline-flex items-center gap-1.5 text-sm text-slate-600 hover:text-slate-900"
          >
            <ArrowLeft className="h-4 w-4" /> Back to welcome
          </Link>
          <Button
            size="sm"
            variant="outline"
            onClick={() => window.print()}
            className="bg-white"
          >
            <Printer className="h-4 w-4" /> Print poster
          </Button>
        </div>
      </div>

      <div className="max-w-4xl mx-auto p-4 sm:p-8">
        <div className="print:hidden mb-4 text-sm text-slate-600">
          <p>
            Pin this in the lobby. Residents scan with their phone camera and
            land on the welcome screen.
          </p>
        </div>

        <div
          className="bg-white rounded-2xl shadow-xl print:shadow-none print:rounded-none aspect-[1/1.414] flex flex-col"
          id="poster"
        >
          <div className="flex-1 flex flex-col items-center justify-between p-8 sm:p-12 text-center">
            <div className="flex items-center gap-3">
              <div className="h-12 w-12 rounded-xl bg-brand-600 grid place-items-center text-white">
                <Building2 className="h-6 w-6" />
              </div>
              <div className="leading-tight text-left">
                <p className="text-xs uppercase tracking-wider text-brand-700 font-semibold">
                  Hard Waste Hub
                </p>
                <p className="text-sm text-slate-600">{building.name}</p>
              </div>
            </div>

            <div>
              <p className="text-xs uppercase tracking-[0.2em] text-brand-700 font-semibold">
                Building noticeboard
              </p>
              <h1 className="mt-3 text-3xl sm:text-5xl font-bold text-slate-900 leading-tight max-w-md">
                Share, swap & save.
              </h1>
              <p className="mt-4 text-base sm:text-lg text-slate-700 max-w-md mx-auto leading-snug">
                Don't dump it. List it for your neighbours, or pool your unused
                council waste entitlement.
              </p>
            </div>

            <div className="rounded-2xl bg-white p-4 border-2 border-slate-900">
              {qrSrc ? (
                <img
                  src={qrSrc}
                  alt={`QR code linking to ${target}`}
                  className="h-56 w-56 sm:h-64 sm:w-64 block"
                />
              ) : (
                <div className="h-56 w-56 sm:h-64 sm:w-64 grid place-items-center bg-slate-100 rounded-lg">
                  <QrCode className="h-10 w-10 text-slate-400" />
                </div>
              )}
            </div>

            <div>
              <p className="text-xl sm:text-2xl font-semibold text-slate-900">
                Scan with your phone camera
              </p>
              <p className="text-sm text-slate-500 mt-1 break-all max-w-md mx-auto">
                {target || '—'}
              </p>
            </div>

            <div className="text-[11px] uppercase tracking-wider text-slate-400 font-semibold">
              {building.address}
            </div>
          </div>
        </div>

        <div className="print:hidden mt-4 text-xs text-slate-500">
          <p>
            Tip: on a phone, open the QR poster page directly and use your
            laptop camera to scan it — the prototype works end-to-end without
            any printing.
          </p>
        </div>
      </div>
    </div>
  )
}
