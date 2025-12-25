'use client'

import { useState, useRef } from 'react'
import { Upload, FileText, Check, X, Loader2 } from 'lucide-react'
import { api } from '@/lib/api'

interface ImportButtonProps {
  onImportComplete?: (pathwayId: string) => void
}

interface ImportPreview {
  pathway: {
    title: string
    subtitle?: string
    level: string
    duration_weeks: number
    modules_count: number
    exercises_count: number
  }
  isValid: boolean
  errors: string[]
}

export default function ImportButton({ onImportComplete }: ImportButtonProps) {
  const [isImporting, setIsImporting] = useState(false)
  const [isDragging, setIsDragging] = useState(false)
  const [preview, setPreview] = useState<ImportPreview | null>(null)
  const [importStatus, setImportStatus] = useState<'idle' | 'success' | 'error'>('idle')
  const [importError, setImportError] = useState<string | null>(null)
  const fileInputRef = useRef<HTMLInputElement>(null)

  const handleFileSelect = (file: File) => {
    const reader = new FileReader()
    reader.onload = (e) => {
      try {
        const content = e.target?.result as string
        const data = JSON.parse(content)

        // Validate basic structure
        const errors: string[] = []
        
        if (!data.pathway) {
          errors.push('Missing pathway data')
        } else {
          if (!data.pathway.title) errors.push('Pathway title is required')
          if (!data.pathway.level) errors.push('Pathway level is required')
          if (!data.pathway.duration_weeks) errors.push('Pathway duration is required')
          if (!data.pathway.modules || !Array.isArray(data.pathway.modules)) {
            errors.push('Pathway must have modules array')
          }
        }

        // Count modules and exercises
        let modulesCount = 0
        let exercisesCount = 0
        
        if (data.pathway?.modules) {
          modulesCount = data.pathway.modules.length
          data.pathway.modules.forEach((module: any) => {
            if (module.exercises && Array.isArray(module.exercises)) {
              exercisesCount += module.exercises.length
            }
          })
        }

        setPreview({
          pathway: {
            title: data.pathway?.title || 'Unknown',
            subtitle: data.pathway?.subtitle,
            level: data.pathway?.level || 'Unknown',
            duration_weeks: data.pathway?.duration_weeks || 0,
            modules_count: modulesCount,
            exercises_count: exercisesCount,
          },
          isValid: errors.length === 0,
          errors,
        })
      } catch (error) {
        setPreview({
          pathway: {
            title: 'Invalid JSON',
            level: 'Unknown',
            duration_weeks: 0,
            modules_count: 0,
            exercises_count: 0,
          },
          isValid: false,
          errors: ['Invalid JSON file'],
        })
      }
    }
    reader.readAsText(file)
  }

  const handleFileChange = (e: React.ChangeEvent<HTMLInputElement>) => {
    const file = e.target.files?.[0]
    if (file) {
      handleFileSelect(file)
    }
  }

  const handleDragOver = (e: React.DragEvent) => {
    e.preventDefault()
    setIsDragging(true)
  }

  const handleDragLeave = (e: React.DragEvent) => {
    e.preventDefault()
    setIsDragging(false)
  }

  const handleDrop = (e: React.DragEvent) => {
    e.preventDefault()
    setIsDragging(false)
    
    const file = e.dataTransfer.files?.[0]
    if (file && file.type === 'application/json') {
      handleFileSelect(file)
    } else {
      alert('Please drop a JSON file')
    }
  }

  const handleImport = async () => {
    if (!preview || !fileInputRef.current?.files?.[0]) return

    try {
      setIsImporting(true)
      setImportStatus('idle')
      setImportError(null)

      const file = fileInputRef.current.files[0]
      const content = await file.text()
      const data = JSON.parse(content)

      // Import the pathway
      const importedPathway = await api.post('/content-creator/pathways/import', {
        pathway: data.pathway || data, // Support both export format and raw pathway
        status: 'draft',
      })

      setImportStatus('success')
      
      // Reset after successful import
      setTimeout(() => {
        setPreview(null)
        setImportStatus('idle')
        if (fileInputRef.current) {
          fileInputRef.current.value = ''
        }
      }, 2000)

      if (onImportComplete) {
        onImportComplete(importedPathway.id)
      }
    } catch (error: any) {
      setImportStatus('error')
      setImportError(error.message || 'Failed to import pathway')
    } finally {
      setIsImporting(false)
    }
  }

  const handleCancel = () => {
    setPreview(null)
    setImportStatus('idle')
    setImportError(null)
    if (fileInputRef.current) {
      fileInputRef.current.value = ''
    }
  }

  return (
    <div className="relative">
      <button
        onClick={() => fileInputRef.current?.click()}
        className="flex items-center gap-2 px-4 py-2 bg-green-500 text-white rounded-lg hover:bg-green-600 transition-colors"
      >
        <Upload className="w-4 h-4" />
        Import Pathway
      </button>

      <input
        ref={fileInputRef}
        type="file"
        accept=".json,application/json"
        onChange={handleFileChange}
        className="hidden"
      />

      {(preview || importStatus !== 'idle') && (
        <div className="absolute top-full left-0 mt-2 w-96 bg-bg-elevated border border-border-default rounded-lg shadow-xl z-50">
          <div className="p-4">
            {importStatus === 'success' ? (
              <div className="text-center">
                <div className="w-12 h-12 bg-green-100 rounded-full flex items-center justify-center mx-auto mb-3">
                  <Check className="w-6 h-6 text-green-600" />
                </div>
                <h3 className="font-semibold text-text-primary mb-2">Import Successful!</h3>
                <p className="text-sm text-text-secondary">
                  Pathway imported successfully and saved as draft.
                </p>
              </div>
            ) : importStatus === 'error' ? (
              <div className="text-center">
                <div className="w-12 h-12 bg-red-100 rounded-full flex items-center justify-center mx-auto mb-3">
                  <X className="w-6 h-6 text-red-600" />
                </div>
                <h3 className="font-semibold text-text-primary mb-2">Import Failed</h3>
                <p className="text-sm text-red-500 mb-3">{importError}</p>
                <button
                  onClick={handleCancel}
                  className="px-4 py-2 bg-bg-primary border border-border-default rounded-lg hover:bg-bg-hover transition-colors"
                >
                  Try Again
                </button>
              </div>
            ) : preview ? (
              <>
                <div className="flex items-center justify-between mb-4">
                  <h3 className="font-semibold text-text-primary">Import Preview</h3>
                  <button
                    onClick={handleCancel}
                    className="p-1 hover:bg-bg-hover rounded"
                  >
                    <X className="w-4 h-4 text-text-secondary" />
                  </button>
                </div>

                {/* Validation Status */}
                <div className={`p-3 rounded-lg mb-4 ${preview.isValid ? 'bg-green-500/10 border border-green-500/20' : 'bg-red-500/10 border border-red-500/20'}`}>
                  <div className="flex items-center gap-2 mb-2">
                    {preview.isValid ? (
                      <>
                        <Check className="w-4 h-4 text-green-500" />
                        <span className="font-medium text-green-500">Valid JSON</span>
                      </>
                    ) : (
                      <>
                        <X className="w-4 h-4 text-red-500" />
                        <span className="font-medium text-red-500">Validation Errors</span>
                      </>
                    )}
                  </div>
                  
                  {preview.errors.length > 0 && (
                    <ul className="text-sm text-red-500 space-y-1">
                      {preview.errors.map((error, index) => (
                        <li key={index}>• {error}</li>
                      ))}
                    </ul>
                  )}
                </div>

                {/* Pathway Preview */}
                <div className="space-y-3 mb-4">
                  <div>
                    <h4 className="text-sm font-medium text-text-secondary mb-1">Pathway</h4>
                    <div className="bg-bg-primary border border-border-default rounded p-3">
                      <div className="font-medium text-text-primary">{preview.pathway.title}</div>
                      {preview.pathway.subtitle && (
                        <div className="text-sm text-text-secondary mt-1">{preview.pathway.subtitle}</div>
                      )}
                      <div className="flex items-center gap-3 mt-2 text-xs text-text-tertiary">
                        <span>Level: {preview.pathway.level}</span>
                        <span>Duration: {preview.pathway.duration_weeks} weeks</span>
                      </div>
                    </div>
                  </div>

                  <div className="grid grid-cols-2 gap-3">
                    <div className="bg-blue-500/10 border border-blue-500/20 rounded p-3">
                      <div className="text-2xl font-bold text-blue-500">{preview.pathway.modules_count}</div>
                      <div className="text-xs text-text-secondary">Modules</div>
                    </div>
                    <div className="bg-purple-500/10 border border-purple-500/20 rounded p-3">
                      <div className="text-2xl font-bold text-purple-500">{preview.pathway.exercises_count}</div>
                      <div className="text-xs text-text-secondary">Exercises</div>
                    </div>
                  </div>
                </div>

                {/* Import Button */}
                <button
                  onClick={handleImport}
                  disabled={!preview.isValid || isImporting}
                  className="w-full flex items-center justify-center gap-2 px-4 py-3 bg-green-500 text-white rounded-lg hover:bg-green-600 transition-colors disabled:opacity-50 disabled:cursor-not-allowed"
                >
                  {isImporting ? (
                    <>
                      <Loader2 className="w-4 h-4 animate-spin" />
                      Importing...
                    </>
                  ) : (
                    <>
                      <Upload className="w-4 h-4" />
                      Import as Draft
                    </>
                  )}
                </button>
              </>
            ) : (
              // Drag & Drop Area
              <div
                onDragOver={handleDragOver}
                onDragLeave={handleDragLeave}
                onDrop={handleDrop}
                className={`p-8 border-2 border-dashed rounded-lg text-center cursor-pointer transition-colors ${
                  isDragging 
                    ? 'border-accent-primary bg-accent-primary/5' 
                    : 'border-border-default hover:border-accent-primary/50'
                }`}
                onClick={() => fileInputRef.current?.click()}
              >
                <FileText className="w-12 h-12 text-text-tertiary mx-auto mb-3" />
                <p className="text-text-primary font-medium mb-1">Drop JSON file here</p>
                <p className="text-sm text-text-secondary">or click to browse</p>
                <p className="text-xs text-text-tertiary mt-3">Supports pathway export files</p>
              </div>
            )}
          </div>
        </div>
      )}
    </div>
  )
}
