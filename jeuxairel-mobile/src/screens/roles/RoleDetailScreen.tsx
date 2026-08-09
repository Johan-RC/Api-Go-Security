import React, { useCallback, useState } from 'react';
import { StyleSheet, View, Alert, Pressable, ActivityIndicator } from 'react-native';
import type { NativeStackScreenProps } from '@react-navigation/native-stack';
import { Feather } from '@expo/vector-icons';
import { useFocusEffect } from '@react-navigation/native';
import { Screen, Text, Card, Input, Button, Badge, Select } from '@/components';
import { useToast } from '@/components/Toast';
import { useAuth } from '@/contexts/AuthContext';
import {
  getRole,
  updateRole,
  listRoleFeatures,
  assignRoleFeatures,
  listFeatures,
  deleteRole,
} from '@/api/endpoints';
import { getErrorMessage } from '@/api/errors';
import { isRequired } from '@/validation';
import { colors, spacing } from '@/theme';
import type { RoleResponse, FeatureResponse, FeatureAssignment, ScopeType } from '@/types/models';
import type { RolesStackParamList } from '@/navigation/types';

type Props = NativeStackScreenProps<RolesStackParamList, 'RoleDetail'>;

const scopeOptions: Array<{ label: string; value: ScopeType }> = [
  { label: 'Global', value: 'GLOBAL' },
  { label: 'Centro de formación', value: 'TRAINING_CENTER' },
  { label: 'Área', value: 'AREA' },
  { label: 'Fichas propias', value: 'OWN_FICHAS' },
  { label: 'Horario propio', value: 'OWN_SCHEDULE' },
  { label: 'Perfil propio', value: 'OWN_PROFILE' },
  { label: 'Ficha propia (aprendiz)', value: 'OWN_FICHA_AS_LEARNER' },
];

