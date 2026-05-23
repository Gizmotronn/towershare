import { useEffect } from 'react'
import { Outlet } from 'react-router-dom'
import { useAuth } from '../context/AuthContext'
import { useAppData } from '../context/AppDataContext'

const DEFAULT_PERSONA_ID = 'r_amelia'

export function ProtectedRoute() {
  const { currentUserId, loginAs } = useAuth()
  const { currentUser } = useAppData()

  useEffect(() => {
    if (!currentUserId) loginAs(DEFAULT_PERSONA_ID)
  }, [currentUserId, loginAs])

  if (!currentUser) return null
  return <Outlet />
}
