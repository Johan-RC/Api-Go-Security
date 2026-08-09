import AsyncStorage from '@react-native-async-storage/async-storage';

const ACCESS_KEY = '@jeuxairel/access_token';
const REFRESH_KEY = '@jeuxairel/refresh_token';

export interface StoredTokens {
  accessToken: string | null;
  refreshToken: string | null;
}

export async function saveTokens(accessToken: string | null, refreshToken: string | null) {
  await AsyncStorage.multiSet([
    [ACCESS_KEY, accessToken ?? ''],
    [REFRESH_KEY, refreshToken ?? ''],
  ]);
}

export async function loadTokens(): Promise<StoredTokens> {
  const values = await AsyncStorage.multiGet([ACCESS_KEY, REFRESH_KEY]);
  const [accessToken, refreshToken] = values.map(([, value]) => value || null);
  return { accessToken, refreshToken };
}

export async function clearTokens() {
  await AsyncStorage.multiRemove([ACCESS_KEY, REFRESH_KEY]);
}