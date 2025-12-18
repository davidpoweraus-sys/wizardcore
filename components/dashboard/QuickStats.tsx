"use client"

import { BookOpen, Target, Clock, TrendingUp } from 'lucide-react'
import { useEffect, useState } from 'react'
import { api } from '@/lib/api'

interface Stats {
  active_courses: number
  completion_rate: number
  study_time: number
  xp_earned: number
  total_xp: number
  xp_this_week: number
  current_streak: number
  modules_completed: number
  modules_total: number
}

export default function QuickStats() {
  const [stats, setStats] = useState<Stats | null>(null)
  const [loading, setLoading] = useState(true)
  const [error, setError] = useState<string | null>(null)

  useEffect(() => {
    async function fetchStats() {
      try {
        const data = await api.get<{ stats: Stats }>('/users/me/stats')
        setStats(data.stats)
      } catch (err: any) {
        console.error('Failed to fetch stats:', err)
        setError(err.message || 'Failed to load stats')
      } finally {
        setLoading(false)
      }
    }
    fetchStats()
  }, [])

  const statsData = stats ? [
    {
      title: 'Active Courses',
      value: stats.active_courses.toString(),
      change: '+1 this week',
      icon: BookOpen,
      color: 'from-neon-cyan to-neon-lavender',
    },
    {
      title: 'Completion Rate',
      value: `${stats.completion_rate}%`,
      change: '+5% from last month',
      icon: Target,
      color: 'from-neon-pink to-neon-purple',
    },
    {
      title: 'Study Time',
      value: `${stats.study_time}h`,
      change: '12h this week',
      icon: Clock,
      color: 'from-green-400 to-cyan-400',
    },
    {
      title: 'XP Earned',
      value: stats.xp_earned.toLocaleString(),
      change: '+320 today',
      icon: TrendingUp,
      color: 'from-yellow-400 to-orange-400',
    },
  ] : []

  if (loading) {
    return (
      <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-4 gap-4">
        {[...Array(4)].map((_, i) => (
          <div
            key={i}
            className="bg-bg-elevated border border-border-default rounded-2xl p-6 animate-pulse"
          >
            <div className="h-4 bg-gray-700 rounded w-1/2 mb-4"></div>
            <div className="h-8 bg-gray-700 rounded w-3/4 mb-2"></div>
            <div className="h-3 bg-gray-700 rounded w-1/3"></div>
          </div>
        ))}
      </div>
    )
  }

  if (error) {
    return (
      <div className="bg-bg-elevated border border-border-default rounded-2xl p-6 text-center">
        <p className="text-text-secondary">Failed to load stats: {error}</p>
        <button
          onClick={() => window.location.reload()}
          className="mt-2 px-4 py-2 bg-primary text-white rounded-lg hover:bg-primary-dark"
        >
          Retry
        </button>
      </div>
    )
  }

  return (
    <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-4 gap-4">
      {statsData.map((stat, index) => {
        const Icon = stat.icon
        return (
          <div
            key={index}
            className="bg-bg-elevated border border-border-default rounded-2xl p-6"
          >
            <div className="flex items-center justify-between">
              <div>
                <p className="text-sm text-text-secondary">{stat.title}</p>
                <p className="text-3xl font-bold text-text-primary mt-2">{stat.value}</p>
                <p className="text-xs text-text-muted mt-1">{stat.change}</p>
              </div>
              <div className={`p-3 rounded-xl bg-gradient-to-br ${stat.color}`}>
                <Icon className="w-6 h-6 text-white" />
              </div>
            </div>
          </div>
        )
      })}
    </div>
  )
}