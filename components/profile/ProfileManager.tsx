'use client'

import { useState } from 'react'
import { User, Mail, Globe, Lock, Bell, Eye, Palette, Save } from 'lucide-react'

export default function ProfileManager() {
  const [profile, setProfile] = useState({
    name: 'Alex Chen',
    email: 'alex@example.com',
    bio: 'Security researcher and aspiring exploit developer. Currently learning Python for offensive security.',
    location: 'San Francisco, CA',
    website: 'https://alexchen.dev',
    github: 'alexchen',
    twitter: '@alexchen',
  })

  const [preferences, setPreferences] = useState({
    emailNotifications: true,
    pushNotifications: false,
    darkMode: true,
    publicProfile: true,
    showProgress: true,
    language: 'en',
  })

  const handleProfileChange = (field: string, value: string) => {
    setProfile({ ...profile, [field]: value })
  }

  const handlePreferenceChange = (field: string, value: boolean | string) => {
    setPreferences({ ...preferences, [field]: value })
  }

  const handleSave = () => {
    console.log('Saving profile:', profile)
    console.log('Saving preferences:', preferences)
    // In a real app, this would call Supabase
  }

  return (
    <div className="grid grid-cols-1 lg:grid-cols-3 gap-6">
      {/* Left Column: Profile Info */}
      <div className="lg:col-span-2 space-y-6">
        {/* Personal Information */}
        <div className="bg-bg-elevated border border-border-default rounded-2xl p-6">
          <h2 className="text-xl font-bold text-text-primary mb-6 flex items-center gap-2">
            <User className="w-5 h-5" />
            Personal Information
          </h2>

          <div className="space-y-6">
            <div className="grid grid-cols-1 md:grid-cols-2 gap-6">
              <div>
                <label className="block text-sm font-medium text-text-secondary mb-2">
                  Full Name
                </label>
                <input
                  type="text"
                  value={profile.name}
                  onChange={(e) => handleProfileChange('name', e.target.value)}
                  className="w-full px-4 py-3 bg-bg-tertiary border border-border-subtle rounded-lg text-text-primary placeholder-text-muted focus:outline-none focus:ring-2 focus:ring-neon-lavender"
                />
              </div>

              <div>
                <label className="block text-sm font-medium text-text-secondary mb-2">
                  Email Address
                </label>
                <div className="flex items-center">
                  <input
                    type="email"
                    value={profile.email}
                    onChange={(e) => handleProfileChange('email', e.target.value)}
                    className="flex-1 px-4 py-3 bg-bg-tertiary border border-border-subtle rounded-lg text-text-primary placeholder-text-muted focus:outline-none focus:ring-2 focus:ring-neon-lavender"
                  />
                  <Mail className="w-5 h-5 text-text-muted ml-3" />
                </div>
              </div>
            </div>

            <div>
              <label className="block text-sm font-medium text-text-secondary mb-2">
                Bio
              </label>
              <textarea
                value={profile.bio}
                onChange={(e) => handleProfileChange('bio', e.target.value)}
                rows={3}
                className="w-full px-4 py-3 bg-bg-tertiary border border-border-subtle rounded-lg text-text-primary placeholder-text-muted focus:outline-none focus:ring-2 focus:ring-neon-lavender"
              />
              <p className="text-xs text-text-muted mt-2">
                Tell us about yourself and your security interests.
              </p>
            </div>

            <div className="grid grid-cols-1 md:grid-cols-3 gap-6">
              <div>
                <label className="block text-sm font-medium text-text-secondary mb-2">
                  Location
                </label>
                <input
                  type="text"
                  value={profile.location}
                  onChange={(e) => handleProfileChange('location', e.target.value)}
                  className="w-full px-4 py-3 bg-bg-tertiary border border-border-subtle rounded-lg text-text-primary placeholder-text-muted focus:outline-none focus:ring-2 focus:ring-neon-lavender"
                />
              </div>

              <div>
                <label className="block text-sm font-medium text-text-secondary mb-2">
                  Website
                </label>
                <div className="flex items-center">
                  <input
                    type="url"
                    value={profile.website}
                    onChange={(e) => handleProfileChange('website', e.target.value)}
                    className="flex-1 px-4 py-3 bg-bg-tertiary border border-border-subtle rounded-lg text-text-primary placeholder-text-muted focus:outline-none focus:ring-2 focus:ring-neon-lavender"
                  />
                  <Globe className="w-5 h-5 text-text-muted ml-3" />
                </div>
              </div>

              <div>
                <label className="block text-sm font-medium text-text-secondary mb-2">
                  GitHub Username
                </label>
                <input
                  type="text"
                  value={profile.github}
                  onChange={(e) => handleProfileChange('github', e.target.value)}
                  className="w-full px-4 py-3 bg-bg-tertiary border border-border-subtle rounded-lg text-text-primary placeholder-text-muted focus:outline-none focus:ring-2 focus:ring-neon-lavender"
                />
              </div>
            </div>
          </div>
        </div>

        {/* Security */}
        <div className="bg-bg-elevated border border-border-default rounded-2xl p-6">
          <h2 className="text-xl font-bold text-text-primary mb-6 flex items-center gap-2">
            <Lock className="w-5 h-5" />
            Security
          </h2>

          <div className="space-y-4">
            <button className="w-full flex items-center justify-between p-4 border border-border-subtle rounded-lg hover:bg-bg-tertiary transition">
              <div className="flex items-center">
                <Lock className="w-5 h-5 text-text-secondary mr-3" />
                <div>
                  <div className="font-medium text-text-primary">Change Password</div>
                  <div className="text-sm text-text-secondary">Update your password regularly</div>
                </div>
              </div>
              <div className="text-neon-cyan">→</div>
            </button>

            <button className="w-full flex items-center justify-between p-4 border border-border-subtle rounded-lg hover:bg-bg-tertiary transition">
              <div className="flex items-center">
                <Eye className="w-5 h-5 text-text-secondary mr-3" />
                <div>
                  <div className="font-medium text-text-primary">Two-Factor Authentication</div>
                  <div className="text-sm text-text-secondary">Add an extra layer of security</div>
                </div>
              </div>
              <div className="text-text-muted">Not enabled</div>
            </button>

            <button className="w-full flex items-center justify-between p-4 border border-border-subtle rounded-lg hover:bg-bg-tertiary transition">
              <div className="flex items-center">
                <Globe className="w-5 h-5 text-text-secondary mr-3" />
                <div>
                  <div className="font-medium text-text-primary">Connected Accounts</div>
                  <div className="text-sm text-text-secondary">GitHub, Google, etc.</div>
                </div>
              </div>
              <div className="text-neon-cyan">Manage</div>
            </button>
          </div>
        </div>
      </div>

      {/* Right Column: Preferences & Actions */}
      <div className="space-y-6">
        {/* Preferences */}
        <div className="bg-bg-elevated border border-border-default rounded-2xl p-6">
          <h2 className="text-xl font-bold text-text-primary mb-6 flex items-center gap-2">
            <Bell className="w-5 h-5" />
            Preferences
          </h2>

          <div className="space-y-4">
            {[
              { key: 'emailNotifications', label: 'Email Notifications', description: 'Receive updates about your progress' },
              { key: 'pushNotifications', label: 'Push Notifications', description: 'Get browser notifications' },
              { key: 'publicProfile', label: 'Public Profile', description: 'Allow others to see your profile' },
              { key: 'showProgress', label: 'Show Progress', description: 'Display your progress on leaderboard' },
            ].map((item) => (
              <div key={item.key} className="flex items-center justify-between">
                <div>
                  <div className="font-medium text-text-primary">{item.label}</div>
                  <div className="text-sm text-text-secondary">{item.description}</div>
                </div>
                <button
                  onClick={() => handlePreferenceChange(item.key, !preferences[item.key as keyof typeof preferences])}
                  className={`w-12 h-6 rounded-full transition ${preferences[item.key as keyof typeof preferences] ? 'bg-neon-cyan' : 'bg-bg-tertiary'}`}
                >
                  <div
                    className={`w-5 h-5 rounded-full bg-white transform transition ${preferences[item.key as keyof typeof preferences] ? 'translate-x-7' : 'translate-x-1'}`}
                  />
                </button>
              </div>
            ))}

            <div>
              <div className="font-medium text-text-primary mb-2">Theme</div>
              <div className="flex gap-2">
                {[
                  { value: 'light', label: 'Light' },
                  { value: 'dark', label: 'Dark' },
                  { value: 'auto', label: 'Auto' },
                ].map((theme) => (
                  <button
                    key={theme.value}
                    onClick={() => handlePreferenceChange('darkMode', theme.value === 'dark')}
                    className={`px-3 py-2 rounded-lg text-sm ${preferences.darkMode === (theme.value === 'dark') ? 'bg-neon-cyan text-white' : 'bg-bg-tertiary text-text-secondary'}`}
                  >
                    {theme.label}
                  </button>
                ))}
              </div>
            </div>

            <div>
              <div className="font-medium text-text-primary mb-2">Language</div>
              <select
                value={preferences.language}
                onChange={(e) => handlePreferenceChange('language', e.target.value)}
                className="w-full px-4 py-3 bg-bg-tertiary border border-border-subtle rounded-lg text-text-primary focus:outline-none focus:ring-2 focus:ring-neon-lavender"
              >
                <option value="en">English</option>
                <option value="es">Español</option>
                <option value="fr">Français</option>
                <option value="de">Deutsch</option>
                <option value="ja">日本語</option>
              </select>
            </div>
          </div>
        </div>

        {/* Profile Preview */}
        <div className="bg-bg-elevated border border-border-default rounded-2xl p-6">
          <h2 className="text-xl font-bold text-text-primary mb-6 flex items-center gap-2">
            <Palette className="w-5 h-5" />
            Profile Preview
          </h2>

          <div className="border border-border-subtle rounded-xl p-4">
            <div className="flex items-center space-x-4">
              <div className="w-16 h-16 rounded-full bg-gradient-to-r from-neon-pink to-neon-cyan" />
              <div>
                <div className="font-bold text-text-primary">{profile.name}</div>
                <div className="text-sm text-text-secondary">{profile.email}</div>
                <div className="text-xs text-text-muted mt-1">{profile.location}</div>
              </div>
            </div>
            <div className="mt-4 text-sm text-text-secondary">
              {profile.bio}
            </div>
            <div className="mt-4 flex items-center text-sm text-text-muted">
              <Globe className="w-4 h-4 mr-2" />
              {profile.website}
            </div>
          </div>
        </div>

        {/* Actions */}
        <div className="bg-bg-elevated border border-border-default rounded-2xl p-6">
          <button
            onClick={handleSave}
            className="w-full py-3 px-4 bg-gradient-to-r from-neon-cyan to-neon-lavender text-white font-semibold rounded-lg hover:opacity-90 transition flex items-center justify-center gap-2"
          >
            <Save className="w-5 h-5" />
            Save Changes
          </button>
          <button className="w-full py-3 px-4 border border-border-subtle text-text-secondary rounded-lg hover:bg-bg-tertiary transition mt-3">
            Cancel
          </button>
          <button className="w-full py-3 px-4 border border-red-700 text-red-400 rounded-lg hover:bg-red-900/20 transition mt-3">
            Delete Account
          </button>
        </div>
      </div>
    </div>
  )
}