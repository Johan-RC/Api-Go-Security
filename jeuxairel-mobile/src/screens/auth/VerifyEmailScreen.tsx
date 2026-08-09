import React, { useState } from 'react';
import { StyleSheet, View, ScrollView } from 'react-native';
import type { NativeStackScreenProps } from '@react-navigation/native-stack';
import { Feather } from '@expo/vector-icons';
import { AuthLayout } from '@/screens/auth/AuthLayout';
import { Input, Button, Text } from '@/components';
import { useToast } from '@/components/Toast';
import { verifyEmail } from '@/api/endpoints';
import { getErrorMessage } from '@/api/errors';
import { ruleSixDigitCode } from '@/validation';
import { spacing } from '@/theme';
import type { AuthStackParamList } from '@/navigation/types';

type Props = NativeStackScreenProps<AuthStackParamList, 'VerifyEmail'>;

export function VerifyEmailScreen({ route, navigation }: Props) {
  const { show } = useToast();
  const { email, devCode } = route.params;

  const [code, setCode] = useState(devCode ?? '');
  const [error, setError] = useState<string | undefined>();
  const [loading, setLoading] = useState(false);

  async function handleVerify() {
    const errorMsg = ruleSixDigitCode(code);
    if (errorMsg) {
      setError(errorMsg);
      return;
    }
    setError(undefined);

    setLoading(true);
    try {
      await verifyEmail(email, code);
      show('Cuenta verificada correctamente. Ya puedes iniciar sesión.', 'success');
      navigation.replace('Login');
    } catch (err) {
      show(getErrorMessage(err), 'error');
    } finally {
      setLoading(false);
    }
  }

  function handleResend() {
    show('Por ahora puedes pedir un nuevo código desde "Recuperar contraseña".', 'info');
  }

  return (
    <AuthLayout
      icon="mail"
      title="Verifica tu correo"
      subtitle={`Ingresa el código de 6 dígitos que enviamos a ${email}`}
    >
      <Input
        label="Código de verificación"
        placeholder="Código de 6 dígitos"
        keyboardType="number-pad"
        maxLength={6}
        value={code}
        onChangeText={(t) => {
          setCode(t);
          if (error) setError(undefined);
        }}
        error={error}
      />

      <Button label="Verificar cuenta" loading={loading} onPress={handleVerify} style={{ marginTop: spacing.sm }} />
      <Button variant="ghost" label="Reenviar código" onPress={handleResend} />

      <Text variant="caption" color="textSecondary" align="center" style={{ marginTop: spacing.lg }}>
        Si no encuentras el correo, revisa también la carpeta de spam.
      </Text>
      <View>
        <Button variant="outline" label="Ir a iniciar sesión" onPress={() => navigation.replace('Login')} />
      </View>
    </AuthLayout>
  );
}

const styles = StyleSheet.create({});