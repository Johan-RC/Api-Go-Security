import React, { useCallback, useEffect, useMemo, useState } from 'react';
import { StyleSheet, View, Modal, Pressable, ActivityIndicator, ScrollView, Alert } from 'react-native';
import { Feather } from '@expo/vector-icons';
import { Screen, Text, Card, Input, Button, Badge, Select } from '@/components';
import { useToast } from '@/components/Toast';
import { useAuth } from '@/contexts/AuthContext';
import {
  listModules,
  listFeatures,
  createModule,
  createFeature,
  deleteModule,
  deleteFeature,
} from '@/api/endpoints';
import { getErrorMessage } from '@/api/errors';
import { isRequired, maxLength } from '@/validation';
import { colors, spacing, radius, shadows } from '@/theme';
import type { ModuleResponse, FeatureResponse, ActionLevel } from '@/types/models';

const actionTones: Record<ActionLevel, 'success' | 'info' | 'warning' | 'danger' | 'primary'> = {
  READ: 'info',
  WRITE: 'primary',
  DELETE: 'danger',
  PUBLISH: 'warning',
  APPROVE: 'success',
};

export function CatalogScreen() {
  const { show } = useToast();
  const { can } = useAuth();
  const canManage = can('IDENTITY_ROLE_MANAGE', 'WRITE');

  const [modules, setModules] = useState<ModuleResponse[]>([]);
  const [features, setFeatures] = useState<FeatureResponse[]>([]);
  const [loading, setLoading] = useState(true);

  const [moduleModal, setModuleModal] = useState(false);
  const [featureModal, setFeatureModal] = useState(false);

  const load = useCallback(async () => {
    try {
      const [m, f] = await Promise.allSettled([
        listModules(1, 100),
        listFeatures(1, 100),
      ]);
      if (m.status === 'fulfilled') setModules(m.value.items);
      if (f.status === 'fulfilled') setFeatures(f.value.items);
    } catch (err) {
      show(getErrorMessage(err), 'error');
    } finally {
      setLoading(false);
    }
  }, [show]);

  useEffect(() => {
    load();
  }, [load]);

  const featuresByModule = useMemo(() => {
    const map: Record<string, FeatureResponse[]> = {};
    for (const f of features) {
      (map[f.module_id] ??= []).push(f);
    }
    return map;
  }, [features]);

  if (loading) {
    return (
      <Screen>
        <ActivityIndicator color={colors.primary} style={styles.loader} />
      </Screen>
    );
  }

  return (
    <Screen scroll={false}>
      <ScrollView style={styles.scroll} contentContainerStyle={styles.content} showsVerticalScrollIndicator={false}>
        <View style={styles.actions}>
          {canManage ? (
            <>
              <Button variant="secondary" label="Nuevo módulo" icon={<Feather name="plus" size={16} color={colors.primary} />} onPress={() => setModuleModal(true)} style={styles.actionBtn} />
              <Button variant="secondary" label="Nueva funcionalidad" icon={<Feather name="plus" size={16} color={colors.primary} />} onPress={() => setFeatureModal(true)} style={styles.actionBtn} />
            </>
          ) : null}
        </View>

        {modules.length === 0 ? (
          <Text variant="body" color="textSecondary" align="center" style={{ marginTop: spacing['2xl'] }}>
            No hay módulos todavía.
          </Text>
        ) : (
          modules.map((m) => (
            <Card key={m.id} style={styles.moduleCard}>
              <View style={styles.moduleHeader}>
                <View style={styles.flex1}>
                  <Text variant="h3">{m.name}</Text>
                  <Text variant="caption" color="textSecondary">{m.code} · Orden {m.display_order}</Text>
                </View>
                <Badge tone={m.is_active ? 'success' : 'neutral'} label={m.is_active ? 'Activo' : 'Inactivo'} />
                {canManage ? (
                  <Pressable onPress={() => confirmDeleteModule(m.id, m.name)} hitSlop={8}>
                    <Feather name="trash-2" size={16} color={colors.danger} />
                  </Pressable>
                ) : null}
              </View>

              {(featuresByModule[m.id] ?? []).map((f) => (
                <View key={f.id} style={styles.featureRow}>
                  <View style={styles.flex1}>
                    <Text variant="body">{f.name}</Text>
                    <Text variant="caption" color="textSecondary">{f.code}</Text>
                  </View>
                  <Badge tone={actionTones[f.action_level]} label={f.action_level} />
                  {canManage ? (
                    <Pressable onPress={() => confirmDeleteFeature(f.id)} hitSlop={8}>
                      <Feather name="trash-2" size={14} color={colors.textMuted} />
                    </Pressable>
                  ) : null}
                </View>
              ))}
              {(featuresByModule[m.id] ?? []).length === 0 ? (
                <Text variant="caption" color="textMuted">Sin funcionalidades.</Text>
              ) : null}
            </Card>
          ))
        )}
      </ScrollView>

      <ModuleForm
        visible={moduleModal}
        onClose={() => setModuleModal(false)}
        onCreate={async (data) => {
          await createModule(data);
          show('Módulo creado.', 'success');
          setModuleModal(false);
          load();
        }}
      />
      <FeatureForm
        visible={featureModal}
        modules={modules}
        onClose={() => setFeatureModal(false)}
        onCreate={async (data) => {
          await createFeature(data);
          show('Funcionalidad creada.', 'success');
          setFeatureModal(false);
          load();
        }}
      />
    </Screen>
  );

  function confirmDeleteModule(id: string, name: string) {
    Alert.alert('Eliminar módulo', `¿Eliminar el módulo "${name}"?`, [
      { text: 'Cancelar', style: 'cancel' },
      {
        text: 'Eliminar',
        style: 'destructive',
        onPress: () =>
          deleteModule(id)
            .then(() => { show('Módulo eliminado.', 'success'); load(); })
            .catch((err) => show(getErrorMessage(err), 'error')),
      },
    ]);
  }

  function confirmDeleteFeature(id: string) {
    Alert.alert('Eliminar funcionalidad', '¿Eliminar esta funcionalidad?', [
      { text: 'Cancelar', style: 'cancel' },
      {
        text: 'Eliminar',
        style: 'destructive',
        onPress: () =>
          deleteFeature(id)
            .then(() => { show('Funcionalidad eliminada.', 'success'); load(); })
            .catch((err) => show(getErrorMessage(err), 'error')),
      },
    ]);
  }
}

