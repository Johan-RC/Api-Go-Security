import React, { useCallback, useState } from 'react';
import { StyleSheet, View, Alert, Switch, Pressable } from 'react-native';
import type { NativeStackScreenProps } from '@react-navigation/native-stack';
import { Feather } from '@expo/vector-icons';
import { useFocusEffect } from '@react-navigation/native';
import { Screen, Text, Card, Input, Button, Badge, Select } from '@/components';
import { Spinner } from '@/components/Spinner';
import { useToast } from '@/components/Toast';
import { useAuth } from '@/contexts/AuthContext';
import {
  getUser,
  createUser,
  updateUser,
  deleteUser,
  listUserRoles,
  assignRoleToUser,
  removeUserRole,
  listRoles,
} from '@/api/endpoints';
import { getErrorMessage } from '@/api/errors';
import { isEmail, isRequired, minLength } from '@/validation';
import { colors, spacing } from '@/theme';
import { formatDate } from '@/utils/format';
import type { UserResponse, RoleResponse, UserRoleResponse, ActorType } from '@/types/models';
import type { UsersStackParamList } from '@/navigation/types';

type Props = NativeStackScreenProps<UsersStackParamList, 'UserDetail'>;

const actorOptions: Array<{ label: string; value: ActorType }> = [
  { label: 'Usuario', value: 'USER' },
  { label: 'Instructor', value: 'INSTRUCTOR' },
  { label: 'Aprendiz', value: 'LEARNER' },
];

const actorLabel = (v: string) => actorOptions.find((a) => a.value === v)?.label ?? v;

// ---------------------------------------------------------------------------
// Modo creación de usuario
// ---------------------------------------------------------------------------

function CreateUserForm({ navigation }: { navigation: Props['navigation'] }) {
  const { show } = useToast();
  const [firstName, setFirstName] = useState('');
  const [lastName, setLastName] = useState('');
  const [email, setEmail] = useState('');
  const [password, setPassword] = useState('');
  const [actorType, setActorType] = useState<ActorType>('USER');
  const [errors, setErrors] = useState<Record<string, string>>({});
  const [loading, setLoading] = useState(false);

  async function handleSubmit() {
    const next: Record<string, string> = {};
    if (!isRequired(firstName)) next.firstName = 'El nombre es obligatorio.';
    if (!isRequired(lastName)) next.lastName = 'El apellido es obligatorio.';
    if (!isRequired(email)) next.email = 'El correo es obligatorio.';
    else if (!isEmail(email)) next.email = 'Ingresa un correo válido.';
    if (!isRequired(password)) next.password = 'La contraseña es obligatoria.';
    else if (!minLength(password, 8)) next.password = 'Mínimo 8 caracteres.';

    setErrors(next);
    if (Object.keys(next).length > 0) return;

    setLoading(true);
    try {
      await createUser({ email, password, first_name: firstName, last_name: lastName, actor_type: actorType });
      show('Usuario creado correctamente.', 'success');
      navigation.goBack();
    } catch (err) {
      show(getErrorMessage(err), 'error');
    } finally {
      setLoading(false);
    }
  }

  return (
    <Screen scroll>
      <Text variant="h3" style={{ marginBottom: spacing.lg }}>Nuevo usuario</Text>
      <Input label="Nombre *" value={firstName} onChangeText={setFirstName} error={errors.firstName} />
      <Input label="Apellido *" value={lastName} onChangeText={setLastName} error={errors.lastName} />
      <Input
        label="Correo electrónico *"
        value={email}
        onChangeText={setEmail}
        autoCapitalize="none"
        keyboardType="email-address"
        error={errors.email}
      />
      <Input
        label="Contraseña *"
        value={password}
        onChangeText={setPassword}
        secureTextEntry
        autoCapitalize="none"
        error={errors.password}
      />
      <Select label="Tipo de usuario" value={actorType} options={actorOptions} onChange={setActorType} />
      <Button label="Crear usuario" loading={loading} onPress={handleSubmit} style={{ marginTop: spacing.sm }} />
    </Screen>
  );
}

// ---------------------------------------------------------------------------
// Modo detalle de usuario
// ---------------------------------------------------------------------------

