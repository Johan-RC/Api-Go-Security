import React from 'react';
import { StyleSheet, View, Pressable, Alert } from 'react-native';
import type { NativeStackScreenProps } from '@react-navigation/native-stack';
import { Feather } from '@expo/vector-icons';
import { Screen, Text, Card, Avatar, Badge } from '@/components';
import { useAuth } from '@/contexts/AuthContext';
import { colors, spacing } from '@/theme';
import { initials } from '@/utils/format';
import type { MoreStackParamList, MainTabParamList } from '@/navigation/types';

type Props = NativeStackScreenProps<MoreStackParamList, 'Profile'>;

export function ProfileScreen({ navigation }: Props) {
  const { user, logout, can } = useAuth();

  const goToTab = (tab: keyof MainTabParamList, screen?: string) => {
    const parent = navigation.getParent() as { navigate: (t: string, p?: unknown) => void } | undefined;
    parent?.navigate(tab, screen ? { screen } : undefined);
  };

  function confirmLogout() {
    Alert.alert('Cerrar sesión', '¿Deseas cerrar tu sesión?', [
      { text: 'Cancelar', style: 'cancel' },
      { text: 'Cerrar sesión', style: 'destructive', onPress: () => logout() },
    ]);
  }

  return (
    <Screen scroll>
      <Card>
        <View style={styles.profileHeader}>
          <Avatar text={initials(user?.email ?? 'U')} size={64} />
          <View style={styles.flex1}>
            <Text variant="h3" numberOfLines={1}>{user?.email ?? 'Usuario'}</Text>
            <Badge tone={user?.actorType === 'INSTRUCTOR' ? 'info' : user?.actorType === 'LEARNER' ? 'success' : 'primary'} label={actorLabel(user?.actorType)} style={{ alignSelf: 'flex-start', marginTop: spacing.xs }} />
          </View>
        </View>
      </Card>

      <Text variant="h3" style={styles.sectionTitle}>Seguridad</Text>
      <Card style={styles.menu}>
        <MenuRow icon="smartphone" label="Sesiones activas" hint="Dispositivos conectados" onPress={() => goToTab('HomeTab', 'Sessions')} />
        {can('AUDIT_LOG_VIEW', 'READ') ? (
          <MenuRow icon="file-text" label="Auditoría de accesos" hint="Registro de inicios de sesión" onPress={() => goToTab('HomeTab', 'Audit')} last />
        ) : null}
      </Card>

      <ButtonLogout onPress={confirmLogout} />
    </Screen>
  );
}

function actorLabel(type?: string): string {
  switch (type) {
    case 'INSTRUCTOR': return 'Instructor';
    case 'LEARNER': return 'Aprendiz';
    default: return 'Usuario';
  }
}

function MenuRow({
  icon,
  label,
  hint,
  onPress,
  last = false,
}: {
  icon: keyof typeof Feather.glyphMap;
  label: string;
  hint?: string;
  onPress: () => void;
  last?: boolean;
}) {
  return (
    <Pressable
      style={[styles.menuRow, last ? styles.menuRowLast : null]}
      onPress={onPress}
      android_ripple={{ color: colors.border }}
    >
      <View style={styles.menuIcon}>
        <Feather name={icon} size={18} color={colors.primary} />
      </View>
      <View style={styles.flex1}>
        <Text variant="body">{label}</Text>
        {hint ? <Text variant="caption" color="textSecondary">{hint}</Text> : null}
      </View>
      <Feather name="chevron-right" size={18} color={colors.textMuted} />
    </Pressable>
  );
}

function ButtonLogout({ onPress }: { onPress: () => void }) {
  return (
    <Pressable style={styles.logout} onPress={onPress}>
      <Feather name="log-out" size={18} color={colors.danger} />
      <Text variant="body" color="danger" bold>Cerrar sesión</Text>
    </Pressable>
  );
}

const styles = StyleSheet.create({
  profileHeader: { flexDirection: 'row', alignItems: 'center', gap: spacing.lg },
  flex1: { flex: 1 },
  sectionTitle: { marginTop: spacing.xl, marginBottom: spacing.md },
  menu: { padding: 0, overflow: 'hidden' },
  menuRow: {
    flexDirection: 'row',
    alignItems: 'center',
    gap: spacing.md,
    padding: spacing.lg,
    borderBottomWidth: 1,
    borderBottomColor: colors.border,
  },
  menuRowLast: { borderBottomWidth: 0 },
  menuIcon: {
    width: 36,
    height: 36,
    borderRadius: 10,
    backgroundColor: colors.primarySoft,
    alignItems: 'center',
    justifyContent: 'center',
  },
  logout: {
    flexDirection: 'row',
    alignItems: 'center',
    justifyContent: 'center',
    gap: spacing.sm,
    paddingVertical: spacing.lg,
    marginTop: spacing.xl,
  },
});