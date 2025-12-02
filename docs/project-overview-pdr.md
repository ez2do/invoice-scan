# Project Overview & Product Development Requirements (PDR)

**Project Name**: Invoice Scan MVP
**Version**: 1.0.0
**Last Updated**: 2025-01-26
**Status**: Active Development

## Executive Summary

Invoice Scan MVP is a Progressive Web Application (PWA) that enables users to scan physical invoices using their mobile device camera and automatically extract structured data using AI-powered vision recognition. The application is designed specifically for Vietnamese invoices (both printed and handwritten) and provides a flexible data extraction system that adapts to various invoice formats without requiring predefined schemas.

## Project Purpose

### Vision
Enable finance and accounting staff to quickly and accurately extract data from physical invoices, eliminating manual transcription errors and reducing processing time.

### Mission
Provide a mobile-first, installable PWA that:
- Captures invoice images using device cameras
- Extracts structured data using Gemini Vision API
- Displays extracted data alongside original images for verification
- Supports flexible invoice formats without rigid schemas
- Optimizes for Vietnamese language content (printed and handwritten)

### Value Proposition
- **10x Faster Processing**: Automated extraction vs manual data entry
- **Higher Accuracy**: AI-powered recognition reduces transcription errors
- **Mobile-First**: No app store required, works on any modern smartphone
- **Cost-Effective**: Low per-invoice processing cost (~$0.005 per extraction)
- **Flexible**: Adapts to various invoice formats automatically

## Target Users

### Primary Users
1. **Finance Staff**: Internal team members processing invoices regularly
2. **Accounting Personnel**: Staff who need to digitize invoice data
3. **Administrative Staff**: Users handling invoice documentation

### User Personas

**Persona 1: Finance Staff Member**
- **Needs**: Fast, accurate invoice data extraction
- **Pain Points**: Manual typing errors, time-consuming data entry, difficulty reading handwritten Vietnamese
- **Solution**: Camera-based scanning with AI extraction and visual verification

**Persona 2: Accounting Manager**
- **Needs**: Consistent data format, verification capability
- **Pain Points**: Inconsistent manual entry, hard to verify accuracy
- **Solution**: Structured extraction with side-by-side image comparison

## Key Features & Capabilities

### 1. Camera-Based Invoice Capture
- Native camera access via browser APIs
- Real-time preview with frame guide
- Image capture and compression
- Support for both front and back cameras
- Mobile-optimized interface

### 2. AI-Powered Data Extraction
- Gemini 1.5 Flash Vision API integration
- Flexible schema extraction (key-value pairs, tables, summaries)
- Vietnamese language support (printed and handwritten)
- Confidence scoring for extracted fields
- Error handling and retry mechanisms

### 3. Progressive Web App (PWA)
- Installable on mobile devices
- Offline-capable service worker
- App-like experience without app stores
- Responsive design (320px+ screens)
- Dark mode support

### 4. Data Verification Interface
- Side-by-side image and data display
- Editable extracted fields
- Organized sections (key-value, table, summary)
- Visual comparison for accuracy checking

### 5. Flexible Data Model
- Dynamic key-value pair extraction
- Variable-column table support
- Summary section extraction
- No fixed schema requirements
- Adapts to invoice format automatically

## Technical Requirements

### Functional Requirements

**FR1: Camera Access**
- Access device camera via MediaDevices API
- Support environment-facing camera (back camera)
- Handle camera permission requests
- Provide visual feedback during capture

**FR2: Image Processing**
- Capture image from video stream
- Compress images to reduce upload size
- Convert to base64 for API transmission
- Validate image quality before processing

**FR3: Data Extraction**
- Send image to backend API
- Process through Gemini Vision API
- Extract structured data (key-value, table, summary)
- Return formatted JSON response
- Handle extraction errors gracefully

**FR4: Data Display**
- Show extracted data in organized sections
- Display original image alongside data
- Enable field editing
- Support scrolling for long invoices
- Responsive layout for mobile screens

