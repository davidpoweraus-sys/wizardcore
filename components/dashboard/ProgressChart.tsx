'use client'

import { TrendingUp } from 'lucide-react'
import { useEffect, useState } from 'react'
import { api } from '@/lib/api'

interface DailyActivity {
  day: string
  value: number
  hours: number
}

interface WeeklyActivity {
  weekly_data: DailyActivity[]
  avg_daily_time_minutes: number
  completion_rate: number
  current_streak: number
  trend_percentage: number
}

export default function ProgressChart() {
  const [weeklyData, setWeeklyData] = useState<DailyActivity[]>([])
  const [avgDailyTime, setAvgDailyTime] = useState(0)
  const [completionRate, setCompletionRate] = useState(0)
  const [currentStreak, setCurrentStreak] = useState(0)
  const [trendPercentage, setTrendPercentage] = useState(0)
  const [loading, setLoading] = useState(true)
  const [error, setError] = useState<string | null>(null)

  useEffect(() => {
    async function fetchWeeklyActivity() {
      try {
        const data = await api.get<WeeklyActivity>('/users/me/activity/weekly')
        // Handle null or undefined weekly_data
        setWeeklyData(data.weekly_data || [])
        setAvgDailyTime(data.avg_daily_time_minutes || 0)
        setCompletionRate(data.completion_rate || 0)
        setCurrentStreak(data.current_streak || 0)
        setTrendPercentage(data.trend_percentage || 0)
      } catch (err: any) {
        console.error('Failed to fetch weekly activity:', err)
        setError(err.message || 'Failed to load weekly activity')
      } finally {
        setLoading(false)
      }
    }
    fetchWeeklyActivity()
  }, [])

  const maxValue = weeklyData && weeklyData.length > 0 ? Math.max(...weeklyData.map(d => d.value)) : 100

  if (loading) {
    return (
      <div className="bg-bg-elevated border border-border-default rounded-2xl p-6">
        <div className="h-6 bg-gray-700 rounded w-1/3 mb-4"></div>
        <div className="h-40 flex items-end justify-between">
          {[...Array(7)].map((_, i) => (
            <div key={i} className="w-8 bg-gray-700 rounded-t-lg animate-pulse" style={{ height: `${Math.random() * 80 + 20}%` }} />
          ))}
        </div>
        <div className="mt-6 pt-6 border-t border-border-subtle">
          <div className="flex justify-between">
            {[...Array(3)].map((_, i) => (
              <div key={i}>
                <div className="h-4 bg-gray-700 rounded w-20 mb-2"></div>
                <div className="h-6 bg-gray-700 rounded w-12"></div>
              </div>
            ))}
          </div>
        </div>
      </div>
    )
  }

  if (error) {
    return (
      <div className="bg-bg-elevated border border-border-default rounded-2xl p-6">
        <div className="text-center text-text-secondary">
          <p>Failed to load weekly activity: {error}</p>
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

  return (
    <div className="bg-bg-elevated border border-border-default rounded-2xl p-6">
      <div className="flex items-center justify-between mb-6">
        <div>
          <h2 className="text-xl font-semibold text-text-primary">Weekly Progress</h2>
          <p className="text-sm text-text-secondary">Your daily activity</p>
        </div>
        <div className="flex items-center text-neon-cyan">
          <TrendingUp className="w-5 h-5 mr-2" />
          <span className="font-bold">+{trendPercentage}%</span>
        </div>
      </div>

      <div className="flex items-end justify-between h-40">
        {weeklyData.map((item, index) => (
          <div key={index} className="flex flex-col items-center">
            <div
              className="w-8 rounded-t-lg bg-gradient-to-t from-neon-cyan to-neon-lavender"
              style={{ height: `${(item.value / maxValue) * 100}%` }}
            />
            <div className="mt-2 text-xs text-text-secondary">{item.day}</div>
            <div className="text-xs font-medium text-text-primary">{item.value}%</div>
          </div>
        ))}
      </div>

      <div className="mt-6 pt-6 border-t border-border-subtle">
        <div className="flex justify-between text-sm">
          <div>
            <div className="text-text-secondary">Avg. Daily Time</div>
            <div className="text-text-primary font-bold">
              {Math.floor(avgDailyTime / 60)}h {avgDailyTime % 60}m
            </div>
          </div>
          <div>
            <div className="text-text-secondary">Completion Rate</div>
            <div className="text-text-primary font-bold">{completionRate}%</div>
          </div>
          <div>
            <div className="text-text-secondary">Streak</div>
            <div className="text-text-primary font-bold">{currentStreak} days</div>
          </div>
        </div>
      </div>
    </div>
  )
}