"use client"

import { useEffect, useState } from 'react'
import { api } from '@/lib/api'
import PathwayCard from './PathwayCard'

interface PathwayProgress {
  pathway_id: string
  title: string
  progress_percentage: number
  completed_modules: number
  total_modules: number
  xp_earned: number
  streak_days: number
  last_activity?: string
}

interface ProgressResponse {
  pathways: PathwayProgress[]
  totals: any
}

export default function PathwayProgressList() {
  const [pathways, setPathways] = useState<PathwayProgress[]>([])
  const [loading, setLoading] = useState(true)
  const [error, setError] = useState<string | null>(null)

  useEffect(() => {
    async function fetchPathways() {
      try {
        const data = await api.get<ProgressResponse>('/users/me/progress')
        // Handle null or undefined pathways
        setPathways(data.pathways || [])
      } catch (err: any) {
        console.error('Failed to fetch pathway progress:', err)
        setError(err.message || 'Failed to load pathways')
      } finally {
        setLoading(false)
      }
    }
    fetchPathways()
  }, [])

  if (loading) {
    return (
      <div className="grid grid-cols-1 md:grid-cols-2 gap-4">
        {[...Array(4)].map((_, i) => (
          <div
            key={i}
            className="bg-bg-tertiary border border-border-subtle rounded-xl p-5 animate-pulse"
          >
            <div className="h-4 bg-gray-700 rounded w-3/4 mb-4"></div>
            <div className="h-3 bg-gray-700 rounded w-1/2 mb-2"></div>
            <div className="h-2 bg-gray-700 rounded w-full"></div>
          </div>
        ))}
      </div>
    )
  }

  if (error) {
    return (
      <div className="bg-bg-elevated border border-border-default rounded-2xl p-6 text-center">
        <p className="text-text-secondary">Failed to load pathways: {error}</p>
        <button
          onClick={() => window.location.reload()}
          className="mt-2 px-4 py-2 bg-primary text-white rounded-lg hover:bg-primary-dark"
        >
          Retry
        </button>
      </div>
    )
  }

  if (pathways.length === 0) {
    return (
      <div className="bg-bg-elevated border border-border-default rounded-2xl p-6 text-center">
        <p className="text-text-secondary">You are not enrolled in any pathways yet.</p>
        <button
          onClick={() => window.location.href = '/dashboard/pathways'}
          className="mt-2 px-4 py-2 bg-primary text-white rounded-lg hover:bg-primary-dark"
        >
          Browse Pathways
        </button>
      </div>
    )
  }

  // Map pathway to icon and color based on title or some logic
  const getIconAndColor = (title: string) => {
    if (title.toLowerCase().includes('python')) {
      return { icon: '🐍', color: 'from-green-400 to-cyan-400' }
    }
    if (title.toLowerCase().includes('c') || title.toLowerCase().includes('systems')) {
      return { icon: '🔧', color: 'from-blue-400 to-purple-400' }
    }
    if (title.toLowerCase().includes('assembly')) {
      return { icon: '⚙️', color: 'from-red-400 to-pink-400' }
    }
    if (title.toLowerCase().includes('javascript') || title.toLowerCase().includes('web')) {
      return { icon: '🌐', color: 'from-yellow-400 to-orange-400' }
    }
    // default
    return { icon: '📚', color: 'from-gray-400 to-gray-600' }
  }

  return (
    <div className="grid grid-cols-1 md:grid-cols-2 gap-4">
      {pathways.map((pathway) => {
        const { icon, color } = getIconAndColor(pathway.title)
        return (
          <PathwayCard
            key={pathway.pathway_id}
            title={pathway.title}
            description={`${pathway.completed_modules}/${pathway.total_modules} modules completed`}
            progress={pathway.progress_percentage}
            icon={icon}
            color={color}
          />
        )
      })}
    </div>
  )
}