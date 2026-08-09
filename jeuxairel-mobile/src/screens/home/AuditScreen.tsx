import React from 'react';
import { StyleSheet, FlatList, RefreshControl, ActivityIndicator } from 'react-native';
import { Screen, Text, Badge, Row, EmptyState } from '@/components';
import { usePaginatedList } from '@/hooks/usePaginatedList';
import { listAuditLogins } from '@/api/endpoints';
import { getErrorMessage } from '@/api/errors';
import { colors, spacing } from '@/theme';
import { formatDateTime } from '@/utils/format';
import type { AuditLoginResponse, AuditOutcome } from '@/types/models';

const outcomeTone: Record<AuditOutcome, 'success' | 'danger' | 'warning'> = {
  SUCCESS: 'success',
  INVALID_PASSWORD: 'danger',
  USER_NOT_FOUND: 'warning',
  ACCOUNT_LOCKED: 'danger',
  TOKEN_EXPIRED: 'warning',
};

const outcomeLabel: Record<AuditOutcome, string> = {
  SUCCESS: 'Acceso exitoso',
  INVALID_PASSWORD: 'Contraseña incorrecta',
  USER_NOT_FOUND: 'Usuario no encontrado',
  ACCOUNT_LOCKED: 'Cuenta bloqueada',
  TOKEN_EXPIRED: 'Token expirado',
};

export function AuditScreen() {
  const { items, loading, refreshing, error, refresh, loadMore, hasMore } = usePaginatedList<AuditLoginResponse>(listAuditLogins, { pageSize: 20 });

  if (loading && items.length === 0) {
    return (
      <Screen>
        <ActivityIndicator color={colors.primary} size="large" style={styles.spinner} />
      </Screen>
    );
  }

  return (
    <Screen scroll={false}>
      <FlatList
        data={items}
        keyExtractor={(item) => item.id}
        style={styles.list}
        contentContainerStyle={styles.content}
        refreshControl={<RefreshControl refreshing={refreshing} onRefresh={refresh} tintColor={colors.primary} />}
        ListEmptyComponent={
          <EmptyState icon="file-text" title="Sin registros" subtitle="La auditoría de accesos se mostrará aquí." />
        }
        onEndReached={loadMore}
        onEndReachedThreshold={0.4}
        ListFooterComponent={
          error ? (
            <Text variant="caption" color="danger" align="center" style={styles.footer}>
              {getErrorMessage(error)}
            </Text>
          ) : hasMore ? (
            <ActivityIndicator color={colors.primary} style={styles.footer} />
          ) : null
        }
        renderItem={({ item }) => (
          <Row
            title={item.email_attempted}
            subtitle={`${formatDateTime(item.attempted_at)} · ${item.ip_address ?? 'IP desconocida'}`}
            left={<Badge tone={outcomeTone[item.outcome] ?? 'warning'} label={outcomeLabel[item.outcome] ?? item.outcome} />}
          />
        )}
      />
    </Screen>
  );
}

const styles = StyleSheet.create({
  list: { flex: 1 },
  content: { padding: spacing.lg, gap: spacing.sm },
  spinner: { marginTop: spacing['2xl'] },
  footer: { paddingVertical: spacing.lg },
});