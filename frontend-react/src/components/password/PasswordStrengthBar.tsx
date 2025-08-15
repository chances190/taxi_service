import { LinearProgress, Stack, Typography } from '@mui/material';

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
  const value = scorePassword(password);
  let color: 'error' | 'warning' | 'info' | 'success' = 'error';
  if (value > 70) color = 'success';
  else if (value > 50) color = 'info';
  else if (value > 30) color = 'warning';

  return (
    <Stack spacing={0.5}>
      <LinearProgress variant="determinate" value={value} color={color} sx={{ height: 8, borderRadius: 1 }} />
      <Typography variant="caption" color="text.secondary">Força da senha: {value}%</Typography>
    </Stack>
  );
}
