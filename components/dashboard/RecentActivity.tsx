'use client'

import { Check, Code, BookOpen, Trophy, Clock, Loader } from 'pixelarticons/fonts/react'
import { useEffect, useState } from 'react'
import { api } from '@/lib/api'

interface Activity {
  id: string
  type: 'completion' | 'practice' | 'reading' | 'achievement' | 'streak' | 'other'
  title: string
  description: string
  time: string // ISO timestamp or relative
  icon: string // icon name mapping
  color: string
}

const iconMap: Record<string, any> = {
  completion: Check,
  practice: Code,
  reading: BookOpen,
  achievement: Trophy,
  streak: Clock,
  other: Check,
}

const colorMap: Record<string, string> = {
  completion: 'text-green-400',
  practice: 'text-blue-400',
  reading: 'text-purple-400',
  achievement: 'text-yellow-400',
  streak: 'text-neon-cyan',
  other: 'text-gray-400',
}

export default function RecentActivity() {
  const [activities, setActivities] = useState<Activity[]>([])
  const [loading, setLoading] = useState(true)
  const [error, setError] = useState<string | null>(null)

  useEffect(() => {
    async function fetchActivities() {
      try {
        const data = await api.get<{ activities: Activity[] }>('/users/me/activities')
        setActivities(data.activities)
      } catch (err: any) {
        console.error('Failed to fetch activities:', err)
        setError(err.message || 'Failed to load activities')
      } finally {
        setLoading(false)
      }
    }
    fetchActivities()
  }, [])

  if (loading) {
    return (
      <div className="bg-bg-elevated border border-border-default rounded-2xl p-6">
        <div className="flex items-center justify-between mb-6">
          <h2 className="text-xl font-semibold text-text-primary">Recent Activity</h2>
          <button className="text-sm text-neon-lavender hover:underline">View All</button>
        </div>
        <div className="space-y-4">
          {[...Array(5)].map((_, i) => (
            <div key={i} className="flex items-center p-4 border border-border-subtle rounded-xl animate-pulse">
              <div className="w-12 h-12 bg-gray-700 rounded-lg"></div>
              <div className="ml-4 flex-1">
                <div className="h-4 bg-gray-700 rounded w-3/4 mb-2"></div>
                <div className="h-3 bg-gray-700 rounded w-1/2"></div>
              </div>
              <div className="h-3 bg-gray-700 rounded w-16"></div>
            </div>
          ))}
        </div>
      </div>
    )
  }

  if (error) {
    return (
      <div className="bg-bg-elevated border border-border-default rounded-2xl p-6">
        <div className="flex items-center justify-between mb-6">
          <h2 className="text-xl font-semibold text-text-primary">Recent Activity</h2>
          <button className="text-sm text-neon-lavender hover:underline">View All</button>
        </div>
        <div className="text-center text-text-secondary py-8">
          <p>Failed to load activities: {error}</p>
          <button
            onClick={() => window.location.reload()}
            className="mt-2 px-4 py-2 bg-primary text-white rounded-lg hover:bg-primary-dark"
          >
            Retry
          </button>
        </div>
      </div>
    )
  }

  if (activities.length === 0) {
    return (
      <div className="bg-bg-elevated border border-border-default rounded-2xl p-6">
        <div className="flex items-center justify-between mb-6">
          <h2 className="text-xl font-semibold text-text-primary">Recent Activity</h2>
          <button className="text-sm text-neon-lavender hover:underline">View All</button>
        </div>
        <div className="text-center text-text-secondary py-8">
          <p>No recent activity to show.</p>
          <p className="text-sm">Complete an exercise to see your activity here.</p>
        </div>
      </div>
    )
  }

  return (
    <div className="bg-bg-elevated border border-border-default rounded-2xl p-6">
      <div className="flex items-center justify-between mb-6">
        <h2 className="text-xl font-semibold text-text-primary">Recent Activity</h2>
        <button className="text-sm text-neon-lavender hover:underline">View All</button>
      </div>

      <div className="space-y-4">
        {activities.map((activity) => {
          const Icon = iconMap[activity.type] || Check
          const color = colorMap[activity.type] || 'text-gray-400'
          return (
            <div
              key={activity.id}
              className="flex items-center p-4 border border-border-subtle rounded-xl hover:bg-bg-tertiary transition"
            >
              <div className={`p-3 rounded-lg ${color} bg-opacity-20`}>
                <Icon className="w-5 h-5" />
              </div>
              <div className="ml-4 flex-1">
                <h3 className="font-medium text-text-primary">{activity.title}</h3>
                <p className="text-sm text-text-secondary">{activity.description}</p>
              </div>
              <div className="text-sm text-text-muted">{activity.time}</div>
            </div>
          )
        })}
      </div>
    </div>
  )
}