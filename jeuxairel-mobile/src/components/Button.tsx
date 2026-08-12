import React, { useRef } from 'react';
import {
  StyleSheet,
  TouchableOpacity,
  ActivityIndicator,
  TouchableOpacityProps,
  View,
  Animated,
} from 'react-native';
import { colors, radius, shadows, spacing } from '@/theme';
import { pointer, nativeDriver } from '@/utils/web';
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
  onPressIn,
  onPressOut,
  ...props
}: ButtonProps) {
  const palette = variants[variant];
  const scale = useRef(new Animated.Value(1)).current;

  function handlePressIn(e: Parameters<NonNullable<TouchableOpacityProps['onPressIn']>>[0]) {
    Animated.spring(scale, { toValue: 0.96, speed: 40, bounciness: 0, useNativeDriver: nativeDriver }).start();
    onPressIn?.(e);
  }

  function handlePressOut(e: Parameters<NonNullable<TouchableOpacityProps['onPressOut']>>[0]) {
    Animated.spring(scale, { toValue: 1, speed: 40, bounciness: 0, useNativeDriver: nativeDriver }).start();
    onPressOut?.(e);
  }

  return (
    <TouchableOpacity
      activeOpacity={0.8}
      disabled={disabled || loading}
      onPressIn={handlePressIn}
      onPressOut={handlePressOut}
      style={[
        styles.base,
        pointer,
        { backgroundColor: palette.bg },
        variant === 'primary' ? styles.primaryShadow : null,
        disabled || loading ? styles.disabled : null,
        style,
      ]}
      {...props}
    >
      <Animated.View style={{ transform: [{ scale }] }}>
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
      </Animated.View>
    </TouchableOpacity>
  );
}

const variants: Record<Variant, { bg: string; fg: string }> = {
  primary: { bg: colors.primary, fg: '#FFFFFF' },
  secondary: { bg: colors.primarySoft, fg: colors.primaryLight },
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
    minHeight: 50,
  },
  primaryShadow: { ...shadows.md },
  content: { flexDirection: 'row', alignItems: 'center', gap: spacing.sm },
  disabled: { opacity: 0.5 },
});