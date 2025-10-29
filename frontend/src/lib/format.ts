import type { Job } from '../types/job'

export const formatSalary = (job: Job): string => {
  if (job.salary_min_usd && job.salary_max_usd) {
    return `$${Math.round(job.salary_min_usd / 1000)}k-$${Math.round(job.salary_max_usd / 1000)}k`
  }
  if (job.salary_min_usd) {
    return `$${Math.round(job.salary_min_usd / 1000)}k+`
  }
  return job.salary_raw || ''
}

export const formatDate = (date: string): string => {
  try {
    const d = new Date(date)
    return d.toLocaleDateString('en-US', { month: 'short', day: 'numeric', year: 'numeric' })
  } catch {
    return date
  }
}

export const getJobTags = (tags: string): string[] => {
  return tags ? tags.split(',').filter(Boolean).map(t => t.trim()) : []
}

