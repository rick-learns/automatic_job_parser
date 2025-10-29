import { Building2, MapPin, DollarSign, Calendar, ExternalLink } from 'lucide-react'
import type { Job } from '../types/job'
import { formatDate, getJobTags } from '../lib/format'

type JobCardProps = {
  job: Job
}

export function JobCard({ job }: JobCardProps) {
  const tags = getJobTags(job.tags)
  const salaryDisplay = job.salary_raw || (job.salary_min_usd && job.salary_max_usd 
    ? `$${Math.round(job.salary_min_usd / 1000)}k–$${Math.round(job.salary_max_usd / 1000)}k` 
    : job.salary_min_usd 
      ? `$${Math.round(job.salary_min_usd / 1000)}k+` 
      : '')
  
  return (
    <div className="group relative bg-gradient-to-br from-gray-900/70 to-black/70 backdrop-blur-xl rounded-3xl shadow-2xl border border-white/10 p-8 hover:border-orange-500/50 transition-all duration-300 hover:shadow-orange-500/20">
      {/* Hover glow effect */}
      <div className="absolute inset-0 bg-gradient-to-r from-orange-500/0 via-purple-500/0 to-cyan-500/0 group-hover:from-orange-500/5 group-hover:via-purple-500/5 group-hover:to-cyan-500/5 rounded-3xl transition-all duration-300 pointer-events-none"></div>
      
      <div className="relative flex flex-col lg:flex-row lg:items-start lg:justify-between gap-6">
        <div className="flex-1">
          <div className="mb-5">
            <h3 className="text-3xl font-bold text-white mb-4 tracking-tight group-hover:text-transparent group-hover:bg-gradient-to-r group-hover:from-orange-400 group-hover:to-purple-400 group-hover:bg-clip-text transition-all">
              {job.title}
            </h3>
            <div className="flex items-center gap-3 text-gray-400">
              <Building2 className="w-5 h-5" />
              <span className="font-semibold text-gray-200 text-lg">{job.company}</span>
              <span className="text-gray-600">•</span>
              <span className="text-sm px-3 py-1.5 bg-white/5 rounded-lg font-semibold border border-white/10">
                {job.source}
              </span>
            </div>
          </div>
          <div className="flex flex-wrap gap-6 text-sm text-gray-400 mb-5">
            <div className="flex items-center gap-2">
              <MapPin className="w-5 h-5 text-purple-400" />
              <span className="text-gray-300 text-base">{job.location}</span>
              {job.is_remote_us && (
                <span className="ml-1 px-3 py-1 bg-gradient-to-r from-cyan-500/20 to-cyan-600/20 text-cyan-300 text-xs rounded-full font-bold border border-cyan-500/30">
                  Remote US
                </span>
              )}
            </div>
            {salaryDisplay && (
              <div className="flex items-center gap-2">
                <DollarSign className="w-5 h-5 text-green-400" />
                <span className="font-bold text-white text-base">
                  {salaryDisplay}
                </span>
              </div>
            )}
            <div className="flex items-center gap-2">
              <Calendar className="w-5 h-5 text-orange-400" />
              <span className="text-gray-300 text-base">Added {formatDate(job.discovered_date)}</span>
            </div>
          </div>

          {tags.length > 0 && (
            <div className="flex flex-wrap gap-2">
              {tags.map((tag, i) => (
                <span
                  key={i}
                  className="px-3 py-1.5 bg-gradient-to-r from-purple-500/10 to-orange-500/10 text-purple-300 text-sm rounded-xl font-semibold border border-purple-500/20 hover:border-purple-500/40 transition-colors"
                >
                  {tag}
                </span>
              ))}
            </div>
          )}
        </div>
        
        <a
          href={job.url}
          target="_blank"
          rel="noopener noreferrer"
          className="flex items-center justify-center gap-2 px-8 py-4 bg-gradient-to-r from-orange-500 to-orange-600 hover:from-orange-600 hover:to-orange-700 text-white rounded-2xl transition-all font-bold whitespace-nowrap shadow-lg shadow-orange-500/30 text-base border border-orange-400/20 hover:shadow-orange-500/50 hover:scale-105"
        >
          Apply Now
          <ExternalLink className="w-5 h-5" />
        </a>
      </div>
    </div>
  )
}
