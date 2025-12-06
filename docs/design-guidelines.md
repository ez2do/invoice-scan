# Design Guidelines: Invoice Scan MVP

_Generated: December 2024_

## Design Style

**Refined Financial Minimalism** — A clean, professional aesthetic that conveys trust and precision, appropriate for a financial/invoice application. Emphasis on clarity, consistency, and subtle depth.

## Color Palette

### Primary Colors
- **Primary 50**: `#eef2ff` - Lightest tint, backgrounds
- **Primary 100**: `#e0e7ff` - Light backgrounds, badges
- **Primary 200**: `#c7d2fe` - Selection highlights
- **Primary 500**: `#6366f1` - Hover states
- **Primary 600**: `#4f46e5` - Primary actions, buttons (DEFAULT)
- **Primary 700**: `#4338ca` - Hover on primary buttons
- **Primary 800**: `#3730a3` - Active/pressed states
- **Primary 900**: `#312e81` - Dark mode accents

### Accent Colors (Teal)
- **Accent 400**: `#2dd4bf` - Success indicators
- **Accent 500**: `#14b8a6` - Success badges (DEFAULT)
- **Accent 600**: `#0d9488` - Dark mode success

### Semantic Colors
- **Success**: `#10b981` - Completed states
- **Warning**: `#f59e0b` - Pending/processing states
- **Error**: `#ef4444` - Failed states, validation errors

### Surface Colors (Slate)
- **Surface 50**: `#f8fafc` - Light mode background
- **Surface 100**: `#f1f5f9` - Card backgrounds, secondary surfaces
- **Surface 200**: `#e2e8f0` - Borders, dividers
- **Surface 400**: `#94a3b8` - Placeholder text, icons
- **Surface 500**: `#64748b` - Secondary text
- **Surface 600**: `#475569` - Body text
- **Surface 700**: `#334155` - Dark mode surfaces
- **Surface 800**: `#1e293b` - Dark mode cards
- **Surface 900**: `#0f172a` - Dark mode text
- **Surface 950**: `#020617` - Dark mode background

### Color Psychology
The indigo primary color conveys trust, professionalism, and reliability — essential qualities for financial applications. Teal accents provide a fresh, modern contrast while maintaining professionalism.

## Typography System

### Font Family
- **Primary Font**: Plus Jakarta Sans (Google Fonts)
- **Fallbacks**: system-ui, sans-serif

### Font Hierarchy
| Element | Size | Weight | Usage |
|---------|------|--------|-------|
| Page Title | 18px (text-lg) | 600 (semibold) | Page headers |
| Section Header | 14px (text-sm) | 600 (semibold) | Section titles |
| Body | 14px (text-sm) | 400 (normal) | General content |
| Label | 12px (text-xs) | 500 (medium) | Form labels, metadata |
| Small | 12px (text-xs) | 400 (normal) | Helper text, badges |

### Typography Guidelines
- Line height: Default (1.5 for body text)
- Letter spacing: `tracking-tight` for titles, default for body
- Use `font-semibold` for emphasis, `font-medium` for labels

## Layout Principles

### Container
- Max width: 480px
- Max height: 960px
- Centered with `mx-auto`
- Full height with overflow handling

### Spacing Scale
| Token | Size | Usage |
|-------|------|-------|
| 1 | 4px | Micro spacing |
| 1.5 | 6px | Icon gaps |
| 2 | 8px | Tight spacing |
| 3 | 12px | Card gaps |
| 4 | 16px | Section padding |
| 6 | 24px | Large spacing |
| 8 | 32px | Extra large spacing |

### Responsive Breakpoints
- Mobile-first design (480px max container)
- PWA safe area support: `safe-top`, `safe-bottom`

## Component Styling

### Buttons
| Variant | Description |
|---------|-------------|
| Primary | Indigo background, white text, soft shadow |
| Secondary | Surface-100 background, border, surface text |
| Ghost | Transparent, hover shows surface-100 |
| Danger | Error-500 background, white text |

**Sizes**: `sm` (36px), `md` (44px), `lg` (48px), `xl` (56px)

**States**: 
- Hover: Darker shade, increased shadow
- Active: Scale 0.98, darkest shade
- Disabled: 50% opacity, no pointer events

### Cards
- Border radius: `rounded-2xl` (16px)
- Shadow: `shadow-soft` (subtle layered shadow)
- Padding: 16px
- Border: 1px surface-200 (light) / surface-800 (dark)
- Interactive cards: hover shadow increase, scale 0.98 on active

