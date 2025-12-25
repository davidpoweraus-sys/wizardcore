'use client'

import { useEffect, useState } from 'react'
import { useRouter, useParams } from 'next/navigation'
import PathwayForm, { PathwayFormData } from '@/components/creator/PathwayForm'
import { api } from '@/lib/api'
import { Loader2, ArrowLeft } from 'lucide-react'

interface Pathway {
  id: string
  title: string
  subtitle?: string
  description?: string
  level: string
  duration_weeks: number
  color_gradient?: string
  icon?: string
  sort_order: number
  prerequisites: string[]
  status: string
  created_at: string
  updated_at: string
}

export default function EditPathwayPage() {
  const router = useRouter()
  const params = useParams()
  const pathwayId = params.id as string

  const [loading, setLoading] = useState(true)
  const [pathway, setPathway] = useState<Pathway | null>(null)
  const [error, setError] = useState<string | null>(null)

  useEffect(() => {
    loadPathway()
  }, [pathwayId])

  const loadPathway = async () => {
    try {
      setLoading(true)
      setError(null)

      // Get all pathways and find the one we're editing
      const pathways = await api.get<Pathway[]>('/content-creator/pathways')
      const foundPathway = pathways.find(p => p.id === pathwayId)

      if (!foundPathway) {
        setError('Pathway not found')
        return
      }

      setPathway(foundPathway)
    } catch (err: any) {
      console.error('Failed to load pathway:', err)
      setError(err.message || 'Failed to load pathway')
    } finally {
      setLoading(false)
    }
  }

  const handleSave = async (formData: PathwayFormData) => {
    try {
      const payload = {
        title: formData.title,
        subtitle: formData.subtitle || null,
        description: formData.description || null,
        level: formData.level,
        duration_weeks: formData.duration_weeks,
        color_gradient: formData.color_gradient || null,
        icon: formData.icon || null,
        sort_order: formData.sort_order,
        prerequisites: formData.prerequisites || [],
        status: formData.status,
      }

      await api.put(`/content-creator/pathways/${pathwayId}`, payload)
      
      // Redirect to dashboard with success message
      router.push('/creator/dashboard')
      // In a real app, you might want to show a toast notification
    } catch (error: any) {
      alert(`Failed to update pathway: ${error.message}`)
      throw error
    }
  }

  const handleCancel = () => {
    router.push('/creator/dashboard')
  }

  if (loading) {
    return (
      <div className="flex items-center justify-center min-h-screen bg-bg-primary">
        <div className="flex flex-col items-center gap-4">
          <Loader2 className="w-8 h-8 animate-spin text-accent-primary" />
          <p className="text-text-secondary">Loading pathway...</p>
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
            onClick={() => router.push('/creator/dashboard')}
            className="flex items-center gap-2 px-4 py-2 bg-accent-primary text-white rounded-lg hover:bg-accent-secondary transition-colors"
          >
            <ArrowLeft className="w-4 h-4" />
            Back to Dashboard
          </button>
        </div>
      </div>
    )
  }

  if (!pathway) {
    return (
      <div className="flex items-center justify-center min-h-screen bg-bg-primary">
        <div className="max-w-md p-6 bg-bg-elevated border border-border-default rounded-lg">
          <h2 className="text-xl font-bold text-text-primary mb-2">Pathway Not Found</h2>
          <p className="text-text-secondary mb-4">The pathway you're trying to edit doesn't exist.</p>
          <button
            onClick={() => router.push('/creator/dashboard')}
            className="flex items-center gap-2 px-4 py-2 bg-accent-primary text-white rounded-lg hover:bg-accent-secondary transition-colors"
          >
            <ArrowLeft className="w-4 h-4" />
            Back to Dashboard
          </button>
        </div>
      </div>
    )
  }

  // Convert backend data to form data
  const initialData: PathwayFormData = {
    title: pathway.title,
    subtitle: pathway.subtitle || '',
    description: pathway.description || '',
    level: pathway.level as 'Beginner' | 'Intermediate' | 'Advanced' | 'Expert',
    duration_weeks: pathway.duration_weeks,
    color_gradient: pathway.color_gradient || 'from-blue-500 to-purple-600',
    icon: pathway.icon || '🎓',
    sort_order: pathway.sort_order,
    prerequisites: pathway.prerequisites || [],
    status: pathway.status as 'draft' | 'published',
  }

  return (
    <div className="min-h-screen bg-bg-primary">
      <div className="max-w-6xl mx-auto px-4 py-8">
        {/* Header */}
        <div className="mb-8">
          <button
            onClick={handleCancel}
            className="flex items-center gap-2 text-text-secondary hover:text-text-primary transition-colors mb-4"
          >
            <ArrowLeft className="w-4 h-4" />
            Back to Dashboard
          </button>
          <h1 className="text-3xl font-bold text-text-primary">Edit Pathway</h1>
          <p className="text-text-secondary mt-2">
            Update your pathway details and settings
          </p>
        </div>

        {/* Form */}
        <div className="bg-bg-elevated border border-border-default rounded-lg p-6">
          <PathwayForm
            initialData={initialData}
            onSave={handleSave}
            onCancel={handleCancel}
            isEditing={true}
          />
        </div>

        {/* Last Updated Info */}
        <div className="mt-6 text-sm text-text-tertiary">
          <p>
            Created: {new Date(pathway.created_at).toLocaleDateString()} • 
            Last Updated: {new Date(pathway.updated_at).toLocaleDateString()}
          </p>
        </div>
      </div>
    </div>
  )
}
