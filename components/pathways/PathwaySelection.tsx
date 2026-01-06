'use client'

import { useState, useEffect } from 'react'
import { Lock, Star, TrendingUp, Clock, Users, Loader2, AlertCircle } from 'lucide-react'
import { api } from '@/lib/api'
import { useRouter } from 'next/navigation'

interface Pathway {
  id: string
  title: string
  subtitle?: string
  description?: string
  level: string
  duration_weeks: number
  student_count: number
  rating: number
  module_count: number
  color_gradient?: string
  icon?: string
  is_locked: boolean
  sort_order: number
  prerequisites: string[]
  created_at: string
  updated_at: string
  is_enrolled?: boolean
  progress?: number
}

export default function PathwaySelection() {
  const router = useRouter()
  const [selected, setSelected] = useState<string | null>(null)
  const [pathways, setPathways] = useState<Pathway[]>([])
  const [loading, setLoading] = useState(true)
  const [error, setError] = useState<string | null>(null)
  const [enrolling, setEnrolling] = useState<string | null>(null)
  const [showEnrollModal, setShowEnrollModal] = useState(false)
  const [selectedPathway, setSelectedPathway] = useState<Pathway | null>(null)

  useEffect(() => {
    fetchPathways()
  }, [])

  const fetchPathways = async () => {
    try {
      setLoading(true)
      setError(null)
      const data = await api.get<{ pathways: Pathway[] }>('/pathways')
      setPathways(data.pathways || [])
    } catch (err: any) {
      console.error('Failed to fetch pathways:', err)
      setError(err.message || 'Failed to load pathways')
    } finally {
      setLoading(false)
    }
  }

  const handleEnrollClick = (pathway: Pathway) => {
    if (pathway.is_locked) return
    
    // Check prerequisites
    if (pathway.prerequisites && pathway.prerequisites.length > 0) {
      // TODO: Check if user has completed prerequisites
      // For now, just show a message
      alert(`This pathway requires completing: ${pathway.prerequisites.join(', ')}`)
      return
    }
    
    setSelectedPathway(pathway)
    setShowEnrollModal(true)
  }

  const handleEnrollConfirm = async () => {
    if (!selectedPathway) return
    
    try {
      setEnrolling(selectedPathway.id)
      await api.post(`/pathways/${selectedPathway.id}/enroll`)
      
      // Update local state
      setPathways(prev => prev.map(p => 
        p.id === selectedPathway.id 
          ? { ...p, is_enrolled: true, progress: 0 }
          : p
      ))
      
      // Close modal
      setShowEnrollModal(false)
      setSelectedPathway(null)
      
      // Redirect to first module
      // TODO: Get first module ID and redirect
      router.push('/dashboard/learning')
      
    } catch (err: any) {
      console.error('Failed to enroll:', err)
      alert(err.message || 'Failed to enroll in pathway')
    } finally {
      setEnrolling(null)
    }
  }

  const formatDuration = (weeks: number) => {
    return `${weeks} week${weeks !== 1 ? 's' : ''}`
  }

  const formatStudentCount = (count: number) => {
    if (count >= 1000) {
      return `${(count / 1000).toFixed(1)}k`
    }
    return count.toString()
  }

  const getColorClass = (colorGradient?: string) => {
    if (!colorGradient) return 'from-green-400 to-cyan-400'
    
    // Extract colors from gradient string
    if (colorGradient.includes('#667eea') && colorGradient.includes('#764ba2')) {
      return 'from-blue-500 to-purple-600'
    } else if (colorGradient.includes('#f093fb') && colorGradient.includes('#f5576c')) {
      return 'from-pink-500 to-red-500'
    } else if (colorGradient.includes('#4facfe') && colorGradient.includes('#00f2fe')) {
      return 'from-blue-400 to-cyan-400'
    }
    
    return 'from-green-400 to-cyan-400'
  }

  const getIcon = (icon?: string) => {
    if (!icon) return '📚'
    
    switch (icon) {
      case 'python': return '🐍'
      case 'c': return '⚙️'
      case 'javascript': return '🌐'
      default: return '📚'
    }
  }

  if (loading) {
    return (
      <div className="flex items-center justify-center h-64">
        <Loader2 className="w-8 h-8 animate-spin text-neon-cyan" />
        <span className="ml-3 text-text-secondary">Loading pathways...</span>
      </div>
    )
  }

  if (error) {
    return (
      <div className="flex flex-col items-center justify-center h-64">
        <AlertCircle className="w-12 h-12 text-red-400 mb-4" />
        <h3 className="text-lg font-semibold text-text-primary mb-2">Failed to load pathways</h3>
        <p className="text-text-secondary">{error}</p>
        <button 
          onClick={fetchPathways}
          className="mt-4 px-4 py-2 bg-neon-cyan text-black rounded-lg font-medium hover:bg-neon-cyan/90"
        >
          Try Again
        </button>
      </div>
    )
  }

  return (
    <div>
      <div className="mb-8">
        <h2 className="text-2xl font-bold text-text-primary">Available Pathways</h2>
        <p className="text-text-secondary mt-2">
          Each pathway is a complete course with projects, labs, and a capstone. Start with Python or jump to your skill level.
        </p>
      </div>

      <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-6">
        {pathways.map((pathway) => (
          <div
            key={pathway.id}
            className={`bg-bg-elevated border rounded-2xl p-6 flex flex-col h-full transition-all ${selected === pathway.id ? 'border-neon-cyan ring-2 ring-neon-cyan/20' : 'border-border-default hover:border-neon-lavender'} ${pathway.is_locked ? 'opacity-80' : ''}`}
            onClick={() => !pathway.is_locked && setSelected(pathway.id)}
          >
            <div className="flex justify-between items-start">
              <div>
                <div className={`w-14 h-14 rounded-xl bg-gradient-to-br ${getColorClass(pathway.color_gradient)} flex items-center justify-center text-2xl mb-4`}>
                  {getIcon(pathway.icon)}
                </div>
                <h3 className="text-xl font-bold text-text-primary">{pathway.title}</h3>
                <p className="text-sm text-text-secondary mt-1">{pathway.subtitle || 'Master offensive security skills'}</p>
              </div>
              {pathway.is_locked && (
                <div className="p-2 bg-bg-tertiary rounded-lg">
                  <Lock className="w-5 h-5 text-text-muted" />
                </div>
              )}
            </div>

            <p className="text-text-secondary mt-4 flex-1">{pathway.description || 'Learn essential security skills through hands-on exercises and projects.'}</p>

            <div className="mt-6 space-y-3">
              <div className="flex items-center justify-between text-sm">
                <div className="flex items-center text-text-secondary">
                  <Clock className="w-4 h-4 mr-2" />
                  <span>{formatDuration(pathway.duration_weeks)}</span>
                </div>
                <div className="flex items-center text-text-secondary">
                  <Users className="w-4 h-4 mr-2" />
                  <span>{formatStudentCount(pathway.student_count)} students</span>
                </div>
              </div>

              <div className="flex items-center justify-between text-sm">
                <div className="flex items-center">
                  <Star className="w-4 h-4 text-yellow-400 mr-1" />
                  <span className="font-medium text-text-primary">{pathway.rating.toFixed(1)}</span>
                  <span className="text-text-muted ml-1">rating</span>
                </div>
                <div className="px-3 py-1 rounded-full bg-bg-tertiary text-text-secondary text-xs font-medium capitalize">
                  {pathway.level}
                </div>
              </div>

              {pathway.is_enrolled && pathway.progress !== undefined && (
                <div className="pt-2">
                  <div className="flex justify-between text-xs text-text-secondary mb-1">
                    <span>Progress</span>
                    <span>{pathway.progress}%</span>
                  </div>
                  <div className="h-2 bg-bg-tertiary rounded-full overflow-hidden">
                    <div 
                      className="h-full bg-gradient-to-r from-neon-cyan to-neon-lavender rounded-full transition-all duration-300"
                      style={{ width: `${pathway.progress}%` }}
                    />
                  </div>
                </div>
              )}

              <div className="pt-4 border-t border-border-subtle">
                <div className="flex justify-between items-center">
                  <span className="text-sm text-text-secondary">{pathway.module_count} module{pathway.module_count !== 1 ? 's' : ''}</span>
                  {pathway.is_locked ? (
                    <button className="px-4 py-2 bg-bg-tertiary text-text-secondary rounded-lg text-sm font-medium cursor-not-allowed">
                      Locked
                    </button>
                  ) : pathway.is_enrolled ? (
                    <button 
                      onClick={() => router.push('/dashboard/learning')}
                      className="px-4 py-2 bg-gradient-to-r from-green-500 to-emerald-500 text-white rounded-lg text-sm font-medium hover:opacity-90 transition"
                    >
                      Continue Learning
                    </button>
                  ) : (
                    <button 
                      onClick={(e) => {
                        e.stopPropagation()
                        handleEnrollClick(pathway)
                      }}
                      disabled={enrolling === pathway.id}
                      className="px-4 py-2 bg-gradient-to-r from-neon-cyan to-neon-lavender text-white rounded-lg text-sm font-medium hover:opacity-90 transition disabled:opacity-50 disabled:cursor-not-allowed"
                    >
                      {enrolling === pathway.id ? (
                        <>
                          <Loader2 className="w-4 h-4 animate-spin inline mr-2" />
                          Enrolling...
                        </>
                      ) : (
                        'Enroll Now'
                      )}
                    </button>
                  )}
                </div>
              </div>
            </div>
          </div>
        ))}
      </div>

      {/* Enrollment Confirmation Modal */}
      {showEnrollModal && selectedPathway && (
        <div className="fixed inset-0 bg-black/50 flex items-center justify-center z-50 p-4">
          <div className="bg-bg-elevated border border-border-default rounded-2xl p-6 max-w-md w-full">
            <h3 className="text-xl font-bold text-text-primary mb-2">Enroll in {selectedPathway.title}</h3>
            <p className="text-text-secondary mb-6">
              Are you sure you want to enroll in this pathway? You'll have access to all modules, exercises, and projects.
            </p>
            
            <div className="space-y-4 mb-6">
              <div className="flex items-center text-sm">
                <Clock className="w-4 h-4 text-text-muted mr-3" />
                <span className="text-text-secondary">Duration: {formatDuration(selectedPathway.duration_weeks)}</span>
              </div>
              <div className="flex items-center text-sm">
                <Users className="w-4 h-4 text-text-muted mr-3" />
                <span className="text-text-secondary">{formatStudentCount(selectedPathway.student_count)} students enrolled</span>
              </div>
              <div className="flex items-center text-sm">
                <Star className="w-4 h-4 text-text-muted mr-3" />
                <span className="text-text-secondary">Rating: {selectedPathway.rating.toFixed(1)}/5.0</span>
              </div>
            </div>

            <div className="flex gap-3">
              <button
                onClick={() => {
                  setShowEnrollModal(false)
                  setSelectedPathway(null)
                }}
                className="flex-1 px-4 py-3 border border-border-default text-text-secondary rounded-lg font-medium hover:bg-bg-tertiary transition"
              >
                Cancel
              </button>
              <button
                onClick={handleEnrollConfirm}
                disabled={enrolling === selectedPathway.id}
                className="flex-1 px-4 py-3 bg-gradient-to-r from-neon-cyan to-neon-lavender text-white rounded-lg font-medium hover:opacity-90 transition disabled:opacity-50 disabled:cursor-not-allowed"
              >
                {enrolling === selectedPathway.id ? (
                  <>
                    <Loader2 className="w-4 h-4 animate-spin inline mr-2" />
                    Enrolling...
                  </>
                ) : (
                  'Confirm Enrollment'
                )}
              </button>
            </div>
          </div>
        </div>
      )}

      <div className="mt-10 bg-gradient-to-r from-bg-elevated to-bg-tertiary border border-border-default rounded-2xl p-8">
        <div className="flex flex-col md:flex-row items-center">
          <div className="p-4 rounded-xl bg-gradient-to-br from-neon-pink to-neon-purple mb-4 md:mb-0">
            <TrendingUp className="w-8 h-8 text-white" />
          </div>
          <div className="md:ml-6 flex-1">
            <h3 className="text-xl font-bold text-text-primary">Not sure where to start?</h3>
            <p className="text-text-secondary mt-2">
              Take our 5‑minute skill assessment to get a personalized pathway recommendation.
            </p>
          </div>
          <button className="mt-4 md:mt-0 md:ml-auto px-6 py-3 border border-neon-cyan text-neon-cyan rounded-lg font-medium hover:bg-neon-cyan/10 transition">
            Take Assessment
          </button>
        </div>
      </div>
    </div>
  )
}