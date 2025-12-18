import PathwaySelection from '@/components/pathways/PathwaySelection'

export default function PathwaysPage() {
  return (
    <div className="space-y-6">
      <div>
        <h1 className="text-3xl font-bold text-text-primary">Choose Your Learning Pathway</h1>
        <p className="text-text-secondary mt-2">
          Select a programming language to start your journey. Each pathway is designed to take you from beginner to expert.
        </p>
      </div>

      <PathwaySelection />
    </div>
  )
}