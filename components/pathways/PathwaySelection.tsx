'use client'

import { useState } from 'react'
import { Check, Lock, Star, TrendingUp, Clock, Users } from 'lucide-react'

const pathways = [
  {
    id: 1,
    title: 'Python for Offensive Security',
    subtitle: 'The Hacker\'s Swiss Army Knife',
    description: 'Master Python for security tooling, exploit development, and automation.',
    level: 'Beginner',
    duration: '8 weeks',
    students: '1.2k',
    rating: 4.8,
    modules: 5,
    color: 'from-green-400 to-cyan-400',
    icon: '🐍',
    locked: false,
    progress: 0,
  },
  {
    id: 2,
    title: 'C & Assembly: The Exploit Developer\'s Core',
    subtitle: 'Understanding Memory Corruption',
    description: 'Deep dive into low‑level programming, memory corruption, and shellcoding.',
    level: 'Intermediate',
    duration: '10 weeks',
    students: '850',
    rating: 4.9,
    modules: 5,
    color: 'from-blue-400 to-purple-400',
    icon: '⚙️',
    locked: false,
    progress: 0,
  },
  {
    id: 3,
    title: 'JavaScript & Browser Exploitation',
    subtitle: 'The Modern Attack Surface',
    description: 'Exploit browser vulnerabilities, Node.js servers, and deobfuscate malware.',
    level: 'Intermediate',
    duration: '6 weeks',
    students: '720',
    rating: 4.7,
    modules: 4,
    color: 'from-yellow-400 to-orange-400',
    icon: '🌐',
    locked: false,
    progress: 0,
  },
  {
    id: 4,
    title: 'SQL & Database Exploitation',
    subtitle: 'Beyond \' OR 1=1--',
    description: 'Advanced SQL injection, database‑specific attacks, and post‑exploitation pivoting.',
    level: 'Intermediate',
    duration: '7 weeks',
    students: '640',
    rating: 4.6,
    modules: 4,
    color: 'from-red-400 to-pink-400',
    icon: '🗃️',
    locked: false,
    progress: 0,
  },
  {
    id: 5,
    title: 'Reverse Engineering for Exploit Development',
    subtitle: 'The Ultimate Synthesis',
    description: 'Static/dynamic analysis, Windows/Linux internals, and real‑world bug hunting.',
    level: 'Advanced',
    duration: '12 weeks',
    students: '480',
    rating: 4.9,
    modules: 5,
    color: 'from-purple-400 to-indigo-400',
    icon: '🔍',
    locked: true,
    progress: 0,
  },
  {
    id: 6,
    title: 'Rootkit Development',
    subtitle: 'The Art of Staying Hidden',
    description: 'Kernel‑level rootkits, evasion techniques, and covert communications.',
    level: 'Expert',
    duration: '14 weeks',
    students: '320',
    rating: 5.0,
    modules: 10,
    color: 'from-gray-700 to-black',
    icon: '👻',
    locked: true,
    progress: 0,
  },
]

export default function PathwaySelection() {
  const [selected, setSelected] = useState<number | null>(null)

  return (
    <div>
      <div className="mb-8">
        <h2 className="text-2xl font-bold text-text-primary">Available Pathways</h2>
        <p className="text-text-secondary mt-2">
          Each pathway is a complete course with projects, labs, and a capstone. Start with Python or jump to your skill level.
        </p>
      </div>

      <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-6">
        {pathways.map((pathway) => (
          <div
            key={pathway.id}
            className={`bg-bg-elevated border rounded-2xl p-6 flex flex-col h-full transition-all ${selected === pathway.id ? 'border-neon-cyan ring-2 ring-neon-cyan/20' : 'border-border-default hover:border-neon-lavender'} ${pathway.locked ? 'opacity-80' : ''}`}
            onClick={() => !pathway.locked && setSelected(pathway.id)}
          >
            <div className="flex justify-between items-start">
              <div>
                <div className={`w-14 h-14 rounded-xl bg-gradient-to-br ${pathway.color} flex items-center justify-center text-2xl mb-4`}>
                  {pathway.icon}
                </div>
                <h3 className="text-xl font-bold text-text-primary">{pathway.title}</h3>
                <p className="text-sm text-text-secondary mt-1">{pathway.subtitle}</p>
              </div>
              {pathway.locked && (
                <div className="p-2 bg-bg-tertiary rounded-lg">
                  <Lock className="w-5 h-5 text-text-muted" />
                </div>
              )}
            </div>

            <p className="text-text-secondary mt-4 flex-1">{pathway.description}</p>

            <div className="mt-6 space-y-3">
              <div className="flex items-center justify-between text-sm">
                <div className="flex items-center text-text-secondary">
                  <Clock className="w-4 h-4 mr-2" />
                  <span>{pathway.duration}</span>
                </div>
                <div className="flex items-center text-text-secondary">
                  <Users className="w-4 h-4 mr-2" />
                  <span>{pathway.students} students</span>
                </div>
              </div>

              <div className="flex items-center justify-between text-sm">
                <div className="flex items-center">
                  <Star className="w-4 h-4 text-yellow-400 mr-1" />
                  <span className="font-medium text-text-primary">{pathway.rating}</span>
                  <span className="text-text-muted ml-1">rating</span>
                </div>
                <div className="px-3 py-1 rounded-full bg-bg-tertiary text-text-secondary text-xs font-medium">
                  {pathway.level}
                </div>
              </div>

              <div className="pt-4 border-t border-border-subtle">
                <div className="flex justify-between items-center">
                  <span className="text-sm text-text-secondary">{pathway.modules} modules</span>
                  {pathway.locked ? (
                    <button className="px-4 py-2 bg-bg-tertiary text-text-secondary rounded-lg text-sm font-medium cursor-not-allowed">
                      Locked
                    </button>
                  ) : (
                    <button className="px-4 py-2 bg-gradient-to-r from-neon-cyan to-neon-lavender text-white rounded-lg text-sm font-medium hover:opacity-90 transition">
                      Enroll Now
                    </button>
                  )}
                </div>
              </div>
            </div>
          </div>
        ))}
      </div>

      <div className="mt-10 bg-gradient-to-r from-bg-elevated to-bg-tertiary border border-border-default rounded-2xl p-8">
        <div className="flex items-center">
          <div className="p-4 rounded-xl bg-gradient-to-br from-neon-pink to-neon-purple">
            <TrendingUp className="w-8 h-8 text-white" />
          </div>
          <div className="ml-6">
            <h3 className="text-xl font-bold text-text-primary">Not sure where to start?</h3>
            <p className="text-text-secondary mt-2">
              Take our 5‑minute skill assessment to get a personalized pathway recommendation.
            </p>
          </div>
          <button className="ml-auto px-6 py-3 border border-neon-cyan text-neon-cyan rounded-lg font-medium hover:bg-neon-cyan/10 transition">
            Take Assessment
          </button>
        </div>
      </div>
    </div>
  )
}