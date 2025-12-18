'use client'

import { Trophy, Flame, Zap, TrendingUp, Users, Globe } from 'lucide-react'

const leaderboardData = [
  { rank: 1, name: 'Alex Chen', xp: 12540, streak: 21, badges: 12, country: 'US', change: 'up' },
  { rank: 2, name: 'Sam Rivera', xp: 11200, streak: 18, badges: 10, country: 'CA', change: 'up' },
  { rank: 3, name: 'Jordan Lee', xp: 9850, streak: 14, badges: 9, country: 'UK', change: 'down' },
  { rank: 4, name: 'Taylor Kim', xp: 8760, streak: 12, badges: 8, country: 'KR', change: 'up' },
  { rank: 5, name: 'Casey Morgan', xp: 7540, streak: 9, badges: 7, country: 'AU', change: 'same' },
  { rank: 6, name: 'Riley Patel', xp: 6320, streak: 7, badges: 6, country: 'IN', change: 'up' },
  { rank: 7, name: 'You', xp: 2850, streak: 7, badges: 4, country: 'GB', change: 'up' },
  { rank: 8, name: 'Drew Zhang', xp: 2450, streak: 5, badges: 3, country: 'CN', change: 'down' },
  { rank: 9, name: 'Morgan West', xp: 2100, streak: 4, badges: 3, country: 'DE', change: 'same' },
  { rank: 10, name: 'Blake Soto', xp: 1850, streak: 3, badges: 2, country: 'ES', change: 'up' },
]

const timeframes = [
  { id: 'all', label: 'All Time' },
  { id: 'month', label: 'This Month' },
  { id: 'week', label: 'This Week' },
  { id: 'python', label: 'Python' },
  { id: 'c', label: 'C & Assembly' },
]

