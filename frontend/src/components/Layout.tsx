import { NavLink, Outlet, Link, useNavigate } from 'react-router-dom'
import { Building2, LogOut, User } from 'lucide-react'
import { useAuth } from '../context/AuthContext'
import { useAppData } from '../context/AppDataContext'
import { cn, formatM3 } from '../lib/utils'

const navItems = [
  { to: '/dashboard', label: 'Dashboard', end: true },
  { to: '/marketplace', label: 'Give-aways' },
  { to: '/entitlements', label: 'Share m³' },
  { to: '/collection-day', label: 'Collection day' },
  { to: '/profile', label: 'Profile' },
]

export function Layout() {
  const { logout } = useAuth()
  const { currentUser, building } = useAppData()
  const navigate = useNavigate()

  const handleLogout = () => {
    logout()
    navigate('/login')
  }

  return (
    <div className="min-h-full flex flex-col">
      <header className="bg-white border-b border-slate-200 sticky top-0 z-10">
        <div className="max-w-6xl mx-auto px-4 sm:px-6 lg:px-8">
          <div className="flex items-center justify-between h-16">
            <Link to="/dashboard" className="flex items-center gap-2.5">
              <div className="h-9 w-9 rounded-lg bg-brand-600 grid place-items-center text-white">
                <Building2 className="h-5 w-5" />
              </div>
              <div className="leading-tight">
                <p className="font-semibold text-slate-900 text-sm">
                  {building.name}
                </p>
                <p className="text-xs text-slate-500">Hard Waste Hub</p>
              </div>
            </Link>

            <nav className="hidden md:flex items-center gap-1">
              {navItems.map((item) => (
                <NavLink
                  key={item.to}
                  to={item.to}
                  end={item.end}
                  className={({ isActive }) =>
                    cn(
                      'px-3 py-1.5 rounded-md text-sm font-medium transition-colors',
                      isActive
                        ? 'bg-brand-50 text-brand-700'
                        : 'text-slate-600 hover:text-slate-900 hover:bg-slate-100',
                    )
                  }
                >
                  {item.label}
                </NavLink>
              ))}
            </nav>

            <div className="flex items-center gap-2">
              {currentUser && (
                <div className="hidden sm:flex items-center gap-2 rounded-full bg-slate-50 border border-slate-200 pl-2 pr-3 py-1">
                  <div className="h-6 w-6 rounded-full bg-brand-600 text-white grid place-items-center text-xs font-semibold">
                    {currentUser.name
                      .split(' ')
                      .map((p) => p[0])
                      .slice(0, 2)
                      .join('')}
                  </div>
                  <div className="leading-tight">
                    <p className="text-xs font-semibold text-slate-900">
                      {currentUser.name}
                    </p>
                    <p className="text-[10px] text-slate-500">
                      {currentUser.aptNumber} ·{' '}
                      {formatM3(currentUser.entitlementRemainingM3)} left
                    </p>
                  </div>
                </div>
              )}
              <button
                onClick={handleLogout}
                className="p-2 rounded-md text-slate-500 hover:bg-slate-100 hover:text-slate-700"
                aria-label="Log out"
                title="Log out"
              >
                <LogOut className="h-4 w-4" />
              </button>
            </div>
          </div>

          <nav className="md:hidden flex items-center gap-1 pb-2 -mx-1 overflow-x-auto">
            {navItems.map((item) => (
              <NavLink
                key={item.to}
                to={item.to}
                end={item.end}
                className={({ isActive }) =>
                  cn(
                    'px-3 py-1.5 rounded-md text-sm font-medium whitespace-nowrap',
                    isActive
                      ? 'bg-brand-50 text-brand-700'
                      : 'text-slate-600',
                  )
                }
              >
                {item.label}
              </NavLink>
            ))}
          </nav>
        </div>
      </header>

      <main className="flex-1 max-w-6xl w-full mx-auto px-4 sm:px-6 lg:px-8 py-6 sm:py-8">
        <Outlet />
      </main>

      <footer className="border-t border-slate-200 py-6 text-center text-xs text-slate-400">
        <div className="max-w-6xl mx-auto px-4">
          <p>
            <User className="inline h-3 w-3 -mt-0.5 mr-1" />
            A prototype for {building.name} — {building.address}
          </p>
          <p className="mt-1">
            Built for Melbourne residents to share hard-waste collection
            entitlements and keep usable furniture out of landfill.
          </p>
          <p className="mt-2 space-x-3">
            <Link
              to="/qr-poster"
              className="text-slate-400 hover:text-brand-700 underline-offset-2 hover:underline"
            >
              View lobby QR poster
            </Link>
            <span className="text-slate-300">·</span>
            <Link
              to="/"
              className="text-slate-400 hover:text-brand-700 underline-offset-2 hover:underline"
            >
              Demo navigator
            </Link>
          </p>
        </div>
      </footer>
    </div>
  )
}
