import React from 'react';
import { StyleSheet, View, TouchableOpacity, ViewStyle } from 'react-native';
import { Feather } from '@expo/vector-icons';
import { colors, spacing } from '@/theme';
import { Text } from '@/components/Text';

interface RowProps {
  title: string;
  subtitle?: string | null;
  description?: string | null;
  left?: React.ReactNode;
  right?: React.ReactNode;
  onPress?: () => void;
  style?: ViewStyle;
}

export function Row({ title, subtitle, description, left, right, onPress, style }: RowProps) {
  const content = (
    <>
      {left ? <View style={styles.left}>{left}</View> : null}
      <View style={styles.center}>
        <Text variant="body" numberOfLines={1}>
          {title}
        </Text>
        {subtitle ? (
          <Text variant="caption" color="textSecondary" numberOfLines={1}>
            {subtitle}
          </Text>
        ) : null}
        {description ? (
          <Text variant="caption" color="textMuted" numberOfLines={1}>
            {description}
          </Text>
        ) : null}
      </View>
      {right ? <View style={styles.right}>{right}</View> : null}
      {onPress ? <Feather name="chevron-right" size={18} color={colors.textMuted} /> : null}
    </>
  );

  if (onPress) {
    return (
      <TouchableOpacity activeOpacity={0.7} style={[styles.row, style]} onPress={onPress}>
        {content}
      </TouchableOpacity>
    );
  }
  return <View style={[styles.row, style]}>{content}</View>;
}

const styles = StyleSheet.create({
  row: {
    flexDirection: 'row',
    alignItems: 'center',
    paddingVertical: spacing.md,
    gap: spacing.md,
  },
  left: { flexShrink: 0 },
  center: { flex: 1, gap: 2 },
  right: { flexShrink: 0, alignItems: 'flex-end' },
});