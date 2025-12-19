'use client'

import { useState, useEffect } from 'react'
import Image from 'next/image'
import { Bell, Search, User as UserIcon } from 'lucide-react'
import { createClient } from '@/lib/supabase/client'

export default function Header() {
  const [user, setUser] = useState<any>(null)
  const supabase = createClient()

  useEffect(() => {
    const getUser = async () => {
      const { data: { user } } = await supabase.auth.getUser()
      setUser(user)
    }
    getUser()
  }, [])

  return (
    <header className="sticky top-0 z-10 bg-bg-elevated/80 backdrop-blur-md border-b border-border-default p-4">
      <div className="flex items-center justify-between">
        <div className="flex items-center gap-4">
          <Image
            src="/wizard_logo.png"
            alt="WizardCore Logo"
            width={40}
            height={40}
            className="drop-shadow-[0_0_10px_rgba(138,43,226,0.4)]"
          />
          <h1 className="text-xl font-bold bg-gradient-to-r from-neon-cyan to-neon-lavender bg-clip-text text-transparent hidden sm:block">
            WizardCore
          </h1>
        </div>

        <div className="flex-1 max-w-xl mx-4">
          <div className="relative">
            <Search className="absolute left-3 top-1/2 transform -translate-y-1/2 w-5 h-5 text-text-muted" />
            <input
              type="search"
              placeholder="Search courses, modules, or help..."
              className="w-full pl-10 pr-4 py-2 bg-bg-tertiary border border-border-subtle rounded-lg text-text-primary placeholder-text-muted focus:outline-none focus:ring-2 focus:ring-neon-lavender"
            />
          </div>
        </div>

        <div className="flex items-center space-x-4">
          <button className="relative p-2 rounded-lg hover:bg-bg-tertiary text-text-secondary">
            <Bell className="w-5 h-5" />
            <span className="absolute top-1 right-1 w-2 h-2 bg-neon-pink rounded-full" />
          </button>

          <div className="h-8 w-px bg-border-subtle" />

          <div className="flex items-center space-x-3">
            <div className="w-10 h-10 rounded-full bg-gradient-to-r from-neon-cyan to-neon-lavender flex items-center justify-center">
              {user?.email?.charAt(0).toUpperCase() || <UserIcon className="w-5 h-5 text-white" />}
            </div>
            <div className="hidden md:block">
              <p className="text-sm font-medium text-text-primary">
                {user?.email?.split('@')[0] || 'Guest'}
              </p>
              <p className="text-xs text-text-secondary">
                {user?.email || 'Not signed in'}
              </p>
            </div>
          </div>
        </div>
      </div>
    </header>
  )
}