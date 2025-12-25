'use client'

import { useState } from 'react'
import { Save, X } from 'lucide-react'

export interface ModuleFormData {
  pathway_id: string
  title: string
  description: string
  sort_order: number
  estimated_hours: number
  xp_reward: number
  status: 'draft' | 'published'
}

interface ModuleFormProps {
  pathwayId: string
  pathwayTitle?: string
  initialData?: Partial<ModuleFormData>
  onSave: (data: ModuleFormData) => Promise<void>
  onCancel?: () => void
  isEditing?: boolean
}

export default function ModuleForm({
  pathwayId,
  pathwayTitle,
  initialData,
  onSave,
  onCancel,
  isEditing = false,
}: ModuleFormProps) {
  const [isSaving, setIsSaving] = useState(false)
  const [formData, setFormData] = useState<ModuleFormData>({
    pathway_id: pathwayId,
    title: '',
    description: '',
    sort_order: 0,
    estimated_hours: 4,
    xp_reward: 100,
    status: 'draft',
    ...initialData,
  })

  const updateField = <K extends keyof ModuleFormData>(
    field: K,
    value: ModuleFormData[K]
  ) => {
    setFormData((prev) => ({ ...prev, [field]: value }))
  }

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault()

    // Validation
    if (!formData.title.trim()) {
      alert('Please enter a module title')
      return
    }

    if (formData.estimated_hours < 1) {
      alert('Estimated hours must be at least 1')
      return
    }

    if (formData.xp_reward < 0) {
      alert('XP reward cannot be negative')
      return
    }

    setIsSaving(true)
    try {
      await onSave(formData)
    } catch (error: any) {
      alert(`Failed to save module: ${error.message}`)
    } finally {
      setIsSaving(false)
    }
  }

  return (
    <form onSubmit={handleSubmit} className="space-y-6">
      {/* Header */}
      <div className="flex items-center justify-between">
        <div>
          <h2 className="text-2xl font-bold text-text-primary">
            {isEditing ? 'Edit Module' : 'Create New Module'}
          </h2>
          {pathwayTitle && (
            <p className="text-sm text-text-secondary mt-1">
              For pathway: <span className="font-medium">{pathwayTitle}</span>
            </p>
          )}
        </div>

        <div className="flex items-center gap-3">
          {onCancel && (
            <button
              type="button"
              onClick={onCancel}
              className="flex items-center gap-2 px-4 py-2 bg-bg-primary border border-border-default rounded-lg hover:bg-bg-hover transition-colors text-text-primary"
            >
              <X className="w-4 h-4" />
              Cancel
            </button>
          )}

          <button
            type="submit"
            disabled={isSaving}
            className="flex items-center gap-2 px-6 py-2 bg-accent-primary text-white rounded-lg hover:bg-accent-secondary transition-colors disabled:opacity-50"
          >
            <Save className="w-4 h-4" />
            {isSaving ? 'Saving...' : isEditing ? 'Update Module' : 'Create Module'}
          </button>
        </div>
      </div>

      {/* Basic Information */}
      <div className="bg-bg-elevated border border-border-default rounded-lg p-6">
        <h3 className="text-lg font-semibold text-text-primary mb-4">Module Information</h3>

        <div className="space-y-4">
          {/* Title */}
          <div>
            <label className="block text-sm font-medium text-text-secondary mb-2">
              Module Title *
            </label>
            <input
              type="text"
              value={formData.title}
              onChange={(e) => updateField('title', e.target.value)}
              placeholder="e.g., The Hacker's Toolkit (Core Python)"
              className="w-full px-4 py-2 bg-bg-primary border border-border-default rounded-lg text-text-primary focus:outline-none focus:ring-2 focus:ring-accent-primary"
              required
            />
          </div>

          {/* Description */}
          <div>
            <label className="block text-sm font-medium text-text-secondary mb-2">
              Description
            </label>
            <textarea
              value={formData.description}
              onChange={(e) => updateField('description', e.target.value)}
              placeholder="What will students learn in this module?&#10;&#10;Topics covered:&#10;- File I/O and string manipulation&#10;- Regular expressions&#10;- Subprocess for tool chaining"
              className="w-full px-4 py-2 bg-bg-primary border border-border-default rounded-lg text-text-primary focus:outline-none focus:ring-2 focus:ring-accent-primary"
              rows={6}
            />
            <p className="text-xs text-text-tertiary mt-1">
              Provide a detailed overview of module content
            </p>
          </div>
        </div>
      </div>

      {/* Module Settings */}
      <div className="bg-bg-elevated border border-border-default rounded-lg p-6">
        <h3 className="text-lg font-semibold text-text-primary mb-4">Module Settings</h3>

        <div className="grid grid-cols-1 md:grid-cols-3 gap-4">
          {/* Sort Order */}
          <div>
            <label className="block text-sm font-medium text-text-secondary mb-2">
              Sort Order
            </label>
            <input
              type="number"
              value={formData.sort_order}
              onChange={(e) => updateField('sort_order', parseInt(e.target.value) || 0)}
              min="0"
              className="w-full px-4 py-2 bg-bg-primary border border-border-default rounded-lg text-text-primary focus:outline-none focus:ring-2 focus:ring-accent-primary"
            />
            <p className="text-xs text-text-tertiary mt-1">Display order in module list</p>
          </div>

          {/* Estimated Hours */}
          <div>
            <label className="block text-sm font-medium text-text-secondary mb-2">
              Estimated Hours *
            </label>
            <input
              type="number"
              value={formData.estimated_hours}
              onChange={(e) => updateField('estimated_hours', parseInt(e.target.value) || 1)}
              min="1"
              className="w-full px-4 py-2 bg-bg-primary border border-border-default rounded-lg text-text-primary focus:outline-none focus:ring-2 focus:ring-accent-primary"
            />
            <p className="text-xs text-text-tertiary mt-1">Time to complete</p>
          </div>

          {/* XP Reward */}
          <div>
            <label className="block text-sm font-medium text-text-secondary mb-2">
              XP Reward
            </label>
            <input
              type="number"
              value={formData.xp_reward}
              onChange={(e) => updateField('xp_reward', parseInt(e.target.value) || 0)}
              min="0"
              step="50"
              className="w-full px-4 py-2 bg-bg-primary border border-border-default rounded-lg text-text-primary focus:outline-none focus:ring-2 focus:ring-accent-primary"
            />
            <p className="text-xs text-text-tertiary mt-1">XP for completing module</p>
          </div>
        </div>

        {/* Info Cards */}
        <div className="mt-4 p-4 bg-accent-primary/10 border border-accent-primary/30 rounded-lg">
          <div className="flex items-start gap-3">
            <div className="text-2xl">💡</div>
            <div className="flex-1">
              <h4 className="font-medium text-text-primary mb-1">Module Guidelines</h4>
              <ul className="text-sm text-text-secondary space-y-1">
                <li>• Modules should focus on a single major topic or skill</li>
                <li>• Include 3-8 exercises per module for best learning experience</li>
                <li>• Estimate hours conservatively - students vary in pace</li>
                <li>• XP rewards typically range from 100-500 per module</li>
              </ul>
            </div>
          </div>
        </div>
      </div>

      {/* Publishing Status */}
      <div className="bg-bg-elevated border border-border-default rounded-lg p-6">
        <h3 className="text-lg font-semibold text-text-primary mb-4">Publishing Status</h3>

        <div className="space-y-3">
          <label className="flex items-center gap-3 p-3 border border-border-default rounded-lg cursor-pointer hover:bg-bg-hover transition-colors">
            <input
              type="radio"
              name="status"
              value="draft"
              checked={formData.status === 'draft'}
              onChange={() => updateField('status', 'draft')}
              className="w-4 h-4 text-accent-primary focus:ring-accent-primary"
            />
            <div>
              <div className="font-medium text-text-primary">Draft</div>
              <div className="text-sm text-text-secondary">
                Save as draft - add exercises before publishing
              </div>
            </div>
          </label>

          <label className="flex items-center gap-3 p-3 border border-border-default rounded-lg cursor-pointer hover:bg-bg-hover transition-colors">
            <input
              type="radio"
              name="status"
              value="published"
              checked={formData.status === 'published'}
              onChange={() => updateField('status', 'published')}
              className="w-4 h-4 text-accent-primary focus:ring-accent-primary"
            />
            <div>
              <div className="font-medium text-text-primary">Published</div>
              <div className="text-sm text-text-secondary">
                Make available to students immediately
              </div>
            </div>
          </label>
        </div>
      </div>

      {/* Preview Card */}
      <div className="bg-bg-elevated border border-border-default rounded-lg p-6">
        <h3 className="text-lg font-semibold text-text-primary mb-4">Preview</h3>
        
        <div className="border border-border-default rounded-lg p-4 bg-bg-primary">
          <div className="flex items-start justify-between">
            <div className="flex-1">
              <div className="flex items-center gap-2 mb-2">
                <span className="text-sm font-medium text-accent-primary">
                  Module {formData.sort_order + 1}
                </span>
                {formData.status === 'draft' && (
                  <span className="text-xs px-2 py-0.5 bg-yellow-500/20 text-yellow-400 rounded">
                    Draft
                  </span>
                )}
              </div>
              <h4 className="text-lg font-semibold text-text-primary mb-1">
                {formData.title || 'Module Title'}
              </h4>
              {formData.description && (
                <p className="text-sm text-text-secondary line-clamp-2">
                  {formData.description}
                </p>
              )}
            </div>
          </div>

          <div className="flex items-center gap-4 mt-4 pt-4 border-t border-border-subtle">
            <div className="flex items-center gap-1.5 text-sm text-text-secondary">
              <span>⏱️</span>
              <span>{formData.estimated_hours}h</span>
            </div>
            <div className="flex items-center gap-1.5 text-sm text-text-secondary">
              <span>⭐</span>
              <span>{formData.xp_reward} XP</span>
            </div>
            <div className="flex items-center gap-1.5 text-sm text-text-secondary">
              <span>📝</span>
              <span>0 exercises</span>
            </div>
          </div>
        </div>
      </div>
    </form>
  )
}
