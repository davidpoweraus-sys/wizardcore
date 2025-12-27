import LearningEnvironment from '@/components/learning/LearningEnvironment'

export default function LearningPage() {
  // For now, use a sample exercise ID from the seed data
  // In a real implementation, this would come from URL params or user progress
  const sampleExerciseId = 'dddddddd-dddd-dddd-dddd-dddddddddddd' // Hello World exercise
  
  return (
    <div className="h-screen">
      <LearningEnvironment exerciseId={sampleExerciseId} />
    </div>
  )
}