// ---------------------------------------------------------------------------
// Formularios (modales)
// ---------------------------------------------------------------------------

function ModuleForm({
  visible,
  onClose,
  onCreate,
}: {
  visible: boolean;
  onClose: () => void;
  onCreate: (data: { code: string; name: string; description?: string; display_order: number }) => Promise<void>;
}) {
  const { show } = useToast();
  const [code, setCode] = useState('');
  const [name, setName] = useState('');
  const [order, setOrder] = useState('1');
  const [errors, setErrors] = useState<Record<string, string>>({});
  const [loading, setLoading] = useState(false);

  function close() {
    setCode('');
    setName('');
    setOrder('1');
    setErrors({});
    onClose();
  }

  async function submit() {
    const next: Record<string, string> = {};
    if (!isRequired(code)) next.code = 'El código es obligatorio.';
    else if (!maxLength(code, 30)) next.code = 'Máximo 30 caracteres.';
    if (!isRequired(name)) next.name = 'El nombre es obligatorio.';
    else if (!maxLength(name, 100)) next.name = 'Máximo 100 caracteres.';
    setErrors(next);
    if (Object.keys(next).length > 0) return;

    setLoading(true);
    try {
      await onCreate({ code, name, display_order: parseInt(order, 10) || 1 });
      close();
    } catch (err) {
      show(getErrorMessage(err), 'error');
    } finally {
      setLoading(false);
    }
  }

  return (
    <Modal visible={visible} transparent animationType="slide" onRequestClose={onClose}>
      <View style={styles.backdrop}>
        <View style={styles.sheet}>
          <Text variant="h3">Nuevo módulo</Text>
          <Input label="Código *" placeholder="IDENTITY" autoCapitalize="characters" value={code} onChangeText={setCode} error={errors.code} style={{ marginTop: spacing.md }} />
          <Input label="Nombre *" placeholder="Gestión de identidades" value={name} onChangeText={setName} error={errors.name} />
          <Input label="Orden de visualización" placeholder="1" keyboardType="number-pad" value={order} onChangeText={setOrder} />
          <Button label="Crear módulo" loading={loading} onPress={submit} style={{ marginTop: spacing.sm }} />
          <Button variant="ghost" label="Cancelar" onPress={onClose} />
        </View>
      </View>
    </Modal>
  );
}

