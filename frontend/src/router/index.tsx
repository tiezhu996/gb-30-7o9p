import { createBrowserRouter, Navigate } from 'react-router-dom'
import App from '@/App'
import Home from '@/pages/Home'
import PetList from '@/pages/PetList'
import PetDetail from '@/pages/PetDetail'
import OrgList from '@/pages/OrgList'
import OrgDetail from '@/pages/OrgDetail'
import Apply from '@/pages/Apply'
import Applications from '@/pages/Applications'
import Reviews from '@/pages/Reviews'
import Community from '@/pages/Community'
import PostDetail from '@/pages/PostDetail'
import Donate from '@/pages/Donate'
import Profile from '@/pages/Profile'
import Login from '@/pages/Login'
import RoleGuard from '@/components/common/RoleGuard'

export const router = createBrowserRouter([
  {
    path: '/',
    element: <App />,
    children: [
      { index: true, element: <Home /> },
      { path: 'pets', element: <PetList /> },
      { path: 'pets/:id', element: <PetDetail /> },
      { path: 'orgs', element: <OrgList /> },
      { path: 'orgs/:id', element: <OrgDetail /> },
      { path: 'apply/:petId', element: <RoleGuard roles={['user']}><Apply /></RoleGuard> },
      { path: 'applications', element: <RoleGuard roles={['user', 'org', 'admin']}><Applications /></RoleGuard> },
      { path: 'reviews', element: <RoleGuard roles={['user', 'org']}><Reviews /></RoleGuard> },
      { path: 'community', element: <Community /> },
      { path: 'community/:id', element: <PostDetail /> },
      { path: 'donate', element: <RoleGuard roles={['user']}><Donate /></RoleGuard> },
      { path: 'profile', element: <RoleGuard roles={['user', 'org', 'admin']}><Profile /></RoleGuard> },
      { path: 'login', element: <Login /> },
      { path: '*', element: <Navigate to="/" replace /> },
    ],
  },
])
