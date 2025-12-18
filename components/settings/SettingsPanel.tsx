'use client'

import { useState } from 'react'
import { Bell, Eye, Globe, Lock, Palette, Download, Printer, Share2 } from 'lucide-react'

export default function SettingsPanel() {
  const [settings, setSettings] = useState({
    theme: 'dark',
    language: 'en',
    emailNotifications: true,
    pushNotifications: false,
    publicProfile: true,
    showProgress: true,
    autoSave: true,
    soundEffects: true,
    twoFactor: false,
  })

  const handleToggle = (key: string) => {
    setSettings({ ...settings, [key]: !settings[key as keyof typeof settings] })
  }

  const handleSelect = (key: string, value: string) => {
    setSettings({ ...settings, [key]: value })
  }

  const certificates = [
    { id: 1, title: 'Python Fundamentals', date: '2025-11-15', url: '#', verified: true },
    { id: 2, title: 'C & Assembly Basics', date: '2025-12-01', url: '#', verified: true },
    { id: 3, title: 'Web Security Essentials', date: '2025-12-10', url: '#', verified: false },
  ]

  return (
    <div className="grid grid-cols-1 lg:grid-cols-3 gap-6">
      {/* Left Column: Preferences */}
      <div className="lg:col-span-2 space-y-6">
        {/* Appearance */}
        <div className="bg-bg-elevated border border-border-default rounded-2xl p-6">
          <h2 className="text-xl font-bold text-text-primary mb-6 flex items-center gap-2">
            <Palette className="w-5 h-5" />
            Appearance
          </h2>

          <div className="space-y-6">
            <div>
              <label className="block text-sm font-medium text-text-secondary mb-3">
                Theme
              </label>
              <div className="flex gap-3">
                {[
                  { value: 'light', label: 'Light' },
                  { value: 'dark', label: 'Dark' },
                  { value: 'auto', label: 'Auto' },
                ].map((theme) => (
                  <button
                    key={theme.value}
                    onClick={() => handleSelect('theme', theme.value)}
                    className={`px-4 py-3 rounded-lg border ${settings.theme === theme.value ? 'border-neon-cyan bg-neon-cyan/10 text-neon-cyan' : 'border-border-subtle bg-bg-tertiary text-text-secondary'}`}
                  >
                    {theme.label}
                  </button>
                ))}
              </div>
            </div>

            <div>
              <label className="block text-sm font-medium text-text-secondary mb-3">
                Language
              </label>
              <select
                value={settings.language}
                onChange={(e) => handleSelect('language', e.target.value)}
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

        {/* Notifications */}
        <div className="bg-bg-elevated border border-border-default rounded-2xl p-6">
          <h2 className="text-xl font-bold text-text-primary mb-6 flex items-center gap-2">
            <Bell className="w-5 h-5" />
            Notifications
          </h2>

          <div className="space-y-4">
            {[
              { key: 'emailNotifications', label: 'Email Notifications', description: 'Receive updates about your progress and new content' },
              { key: 'pushNotifications', label: 'Push Notifications', description: 'Get browser notifications for duels and reminders' },
              { key: 'soundEffects', label: 'Sound Effects', description: 'Play sounds for completions and achievements' },
            ].map((item) => (
              <div key={item.key} className="flex items-center justify-between">
                <div>
                  <div className="font-medium text-text-primary">{item.label}</div>
                  <div className="text-sm text-text-secondary">{item.description}</div>
                </div>
                <button
                  onClick={() => handleToggle(item.key)}
                  className={`w-12 h-6 rounded-full transition ${settings[item.key as keyof typeof settings] ? 'bg-neon-cyan' : 'bg-bg-tertiary'}`}
                >
                  <div
                    className={`w-5 h-5 rounded-full bg-white transform transition ${settings[item.key as keyof typeof settings] ? 'translate-x-7' : 'translate-x-1'}`}
                  />
                </button>
              </div>
            ))}
          </div>
        </div>

        {/* Privacy */}
        <div className="bg-bg-elevated border border-border-default rounded-2xl p-6">
          <h2 className="text-xl font-bold text-text-primary mb-6 flex items-center gap-2">
            <Eye className="w-5 h-5" />
            Privacy & Visibility
          </h2>

          <div className="space-y-4">
            {[
              { key: 'publicProfile', label: 'Public Profile', description: 'Allow other learners to view your profile' },
              { key: 'showProgress', label: 'Show Progress', description: 'Display your learning progress on leaderboards' },
              { key: 'autoSave', label: 'Auto‑Save Exercises', description: 'Automatically save your code as you type' },
            ].map((item) => (
              <div key={item.key} className="flex items-center justify-between">
                <div>
                  <div className="font-medium text-text-primary">{item.label}</div>
                  <div className="text-sm text-text-secondary">{item.description}</div>
                </div>
                <button
                  onClick={() => handleToggle(item.key)}
                  className={`w-12 h-6 rounded-full transition ${settings[item.key as keyof typeof settings] ? 'bg-neon-cyan' : 'bg-bg-tertiary'}`}
                >
                  <div
                    className={`w-5 h-5 rounded-full bg-white transform transition ${settings[item.key as keyof typeof settings] ? 'translate-x-7' : 'translate-x-1'}`}
                  />
                </button>
              </div>
            ))}
          </div>
        </div>
      </div>

      {/* Right Column: Security & Certificates */}
      <div className="space-y-6">
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
                  <div className="text-sm text-text-secondary">Update your password</div>
                </div>
              </div>
              <div className="text-neon-cyan">→</div>
            </button>

            <div className="flex items-center justify-between p-4 border border-border-subtle rounded-lg">
              <div className="flex items-center">
                <Globe className="w-5 h-5 text-text-secondary mr-3" />
                <div>
                  <div className="font-medium text-text-primary">Two‑Factor Authentication</div>
                  <div className="text-sm text-text-secondary">Add an extra layer of security</div>
                </div>
              </div>
              <button
                onClick={() => handleToggle('twoFactor')}
                className={`w-12 h-6 rounded-full transition ${settings.twoFactor ? 'bg-neon-cyan' : 'bg-bg-tertiary'}`}
              >
                <div
                  className={`w-5 h-5 rounded-full bg-white transform transition ${settings.twoFactor ? 'translate-x-7' : 'translate-x-1'}`}
                />
              </button>
            </div>

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

        {/* Certificates */}
        <div className="bg-bg-elevated border border-border-default rounded-2xl p-6">
          <h2 className="text-xl font-bold text-text-primary mb-6">Your Certificates</h2>

          <div className="space-y-4">
            {certificates.map((cert) => (
              <div key={cert.id} className="border border-border-subtle rounded-xl p-4">
                <div className="flex items-center justify-between mb-3">
                  <div>
                    <div className="font-bold text-text-primary">{cert.title}</div>
                    <div className="text-sm text-text-secondary">Issued: {cert.date}</div>
                  </div>
                  {cert.verified ? (
                    <div className="px-3 py-1 rounded-full bg-green-900/30 text-green-300 text-sm">
                      Verified
                    </div>
                  ) : (
                    <div className="px-3 py-1 rounded-full bg-yellow-900/30 text-yellow-300 text-sm">
                      Pending
                    </div>
                  )}
                </div>
                <div className="flex gap-2">
                  <button className="flex-1 py-2 px-3 bg-bg-tertiary text-text-secondary rounded-lg text-sm hover:bg-bg-secondary transition flex items-center justify-center gap-2">
                    <Download className="w-4 h-4" />
                    Download
                  </button>
                  <button className="flex-1 py-2 px-3 bg-bg-tertiary text-text-secondary rounded-lg text-sm hover:bg-bg-secondary transition flex items-center justify-center gap-2">
                    <Printer className="w-4 h-4" />
                    Print
                  </button>
                  <button className="flex-1 py-2 px-3 bg-bg-tertiary text-text-secondary rounded-lg text-sm hover:bg-bg-secondary transition flex items-center justify-center gap-2">
                    <Share2 className="w-4 h-4" />
                    Share
                  </button>
                </div>
              </div>
            ))}
          </div>

          <button className="w-full mt-6 py-3 border border-neon-cyan text-neon-cyan rounded-lg hover:bg-neon-cyan/10 transition">
            View All Certificates
          </button>
        </div>

        {/* Danger Zone */}
        <div className="bg-bg-elevated border border-red-700/30 rounded-2xl p-6">
          <h2 className="text-xl font-bold text-red-400 mb-4">Danger Zone</h2>
          <p className="text-sm text-text-secondary mb-4">
            These actions are irreversible. Please proceed with caution.
          </p>
          <div className="space-y-3">
            <button className="w-full py-3 border border-red-700 text-red-400 rounded-lg hover:bg-red-900/20 transition">
              Delete Account
            </button>
            <button className="w-full py-3 border border-red-700 text-red-400 rounded-lg hover:bg-red-900/20 transition">
              Reset All Progress
            </button>
            <button className="w-full py-3 border border-red-700 text-red-400 rounded-lg hover:bg-red-900/20 transition">
              Export All Data
            </button>
          </div>
        </div>
      </div>
    </div>
  )
}