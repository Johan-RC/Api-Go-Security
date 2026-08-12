import React, { useCallback, useEffect, useState } from 'react';
import { StyleSheet, View, TouchableOpacity } from 'react-native';
import type { NativeStackScreenProps } from '@react-navigation/native-stack';
import { LinearGradient } from 'expo-linear-gradient';
import { Feather } from '@expo/vector-icons';
import { Screen, Card, Text, Avatar, Badge, Reveal } from '@/components';
import { useAuth } from '@/contexts/AuthContext';
import { listUsers, listRoles, listFeatures, listActiveSessions } from '@/api/endpoints';
import { colors, spacing } from '@/theme';
import { pointer } from '@/utils/web';
import { initials } from '@/utils/format';
import type { HomeStackParamList } from '@/navigation/types';

type Props = NativeStackScreenProps<HomeStackParamList, 'HomeOverview'>;
type FeatherIcon = keyof typeof Feather.glyphMap;

interface Stats {
  users?: number;
  roles?: number;
  features?: number;
  sessions?: number;
}

function actorLabel(type?: string): string {
  switch (type) {
    case 'INSTRUCTOR': return 'Instructor';
    case 'LEARNER': return 'Aprendiz';
    default: return 'Usuario';
  }
}

function StatCard({ label, value, loading, color }: { label: string; value?: number; loading: boolean; color: string }) {
  return (
    <Card style={styles.statCard}>
      <View style={[styles.statDot, { backgroundColor: color }]} />
      <Text variant="h2" style={styles.statValue}>{loading ? '…' : value ?? '0'}</Text>
      <Text variant="caption" color="textSecondary">{label}</Text>
    </Card>
  );
}

