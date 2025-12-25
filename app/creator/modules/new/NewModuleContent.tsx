'use client'

import { useRouter, useSearchParams } from 'next/navigation'
import ModuleForm from '@/components/creator/ModuleForm'
import { api } from '@/lib/api'

export default function NewModuleContent() {
  const router = useRouter()
  const searchParams = useSearchParams()
  const pathwayId = searchParams.get('pathway_id')
  const pathwayTitle = searchParams.get('pathway_title')

  if (!pathwayId) {
    return (
      <div className="min-h-screen bg-bg-primary flex items-center justify-center">
        <div className="text-center">
          <h1 className="text-2xl font-bold text-text-primary mb-2">
            Pathway ID Required
          </h1>
          <p className="text-text-secondary">
            Please select a pathway to create a module.
          </p>
          <button
            onClick={() => router.push('/creator/pathways')}
            className="mt-4 px-6 py-2 bg-accent-primary text-white rounded-lg hover:bg-accent-secondary"
          >
            Go to Pathways
          </button>
        </div>
      </div>
    )
  }

  const handleSave = async (moduleData: any) => {
    try {
      const payload = {
        pathway_id: pathwayId,
        title: moduleData.title,
        description: moduleData.description || null,
        sort_order: moduleData.sort_order,
        estimated_hours: moduleData.estimated_hours || null,
        xp_reward: moduleData.xp_reward,
        status: moduleData.status,
      }

      const response = await api.post('/content-creator/modules', payload)
      
      alert('Module created successfully!')
      router.push(`/creator/modules/${response.id}`)
    } catch (error: any) {
      console.error('Failed to create module:', error)
      throw error
    }
  }

  return (
    <div className="min-h-screen bg-bg-primary py-8">
      <div className="max-w-4xl mx-auto px-6">
        <ModuleForm
          pathwayId={pathwayId}
          pathwayTitle={pathwayTitle || undefined}
          onSave={handleSave}
          onCancel={() => router.back()}
        />
      </div>
    </div>
  )
}