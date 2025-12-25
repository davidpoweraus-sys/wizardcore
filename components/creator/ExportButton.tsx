'use client'

import { useState } from 'react'
import { Download, Check, Copy, Loader2 } from 'lucide-react'
import { api } from '@/lib/api'

interface ExportButtonProps {
  pathwayId: string
  pathwayTitle: string
  onExportComplete?: () => void
}

export default function ExportButton({ pathwayId, pathwayTitle, onExportComplete }: ExportButtonProps) {
  const [isExporting, setIsExporting] = useState(false)
  const [isCopied, setIsCopied] = useState(false)
  const [exportData, setExportData] = useState<string | null>(null)

  const handleExport = async () => {
    try {
      setIsExporting(true)
      setExportData(null)
      setIsCopied(false)

      const data = await api.get(`/content-creator/pathways/${pathwayId}/export`)
      
      // Format the JSON nicely
      const formattedJson = JSON.stringify(data, null, 2)
      setExportData(formattedJson)

      // Create a downloadable file
      const blob = new Blob([formattedJson], { type: 'application/json' })
      const url = URL.createObjectURL(blob)
      const a = document.createElement('a')
      a.href = url
      a.download = `${pathwayTitle.replace(/[^a-z0-9]/gi, '_').toLowerCase()}_export.json`
      document.body.appendChild(a)
      a.click()
      document.body.removeChild(a)
      URL.revokeObjectURL(url)

      if (onExportComplete) {
        onExportComplete()
      }
    } catch (error: any) {
      alert(`Failed to export pathway: ${error.message}`)
    } finally {
      setIsExporting(false)
    }
  }

  const handleCopyToClipboard = async () => {
    if (!exportData) return

    try {
      await navigator.clipboard.writeText(exportData)
      setIsCopied(true)
      setTimeout(() => setIsCopied(false), 2000)
    } catch (error) {
      alert('Failed to copy to clipboard')
    }
  }

  return (
    <div className="relative">
      <button
        onClick={handleExport}
        disabled={isExporting}
        className="flex items-center gap-2 px-4 py-2 bg-blue-500 text-white rounded-lg hover:bg-blue-600 transition-colors disabled:opacity-50 disabled:cursor-not-allowed"
      >
        {isExporting ? (
          <>
            <Loader2 className="w-4 h-4 animate-spin" />
            Exporting...
          </>
        ) : (
          <>
            <Download className="w-4 h-4" />
            Export Pathway
          </>
        )}
      </button>

      {exportData && (
        <div className="absolute top-full left-0 mt-2 w-96 bg-bg-elevated border border-border-default rounded-lg shadow-xl z-50">
          <div className="p-4">
            <div className="flex items-center justify-between mb-3">
              <h3 className="font-semibold text-text-primary">Export Complete</h3>
              <button
                onClick={handleCopyToClipboard}
                className="flex items-center gap-1 px-3 py-1 text-sm bg-bg-primary border border-border-default rounded hover:bg-bg-hover transition-colors"
              >
                {isCopied ? (
                  <>
                    <Check className="w-3 h-3 text-green-500" />
                    Copied!
                  </>
                ) : (
                  <>
                    <Copy className="w-3 h-3" />
                    Copy JSON
                  </>
                )}
              </button>
            </div>
            <div className="text-xs text-text-secondary mb-2">
              Pathway exported successfully. You can also copy the JSON.
            </div>
            <div className="max-h-60 overflow-y-auto bg-bg-primary border border-border-default rounded p-3">
              <pre className="text-xs text-text-secondary whitespace-pre-wrap">
                {exportData.substring(0, 500)}...
              </pre>
            </div>
          </div>
        </div>
      )}
    </div>
  )
}
