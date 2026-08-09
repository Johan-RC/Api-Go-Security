import React, { useMemo, useState } from 'react';
import {
  StyleSheet,
  View,
  FlatList,
  RefreshControl,
  ActivityIndicator,
  Pressable,
} from 'react-native';
import type { NativeStackScreenProps } from '@react-navigation/native-stack';
import { Feather } from '@expo/vector-icons';
import { Screen, Text, Input, Badge, Avatar, EmptyState } from '@/components';
import { Spinner } from '@/components/Spinner';
import { useAuth } from '@/contexts/AuthContext';
import { usePaginatedList } from '@/hooks/usePaginatedList';
import { listUsers } from '@/api/endpoints';
import { getErrorMessage } from '@/api/errors';
import { colors, spacing, radius, shadows } from '@/theme';
import { initials } from '@/utils/format';
import type { UserResponse } from '@/types/models';
import type { UsersStackParamList } from '@/navigation/types';

type Props = NativeStackScreenProps<UsersStackParamList, 'UserList'>;

export function UsersScreen({ navigation }: Props) {
  const { can } = useAuth();
  const canCreateUser = can('IDENTITY_USER_MANAGE', 'WRITE');
  const [query, setQuery] = useState('');

  const { items, loading, refreshing, error, refresh, loadMore, hasMore } = usePaginatedList<UserResponse>(
    (page, size) => listUsers(page, size),
    { pageSize: 50 },
  );

  const filtered = useMemo(() => {
    const q = query.trim().toLowerCase();
    if (!q) return items;
    return items.filter(
      (u) => u.email.toLowerCase().includes(q) || `${u.first_name} ${u.last_name}`.toLowerCase().includes(q),
    );
  }, [items, query]);

  if (loading && items.length === 0) {
    return (
      <Screen>
        <Spinner />
      </Screen>
    );
  }

  return (
    <Screen scroll={false}>
      <Input
        placeholder="Buscar por correo o nombre…"
        value={query}
        onChangeText={setQuery}
        leftIcon={<Feather name="search" size={18} color={colors.textMuted} />}
        style={{ marginBottom: spacing.sm }}
      />

      <FlatList
        data={filtered}
        keyExtractor={(item) => item.id}
        style={styles.list}
        contentContainerStyle={styles.content}
        refreshControl={<RefreshControl refreshing={refreshing} onRefresh={refresh} tintColor={colors.primary} />}
        onEndReached={loadMore}
        onEndReachedThreshold={0.4}
        ListEmptyComponent={
          <EmptyState icon="users" title="Sin usuarios" subtitle="Crea tu primer usuario con el botón de abajo." />
        }
        renderItem={({ item }) => (
          <Pressable style={styles.userRow} onPress={() => navigation.navigate('UserDetail', { userId: item.id })}>
            <Avatar text={initials(item.first_name, item.last_name)} />
            <View style={styles.userInfo}>
              <Text variant="body" bold numberOfLines={1}>
                {`${item.first_name} ${item.last_name}`}
              </Text>
              <Text variant="caption" color="textSecondary" numberOfLines={1}>
                {item.email}
              </Text>
            </View>
            <Badge tone={item.is_active ? 'success' : 'neutral'} label={item.is_active ? 'Activo' : 'Inactivo'} />
          </Pressable>
        )}
        ListFooterComponent={
          error ? (
            <Text variant="caption" color="danger" align="center" style={styles.footer}>
              {getErrorMessage(error)}
            </Text>
          ) : hasMore ? (
            <ActivityIndicator color={colors.primary} style={styles.footer} />
          ) : null
        }
      />

      {canCreateUser ? (
        <Pressable style={styles.fab} onPress={() => navigation.navigate('UserDetail', { userId: 'new' })}>
          <Feather name="user-plus" size={22} color="#FFFFFF" />
        </Pressable>
      ) : null}
    </Screen>
  );
}

const styles = StyleSheet.create({
  list: { flex: 1 },
  content: { padding: spacing.lg, gap: spacing.sm, paddingBottom: spacing['2xl'] },
  userRow: {
    flexDirection: 'row',
    alignItems: 'center',
    gap: spacing.md,
    backgroundColor: colors.surface,
    borderRadius: radius.md,
    borderWidth: 1,
    borderColor: colors.border,
    padding: spacing.md,
  },
  userInfo: { flex: 1, gap: 2 },
  fab: {
    position: 'absolute',
    right: spacing.xl,
    bottom: spacing.xl,
    width: 56,
    height: 56,
    borderRadius: 28,
    backgroundColor: colors.primary,
    alignItems: 'center',
    justifyContent: 'center',
    ...shadows.md,
  },
  footer: { paddingVertical: spacing.lg },
});