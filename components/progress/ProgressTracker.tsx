'use client'

import { TrendingUp, Target, Award, Calendar, BarChart3, CheckCircle, Clock, Zap } from 'lucide-react'

const pathways = [
  {
    id: 1,
    title: 'Python for Offensive Security',
    progress: 75,
    completedModules: 3,
    totalModules: 5,
    xp: 1250,
    streak: 7,
    lastActivity: '2 hours ago',
  },
  {
    id: 2,
    title: 'C & Assembly',
    progress: 30,
    completedModules: 1,
    totalModules: 5,
    xp: 450,
    streak: 3,
    lastActivity: '1 day ago',
  },
  {
    id: 3,
    title: 'JavaScript & Browser',
    progress: 50,
    completedModules: 2,
    totalModules: 4,
    xp: 800,
    streak: 5,
    lastActivity: '5 hours ago',
  },
  {
    id: 4,
    title: 'SQL & Database',
    progress: 20,
    completedModules: 1,
    totalModules: 4,
    xp: 300,
    streak: 2,
    lastActivity: '3 days ago',
  },
  {
    id: 5,
    title: 'Reverse Engineering',
    progress: 10,
    completedModules: 0,
    totalModules: 5,
    xp: 150,
    streak: 1,
    lastActivity: '1 week ago',
  },
  {
    id: 6,
    title: 'Rootkit Development',
    progress: 0,
    completedModules: 0,
    totalModules: 10,
    xp: 0,
    streak: 0,
    lastActivity: 'Not started',
  },
]

const milestones = [
  { title: 'First Exercise', date: '2025-11-10', xp: 100 },
  { title: 'Python Module 1', date: '2025-11-15', xp: 250 },
  { title: '10 Day Streak', date: '2025-11-25', xp: 500 },
  { title: 'C Module 1', date: '2025-12-01', xp: 300 },
  { title: 'First Project', date: '2025-12-10', xp: 750 },
]

