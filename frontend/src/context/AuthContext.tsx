import { createContext, useContext, useState, type ReactNode } from 'react'
import { RESIDENTS } from '../lib/seed'

type AuthState = {
  currentUserId: string | null
  login: (email: string) => boolean
  loginAs: (residentId: string) => void
  logout: () => void
}

const AuthContext = createContext<AuthState | null>(null)

export function AuthProvider({ children }: { children: ReactNode }) {
  const [currentUserId, setCurrentUserId] = useState<string | null>(null)

  const login = (email: string): boolean => {
    const match = RESIDENTS.find(
      (r) => r.email.toLowerCase() === email.trim().toLowerCase(),
    )
    if (match) {
      setCurrentUserId(match.id)
      return true
    }
    return false
  }

  const loginAs = (residentId: string) => setCurrentUserId(residentId)
  const logout = () => setCurrentUserId(null)

  return (
    <AuthContext.Provider value={{ currentUserId, login, loginAs, logout }}>
      {children}
    </AuthContext.Provider>
  )
}

export function useAuth() {
  const ctx = useContext(AuthContext)
  if (!ctx) throw new Error('useAuth must be used within AuthProvider')
  return ctx
}
