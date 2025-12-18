import ProgressTracker from '@/components/progress/ProgressTracker'

export default function ProgressPage() {
  return (
    <div className="space-y-6">
      <div>
        <h1 className="text-3xl font-bold text-text-primary">Your Learning Progress</h1>
        <p className="text-text-secondary mt-2">
          Track your journey across all pathways, modules, and exercises.
        </p>
      </div>

      <ProgressTracker />
    </div>
  )
}