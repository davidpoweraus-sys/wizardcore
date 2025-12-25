'use client'

import { useEffect, useState } from 'react'
import { useRouter, useParams } from 'next/navigation'
import ModuleForm, { ModuleFormData } from '@/components/creator/ModuleForm'
import { api } from '@/lib/api'
import { Loader2, ArrowLeft } from 'lucide-react'

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

interface Pathway {
  id: string
  title: string
}

export default function EditModulePage() {
  const router = useRouter()
  const params = useParams()
  const moduleId = params.id as string

  const [loading, setLoading] = useState(true)
  const [module, setModule] = useState<Module | null>(null)
  const [pathway, setPathway] = useState<Pathway | null>(null)
  const [error, setError] = useState<string | null>(null)

  useEffect(() => {
    loadModule()
  }, [moduleId])

  const loadModule = async () => {
    try {
      setLoading(true)
      setError(null)

      // Get all modules and find the one we're editing
      const modules = await api.get<Module[]>('/content-creator/modules')
      const foundModule = modules.find(m => m.id === moduleId)

      if (!foundModule) {
        setError('Module not found')
        return
      }

      setModule(foundModule)

      // Get the pathway for this module
      const pathways = await api.get<Pathway[]>('/content-creator/pathways')
      const foundPathway = pathways.find(p => p.id === foundModule.pathway_id)
      setPathway(foundPathway || null)
    } catch (err: any) {
      console.error('Failed to load module:', err)
      setError(err.message || 'Failed to load module')
    } finally {
      setLoading(false)
    }
  }

  const handleSave = async (formData: ModuleFormData) => {
    try {
      const payload = {
        title: formData.title,
        description: formData.description || null,
        sort_order: formData.sort_order,
        estimated_hours: formData.estimated_hours || null,
        xp_reward: formData.xp_reward,
        status: formData.status,
      }

      await api.put(`/content-creator/modules/${moduleId}`, payload)
      
      // Redirect to dashboard with success message
      router.push('/creator/dashboard')
      // In a real app, you might want to show a toast notification
    } catch (error: any) {
      alert(`Failed to update module: ${error.message}`)
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
          <p className="text-text-secondary">Loading module...</p>
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

  if (!module) {
    return (
      <div className="flex items-center justify-center min-h-screen bg-bg-primary">
        <div className="max-w-md p-6 bg-bg-elevated border border-border-default rounded-lg">
          <h2 className="text-xl font-bold text-text-primary mb-2">Module Not Found</h2>
          <p className="text-text-secondary mb-4">The module you're trying to edit doesn't exist.</p>
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
  const initialData: Partial<ModuleFormData> = {
    title: module.title,
    description: module.description || '',
    sort_order: module.sort_order,
    estimated_hours: module.estimated_hours || 1,
    xp_reward: module.xp_reward,
    status: module.status as 'draft' | 'published',
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
          <div className="flex items-center justify-between">
            <div>
              <h1 className="text-3xl font-bold text-text-primary">Edit Module</h1>
              <p className="text-text-secondary mt-2">
                Update your module details and settings
              </p>
              {pathway && (
                <div className="mt-2">
                  <span className="text-sm text-text-tertiary">Pathway: </span>
                  <span className="text-sm text-text-secondary">{pathway.title}</span>
                </div>
              )}
            </div>
          </div>
        </div>

        {/* Form */}
        <div className="bg-bg-elevated border border-border-default rounded-lg p-6">
          <ModuleForm
            pathwayId={module.pathway_id}
            pathwayTitle={pathway?.title}
            initialData={initialData}
            onSave={handleSave}
            onCancel={handleCancel}
            isEditing={true}
          />
        </div>

        {/* Last Updated Info */}
        <div className="mt-6 text-sm text-text-tertiary">
          <p>
            Created: {new Date(module.created_at).toLocaleDateString()} • 
            Last Updated: {new Date(module.updated_at).toLocaleDateString()}
          </p>
        </div>
      </div>
    </div>
  )
}
