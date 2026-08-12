import React from 'react';
import { StyleSheet, View, KeyboardAvoidingView, Platform, ScrollView } from 'react-native';
import { SafeAreaView } from 'react-native-safe-area-context';
import { LinearGradient } from 'expo-linear-gradient';
import { StatusBar } from 'expo-status-bar';
import { colors, radius, shadows, spacing } from '@/theme';
import { Text } from '@/components/Text';
import { Reveal } from '@/components/Reveal';
import { ShieldAnimation } from '@/components/ShieldAnimation';

interface AuthLayoutProps {
  title: string;
  subtitle?: string;
  icon?: string;
  children: React.ReactNode;
}

export function AuthLayout({ title, subtitle, children }: AuthLayoutProps) {
  return (
    <LinearGradient colors={[colors.gradientStart, colors.gradientMid, colors.gradientEnd]} style={styles.safe}>
      <StatusBar style="light" />
      <SafeAreaView style={styles.safe} edges={['top', 'bottom']}>
        <KeyboardAvoidingView
          style={styles.flex}
          behavior={Platform.OS === 'ios' ? 'padding' : undefined}
        >
          <ScrollView
            contentContainerStyle={styles.scrollContent}
            keyboardShouldPersistTaps="handled"
            showsVerticalScrollIndicator={false}
          >
            <Reveal style={styles.brandWrap}>
              <View style={styles.brand}>
                <ShieldAnimation size={84} />
                <Text variant="h1" style={styles.title}>
                  {title}
                </Text>
                {subtitle ? (
                  <Text variant="body" color="textMuted" align="center" style={styles.subtitle}>
                    {subtitle}
                  </Text>
                ) : null}
              </View>
            </Reveal>

            <Reveal delay={150} style={styles.cardWrap}>
              <View style={styles.card}>
                <View style={styles.accentBar} />
                <View style={styles.glow} />
                {children}
              </View>
            </Reveal>

            <Reveal delay={260}>
              <Text variant="caption" color="textMuted" align="center" style={styles.footer}>
                JeuxAirel · Plataforma de identidad y accesos
              </Text>
            </Reveal>
          </ScrollView>
        </KeyboardAvoidingView>
      </SafeAreaView>
    </LinearGradient>
  );
}

const styles = StyleSheet.create({
  safe: { flex: 1 },
  flex: { flex: 1 },
  scrollContent: { flexGrow: 1, justifyContent: 'center', alignItems: 'center', padding: 24, paddingVertical: 40 },
  brandWrap: { width: '100%', maxWidth: 420 },
  cardWrap: { width: '100%', maxWidth: 420 },
  brand: { alignItems: 'center', marginBottom: 20 },
  title: { color: '#FFFFFF', textAlign: 'center' },
  subtitle: { marginTop: 8, color: 'rgba(255,255,255,0.85)', maxWidth: 340 },
  card: {
    width: '100%',
    maxWidth: 420,
    backgroundColor: colors.surface,
    borderRadius: radius.xl,
    padding: spacing.xl,
    ...shadows.md,
    overflow: 'hidden',
  },
  accentBar: {
    position: 'absolute',
    top: 0,
    left: 0,
    right: 0,
    height: 4,
    backgroundColor: colors.primary,
  },
  glow: {
    position: 'absolute',
    top: 32,
    right: -60,
    width: 160,
    height: 160,
    borderRadius: 80,
    backgroundColor: 'rgba(239,68,68,0.10)',
  },
  footer: { marginTop: 20, color: 'rgba(255,255,255,0.75)' },
});