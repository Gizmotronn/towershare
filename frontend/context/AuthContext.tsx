import {
  createContext,
  useContext,
  useState,
  useCallback,
  type ReactNode,
} from 'react'
import type { RecordModel } from 'pocketbase'
import { pb } from '../lib/pb'
import { RESIDENTS } from '../lib/seed'

type AuthState = {
  currentUserId: string | null
  pbUser: RecordModel | null
  isPbAuth: boolean
  login: (email: string, password: string) => Promise<boolean>
  loginAs: (residentId: string) => void
  logout: () => void
}

const AuthContext = createContext<AuthState | null>(null)

export function AuthProvider({ children }: { children: ReactNode }) {
  const [currentUserId, setCurrentUserId] = useState<string | null>(null)
  const [pbUser, setPbUser] = useState<RecordModel | null>(null)

  const login = useCallback(async (email: string, password: string): Promise<boolean> => {
    // Try real PocketBase auth first (only when password is provided)
    if (password) {
      try {
        const auth = await pb.collection('users').authWithPassword(email, password)
        setPbUser(auth.record)
        // Map PB user to seed resident by email for entitlement compat; fall back to PB id
        const seedMatch = RESIDENTS.find(
          (r) => r.email.toLowerCase() === email.trim().toLowerCase(),
        )
        setCurrentUserId(seedMatch?.id ?? auth.record.id)
        return true
      } catch {
        // Fall through to seed lookup
      }
    }
    // Demo fallback: match by email in seed data
    const match = RESIDENTS.find(
      (r) => r.email.toLowerCase() === email.trim().toLowerCase(),
    )
    if (match) {
      setCurrentUserId(match.id)
      return true
    }
    return false
  }, [])

  const loginAs = useCallback((residentId: string) => {
    setPbUser(null)
    setCurrentUserId(residentId)
  }, [])

  const logout = useCallback(() => {
    pb.authStore.clear()
    setPbUser(null)
    setCurrentUserId(null)
  }, [])

  return (
    <AuthContext.Provider
      value={{
        currentUserId,
        pbUser,
        isPbAuth: pbUser !== null,
        login,
        loginAs,
        logout,
      }}
    >
      {children}
    </AuthContext.Provider>
  )
}

export function useAuth() {
  const ctx = useContext(AuthContext)
  if (!ctx) throw new Error('useAuth must be used within AuthProvider')
  return ctx
}
