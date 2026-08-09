import React, { useCallback, useState } from 'react';
import { StyleSheet, FlatList, Alert, RefreshControl } from 'react-native';
import { useFocusEffect } from '@react-navigation/native';
import { Screen, Text, Badge, Row, EmptyState } from '@/components';
import { Spinner } from '@/components/Spinner';
import { useToast } from '@/components/Toast';
import { listActiveSessions, revokeSession } from '@/api/endpoints';
import { getErrorMessage } from '@/api/errors';
import { colors, spacing } from '@/theme';
import { formatDateTime } from '@/utils/format';
import type { SessionResponse } from '@/types/models';

export function SessionsScreen() {
  const { show } = useToast();
  const [sessions, setSessions] = useState<SessionResponse[]>([]);
  const [loading, setLoading] = useState(true);
  const [refreshing, setRefreshing] = useState(false);
  const [revoking, setRevoking] = useState<string | null>(null);

  const load = useCallback(async () => {
    try {
      const items = await listActiveSessions();
      setSessions(items);
    } catch (err) {
      show(getErrorMessage(err), 'error');
    } finally {
      setLoading(false);
      setRefreshing(false);
    }
  }, [show]);

  useFocusEffect(
    useCallback(() => {
      load();
    }, [load]),
  );

  async function handleRevoke(id: string) {
    setRevoking(id);
    try {
      await revokeSession(id);
      setSessions((prev) => prev.filter((s) => s.id !== id));
      show('Sesión revocada.', 'success');
    } catch (err) {
      show(getErrorMessage(err), 'error');
    } finally {
      setRevoking(null);
    }
  }

  function confirmRevoke(id: string) {
    Alert.alert('Revocar sesión', '¿Quieres cerrar esta sesión en el dispositivo?', [
      { text: 'Cancelar', style: 'cancel' },
      { text: 'Revocar', style: 'destructive', onPress: () => handleRevoke(id) },
    ]);
  }

  if (loading) {
    return (
      <Screen>
        <Spinner />
      </Screen>
    );
  }

  return (
    <Screen scroll={false}>
      <FlatList
        data={sessions}
        keyExtractor={(item) => item.id}
        style={styles.list}
        contentContainerStyle={styles.content}
        refreshControl={<RefreshControl refreshing={refreshing} onRefresh={load} tintColor={colors.primary} />}
        ListEmptyComponent={
          <EmptyState icon="smartphone" title="Sin sesiones activas" subtitle="Más dispositivos aparecerán aquí al iniciar sesión." />
        }
        renderItem={({ item }) => (
          <Row
            title={item.ip_address ?? 'Dispositivo desconocido'}
            subtitle={`Creada: ${formatDateTime(item.created_at)}`}
            description={`Expira: ${formatDateTime(item.expires_at)}`}
            left={<Badge tone="success" label="Activa" />}
            right={
              revoking === item.id ? (
                <Spinner size="small" />
              ) : (
                <Text variant="caption" color="danger" bold onPress={() => confirmRevoke(item.id)}>
                  Revocar
                </Text>
              )
            }
          />
        )}
      />
    </Screen>
  );
}

const styles = StyleSheet.create({
  list: { flex: 1 },
  content: { padding: spacing.lg, gap: spacing.sm },
});