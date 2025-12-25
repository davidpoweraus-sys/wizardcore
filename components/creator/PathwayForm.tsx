'use client'

import { useState } from 'react'
import { Save, X } from 'lucide-react'

export interface PathwayFormData {
  title: string
  subtitle: string
  description: string
  level: 'Beginner' | 'Intermediate' | 'Advanced' | 'Expert'
  duration_weeks: number
  color_gradient: string
  icon: string
  sort_order: number
  prerequisites: string[]
  status: 'draft' | 'published'
}

interface PathwayFormProps {
  initialData?: Partial<PathwayFormData>
  onSave: (data: PathwayFormData) => Promise<void>
  onCancel?: () => void
  isEditing?: boolean
}

const COLOR_GRADIENTS = [
  { value: 'from-blue-500 to-purple-600', label: 'Blue to Purple', preview: 'bg-gradient-to-r from-blue-500 to-purple-600' },
  { value: 'from-green-500 to-teal-600', label: 'Green to Teal', preview: 'bg-gradient-to-r from-green-500 to-teal-600' },
  { value: 'from-red-500 to-orange-600', label: 'Red to Orange', preview: 'bg-gradient-to-r from-red-500 to-orange-600' },
  { value: 'from-purple-500 to-pink-600', label: 'Purple to Pink', preview: 'bg-gradient-to-r from-purple-500 to-pink-600' },
  { value: 'from-yellow-500 to-red-600', label: 'Yellow to Red', preview: 'bg-gradient-to-r from-yellow-500 to-red-600' },
  { value: 'from-indigo-500 to-blue-600', label: 'Indigo to Blue', preview: 'bg-gradient-to-r from-indigo-500 to-blue-600' },
  { value: 'from-gray-500 to-slate-600', label: 'Gray to Slate', preview: 'bg-gradient-to-r from-gray-500 to-slate-600' },
]

const ICONS = ['🎓', '💻', '🔒', '🔓', '⚡', '🚀', '🎯', '🔥', '💡', '🛡️', '⚔️', '🔑']

