import React, { useEffect, useRef } from 'react';
import { Animated, StyleProp, ViewStyle } from 'react-native';
import { nativeDriver } from '@/utils/web';

interface RevealProps {
  delay?: number;
  distance?: number;
  style?: StyleProp<ViewStyle>;
  children: React.ReactNode;
}

export function Reveal({ delay = 0, distance = 14, style, children }: RevealProps) {
  const opacity = useRef(new Animated.Value(0)).current;
  const translateY = useRef(new Animated.Value(distance)).current;

  useEffect(() => {
    Animated.sequence([
      Animated.delay(delay),
      Animated.parallel([
        Animated.timing(opacity, { toValue: 1, duration: 450, useNativeDriver: nativeDriver }),
        Animated.spring(translateY, { toValue: 0, speed: 14, bounciness: 6, useNativeDriver: nativeDriver }),
      ]),
    ]).start();
  }, [delay, opacity, translateY, distance]);

  return (
    <Animated.View style={[{ opacity, transform: [{ translateY }] }, style]}>
      {children}
    </Animated.View>
  );
}