import React from 'react';
import { Text as RNText, StyleSheet, TextProps as RNTextProps } from 'react-native';
import { colors, typography, ColorName } from '@/theme';

export interface TextProps extends RNTextProps {
  variant?: keyof typeof typography;
  color?: ColorName;
  align?: 'auto' | 'left' | 'right' | 'center';
  bold?: boolean;
}

export function Text({ variant = 'body', color = 'text', align, bold, style, ...props }: TextProps) {
  return (
    <RNText
      {...props}
      style={[
        styles.base,
        typography[variant],
        { color: colors[color] },
        align ? { textAlign: align } : null,
        bold ? styles.bold : null,
        style,
      ]}
    />
  );
}

const styles = StyleSheet.create({
  base: {},
  bold: { fontWeight: '700' },
});