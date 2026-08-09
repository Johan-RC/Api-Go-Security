import React from 'react';
import { StyleSheet, ActivityIndicator, View, ViewStyle } from 'react-native';
import { colors, spacing } from '@/theme';
import { Text } from '@/components/Text';

interface SpinnerProps {
  label?: string;
  size?: 'small' | 'large';
  style?: ViewStyle;
}

export function Spinner({ label, size = 'large', style }: SpinnerProps) {
  return (
    <View style={[styles.container, style]}>
      <ActivityIndicator size={size} color={colors.primary} />
      {label ? (
        <Text variant="body" color="textSecondary" style={{ marginTop: spacing.md }}>
          {label}
        </Text>
      ) : null}
    </View>
  );
}

const styles = StyleSheet.create({
  container: { flex: 1, alignItems: 'center', justifyContent: 'center' },
});