export default function PathwayForm({
  initialData,
  onSave,
  onCancel,
  isEditing = false,
}: PathwayFormProps) {
  const [isSaving, setIsSaving] = useState(false)
  const [formData, setFormData] = useState<PathwayFormData>({
    title: '',
    subtitle: '',
    description: '',
    level: 'Beginner',
    duration_weeks: 4,
    color_gradient: COLOR_GRADIENTS[0].value,
    icon: ICONS[0],
    sort_order: 0,
    prerequisites: [],
    status: 'draft',
    ...initialData,
  })

  const updateField = <K extends keyof PathwayFormData>(
    field: K,
    value: PathwayFormData[K]
  ) => {
    setFormData((prev) => ({ ...prev, [field]: value }))
  }

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault()

    // Validation
    if (!formData.title.trim()) {
      alert('Please enter a pathway title')
      return
    }

    if (formData.duration_weeks < 1) {
      alert('Duration must be at least 1 week')
      return
    }

    setIsSaving(true)
    try {
      await onSave(formData)
    } catch (error: any) {
      alert(`Failed to save pathway: ${error.message}`)
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
            {isEditing ? 'Edit Pathway' : 'Create New Pathway'}
          </h2>
          <p className="text-sm text-text-secondary mt-1">
            Create a learning pathway with multiple modules
          </p>
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
            {isSaving ? 'Saving...' : isEditing ? 'Update Pathway' : 'Create Pathway'}
          </button>
        </div>
      </div>

      {/* Basic Information */}
      <div className="bg-bg-elevated border border-border-default rounded-lg p-6">
        <h3 className="text-lg font-semibold text-text-primary mb-4">Basic Information</h3>

        <div className="space-y-4">
          {/* Title */}
          <div>
            <label className="block text-sm font-medium text-text-secondary mb-2">
              Pathway Title *
            </label>
            <input
              type="text"
              value={formData.title}
              onChange={(e) => updateField('title', e.target.value)}
              placeholder="e.g., Python for Offensive Security"
              className="w-full px-4 py-2 bg-bg-primary border border-border-default rounded-lg text-text-primary focus:outline-none focus:ring-2 focus:ring-accent-primary"
              required
            />
          </div>

          {/* Subtitle */}
          <div>
            <label className="block text-sm font-medium text-text-secondary mb-2">
              Subtitle
            </label>
            <input
              type="text"
              value={formData.subtitle}
              onChange={(e) => updateField('subtitle', e.target.value)}
              placeholder="e.g., The Hacker's Swiss Army Knife"
              className="w-full px-4 py-2 bg-bg-primary border border-border-default rounded-lg text-text-primary focus:outline-none focus:ring-2 focus:ring-accent-primary"
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
              placeholder="Detailed description of what students will learn..."
              className="w-full px-4 py-2 bg-bg-primary border border-border-default rounded-lg text-text-primary focus:outline-none focus:ring-2 focus:ring-accent-primary"
              rows={4}
            />
          </div>

          {/* Level, Duration, Sort Order */}
          <div className="grid grid-cols-1 md:grid-cols-3 gap-4">
            <div>
              <label className="block text-sm font-medium text-text-secondary mb-2">
                Difficulty Level *
              </label>
              <select
                value={formData.level}
                onChange={(e) => updateField('level', e.target.value as any)}
                className="w-full px-4 py-2 bg-bg-primary border border-border-default rounded-lg text-text-primary focus:outline-none focus:ring-2 focus:ring-accent-primary"
              >
                <option value="Beginner">Beginner</option>
                <option value="Intermediate">Intermediate</option>
                <option value="Advanced">Advanced</option>
                <option value="Expert">Expert</option>
              </select>
            </div>

            <div>
              <label className="block text-sm font-medium text-text-secondary mb-2">
                Duration (weeks) *
              </label>
              <input
                type="number"
                value={formData.duration_weeks}
                onChange={(e) => updateField('duration_weeks', parseInt(e.target.value) || 1)}
                min="1"
                className="w-full px-4 py-2 bg-bg-primary border border-border-default rounded-lg text-text-primary focus:outline-none focus:ring-2 focus:ring-accent-primary"
              />
            </div>

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
              <p className="text-xs text-text-tertiary mt-1">Display order in pathway list</p>
            </div>
          </div>
        </div>
      </div>

      {/* Visual Styling */}
      <div className="bg-bg-elevated border border-border-default rounded-lg p-6">
        <h3 className="text-lg font-semibold text-text-primary mb-4">Visual Styling</h3>

        <div className="space-y-4">
          {/* Icon */}
          <div>
            <label className="block text-sm font-medium text-text-secondary mb-2">
              Icon
            </label>
            <div className="flex flex-wrap gap-2">
              {ICONS.map((icon) => (
                <button
                  key={icon}
                  type="button"
                  onClick={() => updateField('icon', icon)}
                  className={`w-12 h-12 flex items-center justify-center text-2xl rounded-lg border-2 transition-all ${
                    formData.icon === icon
                      ? 'border-accent-primary bg-accent-primary/20'
                      : 'border-border-default bg-bg-primary hover:border-accent-primary/50'
                  }`}
                >
                  {icon}
                </button>
              ))}
            </div>
          </div>

          {/* Color Gradient */}
          <div>
            <label className="block text-sm font-medium text-text-secondary mb-2">
              Color Gradient
            </label>
            <div className="grid grid-cols-2 md:grid-cols-4 gap-3">
              {COLOR_GRADIENTS.map((gradient) => (
                <button
                  key={gradient.value}
                  type="button"
                  onClick={() => updateField('color_gradient', gradient.value)}
                  className={`flex flex-col gap-2 p-3 rounded-lg border-2 transition-all ${
                    formData.color_gradient === gradient.value
                      ? 'border-accent-primary bg-bg-hover'
                      : 'border-border-default bg-bg-primary hover:border-accent-primary/50'
                  }`}
                >
                  <div className={`h-8 rounded ${gradient.preview}`} />
                  <span className="text-xs text-text-secondary">{gradient.label}</span>
                </button>
              ))}
            </div>
          </div>

          {/* Preview Card */}
          <div>
            <label className="block text-sm font-medium text-text-secondary mb-2">
              Preview
            </label>
            <div className={`p-6 rounded-lg bg-gradient-to-r ${formData.color_gradient} text-white`}>
              <div className="flex items-start gap-4">
                <div className="text-4xl">{formData.icon}</div>
                <div>
                  <h4 className="text-xl font-bold">
                    {formData.title || 'Pathway Title'}
                  </h4>
                  {formData.subtitle && (
                    <p className="text-sm opacity-90 mt-1">{formData.subtitle}</p>
                  )}
                  <div className="flex items-center gap-3 mt-3 text-sm">
                    <span className="bg-white/20 px-2 py-1 rounded">{formData.level}</span>
                    <span>{formData.duration_weeks} weeks</span>
                  </div>
                </div>
              </div>
            </div>
          </div>
        </div>
      </div>

      {/* Status */}
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
                Save as draft - only visible to you
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
    </form>
  )
}
