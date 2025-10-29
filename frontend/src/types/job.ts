export interface Job {
  url: string
  title: string
  company: string
  location: string
  salary_raw: string
  salary_min_usd: number | null
  salary_max_usd: number | null
  source: string
  posted_date: string
  discovered_date: string
  is_remote_us: boolean
  tags: string
}

