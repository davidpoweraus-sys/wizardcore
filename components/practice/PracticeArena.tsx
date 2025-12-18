'use client'

import { useState } from 'react'
import { Clock, Sword, Target, Zap, Trophy, Users, BarChart3, Play, Shield, Timer, Star, Award } from 'lucide-react'

const challengeTypes = [
  {
    id: 'speed',
    title: 'Speed Run',
    description: 'Solve as many problems as you can in 5 minutes',
    icon: Clock,
    color: 'from-neon-cyan to-neon-lavender',
    difficulty: 'Medium',
    xp: 200,
  },
  {
    id: 'duel',
    title: 'Live Duel',
    description: 'Compete against another learner in real-time',
    icon: Sword,
    color: 'from-neon-pink to-neon-purple',
    difficulty: 'Hard',
    xp: 500,
  },
  {
    id: 'random',
    title: 'Random Challenge',
    description: 'Get a random exercise from any pathway',
    icon: Target,
    color: 'from-green-400 to-cyan-400',
    difficulty: 'Variable',
    xp: 150,
  },
  {
    id: 'streak',
    title: 'Streak Builder',
    description: 'Maintain a perfect streak across 10 exercises',
    icon: Zap,
    color: 'from-yellow-400 to-orange-400',
    difficulty: 'Hard',
    xp: 300,
  },
]

const practiceAreas = [
  { name: 'Python', exercises: 42, completed: 28, color: 'from-green-400 to-cyan-400' },
  { name: 'C & Assembly', exercises: 38, completed: 12, color: 'from-blue-400 to-purple-400' },
  { name: 'JavaScript', exercises: 35, completed: 18, color: 'from-yellow-400 to-orange-400' },
  { name: 'SQL', exercises: 25, completed: 10, color: 'from-red-400 to-pink-400' },
  { name: 'Reverse Engineering', exercises: 30, completed: 5, color: 'from-purple-400 to-indigo-400' },
  { name: 'Rootkits', exercises: 20, completed: 0, color: 'from-gray-700 to-black' },
]

const recentMatches = [
  { opponent: 'Sam Rivera', result: 'win', score: '3-1', xp: 120, time: '10 min ago' },
  { opponent: 'Jordan Lee', result: 'loss', score: '1-3', xp: 40, time: '1 hour ago' },
  { opponent: 'Taylor Kim', result: 'win', score: '2-0', xp: 100, time: '3 hours ago' },
  { opponent: 'Casey Morgan', result: 'draw', score: '2-2', xp: 80, time: '1 day ago' },
]

