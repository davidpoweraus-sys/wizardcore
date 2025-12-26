# WizardCore - Educational Platform

A modern educational platform for learning cybersecurity and programming through hands-on exercises.

## Features

- **Interactive Learning** - Hands-on coding exercises with real-time feedback
- **Code Execution** - Safe, sandboxed code execution via Judge0
- **Progress Tracking** - Monitor student progress and achievements
- **Role-Based Access Control** - Separate interfaces for students and content creators
- **Pathway System** - Structured learning paths for different topics
- **Practice Arena** - Competitive coding challenges

## Tech Stack

- **Frontend**: Next.js 16 (React 19, TypeScript)
- **Backend**: Go (REST API)
- **Database**: PostgreSQL 15
- **Authentication**: Supabase Auth (GoTrue)
- **Cache**: Redis 7
- **Code Execution**: Judge0
- **Deployment**: Docker Swarm via Dokploy

## Quick Start

### Prerequisites

- Docker & Docker Compose
- [Dokploy](https://dokploy.com) (for production deployment)
- Node.js 18+ (for local development)
- Go 1.21+ (for backend development)

### Local Development

1. **Clone the repository**
   ```bash
   git clone https://github.com/yourusername/wizardcore.git
   cd wizardcore
   ```

2. **Set up environment variables**
   ```bash
   cp .env.production .env
   # Edit .env and replace all CHANGE_THIS values
   ```

3. **Start services**
   ```bash
   docker compose up -d
   ```

4. **Access the application**
   - Frontend: http://localhost:3000
   - Backend API: http://localhost:8080
   - Auth Service: http://localhost:9999
   - Judge0: http://localhost:2358

### Production Deployment with Dokploy

1. **Install Dokploy on your VPS**
   ```bash
   curl -sSL https://dokploy.com/install.sh | sh
   ```

2. **Configure environment in Dokploy**
   - Copy values from `.env.production`
   - Update domains and secrets
   - Set strong passwords

3. **Deploy via Dokploy UI**
   - Create new project
   - Point to your Git repository
   - Select `docker-compose.yml`
   - Enable Docker Swarm mode
   - Deploy

4. **Configure reverse proxy**
   - Set up domains in Dokploy
   - Enable SSL certificates
   - Configure routing

## Environment Variables

See `.env.production` for a complete list of required environment variables.

**Critical variables to change:**
- `POSTGRES_PASSWORD` - Main database password
- `SUPABASE_POSTGRES_PASSWORD` - Auth database password
- `SUPABASE_JWT_SECRET` - JWT signing secret (generate with `openssl rand -base64 64`)
- `REDIS_PASSWORD` - Redis password
- `JUDGE0_API_KEY` - Judge0 API key
- `JUDGE0_POSTGRES_PASSWORD` - Judge0 database password

## Architecture

```
┌─────────────┐     ┌─────────────┐     ┌──────────────┐
│   Next.js   │────▶│   Go API    │────▶│  PostgreSQL  │
│  Frontend   │     │   Backend   │     │   Database   │
└─────────────┘     └─────────────┘     └──────────────┘
       │                   │                     
       │                   ▼                     
       │            ┌─────────────┐              
       │            │    Redis    │              
       │            │    Cache    │              
       │            └─────────────┘              
       │                                         
       ▼                   ▼                     
┌─────────────┐     ┌─────────────┐              
│  Supabase   │     │   Judge0    │              
│    Auth     │     │   Engine    │              
└─────────────┘     └─────────────┘              
```

## Scaling with Docker Swarm

The `docker-compose.yml` is optimized for Docker Swarm with:
- **Rolling updates** - Zero-downtime deployments
- **Health checks** - Auto-restart failed services
- **Replicas** - Frontend (2x), Backend (2x), Judge0 Workers (2x)
- **Placement constraints** - Databases pinned to manager node
- **Overlay networking** - Service discovery across nodes

### Single Node Deployment
Good for: Development, staging, <10K users
```bash
docker swarm init
docker stack deploy -c docker-compose.yml wizardcore
```

### Multi-Node Deployment
Good for: Production, >10K users, high availability

Via Dokploy:
1. Add worker nodes in Dokploy UI
2. Dokploy automatically configures Swarm
3. Services scale across nodes automatically

## Project Structure

```
wizardcore/
├── app/                    # Next.js app directory
│   ├── (auth)/            # Auth pages (login, register)
│   ├── api/               # API routes
│   ├── creator/           # Content creator interface
│   └── dashboard/         # Student dashboard
├── components/            # React components
├── lib/                   # Utilities and libraries
├── wizardcore-backend/    # Go backend
│   ├── cmd/              # Entry points
│   ├── internal/         # Business logic
│   └── pkg/              # Shared packages
├── init-scripts/         # Database initialization
├── docker-compose.yml    # Production deployment
└── .env.production       # Environment template
```

## API Documentation

Backend API runs on port 8080 with the following endpoints:

- `GET /health` - Health check
- `POST /api/auth/*` - Authentication
- `GET /api/pathways` - Learning pathways
- `GET /api/exercises` - Exercise list
- `POST /api/submissions` - Submit code
- `GET /api/progress` - User progress
- `GET /api/leaderboard` - Rankings

See `RBAC_API_QUICK_REFERENCE.md` for complete API documentation.

## Contributing

1. Fork the repository
2. Create a feature branch (`git checkout -b feature/amazing-feature`)
3. Commit your changes (`git commit -m 'Add amazing feature'`)
4. Push to the branch (`git push origin feature/amazing-feature`)
5. Open a Pull Request

## License

This project is licensed under the MIT License - see the LICENSE file for details.

## Support

- Documentation: See `Documentation/` folder
- Issues: GitHub Issues
- Discord: [Join our community](#)

## Security

- Never commit `.env` files
- Rotate secrets regularly
- Use strong passwords (32+ characters)
- Keep dependencies updated
- Enable firewall rules in production

## Acknowledgments

- [Next.js](https://nextjs.org/) - React framework
- [Supabase](https://supabase.com/) - Authentication
- [Judge0](https://judge0.com/) - Code execution
- [Dokploy](https://dokploy.com/) - Deployment platform
