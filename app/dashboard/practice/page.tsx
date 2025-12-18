import PracticeArena from '@/components/practice/PracticeArena'

export default function PracticePage() {
  return (
    <div className="space-y-6">
      <div>
        <h1 className="text-3xl font-bold text-text-primary">Practice Arena</h1>
        <p className="text-text-secondary mt-2">
          Sharpen your skills with random challenges, timed exercises, and competitive duels.
        </p>
      </div>

      <PracticeArena />
    </div>
  )
}