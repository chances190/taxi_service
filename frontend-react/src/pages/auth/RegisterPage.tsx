import { useNavigate } from 'react-router-dom';
import { useForm } from 'react-hook-form';
import { yupResolver } from '@hookform/resolvers/yup';
import * as yup from 'yup';
import { Paper, Typography, Stack, LinearProgress } from '@mui/material';
import AppButton from '../../components/ui/AppButton';
import AppAlert from '../../components/ui/AppAlert';
import FormTextField from '../../components/ui/FormTextField';
import PasswordStrengthBar from '@components/password/PasswordStrengthBar';
import { useMutation } from '@tanstack/react-query';
import api from '@services/api';
import { sanitizeCPF, sanitizeTelefone, sanitizePlaca, sanitizeCNH, sanitizeEmail } from '@shared/format';
import { saveAuth } from '@services/auth';

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

  const onSubmit = (data: RegisterForm) => {
    const payload: RegisterForm = {
      ...data,
      cpf: sanitizeCPF(data.cpf),
      telefone: sanitizeTelefone(data.telefone),
      placa_veiculo: sanitizePlaca(data.placa_veiculo),
      cnh: sanitizeCNH(data.cnh),
      email: sanitizeEmail(data.email)
    };
    mutation.mutate(payload);
  };

  return (
    <Paper sx={{ p: 4 }}>
      <Typography variant="h4" mb={2}>Criar conta</Typography>
      <Stack component="form" spacing={2} onSubmit={handleSubmit(onSubmit)}>
  <FormTextField label="Nome" {...register('nome')} error={!!errors.nome} helperText={errors.nome?.message} />
  <FormTextField label="Data Nascimento (DD/MM/AAAA)" {...register('data_nascimento')} error={!!errors.data_nascimento} helperText={errors.data_nascimento?.message} />
  <FormTextField label="CPF" {...register('cpf')} error={!!errors.cpf} helperText={errors.cpf?.message} />
  <FormTextField label="Telefone" {...register('telefone')} error={!!errors.telefone} helperText={errors.telefone?.message} />
  <FormTextField label="E-mail" {...register('email')} error={!!errors.email} helperText={errors.email?.message} />
  <FormTextField label="CNH" {...register('cnh')} error={!!errors.cnh} helperText={errors.cnh?.message} />
  <FormTextField label="Categoria CNH" {...register('categoria_cnh')} error={!!errors.categoria_cnh} helperText={errors.categoria_cnh?.message} />
  <FormTextField label="Validade CNH (DD/MM/AAAA)" {...register('validade_cnh')} error={!!errors.validade_cnh} helperText={errors.validade_cnh?.message} />
  <FormTextField label="Placa Veículo" {...register('placa_veiculo')} error={!!errors.placa_veiculo} helperText={errors.placa_veiculo?.message} />
  <FormTextField label="Modelo Veículo" {...register('modelo_veiculo')} error={!!errors.modelo_veiculo} helperText={errors.modelo_veiculo?.message} />
  <FormTextField label="Senha" type="password" {...register('senha')} error={!!errors.senha} helperText={errors.senha?.message} />
  <PasswordStrengthBar password={password || ''} />
  <FormTextField label="Confirmar Senha" type="password" {...register('confirmacao_senha')} error={!!errors.confirmacao_senha} helperText={errors.confirmacao_senha?.message} />
        {mutation.isPending && <LinearProgress />}
  {mutation.error && <AppAlert severity="error" show>Falha ao registrar</AppAlert>}
  <AppButton type="submit" variant="contained" size="large" loading={mutation.isPending}>Registrar</AppButton>
  <AppButton variant="text" onClick={() => navigate('/login')}>Já tenho conta</AppButton>
      </Stack>
    </Paper>
  );
}
