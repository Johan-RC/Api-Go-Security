import React from 'react';
import { StyleSheet, View, KeyboardAvoidingView, Platform } from 'react-native';
import { SafeAreaView } from 'react-native-safe-area-context';
import { Feather } from '@expo/vector-icons';
import { colors, spacing } from '@/theme';
import { Text } from '@/components/Text';

interface AuthLayoutProps {
  title: string;
  subtitle?: string;
  icon?: keyof typeof Feather.glyphMap;
  children: React.ReactNode;
}

export function AuthLayout({ title, subtitle, icon = 'shield', children }: AuthLayoutProps) {
  return (
    <SafeAreaView style={styles.safe} edges={['top', 'bottom']}>
      <KeyboardAvoidingView
        style={styles.flex}
        behavior={Platform.OS === 'ios' ? 'padding' : undefined}
      >
        <View style={styles.content}>
          <View style={styles.brand}>
            <View style={styles.logo}>
              <Feather name={icon} size={28} color={colors.primary} />
            </View>
            <Text variant="h1" style={styles.title}>
              {title}
            </Text>
            {subtitle ? (
              <Text variant="body" color="textSecondary" align="center" style={styles.subtitle}>
                {subtitle}
              </Text>
            ) : null}
          </View>

          <View style={styles.form}>{children}</View>
        </View>
      </KeyboardAvoidingView>
    </SafeAreaView>
  );
}

const styles = StyleSheet.create({
  safe: { flex: 1, backgroundColor: colors.background },
  flex: { flex: 1 },
  content: { flex: 1, justifyContent: 'center', padding: 24 },
  title: { marginTop: spacing.md },
  brand: { alignItems: 'center', marginBottom: 28 },
  logo: {
    width: 60,
    height: 60,
    borderRadius: 18,
    backgroundColor: colors.primarySoft,
    alignItems: 'center',
    justifyContent: 'center',
    marginBottom: 16,
  },
  subtitle: { marginTop: 8 },
  form: { width: '100%', maxWidth: 420, alignSelf: 'center' },
});