# Jobsite Frontend

Modern React frontend for the jobsite application, built with Vite, TypeScript, and Tailwind CSS.

## Development

```bash
npm install
npm run dev
```

The dev server will start on `http://localhost:5173` with hot module replacement.

## Building

```bash
npm run build
```

This builds the frontend and outputs to `../jobsite_golang/public/latest/`. The build is configured to preserve existing `jobs.json` and `jobs.csv` files.

## Features

- **Client-side filtering**: Search by keywords, filter by remote-US only
- **URL state management**: Share filtered views via URL parameters
- **Debounced search**: Optimized performance with 300ms debounce
- **Responsive design**: Mobile-first, works on all screen sizes
- **Dark theme**: Modern gradient design
- **Accessibility**: Keyboard navigation, ARIA labels, semantic HTML

## Data Format

The frontend reads `/latest/jobs.json` which must follow this structure:

```typescript
interface Job {
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
```

## Build Process

1. Go backend generates `jobs.json` and `jobs.csv` in `public/latest/`
2. React frontend builds to the same directory
3. `emptyOutDir: false` in vite.config preserves data files

## Tech Stack

- React 19
- TypeScript
- Vite
- Tailwind CSS
- Lucide React (icons)
