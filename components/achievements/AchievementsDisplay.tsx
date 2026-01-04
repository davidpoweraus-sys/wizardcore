'use client'

import { useState, useEffect } from 'react'
import { Trophy, Star, Zap, Target, Flame, Shield, Lock, Check, Calendar, Infinity, Sun, Moon, Share, Sword, Brain, Rocket, Heart, MessageSquare, Compass, Gem, AlertCircle, Loader2, BookOpen, Clock, TrendingUp, Users } from 'lucide-react'
import { api } from '@/lib/api'

interface Achievement {
  id: string
  title: string
  description?: string
  icon?: string
  color_gradient?: string
  rarity: string
  xp_reward: number
  criteria_type: string
  criteria_value?: number
  criteria_metadata?: Record<string, any>
  is_hidden: boolean
  sort_order: number
  created_at: string
  earned: boolean
  progress: number
  earned_date?: string
}

interface UserStats {
  earned_count: number
  total_count: number
  total_xp: number
  global_rank?: number
  streak_days: number
}

interface LeaderboardEntry {
  rank: number
  name: string
  xp: number
  streak: number
  badges: number
  is_current_user?: boolean
}

export default function AchievementsDisplay() {
  const [achievements, setAchievements] = useState<Achievement[]>([])
  const [stats, setStats] = useState<UserStats | null>(null)
  const [leaderboard, setLeaderboard] = useState<LeaderboardEntry[]>([])
  const [loading, setLoading] = useState(true)
  const [error, setError] = useState<string | null>(null)
  const [unlockedAchievement, setUnlockedAchievement] = useState<Achievement | null>(null)

  useEffect(() => {
    fetchAchievementsData()
  }, [])

  const fetchAchievementsData = async () => {
    try {
      setLoading(true)
      setError(null)
      
      // Fetch achievements
      const achievementsResponse = await api.get<{ achievements: Achievement[] }>('/users/me/achievements')
      // Handle null achievements from API
      const achievements = achievementsResponse.achievements || []
      setAchievements(achievements)
      
      // Calculate stats
      const earnedCount = achievements.filter(a => a.earned).length
      const totalCount = achievements.length
      const totalXP = achievements
        .filter(a => a.earned)
        .reduce((sum, a) => sum + a.xp_reward, 0)
      
      setStats({
        earned_count: earnedCount,
        total_count: totalCount,
        total_xp: totalXP,
        streak_days: 7, // TODO: Get from user stats endpoint
        global_rank: 7 // TODO: Get from leaderboard endpoint
      })
      
      // TODO: Fetch leaderboard from API
      // For now, use mock data that will be replaced
      setLeaderboard([
        { rank: 1, name: 'Alex Chen', xp: 12540, streak: 21, badges: 12 },
        { rank: 2, name: 'Sam Rivera', xp: 11200, streak: 18, badges: 10 },
        { rank: 3, name: 'Jordan Lee', xp: 9850, streak: 14, badges: 9 },
        { rank: 4, name: 'Taylor Kim', xp: 8760, streak: 12, badges: 8 },
        { rank: 5, name: 'Casey Morgan', xp: 7540, streak: 9, badges: 7 },
        { rank: 6, name: 'Riley Patel', xp: 6320, streak: 7, badges: 6 },
        { rank: 7, name: 'You', xp: totalXP, streak: 7, badges: earnedCount, is_current_user: true },
        { rank: 8, name: 'Drew Zhang', xp: 2450, streak: 5, badges: 3 },
      ])
      
    } catch (err: any) {
      console.error('Failed to fetch achievements data:', err)
      setError(err.message || 'Failed to load achievements')
    } finally {
      setLoading(false)
    }
  }

  const getIconComponent = (iconName?: string) => {
    if (!iconName) return Star
    
    const iconMap: Record<string, any> = {
      'star': Star,
      'zap': Zap,
      'flame': Flame,
      'bullseye': Target,
      'shield': Shield,
      'lock': Lock,
      'check': Check,
      'calendar': Calendar,
      'infinity': Infinity,
      'sun': Sun,
      'moon': Moon,
      'share': Share,
      'swords': Sword,
      'brain': Brain,
      'rocket': Rocket,
      'heart': Heart,
      'message': MessageSquare,
      'compass': Compass,
      'treasure': Gem,
      'trophy': Trophy,
      'python': Trophy, // Fallback for python icon
    }
    
    return iconMap[iconName] || Star
  }

  const getColorClass = (colorGradient?: string) => {
    if (!colorGradient) return 'from-yellow-400 to-orange-400'
    
    // Extract colors from gradient string
    if (colorGradient.includes('#ffd166') && colorGradient.includes('#f3722c')) {
      return 'from-yellow-400 to-orange-400' // First Steps
    } else if (colorGradient.includes('#3776ab') && colorGradient.includes('#ffd343')) {
      return 'from-blue-600 to-yellow-400' // Python Master
    } else if (colorGradient.includes('#00f5d4') && colorGradient.includes('#9b5de5')) {
      return 'from-cyan-400 to-purple-500' // Speed Coder
    } else if (colorGradient.includes('#ff6b6b') && colorGradient.includes('#ff8e8e')) {
      return 'from-red-400 to-red-300' // Week Warrior
    } else if (colorGradient.includes('#8338ec') && colorGradient.includes('#3a86ff')) {
      return 'from-purple-500 to-blue-500' // Bug Hunter
    } else if (colorGradient.includes('#4361ee') && colorGradient.includes('#7209b7')) {
      return 'from-blue-500 to-purple-700' // Assembly Guru
    } else if (colorGradient.includes('#f72585') && colorGradient.includes('#b5179e')) {
      return 'from-pink-600 to-purple-600' // 100% Club
    } else if (colorGradient.includes('#495057') && colorGradient.includes('#212529')) {
      return 'from-gray-600 to-gray-900' // Rootkit Researcher
    }
    
    return 'from-green-400 to-cyan-400'
  }

  const formatDate = (dateString?: string) => {
    if (!dateString) return ''
    const date = new Date(dateString)
    return date.toLocaleDateString('en-US', { month: 'short', day: 'numeric', year: 'numeric' })
  }

  const calculateProgressPercentage = (achievement: Achievement) => {
    if (achievement.earned) return 100
    if (!achievement.criteria_value) return 0
    return Math.min(Math.round((achievement.progress / achievement.criteria_value) * 100), 100)
  }

  const handleShareAchievement = (achievement: Achievement) => {
    if (!achievement.earned) return
    
    const shareText = `I just unlocked the "${achievement.title}" achievement on WizardCore! 🎉`
    const shareUrl = window.location.href
    
    if (navigator.share) {
      navigator.share({
        title: 'WizardCore Achievement',
        text: shareText,
        url: shareUrl,
      })
    } else {
      // Fallback to copying to clipboard
      navigator.clipboard.writeText(`${shareText} ${shareUrl}`)
      alert('Achievement link copied to clipboard!')
    }
  }

  if (loading) {
    return (
      <div className="flex items-center justify-center h-64">
        <Loader2 className="w-8 h-8 animate-spin text-neon-cyan" />
        <span className="ml-3 text-text-secondary">Loading achievements...</span>
      </div>
    )
  }

  if (error) {
    return (
      <div className="flex flex-col items-center justify-center h-64">
        <AlertCircle className="w-12 h-12 text-red-400 mb-4" />
        <h3 className="text-lg font-semibold text-text-primary mb-2">Failed to load achievements</h3>
        <p className="text-text-secondary">{error}</p>
        <button 
          onClick={fetchAchievementsData}
          className="mt-4 px-4 py-2 bg-neon-cyan text-black rounded-lg font-medium hover:bg-neon-cyan/90"
        >
          Try Again
        </button>
      </div>
    )
  }

  return (
    <div className="space-y-8">
      {/* Stats */}
      {stats && (
        <div className="grid grid-cols-1 md:grid-cols-3 gap-4">
          <div className="bg-bg-elevated border border-border-default rounded-2xl p-6">
            <div className="flex items-center justify-between">
              <div>
                <p className="text-sm text-text-secondary">Badges Earned</p>
                <p className="text-3xl font-bold text-text-primary">{stats.earned_count}/{stats.total_count}</p>
              </div>
              <div className="p-3 rounded-xl bg-gradient-to-br from-yellow-400 to-orange-400">
                <Trophy className="w-6 h-6 text-white" />
              </div>
            </div>
            <div className="mt-4 text-sm text-text-secondary">
              <span className="text-neon-cyan">Keep going!</span>
            </div>
          </div>

          <div className="bg-bg-elevated border border-border-default rounded-2xl p-6">
            <div className="flex items-center justify-between">
              <div>
                <p className="text-sm text-text-secondary">Global Rank</p>
                <p className="text-3xl font-bold text-text-primary">#{stats.global_rank || '--'}</p>
              </div>
              <div className="p-3 rounded-xl bg-gradient-to-br from-neon-cyan to-neon-lavender">
                <Star className="w-6 h-6 text-white" />
              </div>
            </div>
            <div className="mt-4 text-sm text-text-secondary">
              <span>Top 10% of learners</span>
            </div>
          </div>

          <div className="bg-bg-elevated border border-border-default rounded-2xl p-6">
            <div className="flex items-center justify-between">
              <div>
                <p className="text-sm text-text-secondary">Total XP</p>
                <p className="text-3xl font-bold text-text-primary">{stats.total_xp.toLocaleString()}</p>
              </div>
              <div className="p-3 rounded-xl bg-gradient-to-br from-neon-pink to-neon-purple">
                <Zap className="w-6 h-6 text-white" />
              </div>
            </div>
            <div className="mt-4 text-sm text-text-secondary">
              <span>Next rank at 3,000 XP</span>
            </div>
          </div>
        </div>
      )}

      {/* Badges Grid */}
      <div className="bg-bg-elevated border border-border-default rounded-2xl p-6">
        <h2 className="text-xl font-bold text-text-primary mb-6">Badge Collection</h2>
        <div className="grid grid-cols-1 sm:grid-cols-2 md:grid-cols-4 gap-6">
          {achievements.map((achievement) => {
            const Icon = getIconComponent(achievement.icon)
            const progressPercentage = calculateProgressPercentage(achievement)
            
            return (
              <div
                key={achievement.id}
                className={`border rounded-xl p-5 flex flex-col items-center text-center transition-all duration-300 ${achievement.earned ? 'border-neon-cyan hover:scale-105' : 'border-border-subtle opacity-80'}`}
              >
                <div className={`w-20 h-20 rounded-full bg-gradient-to-br ${getColorClass(achievement.color_gradient)} flex items-center justify-center mb-4`}>
                  <Icon className="w-10 h-10 text-white" />
                </div>
                <h3 className="font-bold text-text-primary">{achievement.title}</h3>
                <p className="text-sm text-text-secondary mt-2">{achievement.description || 'Complete the challenge to earn this badge'}</p>
                
                {achievement.earned ? (
                  <div className="mt-4 space-y-2">
                    <div className="text-xs px-3 py-1 rounded-full bg-green-900/30 text-green-300">
                      Earned {formatDate(achievement.earned_date)}
                    </div>
                    <div className="text-xs text-text-muted">{achievement.rarity} • {achievement.xp_reward} XP</div>
                    <button
                      onClick={() => handleShareAchievement(achievement)}
                      className="text-xs px-3 py-1 border border-neon-cyan text-neon-cyan rounded-lg hover:bg-neon-cyan/10 transition"
                    >
                      Share
                    </button>
                  </div>
                ) : (
                  <div className="mt-4 w-full">
                    <div className="text-xs text-text-muted mb-1">
                      Progress: {progressPercentage}% ({achievement.progress}/{achievement.criteria_value || '?'})
                    </div>
                    <div className="h-2 bg-bg-secondary rounded-full overflow-hidden">
                      <div
                        className="h-full rounded-full bg-gradient-to-r from-neon-cyan to-neon-lavender transition-all duration-500"
                        style={{ width: `${progressPercentage}%` }}
                      />
                    </div>
                    <div className="text-xs text-text-muted mt-2">{achievement.rarity} • {achievement.xp_reward} XP</div>
                  </div>
                )}
              </div>
            )
          })}
        </div>
      </div>

      {/* Leaderboard */}
      <div className="bg-bg-elevated border border-border-default rounded-2xl p-6">
        <div className="flex items-center justify-between mb-6">
          <h2 className="text-xl font-bold text-text-primary">Global Leaderboard</h2>
          <div className="flex gap-2">
            <button className="px-4 py-2 bg-bg-tertiary text-text-secondary rounded-lg text-sm">
              This Week
            </button>
            <button className="px-4 py-2 bg-gradient-to-r from-neon-cyan to-neon-lavender text-white rounded-lg text-sm">
              All Time
            </button>
          </div>
        </div>

        <div className="overflow-x-auto">
          <table className="w-full">
            <thead>
              <tr className="border-b border-border-subtle">
                <th className="text-left py-3 text-text-secondary font-medium">Rank</th>
                <th className="text-left py-3 text-text-secondary font-medium">Learner</th>
                <th className="text-left py-3 text-text-secondary font-medium">XP</th>
                <th className="text-left py-3 text-text-secondary font-medium">Streak</th>
                <th className="text-left py-3 text-text-secondary font-medium">Badges</th>
                <th className="text-left py-3 text-text-secondary font-medium">Action</th>
              </tr>
            </thead>
            <tbody>
              {leaderboard.map((row) => (
                <tr
                  key={row.rank}
                  className={`border-b border-border-subtle ${row.is_current_user ? 'bg-gradient-to-r from-neon-cyan/10 to-neon-lavender/10' : ''}`}
                >
                  <td className="py-4">
                    <div className={`w-8 h-8 rounded-full flex items-center justify-center ${row.rank <= 3 ? 'bg-gradient-to-br from-yellow-400 to-orange-400 text-white' : 'bg-bg-tertiary text-text-secondary'}`}>
                      {row.rank}
                    </div>
                  </td>
                  <td className="py-4">
                    <div className="flex items-center">
                      <div className="w-10 h-10 rounded-full bg-gradient-to-r from-neon-pink to-neon-cyan mr-3" />
                      <div>
                        <div className="font-medium text-text-primary">{row.name}</div>
                        {row.is_current_user && (
                          <div className="text-xs text-neon-cyan">That's you!</div>
                        )}
                      </div>
                    </div>
                  </td>
                  <td className="py-4">
                    <div className="font-bold text-text-primary">{row.xp.toLocaleString()}</div>
                    <div className="text-xs text-text-secondary">XP</div>
                  </td>
                  <td className="py-4">
                    <div className="flex items-center">
                      <Flame className="w-4 h-4 text-orange-400 mr-2" />
                      <span className="font-medium text-text-primary">{row.streak} days</span>
                    </div>
                  </td>
                  <td className="py-4">
                    <div className="flex items-center">
                      <Trophy className="w-4 h-4 text-yellow-400 mr-2" />
                      <span className="font-medium text-text-primary">{row.badges}</span>
                    </div>
                  </td>
                  <td className="py-4">
                    <button className="px-4 py-2 border border-border-subtle text-text-secondary rounded-lg hover:bg-bg-tertiary text-sm">
                      View Profile
                    </button>
                  </td>
                </tr>
              ))}
            </tbody>
          </table>
        </div>

        <div className="mt-6 pt-6 border-t border-border-subtle flex justify-between items-center">
          <div className="text-sm text-text-secondary">
            Leaderboard updates daily at midnight UTC
          </div>
          <button className="text-sm text-neon-cyan hover:underline">
            View Full Leaderboard →
          </button>
        </div>
      </div>

      {/* Achievement Unlock Animation */}
      {unlockedAchievement && (
        <div className="fixed inset-0 bg-black/70 flex items-center justify-center z-50 p-4">
          <div className="bg-bg-elevated border border-neon-cyan rounded-2xl p-8 max-w-md w-full text-center animate-pulse">
            <div className={`w-32 h-32 rounded-full bg-gradient-to-br ${getColorClass(unlockedAchievement.color_gradient)} flex items-center justify-center mx-auto mb-6`}>
              {(() => {
                const Icon = getIconComponent(unlockedAchievement.icon)
                return <Icon className="w-16 h-16 text-white" />
              })()}
            </div>
            <h3 className="text-2xl font-bold text-text-primary mb-2">Achievement Unlocked!</h3>
            <h4 className="text-xl font-semibold text-neon-cyan mb-4">{unlockedAchievement.title}</h4>
            <p className="text-text-secondary mb-6">{unlockedAchievement.description}</p>
            <div className="text-lg font-bold text-yellow-400 mb-6">+{unlockedAchievement.xp_reward} XP</div>
            <button
              onClick={() => setUnlockedAchievement(null)}
              className="px-6 py-3 bg-gradient-to-r from-neon-cyan to-neon-lavender text-white rounded-lg font-medium hover:opacity-90"
            >
              Awesome!
            </button>
          </div>
        </div>
      )}
    </div>
  )
}