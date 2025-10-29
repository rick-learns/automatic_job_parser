import { useEffect, useState, useMemo } from 'react'
import type { Job } from './types/job'
import { FilterBar } from './components/FilterBar'
import { JobList } from './components/JobList'
import { LoadingState } from './components/LoadingState'
import { EmptyState } from './components/EmptyState'
import { ErrorState } from './components/ErrorState'
import { useDebounce } from './lib/useDebounce'

function App() {
  const [jobs, setJobs] = useState<Job[]>([])
  const [loading, setLoading] = useState(true)
  const [error, setError] = useState<string | null>(null)
  const [searchQuery, setSearchQuery] = useState('')
  const [remoteOnly, setRemoteOnly] = useState(false)

  // Debounce search query
  const debouncedSearchQuery = useDebounce(searchQuery, 300)

  // Fetch jobs on mount
  useEffect(() => {
    const controller = new AbortController()

    const fetchJobs = async () => {
      try {
        setLoading(true)
        setError(null)

        const res = await fetch('./jobs.json', {
          cache: 'no-store',
          signal: controller.signal,
        })

        if (!res.ok) {
          throw new Error(`HTTP ${res.status}`)
        }

        const data = await res.json()
        setJobs(Array.isArray(data) ? data : [])
      } catch (err) {
        if ((err as any).name !== 'AbortError') {
          setError('Failed to load jobs')
          console.error('Fetch error:', err)
        }
      } finally {
        setLoading(false)
      }
    }

    fetchJobs()

    return () => {
      controller.abort()
    }
  }, [])

  // Initialize filters from URL on mount
  useEffect(() => {
    const params = new URLSearchParams(window.location.search)
    const qParam = params.get('q')
    const remoteParam = params.get('remote')

    if (qParam) setSearchQuery(qParam)
    if (remoteParam === '1') setRemoteOnly(true)
  }, [])

  // Update URL when filters change
  useEffect(() => {
    const params = new URLSearchParams()
    if (debouncedSearchQuery) params.set('q', debouncedSearchQuery)
    if (remoteOnly) params.set('remote', '1')

    const newUrl = `${window.location.pathname}${params.toString() ? `?${params.toString()}` : ''}`
    window.history.replaceState(null, '', newUrl)
  }, [debouncedSearchQuery, remoteOnly])

  // Precompute searchable text once per job
  const searchableJobs = useMemo(() => {
    return jobs.map((job) => {
      const searchableText = [
        job.title,
        job.company,
        job.location,
        job.tags,
        job.source,
        job.salary_raw,
      ]
        .filter(Boolean)
        .join(' ')
        .toLowerCase()

      return { job, searchableText }
    })
  }, [jobs])

  // Filter jobs
  const filteredJobs = useMemo(() => {
    const needle = debouncedSearchQuery.trim().toLowerCase()

    return searchableJobs
      .filter(({ job, searchableText }) => {
        // Remote filter
        if (remoteOnly && !job.is_remote_us) return false

        // Search filter
        if (needle && !searchableText.includes(needle)) return false

        return true
      })
      .map(({ job }) => job)
  }, [searchableJobs, debouncedSearchQuery, remoteOnly])

  const handleRetry = () => {
    window.location.reload()
  }

  return (
    <div className="min-h-screen bg-gradient-to-br from-gray-950 via-gray-900 to-black text-white font-sans relative overflow-hidden">
      {/* Animated gradient overlay blobs */}
      <div className="pointer-events-none absolute inset-0 -z-10">
        <div className="absolute -top-24 -left-16 h-80 w-80 rounded-full bg-purple-500 opacity-25 blur-3xl mix-blend-multiply animate-pulse4" />
        <div className="absolute top-1/3 -right-20 h-96 w-96 rounded-full bg-orange-500 opacity-25 blur-3xl mix-blend-multiply animate-pulse5" />
        <div className="absolute -bottom-24 left-1/4 h-72 w-72 rounded-full bg-cyan-500 opacity-25 blur-3xl mix-blend-multiply animate-pulse6" />
      </div>

      <div className="relative z-10">
        <header className="border-b border-white/10">
          <div className="max-w-7xl mx-auto px-4 py-6">
            <h1 className="text-3xl font-bold">QA/SDET Roles</h1>
            <p className="text-gray-400 mt-1">Remote US + Wichita, KS</p>
          </div>
        </header>

        <main className="max-w-7xl mx-auto px-4 py-8">
          <FilterBar
            searchQuery={searchQuery}
            setSearchQuery={setSearchQuery}
            remoteOnly={remoteOnly}
            setRemoteOnly={setRemoteOnly}
          />

          <div className="mb-4">
            <p className="text-gray-400">
              {loading ? 'Loading...' : `${filteredJobs.length} jobs found`}
            </p>
          </div>

          {loading && <LoadingState />}
          {error && !loading && <ErrorState error={error} onRetry={handleRetry} />}
          {!loading && !error && filteredJobs.length === 0 && <EmptyState />}
          {!loading && !error && filteredJobs.length > 0 && <JobList jobs={filteredJobs} />}
        </main>
      </div>
    </div>
  )
}

export default App
