import React from 'react';
import { ActivityIndicator, View, StyleSheet } from 'react-native';
import { createNativeStackNavigator } from '@react-navigation/native-stack';
import { createBottomTabNavigator } from '@react-navigation/bottom-tabs';
import { Feather } from '@expo/vector-icons';
import { useAuth } from '@/contexts/AuthContext';
import { colors } from '@/theme';

import type {
  AuthStackParamList,
  HomeStackParamList,
  UsersStackParamList,
  RolesStackParamList,
  CatalogStackParamList,
  MoreStackParamList,
  MainTabParamList,
} from '@/navigation/types';

import { LoginScreen } from '@/screens/auth/LoginScreen';
import { RegisterScreen } from '@/screens/auth/RegisterScreen';
import { VerifyEmailScreen } from '@/screens/auth/VerifyEmailScreen';
import { ForgotPasswordScreen } from '@/screens/auth/ForgotPasswordScreen';
import { ResetPasswordScreen } from '@/screens/auth/ResetPasswordScreen';

import { HomeScreen } from '@/screens/home/HomeScreen';
import { SessionsScreen } from '@/screens/home/SessionsScreen';
import { AuditScreen } from '@/screens/home/AuditScreen';

import { UsersScreen } from '@/screens/users/UsersScreen';
import { UserDetailScreen } from '@/screens/users/UserDetailScreen';

import { RolesScreen } from '@/screens/roles/RolesScreen';
import { RoleCreateScreen } from '@/screens/roles/RoleCreateScreen';
import { RoleDetailScreen } from '@/screens/roles/RoleDetailScreen';

import { CatalogScreen } from '@/screens/catalog/CatalogScreen';
import { ProfileScreen } from '@/screens/profile/ProfileScreen';

const AuthStack = createNativeStackNavigator<AuthStackParamList>();
const HomeStack = createNativeStackNavigator<HomeStackParamList>();
const UsersStack = createNativeStackNavigator<UsersStackParamList>();
const RolesStack = createNativeStackNavigator<RolesStackParamList>();
const CatalogStack = createNativeStackNavigator<CatalogStackParamList>();
const MoreStack = createNativeStackNavigator<MoreStackParamList>();
const Tab = createBottomTabNavigator<MainTabParamList>();

const header = {
  headerStyle: { backgroundColor: colors.surface },
  headerShadowVisible: false,
  headerTintColor: colors.primary,
  headerTitleStyle: { color: colors.text, fontWeight: '700' as const },
  headerTitleAlign: 'center' as const,
} as const;

function AuthNavigator() {
  return (
    <AuthStack.Navigator screenOptions={{ headerShown: false }}>
      <AuthStack.Screen name="Login" component={LoginScreen} />
      <AuthStack.Screen name="Register" component={RegisterScreen} />
      <AuthStack.Screen name="VerifyEmail" component={VerifyEmailScreen} />
      <AuthStack.Screen name="ForgotPassword" component={ForgotPasswordScreen} />
      <AuthStack.Screen name="ResetPassword" component={ResetPasswordScreen} />
    </AuthStack.Navigator>
  );
}

function HomeNavigator() {
  return (
    <HomeStack.Navigator screenOptions={header}>
      <HomeStack.Screen name="HomeOverview" component={HomeScreen} options={{ title: 'Inicio' }} />
      <HomeStack.Screen name="Sessions" component={SessionsScreen} options={{ title: 'Sesiones activas' }} />
      <HomeStack.Screen name="Audit" component={AuditScreen} options={{ title: 'Auditoría de accesos' }} />
    </HomeStack.Navigator>
  );
}

function UsersNavigator() {
  return (
    <UsersStack.Navigator screenOptions={header}>
      <UsersStack.Screen name="UserList" component={UsersScreen} options={{ title: 'Usuarios' }} />
      <UsersStack.Screen name="UserDetail" component={UserDetailScreen} options={{ title: 'Detalle de usuario' }} />
    </UsersStack.Navigator>
  );
}