export default function ProgressTracker() {
  const totalXP = pathways.reduce((sum, p) => sum + p.xp, 0)
  const totalCompleted = pathways.reduce((sum, p) => sum + p.completedModules, 0)
  const totalModules = pathways.reduce((sum, p) => sum + p.totalModules, 0)
  const overallProgress = Math.round((totalCompleted / totalModules) * 100)

  return (
    <div className="space-y-8">
      {/* Stats Overview */}
      <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-4 gap-4">
        <div className="bg-bg-elevated border border-border-default rounded-2xl p-6">
          <div className="flex items-center justify-between">
            <div>
              <p className="text-sm text-text-secondary">Total XP</p>
              <p className="text-3xl font-bold text-text-primary">{totalXP.toLocaleString()}</p>
            </div>
            <div className="p-3 rounded-xl bg-gradient-to-br from-neon-cyan to-neon-lavender">
              <Zap className="w-6 h-6 text-white" />
            </div>
          </div>
          <div className="mt-4 text-sm text-text-secondary">
            <TrendingUp className="w-4 h-4 inline mr-2 text-green-400" />
            <span>+320 this week</span>
          </div>
        </div>

        <div className="bg-bg-elevated border border-border-default rounded-2xl p-6">
          <div className="flex items-center justify-between">
            <div>
              <p className="text-sm text-text-secondary">Overall Progress</p>
              <p className="text-3xl font-bold text-text-primary">{overallProgress}%</p>
            </div>
            <div className="p-3 rounded-xl bg-gradient-to-br from-neon-pink to-neon-purple">
              <Target className="w-6 h-6 text-white" />
            </div>
          </div>
          <div className="mt-4">
            <div className="h-2 bg-bg-secondary rounded-full overflow-hidden">
              <div
                className="h-full rounded-full bg-gradient-to-r from-neon-pink to-neon-purple"
                style={{ width: `${overallProgress}%` }}
              />
            </div>
          </div>
        </div>

        <div className="bg-bg-elevated border border-border-default rounded-2xl p-6">
          <div className="flex items-center justify-between">
            <div>
              <p className="text-sm text-text-secondary">Current Streak</p>
              <p className="text-3xl font-bold text-text-primary">7 days</p>
            </div>
            <div className="p-3 rounded-xl bg-gradient-to-br from-green-400 to-cyan-400">
              <Calendar className="w-6 h-6 text-white" />
            </div>
          </div>
          <div className="mt-4 text-sm text-text-secondary">
            <Clock className="w-4 h-4 inline mr-2" />
            <span>Keep going!</span>
          </div>
        </div>

        <div className="bg-bg-elevated border border-border-default rounded-2xl p-6">
          <div className="flex items-center justify-between">
            <div>
              <p className="text-sm text-text-secondary">Modules Completed</p>
              <p className="text-3xl font-bold text-text-primary">{totalCompleted}/{totalModules}</p>
            </div>
            <div className="p-3 rounded-xl bg-gradient-to-br from-yellow-400 to-orange-400">
              <Award className="w-6 h-6 text-white" />
            </div>
          </div>
          <div className="mt-4 text-sm text-text-secondary">
            <CheckCircle className="w-4 h-4 inline mr-2 text-green-400" />
            <span>{Math.round((totalCompleted / totalModules) * 100)}% complete</span>
          </div>
        </div>
      </div>

      {/* Pathway Progress */}
      <div className="bg-bg-elevated border border-border-default rounded-2xl p-6">
        <div className="flex items-center justify-between mb-6">
          <h2 className="text-xl font-bold text-text-primary">Pathway Progress</h2>
          <button className="text-sm text-neon-cyan hover:underline">View Details</button>
        </div>

        <div className="space-y-6">
          {pathways.map((pathway) => (
            <div key={pathway.id} className="border border-border-subtle rounded-xl p-5">
              <div className="flex items-center justify-between mb-4">
                <div>
                  <h3 className="font-bold text-text-primary">{pathway.title}</h3>
                  <p className="text-sm text-text-secondary">
                    {pathway.completedModules} of {pathway.totalModules} modules • {pathway.xp} XP
                  </p>
                </div>
                <div className="text-right">
                  <div className="text-2xl font-bold text-text-primary">{pathway.progress}%</div>
                  <div className="text-xs text-text-muted">complete</div>
                </div>
              </div>

              <div className="mb-4">
                <div className="flex justify-between text-sm text-text-secondary mb-1">
                  <span>Progress</span>
                  <span>{pathway.progress}%</span>
                </div>
                <div className="h-3 bg-bg-secondary rounded-full overflow-hidden">
                  <div
                    className="h-full rounded-full bg-gradient-to-r from-neon-cyan to-neon-lavender"
                    style={{ width: `${pathway.progress}%` }}
                  />
                </div>
              </div>

              <div className="flex items-center justify-between text-sm">
                <div className="flex items-center text-text-secondary">
                  <Calendar className="w-4 h-4 mr-2" />
                  <span>Streak: {pathway.streak} days</span>
                </div>
                <div className="text-text-secondary">
                  Last activity: {pathway.lastActivity}
                </div>
              </div>
            </div>
          ))}
        </div>
      </div>

      {/* Milestones & Timeline */}
      <div className="grid grid-cols-1 lg:grid-cols-3 gap-6">
        <div className="lg:col-span-2">
          <div className="bg-bg-elevated border border-border-default rounded-2xl p-6">
            <h2 className="text-xl font-bold text-text-primary mb-6">Learning Timeline</h2>
            <div className="space-y-4">
              {milestones.map((milestone, idx) => (
                <div key={idx} className="flex items-start">
                  <div className="flex-shrink-0 w-8 h-8 rounded-full bg-gradient-to-r from-neon-pink to-neon-purple flex items-center justify-center text-white font-bold mr-4">
                    {idx + 1}
                  </div>
                  <div className="flex-1 border-b border-border-subtle pb-4">
                    <div className="flex justify-between">
                      <h3 className="font-medium text-text-primary">{milestone.title}</h3>
                      <span className="text-sm text-text-secondary">{milestone.date}</span>
                    </div>
                    <p className="text-sm text-text-secondary mt-1">
                      Earned <span className="text-neon-cyan font-medium">{milestone.xp} XP</span>
                    </p>
                  </div>
                </div>
              ))}
            </div>
          </div>
        </div>

        <div>
          <div className="bg-bg-elevated border border-border-default rounded-2xl p-6">
            <h2 className="text-xl font-bold text-text-primary mb-6">Weekly Activity</h2>
            <div className="space-y-4">
              {[
                { day: 'Mon', hours: 2.5 },
                { day: 'Tue', hours: 1.8 },
                { day: 'Wed', hours: 3.2 },
                { day: 'Thu', hours: 2.0 },
                { day: 'Fri', hours: 4.5 },
                { day: 'Sat', hours: 1.5 },
                { day: 'Sun', hours: 2.8 },
              ].map((item, idx) => (
                <div key={idx} className="flex items-center">
                  <div className="w-12 text-sm text-text-secondary">{item.day}</div>
                  <div className="flex-1 h-4 bg-bg-secondary rounded-full overflow-hidden">
                    <div
                      className="h-full rounded-full bg-gradient-to-r from-neon-cyan to-neon-lavender"
                      style={{ width: `${(item.hours / 5) * 100}%` }}
                    />
                  </div>
                  <div className="w-16 text-right text-sm text-text-primary">{item.hours}h</div>
                </div>
              ))}
            </div>
            <div className="mt-6 pt-6 border-t border-border-subtle">
              <div className="flex items-center justify-between text-sm">
                <span className="text-text-secondary">Total this week</span>
                <span className="text-text-primary font-bold">18.3 hours</span>
              </div>
            </div>
          </div>
        </div>
      </div>
    </div>
  )
}