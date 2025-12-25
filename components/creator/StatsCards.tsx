'use client'

import { 
  BookOpen, 
  Code, 
  FileText, 
  Users, 
  Star, 
  Eye, 
  CheckCircle, 
  Clock,
  TrendingUp,
  FileCheck 
} from 'lucide-react'

interface CreatorStats {
  total_pathways: number
  total_modules: number
  total_exercises: number
  total_students: number
  average_rating: number
  total_ratings: number
  total_views: number
  total_enrollments: number
  total_completions: number
  completion_rate: number
  pending_reviews: number
  published_content: number
  draft_content: number
}

interface StatsCardsProps {
  stats: CreatorStats
}

interface StatCard {
  label: string
  value: string | number
  icon: React.ReactNode
  color: string
  bgColor: string
}

export default function StatsCards({ stats }: StatsCardsProps) {
  const cards: StatCard[] = [
    {
      label: 'Total Pathways',
      value: stats.total_pathways,
      icon: <BookOpen className="w-6 h-6" />,
      color: 'text-blue-500',
      bgColor: 'bg-blue-500/10',
    },
    {
      label: 'Total Modules',
      value: stats.total_modules,
      icon: <FileText className="w-6 h-6" />,
      color: 'text-green-500',
      bgColor: 'bg-green-500/10',
    },
    {
      label: 'Total Exercises',
      value: stats.total_exercises,
      icon: <Code className="w-6 h-6" />,
      color: 'text-purple-500',
      bgColor: 'bg-purple-500/10',
    },
    {
      label: 'Total Students',
      value: stats.total_students,
      icon: <Users className="w-6 h-6" />,
      color: 'text-orange-500',
      bgColor: 'bg-orange-500/10',
    },
    {
      label: 'Average Rating',
      value: stats.average_rating ? `${stats.average_rating.toFixed(1)} ⭐` : 'N/A',
      icon: <Star className="w-6 h-6" />,
      color: 'text-yellow-500',
      bgColor: 'bg-yellow-500/10',
    },
    {
      label: 'Total Views',
      value: stats.total_views,
      icon: <Eye className="w-6 h-6" />,
      color: 'text-cyan-500',
      bgColor: 'bg-cyan-500/10',
    },
    {
      label: 'Enrollments',
      value: stats.total_enrollments,
      icon: <TrendingUp className="w-6 h-6" />,
      color: 'text-indigo-500',
      bgColor: 'bg-indigo-500/10',
    },
    {
      label: 'Completions',
      value: stats.total_completions,
      icon: <CheckCircle className="w-6 h-6" />,
      color: 'text-teal-500',
      bgColor: 'bg-teal-500/10',
    },
    {
      label: 'Completion Rate',
      value: stats.completion_rate ? `${stats.completion_rate.toFixed(0)}%` : '0%',
      icon: <TrendingUp className="w-6 h-6" />,
      color: 'text-pink-500',
      bgColor: 'bg-pink-500/10',
    },
    {
      label: 'Published Content',
      value: stats.published_content,
      icon: <FileCheck className="w-6 h-6" />,
      color: 'text-green-600',
      bgColor: 'bg-green-600/10',
    },
    {
      label: 'Draft Content',
      value: stats.draft_content,
      icon: <FileText className="w-6 h-6" />,
      color: 'text-gray-500',
      bgColor: 'bg-gray-500/10',
    },
    {
      label: 'Pending Reviews',
      value: stats.pending_reviews,
      icon: <Clock className="w-6 h-6" />,
      color: 'text-amber-500',
      bgColor: 'bg-amber-500/10',
    },
  ]

  return (
    <div>
      <h2 className="text-xl font-bold text-text-primary mb-4">Statistics Overview</h2>
      <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 xl:grid-cols-4 gap-4">
        {cards.map((card, index) => (
          <div
            key={index}
            className="bg-bg-elevated border border-border-default rounded-lg p-6 hover:shadow-lg transition-shadow"
          >
            <div className="flex items-start justify-between">
              <div className="flex-1">
                <p className="text-sm text-text-secondary mb-1">{card.label}</p>
                <p className="text-2xl font-bold text-text-primary">{card.value}</p>
              </div>
              <div className={`p-3 rounded-lg ${card.bgColor}`}>
                <div className={card.color}>{card.icon}</div>
              </div>
            </div>
          </div>
        ))}
      </div>
    </div>
  )
}
