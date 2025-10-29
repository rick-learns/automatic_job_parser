import type { Job } from '../types/job'
import { JobCard } from './JobCard'

type JobListProps = {
  jobs: Job[]
}

export function JobList({ jobs }: JobListProps) {
  if (jobs.length === 0) {
    return null
  }

  return (
    <div className="grid grid-cols-1 gap-4 sm:gap-6 lg:gap-7 xl:gap-8">
      {jobs.map((job) => (
        <JobCard key={job.url} job={job} />
      ))}
    </div>
  )
}