export function RoleDetailScreen({ route, navigation }: Props) {
  const roleId = route.params.roleId;

  const { show } = useToast();
  const { can } = useAuth();
  const canManage = can('IDENTITY_ROLE_MANAGE', 'WRITE');
  const canManageFeatures = can('IDENTITY_ROLE_ASSIGN', 'WRITE');
  const [role, setRole] = useState<RoleResponse | null>(null);
  const [featuresCatalog, setFeaturesCatalog] = useState<FeatureResponse[]>([]);
  const [assignments, setAssignments] = useState<FeatureAssignment[]>([]);
  const [loading, setLoading] = useState(true);
  const [saving, setSaving] = useState(false);

  const [editing, setEditing] = useState(false);
  const [displayName, setDisplayName] = useState('');
  const [description, setDescription] = useState('');

  // selector de funcionalidad nueva
  const [pendingFeature, setPendingFeature] = useState<string | null>(null);
  const [pendingScope, setPendingScope] = useState<ScopeType | null>(null);

  const refresh = useCallback(async () => {
    const [rl, ft, feats] = await Promise.allSettled([
      getRole(roleId),
      listRoleFeatures(roleId),
      listFeatures(1, 100),
    ]);
    if (rl.status === 'fulfilled') {
      setRole(rl.value);
      setDisplayName(rl.value.display_name);
      setDescription(rl.value.description ?? '');
    }
    if (ft.status === 'fulfilled') setAssignments(ft.value);
    if (feats.status === 'fulfilled') setFeaturesCatalog(feats.value.items);
  }, [roleId]);

  useFocusEffect(
    useCallback(() => {
      refresh()
        .catch((err) => show(getErrorMessage(err), 'error'))
        .finally(() => setLoading(false));
    }, [refresh, show]),
  );

  if (loading) return <Screen><ActivityIndicator color={colors.primary} style={styles.loader} /></Screen>;
  if (!role) return <Screen><Text color="textSecondary">No se encontró el rol.</Text></Screen>;

  const featureName = (featureId: string) =>
    featuresCatalog.find((f) => f.id === featureId)?.name ?? featureId;

  const availableFeatures = featuresCatalog.filter(
    (f) => !assignments.some((a) => a.feature_id === f.id),
  );

  async function persist(newAssignments: FeatureAssignment[]) {
    setSaving(true);
    try {
      await assignRoleFeatures(roleId, newAssignments);
      setAssignments(newAssignments);
      show('Funcionalidades actualizadas.', 'success');
    } catch (err) {
      show(getErrorMessage(err), 'error');
    } finally {
      setSaving(false);
    }
  }

  function addPending() {
    if (!pendingFeature || !pendingScope) return;
    const next = [...assignments, { feature_id: pendingFeature, scope_type: pendingScope }];
    persist(next);
    setPendingFeature(null);
    setPendingScope(null);
  }

  function removeAssignment(featureId: string) {
    const next = assignments.filter((a) => a.feature_id !== featureId);
    persist(next);
  }

  function handleSave() {
    if (!isRequired(displayName)) {
      show('El nombre para mostrar es obligatorio.', 'error');
      return;
    }
    updateRole(roleId, { display_name: displayName, description: description || undefined })
      .then((updated) => {
        setRole(updated);
        setEditing(false);
        show('Rol actualizado.', 'success');
      })
      .catch((err) => show(getErrorMessage(err), 'error'));
  }

  function confirmDelete() {
    Alert.alert('Eliminar rol', '¿Estás seguro de eliminar este rol?', [
      { text: 'Cancelar', style: 'cancel' },
      {
        text: 'Eliminar',
        style: 'destructive',
        onPress: () =>
          deleteRole(roleId)
            .then(() => {
              show('Rol eliminado.', 'success');
              navigation.goBack();
            })
            .catch((err) => show(getErrorMessage(err), 'error')),
      },
    ]);
  }

  return (
    <Screen scroll>
      {/* Info del rol */}
      <Card>
        <View style={styles.headerRow}>
          <View style={styles.flex1}>
            <Text variant="h3">{role.display_name}</Text>
            <Text variant="caption" color="textSecondary">{role.name}</Text>
          </View>
          <Badge tone={role.is_system_role ? 'info' : 'neutral'} label={role.is_system_role ? 'Sistema' : 'Personal'} />
        </View>

        {editing ? (
          <>
            <Input label="Nombre para mostrar" value={displayName} onChangeText={setDisplayName} style={{ marginTop: spacing.md }} />
            <Input label="Descripción" value={description} onChangeText={setDescription} multiline numberOfLines={2} />
            <View style={styles.btnRow}>
              <Button variant="outline" label="Cancelar" onPress={() => setEditing(false)} style={styles.flex1} />
              <Button label="Guardar" onPress={handleSave} style={styles.flex1} />
            </View>
          </>
        ) : (
          <>
            <Text variant="body" style={{ marginTop: spacing.md }}>
              {role.description || 'Sin descripción.'}
            </Text>
            <Pressable style={styles.editRow} disabled={!canManage} onPress={() => setEditing(true)}>
              <Feather name="edit-2" size={14} color={canManage ? colors.primary : colors.textMuted} />
              <Text variant="caption" color={canManage ? 'primary' : 'textSecondary'} bold>Editar</Text>
            </Pressable>
          </>
        )}
      </Card>

      {/* Funcionalidades */}
      <Text variant="h3" style={styles.sectionTitle}>Funcionalidades ({assignments.length})</Text>
      <Card>
        {assignments.length === 0 ? (
          <Text variant="caption" color="textSecondary" style={{ marginBottom: spacing.md }}>
            Este rol no tiene funcionalidades asignadas.
          </Text>
        ) : (
          assignments.map((a) => (
            <Pressable key={a.feature_id} style={styles.featRow} disabled={!canManageFeatures} onPress={() => canManageFeatures && removeAssignment(a.feature_id)}>
              <View style={styles.flex1}>
                <Text variant="body" style={styles.flex1}>{featureName(a.feature_id)}</Text>
                <Badge tone="primary" label={a.scope_type} style={{ alignSelf: 'flex-start', marginTop: spacing.xs }} />
              </View>
              <Feather name="trash-2" size={15} color={canManageFeatures ? colors.textMuted : colors.border} />
            </Pressable>
          ))
        )}

        {canManageFeatures && availableFeatures.length > 0 ? (
          <>
            <Select
              label="Funcionalidad"
              placeholder="Elige una funcionalidad…"
              value={pendingFeature}
              options={availableFeatures.map((f) => ({ label: `${f.name} (${f.code})`, value: f.id }))}
              onChange={setPendingFeature}
              style={{ marginTop: spacing.md }}
            />
            <Select
              label="Ámbito (scope)"
              placeholder="Elige el ámbito…"
              value={pendingScope}
              options={scopeOptions}
              onChange={setPendingScope}
            />
            <Button label="Agregar funcionalidad" loading={saving} disabled={!pendingFeature || !pendingScope} onPress={addPending} />
          </>
        ) : null}
      </Card>

      {canManage ? (
      <Button variant="danger" label="Eliminar rol" onPress={confirmDelete} style={{ marginTop: spacing.xl }} />
    ) : null}
    </Screen>
  );
}

const styles = StyleSheet.create({
  loader: { marginTop: spacing['2xl'] },
  flex1: { flex: 1 },
  headerRow: { flexDirection: 'row', alignItems: 'center', gap: spacing.md },
  editRow: { flexDirection: 'row', alignItems: 'center', gap: spacing.xs, marginTop: spacing.md },
  sectionTitle: { marginTop: spacing.xl, marginBottom: spacing.md },
  featRow: { flexDirection: 'row', alignItems: 'center', gap: spacing.md, paddingVertical: spacing.sm },
  btnRow: { flexDirection: 'row', gap: spacing.sm, marginTop: spacing.sm },
});