export function UserDetailScreen({ route, navigation }: Props) {
  const userId = route.params.userId;
  const isNew = userId === 'new';

  const { show } = useToast();
  const { can } = useAuth();
  const canManage = can('IDENTITY_USER_MANAGE', 'WRITE');
  const canAssignRoles = can('IDENTITY_ROLE_ASSIGN', 'WRITE');
  const [user, setUser] = useState<UserResponse | null>(null);
  const [roles, setRoles] = useState<UserRoleResponse[]>([]);
  const [availableRoles, setAvailableRoles] = useState<RoleResponse[]>([]);
  const [loading, setLoading] = useState(!isNew);

  const refresh = useCallback(async () => {
    const [u, rl, ur] = await Promise.allSettled([
      getUser(userId),
      listRoles(1, 100),
      listUserRoles(userId),
    ]);
    if (u.status === 'fulfilled') setUser(u.value);
    if (rl.status === 'fulfilled') setAvailableRoles(rl.value.items);
    if (ur.status === 'fulfilled') setRoles(ur.value);
  }, [userId]);

  useFocusEffect(
    useCallback(() => {
      if (isNew) return;
      refresh()
        .catch((err) => show(getErrorMessage(err), 'error'))
        .finally(() => setLoading(false));
    }, [refresh, isNew, show]),
  );

  if (isNew) return <CreateUserForm navigation={navigation} />;
  if (loading) return <Screen><Spinner /></Screen>;
  if (!user) return <Screen><Text color="textSecondary">No se encontró el usuario.</Text></Screen>;

  const assignableRoles = availableRoles.filter((r) => !roles.some((ur) => ur.role_id === r.id));

  return (
    <Screen scroll>
      <Card>
        <View style={styles.headerRow}>
          <View style={styles.flex1}>
            <Text variant="h3">{`${user.first_name} ${user.last_name}`}</Text>
            <Text variant="caption" color="textSecondary">{user.email}</Text>
          </View>
          <Badge tone={user.is_active ? 'success' : 'neutral'} label={user.is_active ? 'Activo' : 'Inactivo'} />
        </View>

        <View style={styles.dataGrid}>
          <InfoCell label="Tipo" value={actorLabel(user.actor_type)} />
          <InfoCell label="Último acceso" value={formatDate(user.last_login_at)} />
          <InfoCell label="Creado" value={formatDate(user.created_at)} />
        </View>

        <View style={styles.switchRow}>
          <Text variant="body">Cuenta activa</Text>
          <Switch value={user.is_active} disabled={!canManage} onValueChange={handleToggleActive} />
        </View>
      </Card>

      <Text variant="h3" style={styles.sectionTitle}>Roles asignados</Text>
      <Card>
        {roles.length === 0 ? (
          <Text variant="caption" color="textSecondary" style={{ marginBottom: spacing.sm }}>
            Sin roles asignados.
          </Text>
        ) : (
          roles.map((ur) => (
            <Pressable
              key={ur.id}
              style={styles.roleRow}
              disabled={!canAssignRoles}
              onPress={() =>
                canAssignRoles &&
                Alert.alert('Quitar rol', `¿Quitar el rol "${ur.display_name}"?`, [
                  { text: 'Cancelar', style: 'cancel' },
                  {
                    text: 'Quitar',
                    style: 'destructive',
                    onPress: () =>
                      removeUserRole(userId, ur.id)
                        .then(() => { show('Rol quitado.', 'success'); refresh(); })
                        .catch((err) => show(getErrorMessage(err), 'error')),
                  },
                ])
              }
            >
              <Feather name="shield" size={16} color={colors.primary} />
              <Text variant="body" style={styles.flex1}>{ur.display_name}</Text>
              {canAssignRoles ? <Feather name="x" size={15} color={colors.textMuted} /> : null}
            </Pressable>
          ))
        )}
        {canAssignRoles ? (
          <Select
            label="Asignar un rol"
            placeholder="Elige un rol…"
            value={null}
            options={assignableRoles.map((r) => ({ label: r.display_name, value: r.id }))}
            onChange={handleAssignRole}
            style={{ marginTop: spacing.md }}
          />
        ) : null}
      </Card>

      {canManage ? (
        <Button variant="danger" label="Eliminar usuario" onPress={confirmDelete} style={{ marginTop: spacing.xl }} />
      ) : null}
    </Screen>
  );

  function handleToggleActive(v: boolean) {
    updateUser(userId, { is_active: v })
      .then((updated) => {
        setUser((prev) => (prev ? { ...prev, is_active: updated.is_active } : prev));
        show('Estado actualizado.', 'success');
      })
      .catch((err) => show(getErrorMessage(err), 'error'));
  }

  function handleAssignRole(roleId: string) {
    const role = availableRoles.find((r) => r.id === roleId);
    assignRoleToUser(userId, roleId)
      .then(() => {
        show(`Rol "${role?.display_name}" asignado.`, 'success');
        refresh();
      })
      .catch((err) => show(getErrorMessage(err), 'error'));
  }

  function confirmDelete() {
    const targetId = user?.id ?? '';
    if (!targetId) return;
    Alert.alert('Eliminar usuario', 'Esta acción desactivará al usuario. ¿Continuar?', [
      { text: 'Cancelar', style: 'cancel' },
      {
        text: 'Eliminar',
        style: 'destructive',
        onPress: () =>
          deleteUser(targetId)
            .then(() => {
              show('Usuario eliminado.', 'success');
              navigation.goBack();
            })
            .catch((err) => show(getErrorMessage(err), 'error')),
      },
    ]);
  }
}

function InfoCell({ label, value }: { label: string; value: string }) {
  return (
    <View>
      <Text variant="caption" color="textSecondary">{label}</Text>
      <Text variant="body">{value}</Text>
    </View>
  );
}

const styles = StyleSheet.create({
  flex1: { flex: 1 },
  headerRow: { flexDirection: 'row', alignItems: 'center', gap: spacing.md },
  dataGrid: { gap: spacing.md, marginTop: spacing.lg, marginBottom: spacing.md },
  switchRow: { flexDirection: 'row', alignItems: 'center', justifyContent: 'space-between', borderTopWidth: 1, borderTopColor: colors.border, paddingTop: spacing.md },
  sectionTitle: { marginTop: spacing.xl, marginBottom: spacing.md },
  roleRow: { flexDirection: 'row', alignItems: 'center', gap: spacing.sm, paddingVertical: spacing.sm },
});