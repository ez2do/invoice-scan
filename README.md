# Invoice Scan MVP

A Progressive Web Application for scanning and extracting data from invoices using AI-powered vision recognition. Optimized for Vietnamese language invoices (both printed and handwritten) with flexible schema extraction.

## 📋 Overview

Invoice Scan MVP enables finance and accounting staff to quickly extract structured data from physical invoices by:
- Capturing invoice images using mobile device cameras
- Automatically extracting data using Google Gemini Vision API
- Displaying extracted data alongside original images for verification
- Supporting various invoice formats without predefined schemas

## 🚀 Quick Start

### Prerequisites

- Node.js 18+ and npm/pnpm
- Modern mobile device with camera (iOS 14+ or Android 8+)
- Google Gemini API key (for backend, when implemented)

### Frontend Setup

```bash
# Navigate to frontend directory
cd frontend

# Install dependencies
npm install

# Start development server
npm run dev

# Build for production
npm run build
```

### Environment Configuration

Create `frontend/.env.local`:
```env
VITE_API_URL=http://localhost:3000/api
```

### Mobile Development

For camera access, HTTPS is required. Use one of these options:

1. **Basic SSL (Development)**: Automatically handled by Vite plugin
2. **Custom SSL**: Place certificates in `frontend/ssl/` directory
3. **Mobile Dev Script**: Use `frontend/mobile-dev-http.sh` for local HTTPS

## 📱 Features

- **Camera Integration**: Native camera access via browser APIs
- **AI-Powered Extraction**: Gemini Vision API for flexible data extraction
- **Progressive Web App**: Installable, offline-capable mobile experience
- **Vietnamese Support**: Optimized for Vietnamese invoices (printed and handwritten)
- **Flexible Schema**: Adapts to various invoice formats automatically
- **Data Verification**: Side-by-side image and data display for accuracy checking
- **Mobile-First Design**: Responsive UI optimized for mobile devices

## 🏗️ Architecture

### Technology Stack

**Frontend**:
- React 19.2.0 with TypeScript
- Vite 7.2.4 for build tooling
- React Router 7.9.6 for navigation
- TanStack Query 5.90.11 for data fetching
- Zustand 5.0.8 for state management
- Tailwind CSS 3.4.18 for styling
- vite-plugin-pwa for PWA support

**Backend** (Planned):
- Go 1.21+ with Gin framework
- Google Gemini Go SDK
- RESTful API design

### Project Structure

```
invoice-scan/
├── frontend/              # React PWA application
│   ├── src/
│   │   ├── components/   # React components
│   │   ├── hooks/        # Custom hooks
│   │   ├── lib/          # API client and utilities
│   │   ├── pages/        # Page components
│   │   ├── stores/       # Zustand stores
│   │   └── types/        # TypeScript definitions
│   └── ...
├── backend/              # Backend API (planned)
└── docs/                 # Project documentation
    ├── project-overview-pdr.md
    ├── codebase-summary.md
    ├── code-standards.md
    └── system-architecture.md
```

## 📚 Documentation

Comprehensive documentation is available in the `docs/` directory:

- **[Project Overview & PDR](./docs/project-overview-pdr.md)**: Product requirements, goals, and success criteria
- **[Codebase Summary](./docs/codebase-summary.md)**: Codebase structure, components, and technologies
- **[Code Standards](./docs/code-standards.md)**: Coding conventions and best practices
- **[System Architecture](./docs/system-architecture.md)**: Technical architecture and design decisions

## 🔄 Development Workflow

### Running Locally

1. **Start Frontend**:
   ```bash
   cd frontend
   npm run dev
   ```

2. **Access Application**:
   - Desktop: `https://localhost:5173`
   - Mobile: Use your local IP address (e.g., `https://192.168.1.100:5173`)

3. **Camera Access**:
   - Grant camera permissions when prompted
   - Position invoice in frame guide
   - Capture and review image

### Building for Production

```bash
cd frontend
npm run build
```

Output will be in `frontend/dist/` directory, ready for deployment.

## 🧪 Testing

### Current Status
- Manual testing approach
- Test suite planned for future implementation

### Planned Testing
- Unit tests for hooks and utilities
- Component tests for pages
- Integration tests for API client
- E2E tests for user flows

## 🔒 Security

### Current Security Measures
- HTTPS required for camera access
- Input validation in components
- Environment variable configuration
- No image storage (privacy-first)

### Planned Security
- CORS configuration on backend
- API key protection
- Request validation
- Error sanitization

## 📊 Performance

### Targets
- Image upload + extraction: < 5 seconds (p95)
- PWA initial load: < 3 seconds on 4G
- Camera to preview: < 500ms

### Optimizations
- Image compression before upload
- Code splitting (vendor, router, query chunks)
- Service worker caching
- Lazy loading (planned)

## 🛣️ Roadmap

### Phase 1: MVP (Current)
- ✅ Frontend PWA setup
- ✅ Camera capture functionality
- ✅ Basic UI components
- ✅ Routing and navigation
- 🔄 Backend API integration
- 🔄 Gemini API integration
- 🔄 Data extraction flow

### Phase 2: Enhancement
- 📋 Image quality detection
- 📋 Batch processing support
- 📋 Export functionality
- 📋 Data persistence

### Phase 3: Advanced Features
- 📋 Multi-language support expansion
- 📋 Advanced table recognition
- 📋 Historical invoice management
- 📋 Integration with accounting software

## 🤝 Contributing

### Code Standards
- Follow TypeScript strict mode
- Adhere to ESLint configuration
- Use conventional commit messages
- Write self-documenting code
- See [Code Standards](./docs/code-standards.md) for details

### Git Workflow
- Use feature branches
- Follow conventional commits
- Keep commits focused and atomic
- Update documentation as needed

## 📝 License

[License information to be added]

## 🙏 Acknowledgments

- Google Gemini API for vision recognition
- React and Vite communities
- PWA documentation and best practices

## 📞 Support

For questions or issues:
- Review documentation in `docs/` directory
- Open an issue for bugs or feature requests

---

**Status**: Active Development | **Version**: 1.0.0 | **Last Updated**: 2025-01-26