**FR5: PWA Functionality**
- Service worker for offline support
- Web app manifest for installation
- App icons and splash screens
- Install prompt support
- Caching strategy for assets

### Non-Functional Requirements

**NFR1: Performance**
- Image upload + extraction: < 5 seconds (p95)
- PWA initial load: < 3 seconds on 4G
- Camera to preview: < 500ms
- Smooth 60fps camera preview

**NFR2: Reliability**
- Graceful error handling
- Retry mechanisms for failed extractions
- Offline fallback messaging
- Network error recovery

**NFR3: Usability**
- Intuitive camera interface
- Clear visual feedback
- Accessible controls
- Mobile-first design
- Support for screen sizes 320px+

**NFR4: Security**
- HTTPS required for camera access
- Secure API communication
- No image storage (privacy-first)
- Input validation and sanitization

**NFR5: Cost Efficiency**
- Cost per extraction < $0.005
- Optimized image compression
- Efficient API usage
- Minimal bandwidth consumption

## Success Metrics

### Functional Metrics
- Camera access success rate: > 95%
- Extraction success rate: > 90% (printed), > 75% (handwritten)
- Vietnamese text accuracy: > 90% (printed), > 75% (handwritten)
- PWA installability: Works on iOS 14+ and Android 8+

### Performance Metrics
- Average extraction time: < 5 seconds
- PWA load time: < 3 seconds (4G)
- Image compression ratio: > 70% reduction
- API response time: < 100ms (excluding Gemini)

### Quality Metrics
- User satisfaction: > 4.0/5.0
- Error rate: < 5%
- Field extraction accuracy: > 90% (printed invoices)
- Data verification time: < 30 seconds per invoice

## Technical Architecture

### Core Components

**1. Frontend Application**
- React 18+ with TypeScript
- Vite build tool
- React Router for navigation
- TanStack Query for data fetching
- Zustand for state management
- Tailwind CSS for styling

**2. Backend API** (Planned)
- RESTful API server
- Gemini Vision API integration
- Image processing pipeline
- Error handling and logging

**3. External Services**
- Google Gemini 1.5 Flash API
- Hosting infrastructure (self-hosted)

### Technology Stack

**Frontend**:
- React 19.2.0
- TypeScript 5.9.3
- Vite 7.2.4
- React Router 7.9.6
- TanStack Query 5.90.11
- Zustand 5.0.8
- Tailwind CSS 3.4.18
- vite-plugin-pwa 1.2.0

**Backend** (Planned):
- Go 1.21+ with Gin framework
- Gemini Go SDK
- Environment-based configuration

**Infrastructure**:
- Self-hosted server
- Docker containerization (planned)
- HTTPS/SSL certificates
- Reverse proxy (nginx, planned)

## Use Cases

### UC1: Scan and Extract Invoice
**Actor**: Finance Staff Member
**Goal**: Extract data from physical invoice
**Flow**:
1. Open PWA on mobile device
2. Navigate to scan page
3. Grant camera permission
4. Position invoice in frame
5. Capture image
6. Review captured image
7. Confirm extraction
8. View extracted data
9. Verify and edit if needed
10. Complete process

**Outcome**: Structured invoice data extracted and verified

### UC2: Verify Extracted Data
**Actor**: Finance Staff Member
**Goal**: Ensure extraction accuracy
**Flow**:
1. View extracted data alongside image
2. Compare each field with original
3. Edit incorrect fields
4. Verify table data accuracy
5. Check summary calculations
6. Confirm completion

**Outcome**: Verified and corrected invoice data

### UC3: Handle Extraction Errors
**Actor**: Finance Staff Member
**Goal**: Recover from extraction failures
**Flow**:
1. Receive extraction error notification
2. Review error message
3. Retry extraction
4. If still failing, check image quality
5. Retake photo if needed
6. Retry extraction

**Outcome**: Successful extraction or clear error guidance

## Constraints & Limitations

### Technical Constraints
- Requires HTTPS for camera access
- Modern browser with camera API support
- Internet connection for API calls
- Gemini API requires Google Cloud account
- Self-hosted infrastructure required

