export function EmptyState() {
  return (
    <div className="flex items-center justify-center min-h-[400px]">
      <div className="text-center">
        <p className="text-gray-400 text-lg">No jobs matched your filters</p>
        <p className="text-gray-500 text-sm mt-2">Try adjusting your search or filters</p>
      </div>
    </div>
  )
}

