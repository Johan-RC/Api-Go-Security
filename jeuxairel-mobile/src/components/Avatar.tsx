import React from 'react';
import { StyleSheet, View, ViewStyle } from 'react-native';
import { colors } from '@/theme';
import { Text } from '@/components/Text';

interface AvatarProps {
  text: string;
  color?: string;
  size?: number;
  style?: ViewStyle;
}

export function Avatar({ text, color = colors.primary, size = 44, style }: AvatarProps) {
  return (
    <View style={[styles.avatar, { backgroundColor: color, width: size, height: size, borderRadius: size / 2 }, style]}>
      <Text variant="label" style={{ color: '#FFFFFF', fontSize: size * 0.35 }} bold>
        {text}
      </Text>
    </View>
  );
}

const styles = StyleSheet.create({
  avatar: {
    alignItems: 'center',
    justifyContent: 'center',
  },
});