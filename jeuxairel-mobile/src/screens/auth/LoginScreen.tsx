import React from 'react';
import { StyleSheet, View } from 'react-native';
import type { NativeStackScreenProps } from '@react-navigation/native-stack';
import { Feather } from '@expo/vector-icons';
import { AuthLayout } from '@/screens/auth/AuthLayout';
import { Input, Button, Text } from '@/components';
import { useToast } from '@/components/Toast';
import { useAuth } from '@/contexts/AuthContext';
import { getErrorMessage } from '@/api/errors';
import { useForm } from '@/hooks/useForm';
import { ruleRequired, ruleEmail } from '@/validation';
import { colors, spacing } from '@/theme';
import type { AuthStackParamList } from '@/navigation/types';

type Props = NativeStackScreenProps<AuthStackParamList, 'Login'>;

interface FormValues {
  email: string;
  password: string;
}

export function LoginScreen({ navigation }: Props) {
  const { login } = useAuth();
  const { show } = useToast();
  const [showPass, setShowPass] = React.useState(false);
  const [loading, setLoading] = React.useState(false);

  const { values, errors, setField, blurField, validate } = useForm<FormValues>({
    initial: { email: '', password: '' },
    validators: {
      email: ruleEmail,
      password: ruleRequired('La contraseña'),
    },
  });

  async function handleSubmit() {
    if (!validate()) return;

    setLoading(true);
    try {
      await login(values.email, values.password);
    } catch (err) {
      show(getErrorMessage(err), 'error');
    } finally {
      setLoading(false);
    }
  }

  return (
    <AuthLayout
      title="Bienvenido de nuevo"
      subtitle="Ingresa a la plataforma de identidad y accesos"
    >
      <Input
        label="Correo electrónico"
        placeholder="tunombre@correo.com"
        autoCapitalize="none"
        keyboardType="email-address"
        autoCorrect={false}
        value={values.email}
        onChangeText={(t) => setField('email', t)}
        onBlur={() => blurField('email')}
        error={errors.email}
        leftIcon={<Feather name="mail" size={18} color={colors.textMuted} />}
      />

      <Input
        label="Contraseña"
        placeholder="••••••••"
        secureTextEntry={!showPass}
        autoCapitalize="none"
        value={values.password}
        onChangeText={(t) => setField('password', t)}
        error={errors.password}
        leftIcon={<Feather name="lock" size={18} color={colors.textMuted} />}
        rightIcon={
          <Feather name={showPass ? 'eye-off' : 'eye'} size={18} color={colors.textMuted} />
        }
        onRightPress={() => setShowPass((v) => !v)}
      />

      <Button label="Iniciar sesión" loading={loading} onPress={handleSubmit} style={{ marginTop: spacing.sm }} />

      <TouchableLink onPress={() => navigation.navigate('ForgotPassword')} label="¿Olvidaste tu contraseña?" />

      <View style={styles.footer}>
        <Text variant="body" color="textSecondary">
          ¿No tienes cuenta?
        </Text>
        <TouchableLink onPress={() => navigation.navigate('Register')} label="Regístrate" />
      </View>
    </AuthLayout>
  );
}

function TouchableLink({ label, onPress }: { label: string; onPress: () => void }) {
  return (
    <Text variant="body" color="primary" align="center" bold style={styles.link} onPress={onPress}>
      {label}
    </Text>
  );
}

const styles = StyleSheet.create({
  link: { marginTop: spacing.lg },
  footer: { flexDirection: 'row', justifyContent: 'center', alignItems: 'center', marginTop: spacing.xl, gap: spacing.xs },
});