'use client'

import { useRouter } from 'next/navigation'
import PathwayForm from '@/components/creator/PathwayForm'
import { api } from '@/lib/api'

export default function NewPathwayPage() {
  const router = useRouter()

  const handleSave = async (pathwayData: any) => {
    try {
      const payload = {
        title: pathwayData.title,
        subtitle: pathwayData.subtitle || null,
        description: pathwayData.description || null,
        level: pathwayData.level,
        duration_weeks: pathwayData.duration_weeks,
        color_gradient: pathwayData.color_gradient || null,
        icon: pathwayData.icon || null,
        sort_order: pathwayData.sort_order,
        prerequisites: pathwayData.prerequisites || [],
        status: pathwayData.status,
      }

      const response = await api.post('/content-creator/pathways', payload)
      
      alert('Pathway created successfully!')
      router.push(`/creator/pathways/${response.id}`)
    } catch (error: any) {
      console.error('Failed to create pathway:', error)
      throw error
    }
  }

  return (
    <div className="min-h-screen bg-bg-primary py-8">
      <div className="max-w-4xl mx-auto px-6">
        <PathwayForm
          onSave={handleSave}
          onCancel={() => router.back()}
        />
      </div>
    </div>
  )
}
