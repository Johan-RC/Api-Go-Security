import React, { forwardRef, useState } from 'react';
import { StyleSheet, TextInput, TextInputProps, View, Pressable, TextStyle } from 'react-native';
import { colors, radius, spacing } from '@/theme';
import { Text } from '@/components/Text';

interface InputProps extends TextInputProps {
  label?: string;
  error?: string;
  leftIcon?: React.ReactNode;
  rightIcon?: React.ReactNode;
  onRightPress?: () => void;
}

// En web, el navegador pinta un recuadro de foco (normalmente negro). 'none'
// no está tipado en RN, pero es válido en CSS y elimina ese recuadro feo.
const focusReset = {
  outlineStyle: 'none',
  outlineWidth: 0,
} as unknown as TextStyle;

export const Input = forwardRef<TextInput, InputProps>(function Input(
  { label, error, leftIcon, rightIcon, onRightPress, style, onFocus, onBlur, ...props },
  ref,
) {
  const [focused, setFocused] = useState(false);

  return (
    <View style={styles.wrapper}>
      {label ? (
        <Text variant="label" color="textSecondary" style={{ marginBottom: spacing.xs }}>
          {label}
        </Text>
      ) : null}
      <View
        style={[
          styles.container,
          error ? styles.containerError : null,
          focused && !error ? styles.containerFocused : null,
        ]}
      >
        {leftIcon ? (
          <View style={[styles.sideLeft, focused ? { opacity: 1 } : null]}>
            {leftIcon}
          </View>
        ) : null}
        <TextInput
          ref={ref}
          placeholderTextColor={colors.textMuted}
          selectionColor={colors.primary}
          cursorColor={colors.primary}
          onFocus={(e) => {
            setFocused(true);
            onFocus?.(e);
          }}
          onBlur={(e) => {
            setFocused(false);
            onBlur?.(e);
          }}
          style={[
            styles.input,
            focusReset,
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
    overflow: 'hidden',
  },
  containerFocused: {
    borderColor: colors.primary,
    backgroundColor: colors.primarySoft,
    shadowColor: colors.primary,
    shadowOpacity: 0.12,
    shadowRadius: 8,
    shadowOffset: { width: 0, height: 2 },
    elevation: 3,
  },
  containerError: { borderColor: colors.danger, backgroundColor: colors.surface },
  sideLeft: { paddingLeft: spacing.md, opacity: 0.85 },
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