export default function PracticeArena() {
  const [selectedChallenge, setSelectedChallenge] = useState<string | null>(null)
  const [timer, setTimer] = useState(300) // 5 minutes in seconds

  const startChallenge = (id: string) => {
    setSelectedChallenge(id)
    // Start timer logic would go here
  }

  return (
    <div className="space-y-8">
      {/* Stats */}
      <div className="grid grid-cols-1 md:grid-cols-4 gap-4">
        <div className="bg-bg-elevated border border-border-default rounded-2xl p-6">
          <div className="flex items-center justify-between">
            <div>
              <p className="text-sm text-text-secondary">Practice Score</p>
              <p className="text-3xl font-bold text-text-primary">1,850</p>
            </div>
            <div className="p-3 rounded-xl bg-gradient-to-br from-neon-cyan to-neon-lavender">
              <Trophy className="w-6 h-6 text-white" />
            </div>
          </div>
          <div className="mt-2 text-sm text-text-secondary">
            <BarChart3 className="w-4 h-4 inline mr-2" />
            <span>Top 15%</span>
          </div>
        </div>

        <div className="bg-bg-elevated border border-border-default rounded-2xl p-6">
          <div className="flex items-center justify-between">
            <div>
              <p className="text-sm text-text-secondary">Duels Won</p>
              <p className="text-3xl font-bold text-text-primary">24</p>
            </div>
            <div className="p-3 rounded-xl bg-gradient-to-br from-neon-pink to-neon-purple">
              <Sword className="w-6 h-6 text-white" />
            </div>
          </div>
          <div className="mt-2 text-sm text-text-secondary">
            <span>Win rate: 68%</span>
          </div>
        </div>

        <div className="bg-bg-elevated border border-border-default rounded-2xl p-6">
          <div className="flex items-center justify-between">
            <div>
              <p className="text-sm text-text-secondary">Avg. Time</p>
              <p className="text-3xl font-bold text-text-primary">2:14</p>
            </div>
            <div className="p-3 rounded-xl bg-gradient-to-br from-green-400 to-cyan-400">
              <Timer className="w-6 h-6 text-white" />
            </div>
          </div>
          <div className="mt-2 text-sm text-text-secondary">
            <span>Per exercise</span>
          </div>
        </div>

        <div className="bg-bg-elevated border border-border-default rounded-2xl p-6">
          <div className="flex items-center justify-between">
            <div>
              <p className="text-sm text-text-secondary">Live Opponents</p>
              <p className="text-3xl font-bold text-text-primary">12</p>
            </div>
            <div className="p-3 rounded-xl bg-gradient-to-br from-yellow-400 to-orange-400">
              <Users className="w-6 h-6 text-white" />
            </div>
          </div>
          <div className="mt-2 text-sm text-text-secondary">
            <span>Ready to duel</span>
          </div>
        </div>
      </div>

      {/* Challenge Types */}
      <div className="bg-bg-elevated border border-border-default rounded-2xl p-6">
        <h2 className="text-xl font-bold text-text-primary mb-6">Challenge Types</h2>
        <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-4 gap-6">
          {challengeTypes.map((challenge) => {
            const Icon = challenge.icon
            return (
              <div
                key={challenge.id}
                className="border border-border-subtle rounded-xl p-5 hover:border-neon-cyan transition"
              >
                <div className={`w-14 h-14 rounded-xl bg-gradient-to-br ${challenge.color} flex items-center justify-center mb-4`}>
                  <Icon className="w-7 h-7 text-white" />
                </div>
                <h3 className="font-bold text-text-primary">{challenge.title}</h3>
                <p className="text-sm text-text-secondary mt-2">{challenge.description}</p>
                
                <div className="flex items-center justify-between mt-4">
                  <div className="text-sm">
                    <span className="text-text-secondary">Difficulty: </span>
                    <span className={`font-medium ${challenge.difficulty === 'Hard' ? 'text-red-400' : challenge.difficulty === 'Medium' ? 'text-yellow-400' : 'text-green-400'}`}>
                      {challenge.difficulty}
                    </span>
                  </div>
                  <div className="text-sm font-bold text-neon-cyan">{challenge.xp} XP</div>
                </div>

                <button
                  onClick={() => startChallenge(challenge.id)}
                  className="w-full mt-4 py-2 bg-gradient-to-r from-neon-cyan to-neon-lavender text-white rounded-lg font-medium hover:opacity-90 transition flex items-center justify-center gap-2"
                >
                  <Play className="w-4 h-4" />
                  Start Challenge
                </button>
              </div>
            )
          })}
        </div>
      </div>

      {/* Practice Areas */}
      <div className="grid grid-cols-1 lg:grid-cols-3 gap-6">
        <div className="lg:col-span-2">
          <div className="bg-bg-elevated border border-border-default rounded-2xl p-6">
            <h2 className="text-xl font-bold text-text-primary mb-6">Practice Areas</h2>
            <div className="space-y-4">
              {practiceAreas.map((area) => {
                const progress = Math.round((area.completed / area.exercises) * 100)
                return (
                  <div key={area.name} className="border border-border-subtle rounded-xl p-4">
                    <div className="flex items-center justify-between mb-3">
                      <div className="flex items-center">
                        <div className={`w-10 h-10 rounded-lg bg-gradient-to-br ${area.color} flex items-center justify-center text-white font-bold mr-4`}>
                          {area.name.charAt(0)}
                        </div>
                        <div>
                          <h3 className="font-bold text-text-primary">{area.name}</h3>
                          <p className="text-sm text-text-secondary">{area.exercises} exercises</p>
                        </div>
                      </div>
                      <div className="text-right">
                        <div className="text-2xl font-bold text-text-primary">{progress}%</div>
                        <div className="text-xs text-text-muted">complete</div>
                      </div>
                    </div>

                    <div className="mb-2">
                      <div className="flex justify-between text-sm text-text-secondary mb-1">
                        <span>Progress</span>
                        <span>{area.completed} / {area.exercises}</span>
                      </div>
                      <div className="h-3 bg-bg-secondary rounded-full overflow-hidden">
                        <div
                          className={`h-full rounded-full bg-gradient-to-r ${area.color}`}
                          style={{ width: `${progress}%` }}
                        />
                      </div>
                    </div>

                    <div className="flex justify-between">
                      <button className="text-sm text-neon-cyan hover:underline">
                        View Exercises
                      </button>
                      <button className="px-3 py-1 bg-bg-tertiary text-text-secondary rounded-lg text-sm">
                        Practice
                      </button>
                    </div>
                  </div>
                )
              })}
            </div>
          </div>
        </div>

        {/* Recent Matches & Leaderboard */}
        <div className="space-y-6">
          <div className="bg-bg-elevated border border-border-default rounded-2xl p-6">
            <h2 className="text-xl font-bold text-text-primary mb-6">Recent Matches</h2>
            <div className="space-y-4">
              {recentMatches.map((match, idx) => (
                <div key={idx} className="border border-border-subtle rounded-lg p-4">
                  <div className="flex items-center justify-between mb-2">
                    <div className="flex items-center">
                      <div className="w-10 h-10 rounded-full bg-gradient-to-r from-neon-pink to-neon-cyan mr-3" />
                      <div>
                        <div className="font-medium text-text-primary">{match.opponent}</div>
                        <div className="text-xs text-text-secondary">{match.time}</div>
                      </div>
                    </div>
                    <div className={`px-3 py-1 rounded-full text-sm font-medium ${match.result === 'win' ? 'bg-green-900/30 text-green-300' : match.result === 'loss' ? 'bg-red-900/30 text-red-300' : 'bg-yellow-900/30 text-yellow-300'}`}>
                      {match.result.toUpperCase()}
                    </div>
                  </div>
                  <div className="flex items-center justify-between text-sm">
                    <div className="text-text-secondary">Score: <span className="font-bold text-text-primary">{match.score}</span></div>
                    <div className="flex items-center text-neon-cyan">
                      <Zap className="w-4 h-4 mr-1" />
                      +{match.xp} XP
                    </div>
                  </div>
                </div>
              ))}
            </div>
          </div>

          {/* Quick Practice */}
          <div className="bg-bg-elevated border border-border-default rounded-2xl p-6">
            <h2 className="text-xl font-bold text-text-primary mb-6">Quick Practice</h2>
            <div className="space-y-4">
              <button className="w-full p-4 border border-border-subtle rounded-xl hover:border-neon-cyan transition text-left">
                <div className="flex items-center justify-between">
                  <div>
                    <div className="font-medium text-text-primary">5-Minute Drill</div>
                    <div className="text-sm text-text-secondary">Solve 3 random Python problems</div>
                  </div>
                  <Star className="w-5 h-5 text-yellow-400" />
                </div>
              </button>

              <button className="w-full p-4 border border-border-subtle rounded-xl hover:border-neon-cyan transition text-left">
                <div className="flex items-center justify-between">
                  <div>
                    <div className="font-medium text-text-primary">Find Opponent</div>
                    <div className="text-sm text-text-secondary">Match with a similar skill learner</div>
                  </div>
                  <Users className="w-5 h-5 text-neon-cyan" />
                </div>
              </button>

              <button className="w-full p-4 border border-border-subtle rounded-xl hover:border-neon-cyan transition text-left">
                <div className="flex items-center justify-between">
                  <div>
                    <div className="font-medium text-text-primary">Daily Challenge</div>
                    <div className="text-sm text-text-secondary">Earn bonus XP</div>
                  </div>
                  <Award className="w-5 h-5 text-neon-pink" />
                </div>
              </button>
            </div>
          </div>
        </div>
      </div>
    </div>
  )
}