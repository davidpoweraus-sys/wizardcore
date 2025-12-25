'use client'

import { useEffect, useState } from 'react'
import { useRouter } from 'next/navigation'
import { api } from '@/lib/api'
import { Plus, Loader2, Upload, Download } from 'lucide-react'
import StatsCards from '@/components/creator/StatsCards'
import ContentList from '@/components/creator/ContentList'
import ImportButton from '@/components/creator/ImportButton'
import ExportButton from '@/components/creator/ExportButton'

interface CreatorStats {
  total_pathways: number
  total_modules: number
  total_exercises: number
  total_students: number
  average_rating: number
  total_ratings: number
  total_views: number
  total_enrollments: number
  total_completions: number
  completion_rate: number
  pending_reviews: number
  published_content: number
  draft_content: number
}

interface Pathway {
  id: string
  title: string
  subtitle?: string
  level: string
  duration_weeks: number
  status: string
  created_at: string
  updated_at: string
  color_gradient?: string
  icon?: string
}

interface Module {
  id: string
  pathway_id: string
  title: string
  description?: string
  sort_order: number
  estimated_hours?: number
  xp_reward: number
  status: string
  created_at: string
  updated_at: string
}

interface Exercise {
  id: string
  module_id: string
  title: string
  difficulty: string
  points: number
  status: string
  created_at: string
  updated_at: string
}

