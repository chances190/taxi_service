import { LinearProgress, Stack, Typography } from '@mui/material';
import { useEffect, useState } from 'react';
import api from '../../services/api';

function scorePassword(pwd: string) {
  let score = 0;
  if (!pwd) return 0;
  const letters: Record<string, number> = {};
  for (const char of pwd) {
    letters[char] = (letters[char] || 0) + 1;
    score += 5.0 / letters[char];
  }
  const variations = {
    digits: /\d/.test(pwd),
    lower: /[a-z]/.test(pwd),
    upper: /[A-Z]/.test(pwd),
    nonWords: /[^\w]/.test(pwd)
  };
  let variationCount = 0;
  for (const check in variations) variationCount += variations[check as keyof typeof variations] ? 1 : 0;
  score += (variationCount - 1) * 10;
  return Math.min(100, Math.floor(score));
}

export default function PasswordStrengthBar({ password }: { password: string }) {
  const [label, setLabel] = useState<string>('');
  const localScore = scorePassword(password);
  let color: 'error' | 'warning' | 'info' | 'success' = 'error';
  if (localScore > 70) color = 'success';
  else if (localScore > 50) color = 'info';
  else if (localScore > 30) color = 'warning';

  useEffect(() => {
    let cancelled = false;
    if (!password) {
      setLabel('');
      return;
    }
    const t = setTimeout(async () => {
      try {
  const r = await api.post('/api/utils/check-password', { senha: password });
        if (!cancelled) {
          setLabel(r.data.forca);
        }
      } catch {
        if (!cancelled) setLabel('');
      }
    }, 300); // debounce
    return () => { cancelled = true; clearTimeout(t); };
  }, [password]);

  return (
    <Stack spacing={0.5}>
      <LinearProgress variant="determinate" value={localScore} color={color} sx={{ height: 8, borderRadius: 1 }} />
      <Typography variant="caption" color="text.secondary">Força da senha: {label || localScore + '%'} </Typography>
    </Stack>
  );
}
