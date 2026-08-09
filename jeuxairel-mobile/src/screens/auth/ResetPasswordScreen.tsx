import React, { useState } from 'react';
import type { NativeStackScreenProps } from '@react-navigation/native-stack';
import { AuthLayout } from '@/screens/auth/AuthLayout';
import { Input, Button } from '@/components';
import { useToast } from '@/components/Toast';
import { resetPassword } from '@/api/endpoints';
import { getErrorMessage } from '@/api/errors';
import { isRequired, minLength } from '@/validation';
import { spacing } from '@/theme';
import type { AuthStackParamList } from '@/navigation/types';

type Props = NativeStackScreenProps<AuthStackParamList, 'ResetPassword'>;

export function ResetPasswordScreen({ route, navigation }: Props) {
  const { show } = useToast();

  const [code, setCode] = useState(route.params?.code ?? '');
  const [newPassword, setNewPassword] = useState('');
  const [confirm, setConfirm] = useState('');
  const [errors, setErrors] = useState<Record<string, string>>({});
  const [loading, setLoading] = useState(false);

  async function handleSubmit() {
    const next: Record<string, string> = {};
    if (!isRequired(code)) next.code = 'El código es obligatorio.';
    else if (!/^\d{6}$/.test(code)) next.code = 'El código debe tener exactamente 6 dígitos.';
    if (!isRequired(newPassword)) next.newPassword = 'La nueva contraseña es obligatoria.';
    else if (!minLength(newPassword, 8))
      next.newPassword = 'La contraseña debe tener al menos 8 caracteres.';
    if (confirm !== newPassword) next.confirm = 'Las contraseñas no coinciden.';

    setErrors(next);
    if (Object.keys(next).length > 0) return;

    setLoading(true);
    try {
      await resetPassword(code, newPassword);
      show('Contraseña actualizada correctamente.', 'success');
      navigation.replace('Login');
    } catch (err) {
      show(getErrorMessage(err), 'error');
    } finally {
      setLoading(false);
    }
  }

  return (
    <AuthLayout
      icon="refresh-cw"
      title="Nueva contraseña"
      subtitle="Ingresa el código de 6 dígitos recibido y define tu nueva contraseña"
    >
      <Input
        label="Código de recuperación"
        placeholder="Código de 6 dígitos recibido por correo"
        keyboardType="number-pad"
        maxLength={6}
        value={code}
        onChangeText={setCode}
        error={errors.code}
      />
      <Input
        label="Nueva contraseña"
        placeholder="Mínimo 8 caracteres"
        secureTextEntry
        autoCapitalize="none"
        value={newPassword}
        onChangeText={setNewPassword}
        error={errors.newPassword}
      />
      <Input
        label="Confirmar contraseña"
        placeholder="Repite la nueva contraseña"
        secureTextEntry
        autoCapitalize="none"
        value={confirm}
        onChangeText={setConfirm}
        error={errors.confirm}
      />
      <Button label="Restablecer contraseña" loading={loading} onPress={handleSubmit} style={{ marginTop: spacing.sm }} />
      <Button variant="ghost" label="Volver al inicio de sesión" onPress={() => navigation.goBack()} />
    </AuthLayout>
  );
}