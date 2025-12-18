'use client'

import { useState } from 'react'
import Link from 'next/link'
import { usePathname } from 'next/navigation'
import {
  Home,
  BookOpen,
  Trophy,
  Users,
  Code2,
  User,
  Settings,
  Award,
  BarChart3,
  LogOut,
} from 'lucide-react'
import { createClient } from '@/lib/supabase/client'

const navItems = [
  { name: 'Dashboard', href: '/dashboard', icon: Home },
  { name: 'Pathways', href: '/dashboard/pathways', icon: BookOpen },
  { name: 'Learning', href: '/dashboard/learning', icon: Code2 },
  { name: 'Progress', href: '/dashboard/progress', icon: BarChart3 },
  { name: 'Practice Arena', href: '/dashboard/practice', icon: Trophy },
  { name: 'Achievements', href: '/dashboard/achievements', icon: Award },
  { name: 'Leaderboard', href: '/dashboard/leaderboard', icon: Users },
  { name: 'Profile', href: '/dashboard/profile', icon: User },
  { name: 'Settings', href: '/dashboard/settings', icon: Settings },
]

export default function Sidebar() {
  const pathname = usePathname()
  const [collapsed, setCollapsed] = useState(false)
  const supabase = createClient()

  const handleSignOut = async () => {
    await supabase.auth.signOut()
    window.location.href = '/login'
  }

  return (
    <aside className={`h-screen bg-bg-elevated border-r border-border-default flex flex-col transition-all duration-300 ${collapsed ? 'w-20' : 'w-64'}`}>
      <div className="p-6 border-b border-border-subtle">
        <div className="flex items-center justify-between">
          {!collapsed && (
            <div className="flex items-center space-x-3">
              <div className="w-8 h-8 bg-gradient-to-r from-neon-pink to-neon-cyan rounded-lg" />
              <h1 className="text-xl font-bold bg-gradient-to-r from-neon-pink to-neon-cyan bg-clip-text text-transparent">
                WizardCore
              </h1>
            </div>
          )}
          {collapsed && (
            <div className="w-8 h-8 bg-gradient-to-r from-neon-pink to-neon-cyan rounded-lg mx-auto" />
          )}
          <button
            onClick={() => setCollapsed(!collapsed)}
            className="p-2 rounded-lg hover:bg-bg-tertiary text-text-secondary"
          >
            {collapsed ? '→' : '←'}
          </button>
        </div>
      </div>

      <nav className="flex-1 p-4 space-y-2">
        {navItems.map((item) => {
          const Icon = item.icon
          const isActive = pathname === item.href || pathname.startsWith(item.href + '/')
          return (
            <Link
              key={item.name}
              href={item.href}
              className={`flex items-center ${collapsed ? 'justify-center' : 'justify-start'} p-3 rounded-lg transition ${isActive
                  ? 'bg-gradient-to-r from-neon-pink/20 to-neon-purple/20 text-neon-cyan border-l-4 border-neon-cyan'
                  : 'hover:bg-bg-tertiary text-text-secondary hover:text-text-primary'
                }`}
            >
              <Icon className="w-5 h-5" />
              {!collapsed && <span className="ml-3 font-medium">{item.name}</span>}
            </Link>
          )
        })}
      </nav>

      <div className="p-4 border-t border-border-subtle">
        <button
          onClick={handleSignOut}
          className={`flex items-center ${collapsed ? 'justify-center' : 'justify-start'} w-full p-3 rounded-lg text-red-400 hover:bg-red-900/20 transition`}
        >
          <LogOut className="w-5 h-5" />
          {!collapsed && <span className="ml-3 font-medium">Sign Out</span>}
        </button>
      </div>
    </aside>
  )
}