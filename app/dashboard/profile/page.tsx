import ProfileManager from '@/components/profile/ProfileManager'

export default function ProfilePage() {
  return (
    <div className="space-y-6">
      <div>
        <h1 className="text-3xl font-bold text-text-primary">Profile Management</h1>
        <p className="text-text-secondary mt-2">
          Update your personal information, preferences, and account settings.
        </p>
      </div>

      <ProfileManager />
    </div>
  )
}