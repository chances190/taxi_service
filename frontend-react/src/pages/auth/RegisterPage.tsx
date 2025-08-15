import { useNavigate } from 'react-router-dom';
import { useForm } from 'react-hook-form';
import { yupResolver } from '@hookform/resolvers/yup';
import * as yup from 'yup';
import { Paper, Typography, TextField, Stack, Button, Alert, LinearProgress } from '@mui/material';
import { useMutation } from '@tanstack/react-query';
import api from '../../services/api';
import PasswordStrengthBar from '../../components/password/PasswordStrengthBar';

interface RegisterForm {
  nome: string;
  email: string;
  password: string;
  confirmPassword: string;
}

const schema = yup.object({
  nome: yup.string().required('Obrigatório'),
  email: yup.string().email('E-mail inválido').required('Obrigatório'),
  password: yup.string().min(8, 'Min 8 caracteres').required('Obrigatório'),
  confirmPassword: yup.string().oneOf([yup.ref('password')], 'Senhas não conferem').required('Obrigatório')
});

export default function RegisterPage() {
  const navigate = useNavigate();
  const { register, handleSubmit, watch, formState: { errors } } = useForm<RegisterForm>({ resolver: yupResolver(schema) });
  const password = watch('password');

  const mutation = useMutation({
    mutationFn: (data: RegisterForm) => api.post('/api/auth/register', data).then(r => r.data as { id: string }),
    onSuccess: (data: { id: string }) => navigate(`/profile/${data.id}`)
  });

  const onSubmit = (data: RegisterForm) => mutation.mutate(data);

  return (
    <Paper sx={{ p: 4 }}>
      <Typography variant="h4" mb={2}>Criar conta</Typography>
      <Stack component="form" spacing={2} onSubmit={handleSubmit(onSubmit)}>
        <TextField label="Nome" fullWidth size="small" {...register('nome')} error={!!errors.nome} helperText={errors.nome?.message} />
        <TextField label="E-mail" fullWidth size="small" {...register('email')} error={!!errors.email} helperText={errors.email?.message} />
        <TextField label="Senha" type="password" fullWidth size="small" {...register('password')} error={!!errors.password} helperText={errors.password?.message} />
        <PasswordStrengthBar password={password || ''} />
        <TextField label="Confirmar Senha" type="password" fullWidth size="small" {...register('confirmPassword')} error={!!errors.confirmPassword} helperText={errors.confirmPassword?.message} />
        {mutation.isPending && <LinearProgress />}
        {mutation.error && <Alert severity="error">Falha ao registrar</Alert>}
        <Button type="submit" variant="contained" size="large" disabled={mutation.isPending}>Registrar</Button>
        <Button variant="text" onClick={() => navigate('/login')}>Já tenho conta</Button>
      </Stack>
    </Paper>
  );
}