export default function LeaderboardTable() {
  return (
    <div className="space-y-6">
      {/* Stats */}
      <div className="grid grid-cols-1 md:grid-cols-4 gap-4">
        <div className="bg-bg-elevated border border-border-default rounded-2xl p-6">
          <div className="flex items-center justify-between">
            <div>
              <p className="text-sm text-text-secondary">Total Learners</p>
              <p className="text-3xl font-bold text-text-primary">2,847</p>
            </div>
            <div className="p-3 rounded-xl bg-gradient-to-br from-neon-cyan to-neon-lavender">
              <Users className="w-6 h-6 text-white" />
            </div>
          </div>
        </div>

        <div className="bg-bg-elevated border border-border-default rounded-2xl p-6">
          <div className="flex items-center justify-between">
            <div>
              <p className="text-sm text-text-secondary">Your Rank</p>
              <p className="text-3xl font-bold text-text-primary">#7</p>
            </div>
            <div className="p-3 rounded-xl bg-gradient-to-br from-neon-pink to-neon-purple">
              <Trophy className="w-6 h-6 text-white" />
            </div>
          </div>
          <div className="mt-2 text-sm text-text-secondary">
            <TrendingUp className="w-4 h-4 inline mr-2 text-green-400" />
            <span>Up 2 places this week</span>
          </div>
        </div>

        <div className="bg-bg-elevated border border-border-default rounded-2xl p-6">
          <div className="flex items-center justify-between">
            <div>
              <p className="text-sm text-text-secondary">Top XP</p>
              <p className="text-3xl font-bold text-text-primary">12,540</p>
            </div>
            <div className="p-3 rounded-xl bg-gradient-to-br from-yellow-400 to-orange-400">
              <Zap className="w-6 h-6 text-white" />
            </div>
          </div>
          <div className="mt-2 text-sm text-text-secondary">Alex Chen</div>
        </div>

        <div className="bg-bg-elevated border border-border-default rounded-2xl p-6">
          <div className="flex items-center justify-between">
            <div>
              <p className="text-sm text-text-secondary">Countries</p>
              <p className="text-3xl font-bold text-text-primary">64</p>
            </div>
            <div className="p-3 rounded-xl bg-gradient-to-br from-green-400 to-cyan-400">
              <Globe className="w-6 h-6 text-white" />
            </div>
          </div>
        </div>
      </div>

      {/* Timeframe Filters */}
      <div className="flex flex-wrap gap-2">
        {timeframes.map((timeframe) => (
          <button
            key={timeframe.id}
            className={`px-4 py-2 rounded-lg text-sm font-medium ${timeframe.id === 'all' ? 'bg-gradient-to-r from-neon-cyan to-neon-lavender text-white' : 'bg-bg-tertiary text-text-secondary hover:bg-bg-secondary'}`}
          >
            {timeframe.label}
          </button>
        ))}
      </div>

      {/* Leaderboard Table */}
      <div className="bg-bg-elevated border border-border-default rounded-2xl overflow-hidden">
        <div className="overflow-x-auto">
          <table className="w-full">
            <thead>
              <tr className="border-b border-border-subtle">
                <th className="text-left py-4 px-6 text-text-secondary font-medium">Rank</th>
                <th className="text-left py-4 px-6 text-text-secondary font-medium">Learner</th>
                <th className="text-left py-4 px-6 text-text-secondary font-medium">XP</th>
                <th className="text-left py-4 px-6 text-text-secondary font-medium">Streak</th>
                <th className="text-left py-4 px-6 text-text-secondary font-medium">Badges</th>
                <th className="text-left py-4 px-6 text-text-secondary font-medium">Country</th>
                <th className="text-left py-4 px-6 text-text-secondary font-medium">Trend</th>
              </tr>
            </thead>
            <tbody>
              {leaderboardData.map((row) => (
                <tr
                  key={row.rank}
                  className={`border-b border-border-subtle ${row.name === 'You' ? 'bg-gradient-to-r from-neon-cyan/10 to-neon-lavender/10' : ''}`}
                >
                  <td className="py-4 px-6">
                    <div className={`w-10 h-10 rounded-full flex items-center justify-center ${row.rank <= 3 ? 'bg-gradient-to-br from-yellow-400 to-orange-400 text-white' : 'bg-bg-tertiary text-text-secondary'}`}>
                      {row.rank}
                    </div>
                  </td>
                  <td className="py-4 px-6">
                    <div className="flex items-center">
                      <div className="w-12 h-12 rounded-full bg-gradient-to-r from-neon-pink to-neon-cyan mr-4" />
                      <div>
                        <div className="font-bold text-text-primary">{row.name}</div>
                        <div className="text-sm text-text-secondary">Level {Math.floor(row.xp / 1000)}</div>
                      </div>
                    </div>
                  </td>
                  <td className="py-4 px-6">
                    <div className="font-bold text-text-primary">{row.xp.toLocaleString()}</div>
                    <div className="text-sm text-text-secondary">XP</div>
                  </td>
                  <td className="py-4 px-6">
                    <div className="flex items-center">
                      <Flame className="w-5 h-5 text-orange-400 mr-2" />
                      <span className="font-bold text-text-primary">{row.streak}</span>
                      <span className="text-sm text-text-secondary ml-1">days</span>
                    </div>
                  </td>
                  <td className="py-4 px-6">
                    <div className="flex items-center">
                      <Trophy className="w-5 h-5 text-yellow-400 mr-2" />
                      <span className="font-bold text-text-primary">{row.badges}</span>
                    </div>
                  </td>
                  <td className="py-4 px-6">
                    <div className="flex items-center">
                      <div className="w-6 h-6 rounded-full bg-gradient-to-r from-blue-400 to-purple-400 mr-2" />
                      <span className="text-text-primary">{row.country}</span>
                    </div>
                  </td>
                  <td className="py-4 px-6">
                    {row.change === 'up' && (
                      <div className="flex items-center text-green-400">
                        <TrendingUp className="w-5 h-5 mr-2" />
                        <span>Rising</span>
                      </div>
                    )}
                    {row.change === 'down' && (
                      <div className="flex items-center text-red-400">
                        <TrendingUp className="w-5 h-5 mr-2 rotate-180" />
                        <span>Falling</span>
                      </div>
                    )}
                    {row.change === 'same' && (
                      <div className="text-text-secondary">—</div>
                    )}
                  </td>
                </tr>
              ))}
            </tbody>
          </table>
        </div>

        {/* Footer */}
        <div className="border-t border-border-subtle p-6 flex justify-between items-center">
          <div className="text-sm text-text-secondary">
            Showing top 10 of 2,847 learners
          </div>
          <div className="flex gap-2">
            <button className="px-4 py-2 border border-border-subtle text-text-secondary rounded-lg hover:bg-bg-tertiary">
              Previous
            </button>
            <button className="px-4 py-2 bg-gradient-to-r from-neon-cyan to-neon-lavender text-white rounded-lg">
              1
            </button>
            <button className="px-4 py-2 border border-border-subtle text-text-secondary rounded-lg hover:bg-bg-tertiary">
              2
            </button>
            <button className="px-4 py-2 border border-border-subtle text-text-secondary rounded-lg hover:bg-bg-tertiary">
              3
            </button>
            <button className="px-4 py-2 border border-border-subtle text-text-secondary rounded-lg hover:bg-bg-tertiary">
              Next
            </button>
          </div>
        </div>
      </div>

      {/* Additional Info */}
      <div className="grid grid-cols-1 lg:grid-cols-3 gap-6">
        <div className="lg:col-span-2 bg-bg-elevated border border-border-default rounded-2xl p-6">
          <h3 className="text-lg font-bold text-text-primary mb-4">How Ranking Works</h3>
          <ul className="space-y-3 text-text-secondary">
            <li className="flex items-start">
              <div className="w-6 h-6 rounded-full bg-neon-cyan/20 text-neon-cyan flex items-center justify-center mr-3 flex-shrink-0">1</div>
              <span>Earn XP by completing exercises, modules, and projects</span>
            </li>
            <li className="flex items-start">
              <div className="w-6 h-6 rounded-full bg-neon-cyan/20 text-neon-cyan flex items-center justify-center mr-3 flex-shrink-0">2</div>
              <span>Maintain daily streaks for bonus multipliers</span>
            </li>
            <li className="flex items-start">
              <div className="w-6 h-6 rounded-full bg-neon-cyan/20 text-neon-cyan flex items-center justify-center mr-3 flex-shrink-0">3</div>
              <span>Earn badges for special achievements</span>
            </li>
            <li className="flex items-start">
              <div className="w-6 h-6 rounded-full bg-neon-cyan/20 text-neon-cyan flex items-center justify-center mr-3 flex-shrink-0">4</div>
              <span>Rank is updated daily based on total XP</span>
            </li>
          </ul>
        </div>

        <div className="bg-bg-elevated border border-border-default rounded-2xl p-6">
          <h3 className="text-lg font-bold text-text-primary mb-4">Your Next Milestone</h3>
          <div className="space-y-4">
            <div>
              <div className="flex justify-between text-sm text-text-secondary mb-1">
                <span>Next Rank: #6</span>
                <span>500 XP needed</span>
              </div>
              <div className="h-3 bg-bg-secondary rounded-full overflow-hidden">
                <div className="h-full rounded-full bg-gradient-to-r from-neon-pink to-neon-purple" style={{ width: '70%' }} />
              </div>
            </div>
            <div className="text-sm text-text-secondary">
              Complete 3 more Python exercises to reach the next rank!
            </div>
            <button className="w-full py-2 bg-gradient-to-r from-neon-cyan to-neon-lavender text-white rounded-lg text-sm font-medium">
              View Recommended Exercises
            </button>
          </div>
        </div>
      </div>
    </div>
  )
}