function FeatureForm({
  visible,
  modules,
  onClose,
  onCreate,
}: {
  visible: boolean;
  modules: ModuleResponse[];
  onClose: () => void;
  onCreate: (data: { module_id: string; code: string; name: string; action_level: ActionLevel }) => Promise<void>;
}) {
  const { show } = useToast();
  const [moduleId, setModuleId] = useState<string | null>(null);
  const [code, setCode] = useState('');
  const [name, setName] = useState('');
  const [actionLevel, setActionLevel] = useState<ActionLevel>('READ');
  const [errors, setErrors] = useState<Record<string, string>>({});
  const [loading, setLoading] = useState(false);

  function close() {
    setModuleId(null);
    setCode('');
    setName('');
    setActionLevel('READ');
    setErrors({});
    onClose();
  }

  async function submit() {
    const next: Record<string, string> = {};
    if (!moduleId) next.moduleId = 'Selecciona un módulo.';
    if (!isRequired(code)) next.code = 'El código es obligatorio.';
    else if (!maxLength(code, 60)) next.code = 'Máximo 60 caracteres.';
    if (!isRequired(name)) next.name = 'El nombre es obligatorio.';
    else if (!maxLength(name, 120)) next.name = 'Máximo 120 caracteres.';
    setErrors(next);
    if (Object.keys(next).length > 0) return;

    setLoading(true);
    try {
      await onCreate({ module_id: moduleId!, code, name, action_level: actionLevel });
      close();
    } catch (err) {
      show(getErrorMessage(err), 'error');
    } finally {
      setLoading(false);
    }
  }

  return (
    <Modal visible={visible} transparent animationType="slide" onRequestClose={onClose}>
      <View style={styles.backdrop}>
        <View style={styles.sheet}>
          <Text variant="h3">Nueva funcionalidad</Text>
          <Select
            label="Módulo *"
            placeholder="Elige un módulo…"
            value={moduleId}
            options={modules.map((m) => ({ label: m.name, value: m.id }))}
            onChange={setModuleId}
            style={{ marginTop: spacing.md }}
          />
          {errors.moduleId ? <Text variant="caption" color="danger">{errors.moduleId}</Text> : null}
          <Input label="Código *" placeholder="USER_VIEW" autoCapitalize="characters" value={code} onChangeText={setCode} error={errors.code} />
          <Input label="Nombre *" placeholder="Ver usuarios" value={name} onChangeText={setName} error={errors.name} />
          <Select
            label="Nivel de acción"
            value={actionLevel}
            options={[
              { label: 'Lectura (READ)', value: 'READ' },
              { label: 'Escritura (WRITE)', value: 'WRITE' },
              { label: 'Eliminación (DELETE)', value: 'DELETE' },
              { label: 'Publicación (PUBLISH)', value: 'PUBLISH' },
              { label: 'Aprobación (APPROVE)', value: 'APPROVE' },
            ]}
            onChange={setActionLevel}
          />
          <Button label="Crear funcionalidad" loading={loading} onPress={submit} style={{ marginTop: spacing.sm }} />
          <Button variant="ghost" label="Cancelar" onPress={onClose} />
        </View>
      </View>
    </Modal>
  );
}

const styles = StyleSheet.create({
  loader: { marginTop: spacing['2xl'] },
  scroll: { flex: 1 },
  content: { padding: spacing.lg, gap: spacing.md, paddingBottom: spacing['2xl'] },
  actions: { flexDirection: 'row', gap: spacing.sm },
  actionBtn: { flex: 1 },
  moduleCard: { gap: spacing.sm },
  moduleHeader: { flexDirection: 'row', alignItems: 'center', gap: spacing.sm, marginBottom: spacing.xs },
  featureRow: { flexDirection: 'row', alignItems: 'center', gap: spacing.sm, paddingVertical: spacing.sm, borderTopWidth: 1, borderTopColor: colors.border },
  flex1: { flex: 1 },
  backdrop: { flex: 1, backgroundColor: colors.overlay, justifyContent: 'flex-end' },
  sheet: { backgroundColor: colors.surface, borderTopLeftRadius: radius.lg, borderTopRightRadius: radius.lg, padding: spacing.lg, paddingBottom: spacing.xl, ...shadows.md },
});