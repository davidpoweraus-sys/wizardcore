import AchievementsDisplay from '@/components/achievements/AchievementsDisplay'

export default function AchievementsPage() {
  return (
    <div className="space-y-6">
      <div>
        <h1 className="text-3xl font-bold text-text-primary">Achievements & Badges</h1>
        <p className="text-text-secondary mt-2">
          Celebrate your learning milestones and showcase your skills.
        </p>
      </div>

      <AchievementsDisplay />
    </div>
  )
}