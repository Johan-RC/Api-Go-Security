import { Platform } from 'react-native';
import type { TextStyle, ViewStyle } from 'react-native';

// `cursor` solo existe en react-native-web; en nativo se ignora sin romper nada.
export const pointer: ViewStyle = { cursor: 'pointer' } as unknown as ViewStyle;

export const pointerText: TextStyle = { cursor: 'pointer' } as unknown as TextStyle;

// En web react-native falla al usar el driver nativo y lo degrada a JS con warning.
export const nativeDriver = Platform.OS !== 'web';