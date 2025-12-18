import LeaderboardTable from '@/components/leaderboard/LeaderboardTable'

export default function LeaderboardPage() {
  return (
    <div className="space-y-6">
      <div>
        <h1 className="text-3xl font-bold text-text-primary">Global Leaderboard</h1>
        <p className="text-text-secondary mt-2">
          Compete with learners worldwide and climb the ranks.
        </p>
      </div>

      <LeaderboardTable />
    </div>
  )
}