import React, { useState } from 'react';
import { StyleSheet, View, Modal, TouchableOpacity, FlatList, ViewStyle } from 'react-native';
import { Feather } from '@expo/vector-icons';
import { SafeAreaView } from 'react-native-safe-area-context';
import { colors, radius, spacing, shadows } from '@/theme';
import { Text } from '@/components/Text';

export interface SelectOption<T extends string | number> {
  label: string;
  value: T;
}

interface SelectProps<T extends string | number> {
  label?: string;
  placeholder?: string;
  value: T | null;
  options: SelectOption<T>[];
  onChange: (value: T) => void;
  error?: string;
  style?: ViewStyle;
}

export function Select<T extends string | number>({
  label,
  placeholder = 'Selecciona una opción',
  value,
  options,
  onChange,
  error,
  style,
}: SelectProps<T>) {
  const [open, setOpen] = useState(false);
  const selected = options.find((o) => o.value === value);

  return (
    <View style={[styles.wrapper, style]}>
      {label ? (
        <Text variant="label" color="textSecondary" style={{ marginBottom: spacing.xs }}>
          {label}
        </Text>
      ) : null}

      <TouchableOpacity activeOpacity={0.8} style={[styles.trigger, error ? styles.triggerError : null]} onPress={() => setOpen(true)}>
        <Text variant="body" style={{ color: selected ? colors.text : colors.textMuted }}>
          {selected?.label ?? placeholder}
        </Text>
        <Feather name="chevron-down" size={18} color={colors.textMuted} />
      </TouchableOpacity>

      {error ? (
        <Text variant="caption" color="danger" style={{ marginTop: spacing.xs }}>
          {error}
        </Text>
      ) : null}

      <Modal visible={open} transparent animationType="slide" onRequestClose={() => setOpen(false)}>
        <View style={styles.backdrop}>
          <SafeAreaView style={styles.sheet}>
            <View style={styles.handle} />
            <Text variant="h3" style={{ paddingHorizontal: spacing.lg, marginBottom: spacing.md }}>
              {label ?? 'Selecciona una opción'}
            </Text>
            <FlatList
              data={options}
              keyExtractor={(item) => String(item.value)}
              contentContainerStyle={{ paddingBottom: spacing.xl }}
              renderItem={({ item }) => {
                const isSelected = item.value === value;
                return (
                  <TouchableOpacity
                    style={[styles.option, isSelected ? styles.optionSelected : null]}
                    onPress={() => {
                      onChange(item.value);
                      setOpen(false);
                    }}
                  >
                    <Text variant="body" bold={isSelected} style={{ color: isSelected ? colors.primary : colors.text, flex: 1 }}>
                      {item.label}
                    </Text>
                    {isSelected ? <Feather name="check" size={18} color={colors.primary} /> : null}
                  </TouchableOpacity>
                );
              }}
            />
            <TouchableOpacity style={styles.cancel} onPress={() => setOpen(false)}>
              <Text variant="label" color="textSecondary">Cancelar</Text>
            </TouchableOpacity>
          </SafeAreaView>
        </View>
      </Modal>
    </View>
  );
}

const styles = StyleSheet.create({
  wrapper: { marginBottom: spacing.md },
  trigger: {
    flexDirection: 'row',
    alignItems: 'center',
    justifyContent: 'space-between',
    backgroundColor: colors.surface,
    borderWidth: 1.5,
    borderColor: colors.border,
    borderRadius: radius.md,
    paddingHorizontal: spacing.md,
    minHeight: 54,
  },
  triggerError: { borderColor: colors.danger },
  backdrop: { flex: 1, backgroundColor: colors.overlay, justifyContent: 'flex-end' },
  sheet: {
    backgroundColor: colors.surface,
    borderTopLeftRadius: radius.lg,
    borderTopRightRadius: radius.lg,
    paddingTop: spacing.sm,
    ...shadows.md,
  },
  handle: {
    alignSelf: 'center',
    width: 40,
    height: 4,
    borderRadius: 2,
    backgroundColor: colors.border,
    marginBottom: spacing.md,
  },
  option: {
    flexDirection: 'row',
    alignItems: 'center',
    paddingHorizontal: spacing.lg,
    paddingVertical: spacing.lg,
    gap: spacing.sm,
  },
  optionSelected: { backgroundColor: colors.primarySoft },
  cancel: {
    alignItems: 'center',
    paddingVertical: spacing.md,
    borderTopWidth: 1,
    borderTopColor: colors.border,
  },
});