export const colors = {
  primary: '#EF4444',
  primaryLight: '#FCA5A5',
  primaryDark: '#DC2626',
  primaryDarker: '#7F1D1D',
  primarySoft: '#241014',
  primaryOutline: '#7F1D1D',
  background: '#0B0B0F',
  surface: '#15151C',
  text: '#F4F4F5',
  textSecondary: '#A3A3AD',
  textMuted: '#6B6B76',
  border: '#2A2A35',
  success: '#22C55E',
  successSoft: '#0E2A18',
  successStrong: '#4ADE80',
  danger: '#EF4444',
  dangerSoft: '#2A1215',
  dangerStrong: '#F87171',
  warning: '#F59E0B',
  warningSoft: '#2A2210',
  warningStrong: '#FBBF24',
  info: '#38BDF8',
  infoSoft: '#0A2430',
  infoStrong: '#7DD3FC',
  accent: '#F2A900',
  accentSoft: '#2A2108',
  overlay: 'rgba(0, 0, 0, 0.6)',
  gradientStart: '#EF4444',
  gradientMid: '#8F1D1D',
  gradientEnd: '#0A0A0A',
} as const;

export type ColorName = keyof typeof colors;

export const spacing = {
  xs: 4,
  sm: 8,
  md: 12,
  lg: 16,
  xl: 24,
  '2xl': 32,
} as const;

export const radius = {
  sm: 8,
  md: 12,
  lg: 16,
  xl: 22,
  pill: 999,
} as const;

export const typography = {
  h1: { fontSize: 28, fontWeight: '800' as const, lineHeight: 36, letterSpacing: -0.5 },
  h2: { fontSize: 22, fontWeight: '700' as const, lineHeight: 28, letterSpacing: -0.3 },
  h3: { fontSize: 18, fontWeight: '600' as const, lineHeight: 24 },
  body: { fontSize: 15, lineHeight: 22 },
  caption: { fontSize: 12, lineHeight: 16 },
  label: { fontSize: 13, fontWeight: '600' as const, lineHeight: 18 },
} as const;

export const shadows = {
  sm: {
    shadowColor: '#000000',
    shadowOpacity: 0.18,
    shadowRadius: 10,
    shadowOffset: { width: 0, height: 3 },
    elevation: 2,
  },
  md: {
    shadowColor: '#000000',
    shadowOpacity: 0.35,
    shadowRadius: 18,
    shadowOffset: { width: 0, height: 8 },
    elevation: 5,
  },
} as const;