export function HomeScreen({ navigation }: Props) {
  const { user, logout, can } = useAuth();
  const [stats, setStats] = useState<Stats>({});
  const [loading, setLoading] = useState(true);

  const canViewUsers = can('IDENTITY_USER_VIEW', 'READ');
  const canViewRoles = can('IDENTITY_ROLE_VIEW', 'READ');
  const canViewFeatures = can('IDENTITY_ROLE_MANAGE', 'WRITE');

  const loadStats = useCallback(async () => {
    setLoading(true);
    try {
      const [users, roles, features, sessions] = await Promise.all([
        canViewUsers ? listUsers(1, 1).then((r) => r.total) : Promise.resolve(undefined),
        canViewRoles ? listRoles(1, 1).then((r) => r.total) : Promise.resolve(undefined),
        canViewFeatures ? listFeatures(1, 1).then((r) => r.total) : Promise.resolve(undefined),
        listActiveSessions().then((r) => r.length),
      ]);
      setStats({ users, roles, features, sessions });
    } catch {
      // Si un recurso falla (usuario nuevo, sin permisos), no rompemos la vista.
      setStats({});
    } finally {
      setLoading(false);
    }
  }, [canViewUsers, canViewRoles, canViewFeatures]);

  useEffect(() => {
    loadStats();
  }, [loadStats]);

  const actions: Array<{ id: string; label: string; icon: FeatherIcon; tint: string; onPress: () => void }> = [];
  if (canViewUsers) {
    actions.push({ id: 'users', label: 'Usuarios', icon: 'users', tint: colors.primary, onPress: () => navigation.getParent()?.navigate('UsersTab') });
  }
  if (canViewRoles) {
    actions.push({ id: 'roles', label: 'Roles', icon: 'shield', tint: colors.info, onPress: () => navigation.getParent()?.navigate('RolesTab') });
  }
  if (canViewFeatures) {
    actions.push({ id: 'catalog', label: 'Módulos', icon: 'grid', tint: colors.success, onPress: () => navigation.getParent()?.navigate('CatalogTab') });
  }
  actions.push({ id: 'sessions', label: 'Sesiones', icon: 'smartphone', tint: colors.warning, onPress: () => navigation.navigate('Sessions') });
  if (can('AUDIT_LOG_VIEW', 'READ')) {
    actions.push({ id: 'audit', label: 'Auditoría', icon: 'file-text', tint: colors.danger, onPress: () => navigation.navigate('Audit') });
  }

  return (
    <Screen scroll>
      <Reveal>
        <LinearGradient
          colors={[colors.gradientStart, colors.gradientMid, colors.gradientEnd]}
          start={{ x: 0, y: 0 }}
          end={{ x: 1, y: 1 }}
          style={styles.banner}
        >
          <View style={styles.header}>
            <Avatar text={initials(user?.email ?? 'U')} color="rgba(255,255,255,0.25)" size={48} />
            <View style={styles.headerText}>
              <Text variant="label" style={{ color: 'rgba(255,255,255,0.85)' }}>Bienvenido</Text>
              <Text variant="h2" numberOfLines={1} style={{ color: '#FFFFFF' }}>{user?.email ?? 'Usuario'}</Text>
            </View>
            <Badge tone={user?.actorType === 'INSTRUCTOR' ? 'info' : user?.actorType === 'LEARNER' ? 'success' : 'primary'} label={actorLabel(user?.actorType)} />
          </View>
        </LinearGradient>
      </Reveal>

      {canViewUsers || canViewRoles || canViewFeatures ? (
        <Reveal delay={120}>
          <View style={styles.statsRow}>
            {canViewUsers ? <StatCard label="Usuarios" value={stats.users} loading={loading} color={colors.primary} /> : null}
            {canViewRoles ? <StatCard label="Roles" value={stats.roles} loading={loading} color={colors.info} /> : null}
            {canViewFeatures ? <StatCard label="Funcionalidades" value={stats.features} loading={loading} color={colors.success} /> : null}
            <StatCard label="Sesiones activas" value={stats.sessions} loading={loading} color={colors.warning} />
          </View>
        </Reveal>
      ) : null}

      <Reveal delay={200}>
        <Text variant="h3" style={styles.sectionTitle}>Accesos rápidos</Text>
        <View style={styles.actionsGrid}>
          {actions.map((action) => (
            <TouchableOpacity key={action.id} style={[styles.actionCard, pointer]} onPress={action.onPress}>
              <View style={[styles.actionIcon, { backgroundColor: action.tint }]}>
                <Feather name={action.icon} size={20} color="#FFFFFF" />
              </View>
              <Text variant="body" bold>{action.label}</Text>
            </TouchableOpacity>
          ))}
        </View>
      </Reveal>

      <Reveal delay={300}>
        <TouchableOpacity onPress={logout} style={[styles.logout, pointer]}>
          <Feather name="log-out" size={18} color={colors.danger} />
          <Text variant="body" color="danger" bold>Cerrar sesión</Text>
        </TouchableOpacity>
      </Reveal>
    </Screen>
  );
}

const styles = StyleSheet.create({
  banner: {
    borderRadius: 22,
    padding: spacing.lg,
    marginBottom: spacing.xl,
  },
  header: { flexDirection: 'row', alignItems: 'center', gap: spacing.md },
  headerText: { flex: 1 },
  statsRow: { flexWrap: 'wrap', flexDirection: 'row', gap: spacing.sm, marginBottom: spacing.xl },
  statCard: { flexGrow: 1, flexBasis: '45%', padding: spacing.md },
  statDot: { width: 8, height: 8, borderRadius: 4, marginBottom: spacing.sm },
  statValue: { marginBottom: 2 },
  sectionTitle: { marginBottom: spacing.md },
  actionsGrid: { flexDirection: 'row', flexWrap: 'wrap', gap: spacing.md, marginBottom: spacing.xl },
  actionCard: {
    flexBasis: '47%',
    flexGrow: 1,
    backgroundColor: colors.surface,
    borderRadius: 16,
    padding: spacing.lg,
    flexDirection: 'row',
    alignItems: 'center',
    gap: spacing.md,
    borderWidth: 1,
    borderColor: colors.border,
  },
  actionIcon: { width: 40, height: 40, borderRadius: 12, alignItems: 'center', justifyContent: 'center' },
  logout: { flexDirection: 'row', alignItems: 'center', justifyContent: 'center', gap: spacing.sm, paddingVertical: spacing.lg, marginTop: spacing.sm },
});