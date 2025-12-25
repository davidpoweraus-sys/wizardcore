import { Suspense } from 'react'
import NewModuleContent from './NewModuleContent'

export const dynamic = 'force-dynamic'

export default function NewModulePage() {
  return (
    <Suspense fallback={
      <div className="min-h-screen bg-bg-primary flex items-center justify-center">
        <div className="text-center">
          <h1 className="text-2xl font-bold text-text-primary mb-2">
            Loading Module Creator...
          </h1>
          <p className="text-text-secondary">
            Please wait while we load the module creation interface.
          </p>
        </div>
      </div>
    }>
      <NewModuleContent />
    </Suspense>
  )
}