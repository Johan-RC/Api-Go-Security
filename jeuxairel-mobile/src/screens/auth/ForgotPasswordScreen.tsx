import React, { useState } from 'react';
import type { NativeStackScreenProps } from '@react-navigation/native-stack';
import { Feather } from '@expo/vector-icons';
import { AuthLayout } from '@/screens/auth/AuthLayout';
import { Input, Button } from '@/components';
import { useToast } from '@/components/Toast';
import { forgotPassword } from '@/api/endpoints';
import { getErrorMessage } from '@/api/errors';
import { isEmail, isRequired } from '@/validation';
import { spacing, colors } from '@/theme';
import type { AuthStackParamList } from '@/navigation/types';

type Props = NativeStackScreenProps<AuthStackParamList, 'ForgotPassword'>;

export function ForgotPasswordScreen({ navigation }: Props) {
  const { show } = useToast();
  const [email, setEmail] = useState('');
  const [error, setError] = useState<string | undefined>();
  const [loading, setLoading] = useState(false);

  async function handleSubmit() {
    if (!isRequired(email)) {
      setError('El correo es obligatorio.');
      return;
    }
    if (!isEmail(email)) {
      setError('Ingresa un correo válido.');
      return;
    }
    setError(undefined);

    setLoading(true);
    try {
      const { code } = await forgotPassword(email);
      if (code) {
        show(`[Desarrollo] Código de recuperación: ${code}`, 'info');
        navigation.replace('ResetPassword', { code });
      } else {
        show('Si el correo existe, recibirás un código de 6 dígitos para restablecer tu contraseña.', 'success');
        navigation.replace('Login');
      }
    } catch (err) {
      show(getErrorMessage(err), 'error');
    } finally {
      setLoading(false);
    }
  }

  return (
    <AuthLayout
      icon="key"
      title="Recuperar contraseña"
      subtitle="Te enviaremos un código de 6 dígitos para restablecer tu acceso"
    >
      <Input
        label="Correo electrónico"
        placeholder="tunombre@correo.com"
        autoCapitalize="none"
        keyboardType="email-address"
        autoCorrect={false}
        value={email}
        onChangeText={setEmail}
        error={error}
        leftIcon={<Feather name="mail" size={18} color={colors.textSecondary} />}
      />
      <Button label="Enviar instrucciones" loading={loading} onPress={handleSubmit} style={{ marginTop: spacing.sm }} />
      <Button variant="ghost" label="Volver al inicio de sesión" onPress={() => navigation.goBack()} />
    </AuthLayout>
  );
}