### Expandable Cards (Line Items)
For displaying tabular data with potentially long values on mobile:
- **Collapsed State**: Shows primary field with truncated preview of other fields
- **Expanded State**: Full vertical layout with all fields visible
- **Auto-expand Textarea**: For long text values (>50 chars), uses auto-sizing textarea
- **Visual Indicator**: Chevron icon rotates 180° when expanded
- **Animation**: `animate-fade-in` on expansion (0.3s ease-out)
- **Item Count Badge**: Shows total items count in section header

```
┌─────────────────────────────────┐
│ 📦 Product Name          ▼    │  ← Collapsed
│    Qty: 2 • Price: 100,000đ   │
└─────────────────────────────────┘

┌─────────────────────────────────┐
│ 📦 Product Name          ▲    │  ← Expanded
├─────────────────────────────────┤
│ Tên sản phẩm                   │
│ ┌─────────────────────────────┐ │
│ │ Long product description... │ │  ← Auto-expanding
│ │ that wraps to multiple lines│ │     textarea
│ └─────────────────────────────┘ │
│                                 │
│ Số lượng                       │
│ ┌─────────────────────────────┐ │
│ │ 2                           │ │  ← Standard input
│ └─────────────────────────────┘ │
└─────────────────────────────────┘
```

### Forms
- Input height: 40-44px
- Border radius: `rounded-xl` (12px)
- Focus: 2px primary ring with 30% opacity
- Error: Error-500 border and ring

### Badges
| Type | Colors |
|------|--------|
| Success | Green-50 bg, green-600 text |
| Warning | Yellow-50 bg, yellow-600 text |
| Error | Red-50 bg, red-600 text |
| Info | Primary-50 bg, primary-600 text |

### Icons
- **Library**: Lucide React
- **Sizes**: 16px (small), 20px (default), 24px (large)
- **Style**: Line icons, consistent stroke width

## Visual Hierarchy

### Emphasis Techniques
1. **Size**: Larger text/icons for primary elements
2. **Weight**: Semibold for important text
3. **Color**: Primary color for actions, surface-500 for secondary
4. **Position**: Important actions at bottom (mobile thumb-friendly)
5. **Spacing**: More whitespace around important elements

### Content Flow
- Mobile-optimized single-column layout
- Sticky headers with blur backdrop
- Floating action buttons for primary actions
- Bottom-aligned action buttons in forms

## Micro-Interactions

### Animation Timing
| Duration | Usage |
|----------|-------|
| 150ms | Fast feedback (hover, active states) |
| 200ms | UI transitions |
| 300ms | Content fade-in |
| 400ms | Page transitions |

### Animations
- `animate-fade-in`: Content appearance
- `animate-fade-in-up`: Staggered list items
- `animate-scale-in`: Modal/popup appearance
- `animate-spin`: Loading spinners

### Easing Curves
- Default: `ease-out` for most transitions
- Spring: `bounce-in` for playful interactions

### Interactive States
- Hover: Background color shift, shadow increase
- Active: Scale 0.95-0.98, darker background
- Focus: 2px ring with offset
- Loading: Spinning indicator, disabled state

## Accessibility

### Contrast Ratios
- Body text: 7:1+ (WCAG AAA)
- Large text: 4.5:1+ (WCAG AA)
- UI elements: 3:1+

### Focus Indicators
- 2px ring with 50% opacity primary color
- 2px offset from element

### Touch Targets
- Minimum size: 44x44px
- Comfortable size: 48x48px
- Icon buttons: 40x40px minimum

## Implementation Notes

### Technologies
- React 19 with TypeScript
- Tailwind CSS 3.4 with custom config
- Lucide React for icons
- Plus Jakarta Sans from Google Fonts

### CSS Architecture
- Utility-first with Tailwind
- Component classes in `@layer components`
- Custom properties for dynamic theming
- Dark mode via `dark:` variant

### Performance Considerations
- Font preloading via Google Fonts URL
- CSS animations prefer `transform` and `opacity`
- Backdrop blur used sparingly
- Lazy loading for images

### Browser Support
- Modern browsers (Chrome, Safari, Firefox, Edge)
- iOS 14+ Safari for PWA
- Android Chrome 90+

## Resources

### Design Files
- Colors defined in `tailwind.config.js`
- Global styles in `src/index.css`

### Related Documentation
- [Project Overview](./project-overview-pdr.md)
- [Code Standards](./code-standards.md)
- [System Architecture](./system-architecture.md)

