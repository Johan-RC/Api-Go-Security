import React, { forwardRef } from 'react';
import { StyleSheet, TextInput, TextInputProps, View, Pressable } from 'react-native';
import { colors, radius, spacing } from '@/theme';
import { Text } from '@/components/Text';

interface InputProps extends TextInputProps {
  label?: string;
  error?: string;
  leftIcon?: React.ReactNode;
  rightIcon?: React.ReactNode;
  onRightPress?: () => void;
}

export const Input = forwardRef<TextInput, InputProps>(function Input(
  { label, error, leftIcon, rightIcon, onRightPress, style, ...props },
  ref,
) {
  return (
    <View style={styles.wrapper}>
      {label ? (
        <Text variant="label" color="textSecondary" style={{ marginBottom: spacing.xs }}>
          {label}
        </Text>
      ) : null}
      <View style={[styles.container, error ? styles.containerError : null]}>
        {leftIcon ? <View style={styles.sideLeft}>{leftIcon}</View> : null}
        <TextInput
          ref={ref}
          placeholderTextColor={colors.textMuted}
          style={[
            styles.input,
            leftIcon ? styles.inputSideLeft : null,
            rightIcon ? styles.inputSideRight : null,
            style,
          ]}
          {...props}
        />
        {rightIcon ? (
          onRightPress ? (
            <Pressable onPress={onRightPress} hitSlop={8} style={styles.sideRight}>
              {rightIcon}
            </Pressable>
          ) : (
            <View style={styles.sideRight}>{rightIcon}</View>
          )
        ) : null}
      </View>
      {error ? (
        <Text variant="caption" color="danger" style={{ marginTop: spacing.xs }}>
          {error}
        </Text>
      ) : null}
    </View>
  );
});

const styles = StyleSheet.create({
  wrapper: { marginBottom: spacing.md },
  container: {
    flexDirection: 'row',
    alignItems: 'center',
    backgroundColor: colors.surface,
    borderWidth: 1.5,
    borderColor: colors.border,
    borderRadius: radius.md,
    minHeight: 54,
  },
  containerError: { borderColor: colors.danger },
  sideLeft: { paddingLeft: spacing.md },
  sideRight: { paddingRight: spacing.md },
  input: {
    flex: 1,
    paddingHorizontal: spacing.md,
    paddingVertical: spacing.md,
    color: colors.text,
    fontSize: 15,
  },
  inputSideLeft: { paddingLeft: spacing.sm },
  inputSideRight: { paddingRight: spacing.sm },
});