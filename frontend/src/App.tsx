import { Route, Routes, Navigate } from 'react-router-dom'
import { ProtectedRoute } from './components/ProtectedRoute'
import { Layout } from './components/Layout'
import { Welcome } from './pages/Welcome'
import { QrPoster } from './pages/QrPoster'
import { DemoIndex } from './pages/DemoIndex'
import { Login } from './pages/Login'
import { Dashboard } from './pages/Dashboard'
import { Marketplace } from './pages/Marketplace'
import { NewListing } from './pages/NewListing'
import { ListingDetail } from './pages/ListingDetail'
import { Entitlements } from './pages/Entitlements'
import { CollectionDay } from './pages/CollectionDay'
import { Profile } from './pages/Profile'

export default function App() {
  return (
    <Routes>
      <Route path="/" element={<DemoIndex />} />
      <Route path="/welcome" element={<Welcome />} />
      <Route path="/qr-poster" element={<QrPoster />} />
      <Route path="/login" element={<Login />} />
      <Route element={<ProtectedRoute />}>
        <Route element={<Layout />}>
          <Route path="/dashboard" element={<Dashboard />} />
          <Route path="/marketplace" element={<Marketplace />} />
          <Route path="/listings/new" element={<NewListing />} />
          <Route path="/listings/:id" element={<ListingDetail />} />
          <Route path="/entitlements" element={<Entitlements />} />
          <Route path="/collection-day" element={<CollectionDay />} />
          <Route path="/profile" element={<Profile />} />
        </Route>
      </Route>
      <Route path="*" element={<Navigate to="/" replace />} />
    </Routes>
  )
}