export default function CreatorDashboard() {
  const router = useRouter()
  const [loading, setLoading] = useState(true)
  const [stats, setStats] = useState<CreatorStats | null>(null)
  const [pathways, setPathways] = useState<Pathway[]>([])
  const [modules, setModules] = useState<Module[]>([])
  const [exercises, setExercises] = useState<Exercise[]>([])
  const [error, setError] = useState<string | null>(null)

  useEffect(() => {
    loadDashboardData()
  }, [])

  const loadDashboardData = async () => {
    try {
      setLoading(true)
      setError(null)

      // Load all data in parallel
      const [statsData, pathwaysData, modulesData, exercisesData] = await Promise.all([
        api.get<CreatorStats>('/content-creator/stats'),
        api.get<Pathway[]>('/content-creator/pathways'),
        api.get<Module[]>('/content-creator/modules'),
        api.get<Exercise[]>('/content-creator/exercises'),
      ])

      setStats(statsData)
      setPathways(pathwaysData)
      setModules(modulesData)
      setExercises(exercisesData)
    } catch (err: any) {
      console.error('Failed to load dashboard data:', err)
      setError(err.message || 'Failed to load dashboard data')
    } finally {
      setLoading(false)
    }
  }

  const handleDeletePathway = async (id: string) => {
    if (!confirm('Are you sure you want to delete this pathway? This will also delete all associated modules and exercises.')) {
      return
    }

    try {
      await api.delete(`/content-creator/pathways/${id}`)
      await loadDashboardData() // Reload all data
    } catch (err: any) {
      alert(`Failed to delete pathway: ${err.message}`)
    }
  }

  const handleDeleteModule = async (id: string) => {
    if (!confirm('Are you sure you want to delete this module? This will also delete all associated exercises.')) {
      return
    }

    try {
      await api.delete(`/content-creator/modules/${id}`)
      await loadDashboardData() // Reload all data
    } catch (err: any) {
      alert(`Failed to delete module: ${err.message}`)
    }
  }

  const handleDeleteExercise = async (id: string) => {
    if (!confirm('Are you sure you want to delete this exercise?')) {
      return
    }

    try {
      await api.delete(`/content-creator/exercises/${id}`)
      await loadDashboardData() // Reload all data
    } catch (err: any) {
      alert(`Failed to delete exercise: ${err.message}`)
    }
  }

  const handleExportPathway = async (pathwayId: string, pathwayTitle: string) => {
    try {
      const data = await api.get(`/content-creator/pathways/${pathwayId}/export`)
      
      // Format the JSON nicely
      const formattedJson = JSON.stringify(data, null, 2)
      
      // Create a downloadable file
      const blob = new Blob([formattedJson], { type: 'application/json' })
      const url = URL.createObjectURL(blob)
      const a = document.createElement('a')
      a.href = url
      a.download = `${pathwayTitle.replace(/[^a-z0-9]/gi, '_').toLowerCase()}_export.json`
      document.body.appendChild(a)
      a.click()
      document.body.removeChild(a)
      URL.revokeObjectURL(url)
    } catch (error: any) {
      alert(`Failed to export pathway: ${error.message}`)
    }
  }

  if (loading) {
    return (
      <div className="flex items-center justify-center min-h-screen bg-bg-primary">
        <div className="flex flex-col items-center gap-4">
          <Loader2 className="w-8 h-8 animate-spin text-accent-primary" />
          <p className="text-text-secondary">Loading your dashboard...</p>
        </div>
      </div>
    )
  }

  if (error) {
    return (
      <div className="flex items-center justify-center min-h-screen bg-bg-primary">
        <div className="max-w-md p-6 bg-bg-elevated border border-border-default rounded-lg">
          <h2 className="text-xl font-bold text-red-500 mb-2">Error</h2>
          <p className="text-text-secondary mb-4">{error}</p>
          <button
            onClick={loadDashboardData}
            className="px-4 py-2 bg-accent-primary text-white rounded-lg hover:bg-accent-secondary transition-colors"
          >
            Retry
          </button>
        </div>
      </div>
    )
  }

  return (
    <div className="min-h-screen bg-bg-primary">
      <div className="max-w-7xl mx-auto px-4 py-8">
        {/* Header */}
        <div className="mb-8">
          <div className="flex items-center justify-between">
            <div>
              <h1 className="text-3xl font-bold text-text-primary mb-2">Content Creator Dashboard</h1>
              <p className="text-text-secondary">
                Manage your pathways, modules, and exercises
              </p>
            </div>
            <div className="flex items-center gap-3">
              <ImportButton onImportComplete={loadDashboardData} />
            </div>
          </div>
        </div>

        {/* Quick Actions */}
        <div className="grid grid-cols-1 md:grid-cols-3 gap-4 mb-8">
          <button
            onClick={() => router.push('/creator/pathways/new')}
            className="flex items-center justify-center gap-3 p-6 bg-gradient-to-r from-blue-500 to-purple-600 text-white rounded-lg hover:from-blue-600 hover:to-purple-700 transition-all shadow-lg"
          >
            <Plus className="w-6 h-6" />
            <span className="font-semibold">Create Pathway</span>
          </button>

          <button
            onClick={() => {
              if (pathways.length === 0) {
                alert('Please create a pathway first')
                return
              }
              router.push(`/creator/modules/new?pathway_id=${pathways[0].id}`)
            }}
            className="flex items-center justify-center gap-3 p-6 bg-gradient-to-r from-green-500 to-teal-600 text-white rounded-lg hover:from-green-600 hover:to-teal-700 transition-all shadow-lg"
          >
            <Plus className="w-6 h-6" />
            <span className="font-semibold">Create Module</span>
          </button>

          <button
            onClick={() => {
              if (modules.length === 0) {
                alert('Please create a module first')
                return
              }
              router.push(`/creator/exercises/new?module_id=${modules[0].id}`)
            }}
            className="flex items-center justify-center gap-3 p-6 bg-gradient-to-r from-red-500 to-orange-600 text-white rounded-lg hover:from-red-600 hover:to-orange-700 transition-all shadow-lg"
          >
            <Plus className="w-6 h-6" />
            <span className="font-semibold">Create Exercise</span>
          </button>
        </div>

        {/* Statistics */}
        {stats && <StatsCards stats={stats} />}

        {/* Content Lists */}
        <div className="space-y-8 mt-8">
          {/* Pathways */}
          <ContentList
            title="Your Pathways"
            items={pathways}
            type="pathway"
            onEdit={(id: string) => router.push(`/creator/pathways/${id}/edit`)}
            onDelete={handleDeletePathway}
            onCreate={() => router.push('/creator/pathways/new')}
            onExport={handleExportPathway}
          />

          {/* Modules */}
          <ContentList
            title="Your Modules"
            items={modules}
            type="module"
            onEdit={(id: string) => router.push(`/creator/modules/${id}/edit`)}
            onDelete={handleDeleteModule}
            onCreate={() => {
              if (pathways.length === 0) {
                alert('Please create a pathway first')
                return
              }
              router.push(`/creator/modules/new?pathway_id=${pathways[0].id}`)
            }}
          />

          {/* Exercises */}
          <ContentList
            title="Your Exercises"
            items={exercises}
            type="exercise"
            onEdit={(id: string) => router.push(`/creator/exercises/${id}/edit`)}
            onDelete={handleDeleteExercise}
            onCreate={() => {
              if (modules.length === 0) {
                alert('Please create a module first')
                return
              }
              router.push(`/creator/exercises/new?module_id=${modules[0].id}`)
            }}
          />
        </div>
      </div>
    </div>
  )
}
