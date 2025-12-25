'use client'

import { useEffect, useState } from 'react'
import { useRouter, useParams } from 'next/navigation'
import ExerciseBuilder from '@/components/creator/ExerciseBuilder'
import { api } from '@/lib/api'
import { Loader2, ArrowLeft } from 'lucide-react'

interface ExerciseWithTests {
  exercise: {
    id: string
    module_id: string
    title: string
    difficulty: string
    points: number
    time_limit_minutes?: number
    sort_order: number
    objectives: string[]
    content?: string
    examples: Record<string, any>
    description?: string
    constraints: string[]
    hints: string[]
    starter_code?: string
    solution_code?: string
    language_id: number
    tags: string[]
    status: string
    created_at: string
    updated_at: string
  }
  test_cases: Array<{
    id: string
    exercise_id: string
    input?: string
    expected_output: string
    is_hidden: boolean
    points: number
    sort_order: number
    created_at: string
  }>
}

interface Module {
  id: string
  title: string
  pathway_id: string
}

interface Pathway {
  id: string
  title: string
}

export default function EditExercisePage() {
  const router = useRouter()
  const params = useParams()
  const exerciseId = params.id as string

  const [loading, setLoading] = useState(true)
  const [exerciseData, setExerciseData] = useState<ExerciseWithTests | null>(null)
  const [module, setModule] = useState<Module | null>(null)
  const [pathway, setPathway] = useState<Pathway | null>(null)
  const [error, setError] = useState<string | null>(null)

  useEffect(() => {
    loadExercise()
  }, [exerciseId])

  const loadExercise = async () => {
    try {
      setLoading(true)
      setError(null)

      // Get exercise with test cases
      const data = await api.get<ExerciseWithTests>(`/content-creator/exercises/${exerciseId}`)
      setExerciseData(data)

      // Get the module for this exercise
      const modules = await api.get<Module[]>('/content-creator/modules')
      const foundModule = modules.find(m => m.id === data.exercise.module_id)
      setModule(foundModule || null)

      if (foundModule) {
        // Get the pathway for this module
        const pathways = await api.get<Pathway[]>('/content-creator/pathways')
        const foundPathway = pathways.find(p => p.id === foundModule.pathway_id)
        setPathway(foundPathway || null)
      }
    } catch (err: any) {
      console.error('Failed to load exercise:', err)
      setError(err.message || 'Failed to load exercise')
    } finally {
      setLoading(false)
    }
  }

  const handleSave = async (formData: any) => {
    try {
      const payload = {
        title: formData.title,
        difficulty: formData.difficulty,
        points: formData.points,
        time_limit_minutes: formData.time_limit_minutes || null,
        sort_order: formData.sort_order,
        objectives: formData.objectives || [],
        content: formData.content || null,
        examples: formData.examples || {},
        description: formData.description || null,
        constraints: formData.constraints || [],
        hints: formData.hints || [],
        starter_code: formData.starter_code || null,
        solution_code: formData.solution_code || null,
        language_id: formData.language_id,
        tags: formData.tags || [],
        status: formData.status,
        test_cases: formData.test_cases.map((tc: any) => ({
          input: tc.input || null,
          expected_output: tc.expected_output,
          is_hidden: tc.is_hidden,
          points: tc.points,
          sort_order: tc.sort_order,
        })),
      }

      await api.put(`/content-creator/exercises/${exerciseId}`, payload)
      
      // Redirect to dashboard with success message
      router.push('/creator/dashboard')
      // In a real app, you might want to show a toast notification
    } catch (error: any) {
      alert(`Failed to update exercise: ${error.message}`)
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
          <p className="text-text-secondary">Loading exercise...</p>
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

  if (!exerciseData) {
    return (
      <div className="flex items-center justify-center min-h-screen bg-bg-primary">
        <div className="max-w-md p-6 bg-bg-elevated border border-border-default rounded-lg">
          <h2 className="text-xl font-bold text-text-primary mb-2">Exercise Not Found</h2>
          <p className="text-text-secondary mb-4">The exercise you're trying to edit doesn't exist.</p>
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

  const { exercise, test_cases } = exerciseData

  // Convert backend data to form data
  const initialData: any = {
    module_id: exercise.module_id,
    title: exercise.title,
    difficulty: exercise.difficulty as 'BEGINNER' | 'INTERMEDIATE' | 'ADVANCED',
    points: exercise.points,
    time_limit_minutes: exercise.time_limit_minutes || undefined,
    sort_order: exercise.sort_order,
    objectives: exercise.objectives || [],
    content: exercise.content || '',
    examples: exercise.examples || {},
    description: exercise.description || '',
    constraints: exercise.constraints || [],
    hints: exercise.hints || [],
    starter_code: exercise.starter_code || '',
    solution_code: exercise.solution_code || '',
    language_id: exercise.language_id,
    tags: exercise.tags || [],
    status: exercise.status as 'draft' | 'published',
    test_cases: test_cases.map(tc => ({
      id: tc.id,
      input: tc.input || '',
      expected_output: tc.expected_output,
      is_hidden: tc.is_hidden,
      points: tc.points,
      sort_order: tc.sort_order,
    })),
  }

  return (
    <div className="min-h-screen bg-bg-primary">
      <div className="max-w-7xl mx-auto px-4 py-8">
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
              <h1 className="text-3xl font-bold text-text-primary">Edit Exercise</h1>
              <p className="text-text-secondary mt-2">
                Update your exercise details, code, and test cases
              </p>
              <div className="mt-2 space-y-1">
                {module && (
                  <div>
                    <span className="text-sm text-text-tertiary">Module: </span>
                    <span className="text-sm text-text-secondary">{module.title}</span>
                  </div>
                )}
                {pathway && (
                  <div>
                    <span className="text-sm text-text-tertiary">Pathway: </span>
                    <span className="text-sm text-text-secondary">{pathway.title}</span>
                  </div>
                )}
              </div>
            </div>
          </div>
        </div>

        {/* Form */}
        <div className="bg-bg-elevated border border-border-default rounded-lg p-6">
          <ExerciseBuilder
            moduleId={exercise.module_id}
            initialData={initialData}
            onSave={handleSave}
          />
        </div>

        {/* Last Updated Info */}
        <div className="mt-6 text-sm text-text-tertiary">
          <p>
            Created: {new Date(exercise.created_at).toLocaleDateString()} • 
            Last Updated: {new Date(exercise.updated_at).toLocaleDateString()} • 
            Test Cases: {test_cases.length}
          </p>
        </div>
      </div>
    </div>
  )
}
