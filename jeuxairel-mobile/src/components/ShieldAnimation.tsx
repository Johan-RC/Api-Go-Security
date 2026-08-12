import React, { useEffect, useRef } from 'react';
import { Animated, Easing, StyleSheet, View } from 'react-native';
import Svg, { Path, Rect, Circle, Ellipse, G, Defs, LinearGradient, RadialGradient, Stop } from 'react-native-svg';
import { nativeDriver } from '@/utils/web';

interface ShieldAnimationProps {
  size?: number;
}

// Escena: un caballero sostiene el escudo y repele ataques mientras protege una computadora.
// viewBox 160x110, `size` es la altura en píxeles.
export function ShieldAnimation({ size = 80 }: ShieldAnimationProps) {
  const height = size;
  const width = Math.round(size * 1.45);
  const s = width / 160;

  // Centro del escudo en coords de la escena (viewBox) y píxeles.
  const cx = 44 * s;
  const cy = 56 * s;

  const pulse = useRef(new Animated.Value(0)).current;
  const ringA = useRef(new Animated.Value(0)).current;
  const ringB = useRef(new Animated.Value(0)).current;
  const pa = useRef(new Animated.Value(0)).current;
  const pb = useRef(new Animated.Value(0)).current;
  const pc = useRef(new Animated.Value(0)).current;

  useEffect(() => {
    const loops: Animated.CompositeAnimation[] = [];

    loops.push(
      Animated.loop(
        Animated.sequence([
          Animated.timing(pulse, { toValue: 1, duration: 900, easing: Easing.inOut(Easing.sin), useNativeDriver: nativeDriver }),
          Animated.timing(pulse, { toValue: 0, duration: 900, easing: Easing.inOut(Easing.sin), useNativeDriver: nativeDriver }),
        ]),
      ),
    );

    const ring = (v: Animated.Value, delay: number) =>
      Animated.loop(
        Animated.sequence([
          Animated.delay(delay),
          Animated.timing(v, { toValue: 1, duration: 1300, easing: Easing.out(Easing.ease), useNativeDriver: nativeDriver }),
          Animated.timing(v, { toValue: 0, duration: 0, useNativeDriver: nativeDriver }),
        ]),
      );
    loops.push(ring(ringA, 0), ring(ringB, 650));

    const proj = (v: Animated.Value, delay: number) =>
      Animated.loop(
        Animated.sequence([
          Animated.delay(delay),
          Animated.timing(v, { toValue: 1, duration: 1200, easing: Easing.inOut(Easing.quad), useNativeDriver: nativeDriver }),
          Animated.delay(450),
        ]),
      );
    loops.push(proj(pa, 0), proj(pb, 620), proj(pc, 1240));

    loops.forEach((l) => l.start());
    return () => loops.forEach((l) => l.stop());
  }, [pulse, ringA, ringB, pa, pb, pc]);

  const shieldKick = pulse.interpolate({ inputRange: [0, 1], outputRange: [1, 1.06] });

  const ringScale = (v: Animated.Value) =>
    v.interpolate({ inputRange: [0, 1], outputRange: [0.3, 2.3] });
  const ringOpacity = (v: Animated.Value) =>
    v.interpolate({ inputRange: [0, 0.85, 1], outputRange: [0.55, 0.2, 0] });

  const dot = Math.max(7, size * 0.08);

  const proj = (v: Animated.Value, fromX: number, fromY: number, toX: number, toY: number) => ({
    x: v.interpolate({ inputRange: [0, 0.85, 1], outputRange: [0, (toX - fromX) * s, (toX - fromX) * s] }),
    y: v.interpolate({ inputRange: [0, 0.85, 1], outputRange: [0, (toY - fromY) * s, (toY - fromY) * s] }),
    o: v.interpolate({ inputRange: [0, 0.12, 0.8, 1], outputRange: [0, 1, 1, 0] }),
  });

  const p1 = proj(pa, 2, 40, 20, 50);
  const p2 = proj(pb, 2, 58, 20, 58);
  const p3 = proj(pc, 2, 76, 20, 66);

  const ringBase = 14 * s;

  return (
    <View style={{ width, height }}>
      {/* Ondas de energía (repelen) */}
      {[
        { v: ringA, color: 'rgba(239,68,68,0.55)' },
        { v: ringB, color: 'rgba(255,255,255,0.45)' },
      ].map((r, i) => (
        <Animated.View
          key={i}
          pointerEvents="none"
          style={{
            position: 'absolute',
            left: cx - ringBase / 2,
            top: cy - ringBase / 2,
            width: ringBase,
            height: ringBase,
            borderRadius: ringBase / 2,
            borderWidth: 2,
            borderColor: r.color,
            opacity: ringOpacity(r.v),
            transform: [{ scale: ringScale(r.v) }],
          }}
        />
      ))}

      {/* Proyectiles (ataques) que chocan contra el escudo */}
      {[
        { p: p1, left: 2 * s, top: 40 * s },
        { p: p2, left: 2 * s, top: 58 * s },
        { p: p3, left: 2 * s, top: 76 * s },
      ].map(({ p, left, top }, i) => (
        <Animated.View
          key={i}
          pointerEvents="none"
          style={{
            position: 'absolute',
            left: left - dot / 2,
            top: top - dot / 2,
            width: dot,
            height: dot,
            borderRadius: dot / 2,
            backgroundColor: '#FB923C',
            borderWidth: 1.5,
            borderColor: 'rgba(239,68,68,0.6)',
            opacity: p.o,
            transform: [{ translateX: p.x }, { translateY: p.y }],
          }}
        />
      ))}

      {/* Escena estática: computadora + caballero */}
      <Svg width={width} height={height} viewBox="0 0 160 110" style={StyleSheet.absoluteFill}>
        <Defs>
          <LinearGradient id="shieldGrad" x1="0" y1="0" x2="0" y2="1">
            <Stop offset="0" stopColor="#F87171" />
            <Stop offset="0.55" stopColor="#DC2626" />
            <Stop offset="1" stopColor="#7F1D1D" />
          </LinearGradient>
          <LinearGradient id="ironGrad" x1="0" y1="0" x2="0" y2="1">
            <Stop offset="0" stopColor="#556577" />
            <Stop offset="1" stopColor="#33405A" />
          </LinearGradient>
          <RadialGradient id="screenGlow" cx="0.5" cy="0.5" r="0.5">
            <Stop offset="0" stopColor="#EF4444" stopOpacity="0.4" />
            <Stop offset="1" stopColor="#EF4444" stopOpacity="0" />
          </RadialGradient>
        </Defs>

        {/* Sombra en el piso */}
        <Ellipse cx="86" cy="102" rx="54" ry="6" fill="rgba(0,0,0,0.35)" />

        <G>
          {/* COMPUTADORA (derecha) — monitor, pantalla, torre/base y regleta */}
          <Rect x="118" y="22" width="36" height="46" rx="20" fill="url(#screenGlow)" />
          <Rect x="118" y="24" width="36" height="44" rx="4" fill="#0F172A" stroke="rgba(255,255,255,0.30)" strokeWidth="1.5" />
          <Rect x="123" y="29" width="26" height="30" rx="2" fill="#141C2B" />
          <Rect x="127" y="40" width="6" height="3" rx="1.5" fill="#DC2626" />
          <Rect x="127" y="46" width="16" height="3" rx="1.5" fill="#33405F" />
          <Rect x="132" y="68" width="8" height="10" fill="#0F172A" />
          <Rect x="124" y="78" width="24" height="4" rx="2" fill="#1E293B" />
        </G>

        {/* CABALLERO (centro-izquierda) */}
        <G>
          {/* piernas */}
          <Rect x="66" y="80" width="9" height="16" rx="3" fill="#2B3646" />
          <Rect x="79" y="80" width="9" height="16" rx="3" fill="#2B3646" />
          <Rect x="63" y="94" width="12" height="6" rx="3" fill="#1F2836" />
          <Rect x="80" y="94" width="12" height="6" rx="3" fill="#1F2836" />
          {/* torso */}
          <Rect x="62" y="52" width="30" height="30" rx="10" fill="url(#ironGrad)" />
          <Rect x="62" y="76" width="30" height="6" rx="3" fill="#7F1D1D" />
          {/* hombreras */}
          <Circle cx="66" cy="57" r="7" fill="#5A6B82" />
          <Circle cx="88" cy="57" r="7" fill="#5A6B82" />
          {/* casco */}
          <Path d="M 62 46 V 38 a 15 15 0 0 1 30 0 V 46 Z" fill="url(#ironGrad)" />
          <Path d="M 63 44 h 28" stroke="#111827" strokeWidth="2" strokeLinecap="round" />
          <Path d="M 64 38 h 24" stroke="#EF4444" strokeWidth="3" strokeLinecap="round" />
          {/* brazo hacia el escudo */}
          <Path d="M 84 62 L 56 58" stroke="#46566B" strokeWidth="8" strokeLinecap="round" />
          <Circle cx="56" cy="58" r="5" fill="#2B3646" />
        </G>
      </Svg>

      {/* ESCUDO (capa superior, con latido) */}
      <Animated.View
        pointerEvents="none"
        style={{
          position: 'absolute',
          left: 12 * s,
          top: 33 * s,
          width: 38 * s,
          height: 44.5 * s,
          transform: [{ scale: shieldKick }],
        }}
      >
        <Svg width={38 * s} height={44.5 * s} viewBox="12 6 76 89">
        <Defs>
          <LinearGradient id="shieldGradOverlay" x1="0" y1="0" x2="0" y2="1">
            <Stop offset="0" stopColor="#F87171" />
            <Stop offset="0.55" stopColor="#DC2626" />
            <Stop offset="1" stopColor="#7F1D1D" />
          </LinearGradient>
        </Defs>
        <Path
          d="M50 6 L 88 20 V 48 C 88 70 71 86 50 95 C 29 86 12 70 12 48 V 20 Z"
          fill="url(#shieldGradOverlay)"
          stroke="rgba(255,255,255,0.85)"
          strokeWidth="2.5"
        />
        <Path
          d="M33 50 L 45 62 L 67 39"
          fill="none"
          stroke="#FFFFFF"
          strokeWidth="5"
          strokeLinecap="round"
          strokeLinejoin="round"
        />
      </Svg>
      </Animated.View>
    </View>
  );
}

const styles = StyleSheet.create({});