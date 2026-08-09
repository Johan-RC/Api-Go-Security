export const colors = {
  primary: '#00843D',
  primaryDark: '#006B31',
  primaryDarker: '#004F25',
  primarySoft: '#E6F4EC',
  primaryOutline: '#A8E0C1',
  background: '#F5F9F7',
  surface: '#FFFFFF',
  text: '#0B2E1E',
  textSecondary: '#47685A',
  textMuted: '#7E948A',
  border: '#DCE8E1',
  success: '#16A34A',
  successSoft: '#DCFCE7',
  successStrong: '#15803D',
  danger: '#DC2626',
  dangerSoft: '#FEE2E2',
  dangerStrong: '#B91C1C',
  warning: '#D97706',
  warningSoft: '#FEF3C7',
  warningStrong: '#B45309',
  info: '#0284C7',
  infoSoft: '#E0F2FE',
  infoStrong: '#0369A1',
  accent: '#F2A900',
  accentSoft: '#FEF8E4',
  overlay: 'rgba(9, 45, 29, 0.45)',
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
  h1: { fontSize: 30, fontWeight: '800' as const, lineHeight: 38, letterSpacing: -0.5 },
  h2: { fontSize: 22, fontWeight: '700' as const, lineHeight: 28, letterSpacing: -0.3 },
  h3: { fontSize: 18, fontWeight: '600' as const, lineHeight: 24 },
  body: { fontSize: 15, lineHeight: 22 },
  caption: { fontSize: 12, lineHeight: 16 },
  label: { fontSize: 13, fontWeight: '600' as const, lineHeight: 18 },
} as const;

export const shadows = {
  sm: {
    shadowColor: '#0B2E1E',
    shadowOpacity: 0.06,
    shadowRadius: 10,
    shadowOffset: { width: 0, height: 3 },
    elevation: 2,
  },
  md: {
    shadowColor: '#0B2E1E',
    shadowOpacity: 0.12,
    shadowRadius: 18,
    shadowOffset: { width: 0, height: 8 },
    elevation: 5,
  },
} as const;