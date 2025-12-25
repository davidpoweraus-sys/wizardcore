'use client'

import { Edit, Trash2, Plus, Clock, CheckCircle, FileText, Download } from 'lucide-react'
import ExportButton from './ExportButton'

interface ContentItem {
  id: string
  title: string
  status: string
  created_at: string
  updated_at: string
  // Pathway-specific
  subtitle?: string
  level?: string
  duration_weeks?: number
  color_gradient?: string
  icon?: string
  // Module-specific
  pathway_id?: string
  sort_order?: number
  estimated_hours?: number
  xp_reward?: number
  // Exercise-specific
  module_id?: string
  difficulty?: string
  points?: number
}

interface ContentListProps {
  title: string
  items: ContentItem[]
  type: 'pathway' | 'module' | 'exercise'
  onEdit: (id: string) => void
  onDelete: (id: string) => void
  onCreate: () => void
  onExport?: (id: string, title: string) => void
}

export default function ContentList({
  title,
  items,
  type,
  onEdit,
  onDelete,
  onCreate,
  onExport,
}: ContentListProps) {
  const getStatusBadge = (status: string) => {
    const statusConfig = {
      draft: { 
        label: 'Draft', 
        color: 'bg-gray-500/20 text-gray-400 border-gray-500/30',
        icon: <FileText className="w-3 h-3" />
      },
      published: { 
        label: 'Published', 
        color: 'bg-green-500/20 text-green-400 border-green-500/30',
        icon: <CheckCircle className="w-3 h-3" />
      },
      under_review: { 
        label: 'Under Review', 
        color: 'bg-yellow-500/20 text-yellow-400 border-yellow-500/30',
        icon: <Clock className="w-3 h-3" />
      },
      archived: { 
        label: 'Archived', 
        color: 'bg-red-500/20 text-red-400 border-red-500/30',
        icon: <Trash2 className="w-3 h-3" />
      },
    }

    const config = statusConfig[status as keyof typeof statusConfig] || statusConfig.draft

    return (
      <span className={`inline-flex items-center gap-1 px-2 py-1 rounded-full text-xs font-medium border ${config.color}`}>
        {config.icon}
        {config.label}
      </span>
    )
  }

  const getDifficultyBadge = (difficulty: string) => {
    const difficultyConfig = {
      BEGINNER: { label: 'Beginner', color: 'bg-green-500/20 text-green-400' },
      INTERMEDIATE: { label: 'Intermediate', color: 'bg-yellow-500/20 text-yellow-400' },
      ADVANCED: { label: 'Advanced', color: 'bg-red-500/20 text-red-400' },
    }

    const config = difficultyConfig[difficulty as keyof typeof difficultyConfig] || difficultyConfig.BEGINNER

    return (
      <span className={`inline-flex items-center px-2 py-1 rounded-full text-xs font-medium ${config.color}`}>
        {config.label}
      </span>
    )
  }

  const formatDate = (dateString: string) => {
    return new Date(dateString).toLocaleDateString('en-US', {
      year: 'numeric',
      month: 'short',
      day: 'numeric',
    })
  }

  const renderPathwayCard = (item: ContentItem) => (
    <div
      key={item.id}
      className="bg-bg-elevated border border-border-default rounded-lg p-6 hover:shadow-lg transition-all"
    >
      <div className="flex items-start justify-between mb-4">
        <div className="flex items-start gap-4 flex-1">
          {item.color_gradient && item.icon && (
            <div className={`w-12 h-12 rounded-lg bg-gradient-to-r ${item.color_gradient} flex items-center justify-center text-2xl`}>
              {item.icon}
            </div>
          )}
          <div className="flex-1">
            <h3 className="text-lg font-semibold text-text-primary mb-1">{item.title}</h3>
            {item.subtitle && (
              <p className="text-sm text-text-secondary mb-2">{item.subtitle}</p>
            )}
            <div className="flex items-center gap-3 flex-wrap">
              {getStatusBadge(item.status)}
              {item.level && (
                <span className="text-xs text-text-tertiary">{item.level}</span>
              )}
              {item.duration_weeks && (
                <span className="text-xs text-text-tertiary">{item.duration_weeks} weeks</span>
              )}
            </div>
          </div>
        </div>
        <div className="flex items-center gap-2">
          {onExport && (
            <button
              onClick={() => onExport(item.id, item.title)}
              className="p-2 bg-bg-primary border border-border-default rounded-lg hover:bg-blue-500 hover:text-white transition-colors text-text-secondary"
              title="Export"
            >
              <Download className="w-4 h-4" />
            </button>
          )}
          <button
            onClick={() => onEdit(item.id)}
            className="p-2 bg-bg-primary border border-border-default rounded-lg hover:bg-accent-primary hover:text-white transition-colors text-text-secondary"
            title="Edit"
          >
            <Edit className="w-4 h-4" />
          </button>
          <button
            onClick={() => onDelete(item.id)}
            className="p-2 bg-bg-primary border border-border-default rounded-lg hover:bg-red-500 hover:text-white transition-colors text-text-secondary"
            title="Delete"
          >
            <Trash2 className="w-4 h-4" />
          </button>
        </div>
      </div>
      <div className="text-xs text-text-tertiary">
        Created: {formatDate(item.created_at)} • Updated: {formatDate(item.updated_at)}
      </div>
    </div>
  )

  const renderModuleCard = (item: ContentItem) => (
    <div
      key={item.id}
      className="bg-bg-elevated border border-border-default rounded-lg p-6 hover:shadow-lg transition-all"
    >
      <div className="flex items-start justify-between mb-4">
        <div className="flex-1">
          <div className="flex items-center gap-3 mb-2">
            <h3 className="text-lg font-semibold text-text-primary">{item.title}</h3>
            {item.sort_order !== undefined && (
              <span className="px-2 py-1 bg-bg-primary border border-border-default rounded text-xs text-text-tertiary">
                Order: {item.sort_order}
              </span>
            )}
          </div>
          <div className="flex items-center gap-3 flex-wrap">
            {getStatusBadge(item.status)}
            {item.estimated_hours && (
              <span className="text-xs text-text-tertiary">⏱️ {item.estimated_hours}h</span>
            )}
            {item.xp_reward !== undefined && (
              <span className="text-xs text-text-tertiary">⭐ {item.xp_reward} XP</span>
            )}
          </div>
        </div>
        <div className="flex items-center gap-2">
          <button
            onClick={() => onEdit(item.id)}
            className="p-2 bg-bg-primary border border-border-default rounded-lg hover:bg-accent-primary hover:text-white transition-colors text-text-secondary"
            title="Edit"
          >
            <Edit className="w-4 h-4" />
          </button>
          <button
            onClick={() => onDelete(item.id)}
            className="p-2 bg-bg-primary border border-border-default rounded-lg hover:bg-red-500 hover:text-white transition-colors text-text-secondary"
            title="Delete"
          >
            <Trash2 className="w-4 h-4" />
          </button>
        </div>
      </div>
      <div className="text-xs text-text-tertiary">
        Created: {formatDate(item.created_at)} • Updated: {formatDate(item.updated_at)}
      </div>
    </div>
  )

  const renderExerciseCard = (item: ContentItem) => (
    <div
      key={item.id}
      className="bg-bg-elevated border border-border-default rounded-lg p-6 hover:shadow-lg transition-all"
    >
      <div className="flex items-start justify-between mb-4">
        <div className="flex-1">
          <h3 className="text-lg font-semibold text-text-primary mb-2">{item.title}</h3>
          <div className="flex items-center gap-3 flex-wrap">
            {getStatusBadge(item.status)}
            {item.difficulty && getDifficultyBadge(item.difficulty)}
            {item.points !== undefined && (
              <span className="text-xs text-text-tertiary">🎯 {item.points} points</span>
            )}
          </div>
        </div>
        <div className="flex items-center gap-2">
          <button
            onClick={() => onEdit(item.id)}
            className="p-2 bg-bg-primary border border-border-default rounded-lg hover:bg-accent-primary hover:text-white transition-colors text-text-secondary"
            title="Edit"
          >
            <Edit className="w-4 h-4" />
          </button>
          <button
            onClick={() => onDelete(item.id)}
            className="p-2 bg-bg-primary border border-border-default rounded-lg hover:bg-red-500 hover:text-white transition-colors text-text-secondary"
            title="Delete"
          >
            <Trash2 className="w-4 h-4" />
          </button>
        </div>
      </div>
      <div className="text-xs text-text-tertiary">
        Created: {formatDate(item.created_at)} • Updated: {formatDate(item.updated_at)}
      </div>
    </div>
  )

  const renderCard = (item: ContentItem) => {
    switch (type) {
      case 'pathway':
        return renderPathwayCard(item)
      case 'module':
        return renderModuleCard(item)
      case 'exercise':
        return renderExerciseCard(item)
      default:
        return null
    }
  }

  return (
    <div>
      <div className="flex items-center justify-between mb-4">
        <div>
          <h2 className="text-xl font-bold text-text-primary">{title}</h2>
          <p className="text-sm text-text-secondary mt-1">
            {items.length} {type}{items.length !== 1 ? 's' : ''}
          </p>
        </div>
        <button
          onClick={onCreate}
          className="flex items-center gap-2 px-4 py-2 bg-accent-primary text-white rounded-lg hover:bg-accent-secondary transition-colors"
        >
          <Plus className="w-4 h-4" />
          Create {type}
        </button>
      </div>

      {items.length === 0 ? (
        <div className="bg-bg-elevated border border-border-default rounded-lg p-12 text-center">
          <p className="text-text-secondary mb-4">
            No {type}s created yet. Create your first {type} to get started!
          </p>
          <button
            onClick={onCreate}
            className="inline-flex items-center gap-2 px-6 py-3 bg-accent-primary text-white rounded-lg hover:bg-accent-secondary transition-colors"
          >
            <Plus className="w-5 h-5" />
            Create Your First {type}
          </button>
        </div>
      ) : (
        <div className="grid grid-cols-1 lg:grid-cols-2 gap-4">
          {items.map(renderCard)}
        </div>
      )}
    </div>
  )
}