function RolesNavigator() {
  return (
    <RolesStack.Navigator screenOptions={header}>
      <RolesStack.Screen name="RoleList" component={RolesScreen} options={{ title: 'Roles' }} />
      <RolesStack.Screen name="RoleCreate" component={RoleCreateScreen} options={{ title: 'Crear rol' }} />
      <RolesStack.Screen name="RoleDetail" component={RoleDetailScreen} options={{ title: 'Detalle de rol' }} />
    </RolesStack.Navigator>
  );
}

function CatalogNavigator() {
  return (
    <CatalogStack.Navigator screenOptions={header}>
      <CatalogStack.Screen name="Modules" component={CatalogScreen} options={{ title: 'Módulos y funcionalidades' }} />
    </CatalogStack.Navigator>
  );
}

function MoreNavigator() {
  return (
    <MoreStack.Navigator screenOptions={header}>
      <MoreStack.Screen name="Profile" component={ProfileScreen} options={{ title: 'Mi perfil' }} />
    </MoreStack.Navigator>
  );
}

type FeatherIcon = keyof typeof Feather.glyphMap;

const tabIcons: Record<keyof MainTabParamList, { active: FeatherIcon; inactive: FeatherIcon }> = {
  HomeTab: { active: 'home', inactive: 'home' },
  UsersTab: { active: 'users', inactive: 'users' },
  RolesTab: { active: 'shield', inactive: 'shield' },
  CatalogTab: { active: 'grid', inactive: 'grid' },
  MoreTab: { active: 'user', inactive: 'user' },
};

function MainTabs() {
  const { can } = useAuth();

  const showUsers = can('IDENTITY_USER_VIEW', 'READ');
  const showRoles = can('IDENTITY_ROLE_VIEW', 'READ');
  const showCatalog = can('IDENTITY_ROLE_MANAGE', 'WRITE');

  return (
    <Tab.Navigator
      screenOptions={({ route }) => ({
        headerShown: false,
        tabBarActiveTintColor: colors.primary,
        tabBarInactiveTintColor: colors.textMuted,
        tabBarStyle: {
          backgroundColor: colors.surface,
          borderTopColor: colors.border,
          height: 60,
          paddingBottom: 6,
          paddingTop: 6,
        },
        tabBarLabelStyle: { fontSize: 11, fontWeight: '600' },
        tabBarIcon: ({ color, focused }) => {
          const icon = tabIcons[route.name as keyof MainTabParamList];
          const name = focused ? icon.active : icon.inactive;
          return (
            <View style={focused ? styles.tabIconActive : null}>
              <Feather name={name} size={22} color={color} />
            </View>
          );
        },
      })}
    >
      <Tab.Screen name="HomeTab" component={HomeNavigator} options={{ tabBarLabel: 'Inicio' }} />
      {showUsers ? (
        <Tab.Screen name="UsersTab" component={UsersNavigator} options={{ tabBarLabel: 'Usuarios' }} />
      ) : null}
      {showRoles ? (
        <Tab.Screen name="RolesTab" component={RolesNavigator} options={{ tabBarLabel: 'Roles' }} />
      ) : null}
      {showCatalog ? (
        <Tab.Screen name="CatalogTab" component={CatalogNavigator} options={{ tabBarLabel: 'Catálogo' }} />
      ) : null}
      <Tab.Screen name="MoreTab" component={MoreNavigator} options={{ tabBarLabel: 'Perfil' }} />
    </Tab.Navigator>
  );
}

export function RootNavigator() {
  const { isAuthenticated, isLoading } = useAuth();

  if (isLoading) {
    return (
      <View style={styles.loading}>
        <ActivityIndicator size="large" color={colors.primary} />
      </View>
    );
  }

  return isAuthenticated ? <MainTabs /> : <AuthNavigator />;
}

const styles = StyleSheet.create({
  loading: { flex: 1, alignItems: 'center', justifyContent: 'center', backgroundColor: colors.background },
  tabIconActive: {
    paddingHorizontal: 14,
    paddingVertical: 3,
    borderRadius: 999,
    backgroundColor: colors.primarySoft,
  },
});