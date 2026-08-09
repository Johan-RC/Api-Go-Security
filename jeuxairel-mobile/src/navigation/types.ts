import type { NavigatorScreenParams } from '@react-navigation/native';

export type AuthStackParamList = {
  Login: undefined;
  Register: undefined;
  VerifyEmail: { email: string; devCode?: string };
  ForgotPassword: undefined;
  ResetPassword: { code?: string } | undefined;
};

export type HomeStackParamList = {
  HomeOverview: undefined;
  Sessions: undefined;
  Audit: undefined;
};

export type UsersStackParamList = {
  UserList: undefined;
  UserDetail: { userId: string };
};

export type RolesStackParamList = {
  RoleList: undefined;
  RoleCreate: undefined;
  RoleDetail: { roleId: string };
};

export type CatalogStackParamList = {
  Modules: undefined;
};

export type MoreStackParamList = {
  Profile: undefined;
};

export type MainTabParamList = {
  HomeTab: NavigatorScreenParams<HomeStackParamList>;
  UsersTab: NavigatorScreenParams<UsersStackParamList>;
  RolesTab: NavigatorScreenParams<RolesStackParamList>;
  CatalogTab: NavigatorScreenParams<CatalogStackParamList>;
  MoreTab: NavigatorScreenParams<MoreStackParamList>;
};