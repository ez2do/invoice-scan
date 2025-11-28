# Invoice Scanner - Frontend

A React-based Progressive Web App for scanning and extracting data from invoices using AI.

## 🚀 Quick Start

```bash
# Install dependencies
npm install

# Start development server
npm run dev

# Build for production
npm run build

# Preview production build
npm run preview
```

## 📱 Features

- **Camera Integration**: Native camera access for invoice scanning
- **AI-Powered Extraction**: Flexible data extraction using Gemini API
- **Progressive Web App**: Installable mobile experience
- **Responsive Design**: Mobile-first design with Tailwind CSS
- **Type Safety**: Full TypeScript support
- **Modern Stack**: React 18, Vite, TanStack Query

## 🏗️ Architecture

```
src/
├── components/          # Reusable UI components
│   ├── ui/             # Base components (Button, Input)
│   ├── camera/         # Camera-specific components
│   └── invoice/        # Invoice data components
├── hooks/              # Custom React hooks
├── lib/                # Utilities and API client
├── pages/              # Page components
├── stores/             # Zustand state management
├── types/              # TypeScript definitions
└── utils/              # Helper functions
```

## 🛠️ Configuration

### Environment Variables

Copy `.env.example` to `.env.local` and configure:

```env
VITE_API_URL=http://localhost:3000/api
```

### PWA Configuration

The app is configured as a PWA with:
- Service Worker for caching
- Web App Manifest
- Offline support
- Installation prompts

## 📊 Data Flow

1. **Home Page** → Welcome screen with scan button
2. **Scan Page** → Camera interface for capturing invoices  
3. **Verify Page** → Display extracted data for verification

## 🎨 Styling

- **Tailwind CSS**: Utility-first CSS framework
- **Mobile-First**: Responsive design starting from 320px
- **Custom Components**: Consistent design system

## 🔧 Development

### Available Scripts

- `npm run dev` - Start development server
- `npm run build` - Build for production
- `npm run preview` - Preview production build
- `npm run lint` - Run ESLint

### State Management

Uses Zustand for lightweight state management:
- Current image storage
- Extracted invoice data
- Loading and error states
- Edit capabilities

### API Integration

TanStack Query for server state:
- Invoice extraction calls
- Error handling and retries
- Background refetching
- Cache management

## 🚀 Next Steps

1. **Backend Integration**: Connect to Node.js API server
2. **Device Testing**: Test camera functionality on mobile devices
3. **API Configuration**: Set correct backend URL in environment
4. **SSL Setup**: Configure HTTPS for camera access in production
5. **Icon Generation**: Create PWA icons for all required sizes

---

✅ **Frontend Setup Complete!** 
Ready to integrate with backend and start scanning invoices.