import { createContext, useContext, useState, type ReactNode } from 'react'
import { RESIDENTS } from '../lib/seed'
import {
  pbAuth,
  setPbSession,
  clearPbSession,
  getPbUser,
  type PbUser,
} from '../lib/pb'

type AuthState = {
  currentUserId: string | null
  pbUser: PbUser | null
  login: (email: string, password: string) => Promise<boolean>
  loginAs: (residentId: string) => void
  logout: () => void
}

const AuthContext = createContext<AuthState | null>(null)

export function AuthProvider({ children }: { children: ReactNode }) {
  const [currentUserId, setCurrentUserId] = useState<string | null>(null)
  const [pbUser, setPbUser] = useState<PbUser | null>(getPbUser())

  const login = async (email: string, password: string): Promise<boolean> => {
    // Try PocketBase first.
    try {
      const { token, record } = await pbAuth(email, password)
      setPbSession(token, record)
      setPbUser(record)
      // Mirror into seed-data userId so the rest of the UI keeps working.
      const match = RESIDENTS.find(
        (r) => r.email.toLowerCase() === email.trim().toLowerCase(),
      )
      setCurrentUserId(match?.id ?? record.id)
      return true
    } catch {
      // PocketBase unavailable or wrong credentials — fall back to seed data.
    }

    // Seed fallback (password not validated in demo mode).
    const match = RESIDENTS.find(
      (r) => r.email.toLowerCase() === email.trim().toLowerCase(),
    )
    if (match) {
      setCurrentUserId(match.id)
      return true
    }
    return false
  }

  const loginAs = (residentId: string) => {
    clearPbSession()
    setPbUser(null)
    setCurrentUserId(residentId)
  }

  const logout = () => {
    clearPbSession()
    setPbUser(null)
    setCurrentUserId(null)
  }

  return (
    <AuthContext.Provider value={{ currentUserId, pbUser, login, loginAs, logout }}>
      {children}
    </AuthContext.Provider>
  )
}

export function useAuth() {
  const ctx = useContext(AuthContext)
  if (!ctx) throw new Error('useAuth must be used within AuthProvider')
  return ctx
}
