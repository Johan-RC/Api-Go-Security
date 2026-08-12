import React, { useCallback, useMemo } from 'react';
import { StyleSheet, FlatList, RefreshControl, ActivityIndicator, Pressable } from 'react-native';
import type { NativeStackScreenProps } from '@react-navigation/native-stack';
import { Feather } from '@expo/vector-icons';
import { useFocusEffect } from '@react-navigation/native';
import { Screen, Text, Badge, Row, EmptyState } from '@/components';
import { Spinner } from '@/components/Spinner';
import { useAuth } from '@/contexts/AuthContext';
import { usePaginatedList } from '@/hooks/usePaginatedList';
import { listRoles } from '@/api/endpoints';
import { getErrorMessage } from '@/api/errors';
import { colors, spacing } from '@/theme';
import { pointer } from '@/utils/web';
import type { RoleResponse } from '@/types/models';
import type { RolesStackParamList } from '@/navigation/types';

type Props = NativeStackScreenProps<RolesStackParamList, 'RoleList'>;

export function RolesScreen({ navigation }: Props) {
  const { can } = useAuth();
  const canCreate = can('IDENTITY_ROLE_MANAGE', 'WRITE');
  const { items, loading, refreshing, error, refresh, loadMore, hasMore } = usePaginatedList<RoleResponse>(
    (page, size) => listRoles(page, size),
    { pageSize: 50 },
  );

  useFocusEffect(
    useCallback(() => {
      refresh();
    }, [refresh]),
  );

  const data = useMemo(() => items, [items]);

  if (loading && items.length === 0) {
    return (
      <Screen>
        <Spinner />
      </Screen>
    );
  }

  return (
    <Screen scroll={false}>
      <FlatList
        data={data}
        keyExtractor={(item) => item.id}
        style={styles.list}
        contentContainerStyle={styles.content}
        refreshControl={<RefreshControl refreshing={refreshing} onRefresh={refresh} tintColor={colors.primary} />}
        onEndReached={loadMore}
        onEndReachedThreshold={0.4}
        ListEmptyComponent={
          <EmptyState icon="shield" title="Sin roles" subtitle="Crea tu primer rol para empezar." />
        }
        renderItem={({ item }) => (
          <Row
            title={item.display_name}
            subtitle={item.name}
            description={item.description ?? undefined}
            left={<Feather name="shield" size={18} color={colors.primary} />}
            right={item.is_system_role ? <Badge tone="info" label="Sistema" /> : <Badge tone="neutral" label="Personalizado" />}
            onPress={() => navigation.navigate('RoleDetail', { roleId: item.id })}
          />
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
    {canCreate ? (
        <Pressable style={[styles.fab, pointer]} onPress={() => navigation.navigate('RoleCreate')}>
          <Feather name="plus" size={22} color="#FFFFFF" />
        </Pressable>
      ) : null}
    </Screen>
  );
}

const styles = StyleSheet.create({
  list: { flex: 1 },
  content: { padding: spacing.lg, gap: spacing.sm, paddingBottom: spacing['2xl'] },
  footer: { paddingVertical: spacing.lg },
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
  },
});