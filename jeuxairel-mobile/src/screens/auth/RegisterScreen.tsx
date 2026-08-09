import React from 'react';
import { StyleSheet, View } from 'react-native';
import type { NativeStackScreenProps } from '@react-navigation/native-stack';
import { Feather } from '@expo/vector-icons';
import { AuthLayout } from '@/screens/auth/AuthLayout';
import { Input, Button, Select, Text } from '@/components';
import { useToast } from '@/components/Toast';
import { register } from '@/api/endpoints';
import { getErrorMessage } from '@/api/errors';
import { useForm } from '@/hooks/useForm';
import { seq, ruleRequired, ruleEmail, ruleName, ruleStrongPassword, ruleConfirm, checkPassword } from '@/validation';
import { spacing, colors } from '@/theme';
import type { AuthStackParamList } from '@/navigation/types';
import type { ActorType } from '@/types/models';

type Props = NativeStackScreenProps<AuthStackParamList, 'Register'>;

const actorOptions: Array<{ label: string; value: ActorType }> = [
  { label: 'Usuario', value: 'USER' },
  { label: 'Instructor', value: 'INSTRUCTOR' },
  { label: 'Aprendiz', value: 'LEARNER' },
];

interface FormValues {
  firstName: string;
  lastName: string;
  email: string;
  password: string;
  confirm: string;
  actorType: ActorType;
}

export function RegisterScreen({ navigation }: Props) {
  const { show } = useToast();
  const [loading, setLoading] = React.useState(false);

  const passwordRef = React.useRef<string>('');
  const { values, errors, setField, blurField, validate } = useForm<FormValues>({
    initial: { firstName: '', lastName: '', email: '', password: '', confirm: '', actorType: 'USER' },
    validators: {
      firstName: ruleName('El nombre'),
      lastName: ruleName('El apellido'),
      email: ruleEmail,
      password: ruleStrongPassword,
      confirm: ruleConfirm('La confirmación', () => passwordRef.current),
    },
  });

  const policy = checkPassword(values.password);
  const policyItems = [
    { key: 'len' as const, label: 'Mínimo 8 caracteres', ok: policy.length },
    { key: 'upper' as const, label: 'Una mayúscula', ok: policy.upper },
    { key: 'lower' as const, label: 'Una minúscula', ok: policy.lower },
    { key: 'digit' as const, label: 'Un número', ok: policy.digit },
    { key: 'special' as const, label: 'Un símbolo', ok: policy.special },
  ];

  async function handleSubmit() {
    if (!validate()) return;

    setLoading(true);
    try {
      const { dev_code } = await register({
        email: values.email,
        password: values.password,
        first_name: values.firstName,
        last_name: values.lastName,
        actor_type: values.actorType,
      });
      show('Cuenta creada. Revisa tu correo para verificar tu cuenta.', 'success');
      navigation.replace('VerifyEmail', { email: values.email, devCode: dev_code });
    } catch (err) {
      show(getErrorMessage(err), 'error');
    } finally {
      setLoading(false);
    }
  }

  return (
    <AuthLayout
      icon="user-plus"
      title="Crear cuenta"
      subtitle="Regístrate para acceder al panel de administración"
    >
      <Input
        label="Nombre"
        placeholder="Nombre"
        autoCapitalize="words"
        value={values.firstName}
        onChangeText={(t) => setField('firstName', t)}
        onBlur={() => blurField('firstName')}
        error={errors.firstName}
      />
      <Input
        label="Apellido"
        placeholder="Apellido"
        autoCapitalize="words"
        value={values.lastName}
        onChangeText={(t) => setField('lastName', t)}
        onBlur={() => blurField('lastName')}
        error={errors.lastName}
      />
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
      />
      <Input
        label="Contraseña"
        placeholder="Mínimo 8 caracteres"
        secureTextEntry
        autoCapitalize="none"
        value={values.password}
        onChangeText={(t) => {
          passwordRef.current = t;
          setField('password', t);
        }}
        onBlur={() => blurField('password')}
        error={errors.password}
      />
      {values.password.length > 0 ? (
        <View style={styles.policy}>
          {policyItems.map((item) => (
            <View key={item.key} style={styles.policyRow}>
              <Feather name={item.ok ? 'check-circle' : 'circle'} size={14} color={item.ok ? colors.success : colors.textMuted} />
              <Text variant="caption" color={item.ok ? 'textSecondary' : 'textMuted'}>
                {item.label}
              </Text>
            </View>
          ))}
        </View>
      ) : null}
      <Input
        label="Confirmar contraseña"
        placeholder="Repite la contraseña"
        secureTextEntry
        autoCapitalize="none"
        value={values.confirm}
        onChangeText={(t) => setField('confirm', t)}
        onBlur={() => blurField('confirm')}
        error={errors.confirm}
      />
      <Select
        label="Tipo de usuario"
        value={values.actorType}
        options={actorOptions}
        onChange={(v) => setField('actorType', v)}
      />

      <Button label="Crear cuenta" loading={loading} onPress={handleSubmit} style={{ marginTop: spacing.sm }} />
      <Button variant="ghost" label="Ya tengo cuenta" onPress={() => navigation.goBack()} />
    </AuthLayout>
  );
}

const styles = StyleSheet.create({
  policy: {
    backgroundColor: colors.background,
    borderRadius: 8,
    padding: spacing.md,
    marginTop: -spacing.sm,
    marginBottom: spacing.md,
    gap: spacing.xs,
  },
  policyRow: { flexDirection: 'row', alignItems: 'center', gap: spacing.sm },
});