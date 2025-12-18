import PathwayProgressList from '@/components/dashboard/PathwayProgressList'
import ProgressChart from '@/components/dashboard/ProgressChart'
import RecentActivity from '@/components/dashboard/RecentActivity'
import QuickStats from '@/components/dashboard/QuickStats'

export default function DashboardPage() {
  return (
    <div className="space-y-6">
      <div>
        <h1 className="text-3xl font-bold text-text-primary">Welcome back, Coder!</h1>
        <p className="text-text-secondary mt-2">
          Continue your journey to master programming. You're making great progress!
        </p>
      </div>

      <QuickStats />

      <div className="grid grid-cols-1 lg:grid-cols-3 gap-6">
        <div className="lg:col-span-2">
          <div className="bg-bg-elevated border border-border-default rounded-2xl p-6">
            <h2 className="text-xl font-semibold text-text-primary mb-4">Your Learning Pathways</h2>
            <PathwayProgressList />
          </div>
        </div>

        <div className="lg:col-span-1">
          <ProgressChart />
        </div>
      </div>

      <div className="grid grid-cols-1 lg:grid-cols-3 gap-6">
        <div className="lg:col-span-2">
          <RecentActivity />
        </div>
        <div className="lg:col-span-1">
          <div className="bg-bg-elevated border border-border-default rounded-2xl p-6">
            <h2 className="text-xl font-semibold text-text-primary mb-4">Upcoming Deadlines</h2>
            <div className="space-y-4">
              <div className="p-3 border border-border-subtle rounded-lg">
                <p className="font-medium text-text-primary">Python Project: Data Analysis</p>
                <p className="text-sm text-text-secondary">Due in 3 days</p>
              </div>
              <div className="p-3 border border-border-subtle rounded-lg">
                <p className="font-medium text-text-primary">C Memory Management Quiz</p>
                <p className="text-sm text-text-secondary">Due tomorrow</p>
              </div>
              <div className="p-3 border border-border-subtle rounded-lg">
                <p className="font-medium text-text-primary">Assembly Lab 2</p>
                <p className="text-sm text-text-secondary">Due next week</p>
              </div>
            </div>
          </div>
        </div>
      </div>
    </div>
  )
}