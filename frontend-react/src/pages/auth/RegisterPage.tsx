import { useNavigate } from 'react-router-dom';
import { useForm } from 'react-hook-form';
import { yupResolver } from '@hookform/resolvers/yup';
import * as yup from 'yup';
import { Paper, Typography, TextField, Stack, Button, Alert, LinearProgress } from '@mui/material';
import { useMutation } from '@tanstack/react-query';
import api from '../../services/api';
import { saveAuth } from '../../services/auth';
import PasswordStrengthBar from '../../components/password/PasswordStrengthBar';

interface RegisterForm {
  nome: string;
  data_nascimento: string;
  cpf: string;
  cnh: string;
  categoria_cnh: string;
  validade_cnh: string;
  placa_veiculo: string;
  modelo_veiculo: string;
  telefone: string;
  email: string;
  senha: string;
  confirmacao_senha: string;
}

const schema = yup.object({
  nome: yup.string().required('Obrigatório'),
  data_nascimento: yup.string().required('Obrigatório'),
  cpf: yup.string().required('Obrigatório'),
  cnh: yup.string().required('Obrigatório'),
  categoria_cnh: yup.string().required('Obrigatório'),
  validade_cnh: yup.string().required('Obrigatório'),
  placa_veiculo: yup.string().required('Obrigatório'),
  modelo_veiculo: yup.string().required('Obrigatório'),
  telefone: yup.string().required('Obrigatório'),
  email: yup.string().email('E-mail inválido').required('Obrigatório'),
  senha: yup.string().min(8, 'Min 8 caracteres').required('Obrigatório'),
  confirmacao_senha: yup.string().oneOf([yup.ref('senha')], 'Senhas não conferem').required('Obrigatório')
});

export default function RegisterPage() {
  const navigate = useNavigate();
  const { register, handleSubmit, watch, formState: { errors } } = useForm<RegisterForm>({ resolver: yupResolver(schema) });
  const password = watch('senha');

  const mutation = useMutation({
  mutationFn: (data: RegisterForm) => api.post('/api/auth/register', data).then((r: { data: { motorista?: { id: string }; message?: string } }) => r.data),
    onSuccess: (data: { motorista?: { id: string } }) => {
      if (data.motorista?.id) {
        saveAuth({ motoristaId: data.motorista.id, role: 'user' });
        navigate(`/documents/${data.motorista.id}/upload`);
      }
    }
  });

  const onSubmit = (data: RegisterForm) => mutation.mutate(data);

  return (
    <Paper sx={{ p: 4 }}>
      <Typography variant="h4" mb={2}>Criar conta</Typography>
      <Stack component="form" spacing={2} onSubmit={handleSubmit(onSubmit)}>
  <TextField label="Nome" fullWidth size="small" {...register('nome')} error={!!errors.nome} helperText={errors.nome?.message} />
  <TextField label="Data Nascimento (DD/MM/AAAA)" fullWidth size="small" {...register('data_nascimento')} error={!!errors.data_nascimento} helperText={errors.data_nascimento?.message} />
  <TextField label="CPF" fullWidth size="small" {...register('cpf')} error={!!errors.cpf} helperText={errors.cpf?.message} />
  <TextField label="Telefone" fullWidth size="small" {...register('telefone')} error={!!errors.telefone} helperText={errors.telefone?.message} />
  <TextField label="E-mail" fullWidth size="small" {...register('email')} error={!!errors.email} helperText={errors.email?.message} />
  <TextField label="CNH" fullWidth size="small" {...register('cnh')} error={!!errors.cnh} helperText={errors.cnh?.message} />
  <TextField label="Categoria CNH" fullWidth size="small" {...register('categoria_cnh')} error={!!errors.categoria_cnh} helperText={errors.categoria_cnh?.message} />
  <TextField label="Validade CNH (DD/MM/AAAA)" fullWidth size="small" {...register('validade_cnh')} error={!!errors.validade_cnh} helperText={errors.validade_cnh?.message} />
  <TextField label="Placa Veículo" fullWidth size="small" {...register('placa_veiculo')} error={!!errors.placa_veiculo} helperText={errors.placa_veiculo?.message} />
  <TextField label="Modelo Veículo" fullWidth size="small" {...register('modelo_veiculo')} error={!!errors.modelo_veiculo} helperText={errors.modelo_veiculo?.message} />
  <TextField label="Senha" type="password" fullWidth size="small" {...register('senha')} error={!!errors.senha} helperText={errors.senha?.message} />
  <PasswordStrengthBar password={password || ''} />
  <TextField label="Confirmar Senha" type="password" fullWidth size="small" {...register('confirmacao_senha')} error={!!errors.confirmacao_senha} helperText={errors.confirmacao_senha?.message} />
        {mutation.isPending && <LinearProgress />}
        {mutation.error && <Alert severity="error">Falha ao registrar</Alert>}
        <Button type="submit" variant="contained" size="large" disabled={mutation.isPending}>Registrar</Button>
        <Button variant="text" onClick={() => navigate('/login')}>Já tenho conta</Button>
      </Stack>
    </Paper>
  );
}
