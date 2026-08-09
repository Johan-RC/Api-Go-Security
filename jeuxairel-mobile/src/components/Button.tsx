import React from 'react';
import {
  StyleSheet,
  TouchableOpacity,
  ActivityIndicator,
  TouchableOpacityProps,
  View,
} from 'react-native';
import { colors, radius, spacing } from '@/theme';
import { Text } from '@/components/Text';

type Variant = 'primary' | 'secondary' | 'danger' | 'ghost' | 'outline';

interface ButtonProps extends TouchableOpacityProps {
  variant?: Variant;
  loading?: boolean;
  label: string;
  icon?: React.ReactNode;
}

export function Button({
  variant = 'primary',
  loading = false,
  label,
  icon,
  disabled,
  style,
  ...props
}: ButtonProps) {
  const palette = variants[variant];

  return (
    <TouchableOpacity
      activeOpacity={0.85}
      disabled={disabled || loading}
      style={[styles.base, { backgroundColor: palette.bg }, disabled || loading ? styles.disabled : null, style]}
      {...props}
    >
      {loading ? (
        <ActivityIndicator color={palette.fg} />
      ) : (
        <View style={styles.content}>
          {icon}
          <Text variant="label" bold style={{ color: palette.fg }}>
            {label}
          </Text>
        </View>
      )}
    </TouchableOpacity>
  );
}

const variants: Record<Variant, { bg: string; fg: string }> = {
  primary: { bg: colors.primary, fg: '#FFFFFF' },
  secondary: { bg: colors.primarySoft, fg: colors.primary },
  danger: { bg: colors.danger, fg: '#FFFFFF' },
  ghost: { bg: 'transparent', fg: colors.primary },
  outline: { bg: colors.surface, fg: colors.textSecondary },
};

const styles = StyleSheet.create({
  base: {
    borderRadius: radius.md,
    paddingVertical: spacing.md,
    paddingHorizontal: spacing.lg,
    alignItems: 'center',
    justifyContent: 'center',
    minHeight: 48,
  },
  content: { flexDirection: 'row', alignItems: 'center', gap: spacing.sm },
  disabled: { opacity: 0.5 },
});