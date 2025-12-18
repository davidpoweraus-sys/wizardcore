interface PathwayCardProps {
  title: string
  description: string
  progress: number
  icon: string
  color: string
}

export default function PathwayCard({
  title,
  description,
  progress,
  icon,
  color,
}: PathwayCardProps) {
  return (
    <div className="bg-bg-tertiary border border-border-subtle rounded-xl p-5 hover:border-neon-lavender transition">
      <div className="flex items-start justify-between">
        <div>
          <div className="flex items-center space-x-3">
            <div className={`w-12 h-12 rounded-lg bg-gradient-to-r ${color} flex items-center justify-center text-2xl`}>
              {icon}
            </div>
            <div>
              <h3 className="font-bold text-text-primary">{title}</h3>
              <p className="text-sm text-text-secondary">{description}</p>
            </div>
          </div>
        </div>
        <div className="text-right">
          <div className="text-2xl font-bold text-text-primary">{progress}%</div>
          <div className="text-xs text-text-muted">complete</div>
        </div>
      </div>

      <div className="mt-4">
        <div className="flex justify-between text-sm text-text-secondary mb-1">
          <span>Progress</span>
          <span>{progress}%</span>
        </div>
        <div className="h-2 bg-bg-secondary rounded-full overflow-hidden">
          <div
            className={`h-full rounded-full bg-gradient-to-r ${color}`}
            style={{ width: `${progress}%` }}
          />
        </div>
      </div>

      <div className="mt-4 flex space-x-2">
        <button className="flex-1 py-2 px-3 bg-bg-secondary hover:bg-bg-tertiary text-text-primary rounded-lg text-sm font-medium transition">
          Continue
        </button>
        <button className="py-2 px-3 border border-border-subtle hover:border-neon-cyan text-text-secondary rounded-lg text-sm transition">
          Details
        </button>
      </div>
    </div>
  )
}