import { Suspense } from 'react'
import NewExerciseContent from './NewExerciseContent'

export const dynamic = 'force-dynamic'

export default function NewExercisePage() {
  return (
    <Suspense fallback={
      <div className="min-h-screen bg-bg-primary flex items-center justify-center">
        <div className="text-center">
          <h1 className="text-2xl font-bold text-text-primary mb-2">
            Loading Exercise Creator...
          </h1>
          <p className="text-text-secondary">
            Please wait while we load the exercise creation interface.
          </p>
        </div>
      </div>
    }>
      <NewExerciseContent />
    </Suspense>
  )
}