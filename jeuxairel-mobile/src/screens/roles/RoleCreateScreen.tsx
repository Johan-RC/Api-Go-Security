import React, { useState } from 'react';
import type { NativeStackScreenProps } from '@react-navigation/native-stack';
import { Screen, Text, Input, Button } from '@/components';
import { useToast } from '@/components/Toast';
import { createRole } from '@/api/endpoints';
import { getErrorMessage } from '@/api/errors';
import { isRequired, minLength, maxLength } from '@/validation';
import { spacing } from '@/theme';
import type { RolesStackParamList } from '@/navigation/types';

type Props = NativeStackScreenProps<RolesStackParamList, 'RoleCreate'>;

export function RoleCreateScreen({ navigation }: Props) {
  const { show } = useToast();
  const [name, setName] = useState('');
  const [displayName, setDisplayName] = useState('');
  const [description, setDescription] = useState('');
  const [errors, setErrors] = useState<Record<string, string>>({});
  const [loading, setLoading] = useState(false);

  async function handleSubmit() {
    const next: Record<string, string> = {};
    if (!isRequired(name)) next.name = 'El nombre (código) es obligatorio.';
    else if (!minLength(name, 2)) next.name = 'Usa al menos 2 caracteres.';
    else if (!maxLength(name, 50)) next.name = 'Máximo 50 caracteres.';

    if (!isRequired(displayName)) next.displayName = 'El nombre para mostrar es obligatorio.';
    else if (!maxLength(displayName, 100)) next.displayName = 'Máximo 100 caracteres.';

    setErrors(next);
    if (Object.keys(next).length > 0) return;

    setLoading(true);
    try {
      await createRole({ name, display_name: displayName, description: description || undefined });
      show('Rol creado.', 'success');
      navigation.goBack();
    } catch (err) {
      show(getErrorMessage(err), 'error');
    } finally {
      setLoading(false);
    }
  }

  return (
    <Screen scroll>
      <Text variant="caption" color="textSecondary" style={{ marginBottom: spacing.lg }}>
        El nombre es un identificador único (ej.: ADMIN), el nombre para mostrar es el que ven los usuarios.
      </Text>

      <Input
        label="Nombre (código) *"
        placeholder="ADMIN"
        autoCapitalize="characters"
        value={name}
        onChangeText={setName}
        error={errors.name}
      />
      <Input
        label="Nombre para mostrar *"
        placeholder="Administrador"
        value={displayName}
        onChangeText={setDisplayName}
        error={errors.displayName}
      />
      <Input
        label="Descripción"
        placeholder="Descripción del rol (opcional)"
        multiline
        numberOfLines={3}
        value={description}
        onChangeText={setDescription}
      />

      <Button label="Crear rol" loading={loading} onPress={handleSubmit} style={{ marginTop: spacing.sm }} />
    </Screen>
  );
}