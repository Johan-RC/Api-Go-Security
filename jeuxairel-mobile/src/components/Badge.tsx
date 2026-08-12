import React from 'react';
import { StyleSheet, View, ViewStyle } from 'react-native';
import { colors, radius } from '@/theme';
import { Text } from '@/components/Text';

export type BadgeTone = 'primary' | 'success' | 'danger' | 'warning' | 'info' | 'neutral';

const tones: Record<BadgeTone, { bg: string; fg: string }> = {
  primary: { bg: colors.primarySoft, fg: colors.primaryLight },
  success: { bg: colors.successSoft, fg: colors.successStrong },
  danger: { bg: colors.dangerSoft, fg: colors.dangerStrong },
  warning: { bg: colors.warningSoft, fg: colors.warningStrong },
  info: { bg: colors.infoSoft, fg: colors.infoStrong },
  neutral: { bg: '#1C1C24', fg: colors.textSecondary },
};

interface BadgeProps {
  tone?: BadgeTone;
  label: string;
  style?: ViewStyle;
}

export function Badge({ tone = 'neutral', label, style }: BadgeProps) {
  const palette = tones[tone];
  return (
    <View style={[styles.badge, { backgroundColor: palette.bg }, style]}>
      <Text variant="caption" style={{ color: palette.fg }} bold>
        {label}
      </Text>
    </View>
  );
}

const styles = StyleSheet.create({
  badge: { paddingHorizontal: 8, paddingVertical: 4, borderRadius: radius.pill },
});