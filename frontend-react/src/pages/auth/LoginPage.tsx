import { useState } from 'react';
import { useNavigate } from 'react-router-dom';
import { useForm } from 'react-hook-form';
import { yupResolver } from '@hookform/resolvers/yup';
import * as yup from 'yup';
import { Paper, Typography, Stack } from '@mui/material';
import AppButton from '../../components/ui/AppButton';
import AppAlert from '../../components/ui/AppAlert';
import FormTextField from '../../components/ui/FormTextField';
import { useMutation } from '@tanstack/react-query';
import api from '@services/api';
import { saveAuth } from '@services/auth';

interface LoginForm {
  email: string;
  password: string;
}

const schema = yup.object({
  email: yup.string().email('E-mail inválido').required('Obrigatório'),
  password: yup.string().required('Obrigatório')
});

export default function LoginPage() {
  const navigate = useNavigate();
  const [serverError, setServerError] = useState('');
  const { register, handleSubmit, formState: { errors } } = useForm<LoginForm>({ resolver: yupResolver(schema) });

  const mutation = useMutation({
    mutationFn: (data: LoginForm) => api.post('/api/auth/login', data).then((r: { data: { motorista: { id: string } } }) => r.data),
    onSuccess: (data: { motorista: { id: string } }) => {
      if (data.motorista?.id) {
        saveAuth({ motoristaId: data.motorista.id, role: 'user' });
        navigate(`/profile/${data.motorista.id}`);
      }
    },
    onError: (err: unknown) => {
      setServerError('Credenciais inválidas');
      console.error(err);
    }
  });

  const onSubmit = (data: LoginForm) => mutation.mutate(data);

  return (
    <Paper sx={{ p: 4, backdropFilter: 'blur(6px)' }}>
      <Typography variant="h4" mb={2}>Entrar</Typography>
      <Stack component="form" onSubmit={handleSubmit(onSubmit)} spacing={2}>
  {serverError && <AppAlert severity="error" show>{serverError}</AppAlert>}
  <FormTextField label="E-mail" type="email" {...register('email')} error={!!errors.email} helperText={errors.email?.message} />
  <FormTextField label="Senha" type="password" {...register('password')} error={!!errors.password} helperText={errors.password?.message} />
  <AppButton type="submit" variant="contained" size="large" loading={mutation.isPending}>Acessar</AppButton>
  <AppButton variant="text" onClick={() => navigate('/register')}>Criar conta</AppButton>
  <AppButton variant="text" onClick={() => navigate('/documents/placeholder/review')} sx={{ alignSelf: 'flex-start' }}>Recuperar conta</AppButton>
      </Stack>
    </Paper>
  );
}
