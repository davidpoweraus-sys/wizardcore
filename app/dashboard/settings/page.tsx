import SettingsPanel from '@/components/settings/SettingsPanel'

export default function SettingsPage() {
  return (
    <div className="space-y-6">
      <div>
        <h1 className="text-3xl font-bold text-text-primary">Settings</h1>
        <p className="text-text-secondary mt-2">
          Manage your account, preferences, and platform settings.
        </p>
      </div>

      <SettingsPanel />
    </div>
  )
}