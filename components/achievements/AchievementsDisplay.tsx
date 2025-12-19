'use client'

import { Trophy, MoonStar as Star, Zap, Bullseye, Zap as Flame, Shield, Lock, Check } from 'pixelarticons/fonts/react'

const badges = [
  {
    id: 1,
    title: 'First Steps',
    description: 'Complete your first exercise',
    icon: Star,
    color: 'from-yellow-400 to-orange-400',
    earned: true,
    date: '2025-11-10',
    rarity: 'Common',
  },
  {
    id: 2,
    title: 'Speed Coder',
    description: 'Solve 10 exercises under 5 minutes each',
    icon: Zap,
    color: 'from-neon-cyan to-neon-lavender',
    earned: true,
    date: '2025-11-25',
    rarity: 'Rare',
  },
  {
    id: 3,
    title: 'Python Master',
    description: 'Complete all Python pathway modules',
    icon: Trophy,
    color: 'from-green-400 to-cyan-400',
    earned: false,
    progress: 75,
    rarity: 'Epic',
  },
  {
    id: 4,
    title: 'Week Warrior',
    description: 'Maintain a 7-day learning streak',
    icon: Flame,
    color: 'from-red-400 to-pink-400',
    earned: true,
    date: '2025-12-01',
    rarity: 'Uncommon',
  },
  {
    id: 5,
    title: 'Bug Hunter',
    description: 'Find and report a security vulnerability',
    icon: Bullseye,
    color: 'from-purple-400 to-indigo-400',
    earned: false,
    progress: 0,
    rarity: 'Legendary',
  },
  {
    id: 6,
    title: 'Assembly Guru',
    description: 'Write functional x64 shellcode',
    icon: Shield,
    color: 'from-blue-400 to-purple-400',
    earned: false,
    progress: 15,
    rarity: 'Epic',
  },
  {
    id: 7,
    title: '100% Club',
    description: 'Achieve perfect score on 20 exercises',
    icon: Check,
    color: 'from-neon-pink to-neon-purple',
    earned: false,
    progress: 45,
    rarity: 'Rare',
  },
  {
    id: 8,
    title: 'Rootkit Researcher',
    description: 'Complete Rootkit Development course',
    icon: Lock,
    color: 'from-gray-700 to-black',
    earned: false,
    progress: 0,
    rarity: 'Mythic',
  },
]

const leaderboard = [
  { rank: 1, name: 'Alex Chen', xp: 12540, streak: 21, badges: 12 },
  { rank: 2, name: 'Sam Rivera', xp: 11200, streak: 18, badges: 10 },
  { rank: 3, name: 'Jordan Lee', xp: 9850, streak: 14, badges: 9 },
  { rank: 4, name: 'Taylor Kim', xp: 8760, streak: 12, badges: 8 },
  { rank: 5, name: 'Casey Morgan', xp: 7540, streak: 9, badges: 7 },
  { rank: 6, name: 'Riley Patel', xp: 6320, streak: 7, badges: 6 },
  { rank: 7, name: 'You', xp: 2850, streak: 7, badges: 4 },
  { rank: 8, name: 'Drew Zhang', xp: 2450, streak: 5, badges: 3 },
]

export default function AchievementsDisplay() {
  const earnedCount = badges.filter(b => b.earned).length
  const totalXP = leaderboard.find(l => l.name === 'You')?.xp || 0

  return (
    <div className="space-y-8">
      {/* Stats */}
      <div className="grid grid-cols-1 md:grid-cols-3 gap-4">
        <div className="bg-bg-elevated border border-border-default rounded-2xl p-6">
          <div className="flex items-center justify-between">
            <div>
              <p className="text-sm text-text-secondary">Badges Earned</p>
              <p className="text-3xl font-bold text-text-primary">{earnedCount}/{badges.length}</p>
            </div>
            <div className="p-3 rounded-xl bg-gradient-to-br from-yellow-400 to-orange-400">
              <Trophy className="w-6 h-6 text-white" />
            </div>
          </div>
          <div className="mt-4 text-sm text-text-secondary">
            <span className="text-neon-cyan">+2 this month</span>
          </div>
        </div>

        <div className="bg-bg-elevated border border-border-default rounded-2xl p-6">
          <div className="flex items-center justify-between">
            <div>
              <p className="text-sm text-text-secondary">Global Rank</p>
              <p className="text-3xl font-bold text-text-primary">#7</p>
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
              <p className="text-3xl font-bold text-text-primary">{totalXP.toLocaleString()}</p>
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

      {/* Badges Grid */}
      <div className="bg-bg-elevated border border-border-default rounded-2xl p-6">
        <h2 className="text-xl font-bold text-text-primary mb-6">Badge Collection</h2>
        <div className="grid grid-cols-1 sm:grid-cols-2 md:grid-cols-4 gap-6">
          {badges.map((badge) => {
            const Icon = badge.icon
            return (
              <div
                key={badge.id}
                className={`border rounded-xl p-5 flex flex-col items-center text-center ${badge.earned ? 'border-neon-cyan' : 'border-border-subtle opacity-70'}`}
              >
                <div className={`w-20 h-20 rounded-full bg-gradient-to-br ${badge.color} flex items-center justify-center mb-4`}>
                  <Icon className="w-10 h-10 text-white" />
                </div>
                <h3 className="font-bold text-text-primary">{badge.title}</h3>
                <p className="text-sm text-text-secondary mt-2">{badge.description}</p>
                
                {badge.earned ? (
                  <div className="mt-4">
                    <div className="text-xs px-3 py-1 rounded-full bg-green-900/30 text-green-300">
                      Earned {badge.date}
                    </div>
                    <div className="text-xs text-text-muted mt-2">{badge.rarity}</div>
                  </div>
                ) : (
                  <div className="mt-4 w-full">
                    <div className="text-xs text-text-muted mb-1">Progress: {badge.progress}%</div>
                    <div className="h-2 bg-bg-secondary rounded-full overflow-hidden">
                      <div
                        className="h-full rounded-full bg-gradient-to-r from-neon-cyan to-neon-lavender"
                        style={{ width: `${badge.progress}%` }}
                      />
                    </div>
                    <div className="text-xs text-text-muted mt-2">{badge.rarity}</div>
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
                  className={`border-b border-border-subtle ${row.name === 'You' ? 'bg-gradient-to-r from-neon-cyan/10 to-neon-lavender/10' : ''}`}
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
                        {row.name === 'You' && (
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
    </div>
  )
}