### Business Constraints
- MVP scope: No data persistence
- Internal tool: No public access
- Single language focus: Vietnamese (English fallback)
- No multi-user authentication in MVP

### Assumptions
- Users have modern smartphones (iOS 14+, Android 8+)
- Stable internet connection during scanning
- Invoice formats follow general structure (header, table, summary)
- Server has outbound internet access

## Risks & Mitigation

### Risk 1: Camera API Limitations
**Impact**: High
**Likelihood**: Medium
**Mitigation**: Progressive enhancement, fallback to file upload, clear error messages

### Risk 2: Gemini API Failures
**Impact**: High
**Likelihood**: Low
**Mitigation**: Retry logic, error handling, fallback messaging, API key rotation

### Risk 3: Poor Image Quality
**Impact**: Medium
**Likelihood**: High
**Mitigation**: Image quality validation, user guidance, retake option, compression optimization

### Risk 4: Vietnamese Handwriting Recognition
**Impact**: Medium
**Likelihood**: Medium
**Mitigation**: Gemini's Vietnamese support, confidence scoring, manual editing capability

### Risk 5: PWA Installation Issues
**Impact**: Low
**Likelihood**: Low
**Mitigation**: Comprehensive manifest, testing across devices, installation guidance

## Future Roadmap

### Phase 1: MVP (Current)
- ✅ Frontend PWA setup
- ✅ Camera capture functionality
- ✅ Basic UI components
- ✅ Routing and navigation
- 🔄 Backend API integration
- 🔄 Gemini API integration
- 🔄 Data extraction flow
- 🔄 Verification interface

### Phase 2: Enhancement
- 📋 Image quality detection
- 📋 Batch processing support
- 📋 Export functionality
- 📋 Data persistence
- 📋 User authentication

### Phase 3: Advanced Features
- 📋 Multi-language support expansion
- 📋 Advanced table recognition
- 📋 OCR confidence visualization
- 📋 Historical invoice management
- 📋 Integration with accounting software

## Dependencies & Integration

### Required Dependencies
- Google Gemini API access
- HTTPS-enabled server
- Modern web browser
- Mobile device with camera

### Optional Dependencies
- Docker for containerization
- Reverse proxy (nginx)
- SSL certificate (Let's Encrypt)
- Monitoring tools

## Compliance & Standards

### Coding Standards
- TypeScript strict mode
- ESLint configuration
- Prettier formatting
- Component-based architecture
- Mobile-first responsive design

### Security Standards
- HTTPS required
- No image storage
- API key protection
- Input validation
- CORS configuration

### Accessibility Standards
- WCAG 2.1 Level AA compliance
- Keyboard navigation support
- Screen reader compatibility
- High contrast mode support
- Touch target sizes (44x44px minimum)

## Glossary

- **PWA**: Progressive Web App - Web application with app-like features
- **Gemini API**: Google's multimodal AI API for vision and language tasks
- **Key-Value Pair**: Data structure with a key (field name) and value (field content)
- **Service Worker**: Background script enabling offline functionality
- **Base64**: Encoding scheme for binary data in text format
- **MediaDevices API**: Browser API for accessing camera and microphone

## Appendix

### Related Documentation
- [Codebase Summary](./codebase-summary.md)
- [Code Standards](./code-standards.md)
- [System Architecture](./system-architecture.md)

### External Resources
- [Gemini API Documentation](https://ai.google.dev/docs)
- [PWA Documentation](https://web.dev/progressive-web-apps/)
- [MediaDevices API](https://developer.mozilla.org/en-US/docs/Web/API/MediaDevices)

## Unresolved Questions

1. **Backend Implementation**: Final decision on Go vs Node.js backend?
2. **Deployment Strategy**: Docker vs direct deployment?
3. **Monitoring**: What monitoring tools to use?
4. **Error Tracking**: Integration with error tracking service?
5. **Analytics**: User analytics implementation?
