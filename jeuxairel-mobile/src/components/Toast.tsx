import React, { createContext, useCallback, useContext, useRef, useState } from 'react';
import { StyleSheet, Animated, View } from 'react-native';
import { Feather } from '@expo/vector-icons';
import { colors, radius, spacing } from '@/theme';
import { Text } from '@/components/Text';

type ToastType = 'success' | 'error' | 'info';

interface ToastState {
  visible: boolean;
  type: ToastType;
  message: string;
}

interface ToastContextValue {
  show: (message: string, type?: ToastType) => void;
}

const ToastContext = createContext<ToastContextValue | undefined>(undefined);

const icons: Record<ToastType, { icon: 'check-circle' | 'alert-circle' | 'info'; color: string }> = {
  success: { icon: 'check-circle', color: colors.success },
  error: { icon: 'alert-circle', color: colors.danger },
  info: { icon: 'info', color: colors.info },
};

export function ToastProvider({ children }: { children: React.ReactNode }) {
  const [toast, setToast] = useState<ToastState>({ visible: false, type: 'info', message: '' });
  const opacity = useRef(new Animated.Value(0)).current;
  const timer = useRef<ReturnType<typeof setTimeout> | null>(null);

  const show = useCallback(
    (message: string, type: ToastType = 'info') => {
      if (timer.current) clearTimeout(timer.current);
      setToast({ visible: true, type, message });

      Animated.timing(opacity, { toValue: 1, duration: 200, useNativeDriver: true }).start();

      timer.current = setTimeout(() => {
        Animated.timing(opacity, { toValue: 0, duration: 250, useNativeDriver: true }).start(() => {
          setToast((prev) => ({ ...prev, visible: false }));
        });
      }, 3200);
    },
    [opacity],
  );

  const palette = icons[toast.type];

  return (
    <ToastContext.Provider value={{ show }}>
      {children}
      {toast.visible ? (
        <Animated.View pointerEvents="none" style={[styles.container, { opacity }]}>
          <View style={styles.toast}>
            <Feather name={palette.icon} size={20} color={palette.color} />
            <Text variant="body" style={{ color: colors.text, flex: 1 }}>
              {toast.message}
            </Text>
          </View>
        </Animated.View>
      ) : null}
    </ToastContext.Provider>
  );
}

export function useToast(): ToastContextValue {
  const ctx = useContext(ToastContext);
  if (!ctx) throw new Error('useToast debe usarse dentro de <ToastProvider>');
  return ctx;
}

const styles = StyleSheet.create({
  container: {
    position: 'absolute',
    top: 0,
    left: 0,
    right: 0,
    alignItems: 'center',
    zIndex: 1000,
  },
  toast: {
    flexDirection: 'row',
    alignItems: 'center',
    gap: spacing.sm,
    backgroundColor: colors.surface,
    borderRadius: radius.md,
    marginTop: 8,
    marginHorizontal: spacing.lg,
    paddingHorizontal: spacing.lg,
    paddingVertical: spacing.md,
    borderWidth: 1,
    borderColor: colors.border,
    maxWidth: 480,
    width: '92%',
  },
});