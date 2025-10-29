type FilterBarProps = {
  searchQuery: string
  setSearchQuery: (q: string) => void
  remoteOnly: boolean
  setRemoteOnly: (remote: boolean) => void
}

export function FilterBar({
  searchQuery,
  setSearchQuery,
  remoteOnly,
  setRemoteOnly,
}: FilterBarProps) {
  return (
    <div className="sticky top-0 z-20 -mx-4 sm:-mx-6 lg:-mx-8 px-4 sm:px-6 lg:px-8 py-3 backdrop-blur-xl bg-black/30 border-b border-white/10 mb-8">
      <div className="flex flex-col md:flex-row gap-4">
        <input
          type="search"
          placeholder="Search by title, company, keyword..."
          value={searchQuery}
          onChange={(e) => setSearchQuery(e.target.value)}
          className="flex-1 px-4 py-2 bg-white/5 border border-white/10 rounded-lg text-white placeholder-gray-400 focus:outline-none focus:ring-2 focus:ring-orange-500/50 focus:border-orange-500/50"
        />
        <label className="flex items-center gap-2 px-4 py-2 bg-white/5 border border-white/10 rounded-lg cursor-pointer hover:bg-white/10 transition-colors">
          <input
            type="checkbox"
            checked={remoteOnly}
            onChange={(e) => setRemoteOnly(e.target.checked)}
            className="w-4 h-4"
          />
          <span className="text-white">Remote-US Only</span>
        </label>
      </div>
    </div>
  )
}
