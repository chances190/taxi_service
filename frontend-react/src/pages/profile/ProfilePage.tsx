import { useParams } from 'react-router-dom';
import React from 'react';
import { useQuery } from '@tanstack/react-query';
import { Avatar, Box, Chip, Divider, Paper, Stack, TextField, Typography, LinearProgress } from '@mui/material';
import AppButton from '../../components/ui/AppButton';
import AppAlert from '../../components/ui/AppAlert';
import FieldRow from '../../components/ui/FieldRow';
import { useQueryClient } from '@tanstack/react-query';
import api from '@services/api';
import { useState } from 'react';

interface Motorista {
  id: string;
  nome: string;
  email: string;
  telefone: string;
  cpf?: string;
  cnh?: string;
  categoria_cnh?: string;
  validade_cnh?: string | Date;
  placa_veiculo?: string;
  modelo_veiculo?: string;
  status?: string;
  foto_perfil_url?: string;
}

export default function ProfilePage() {
  const { id } = useParams();
  const queryClient = useQueryClient();
  const { data, isLoading } = useQuery<{ motorista: Motorista }>({
    queryKey: ['motorista', id],
    queryFn: () => api.get(`/api/profile/${id}`).then((r: { data: { motorista: Motorista } }) => r.data),
    enabled: !!id
  });

  const [editing, setEditing] = useState(false);
  const [telefone, setTelefone] = useState('');
  const [email, setEmail] = useState('');
  const [pwEditing, setPwEditing] = useState(false);
  const [senhaAtual, setSenhaAtual] = useState('');
  const [novaSenha, setNovaSenha] = useState('');
  const [confirmacao, setConfirmacao] = useState('');
  const [msg, setMsg] = useState('');
  const [err, setErr] = useState('');
  const [photoUploading, setPhotoUploading] = useState(false);
  const [deletionRequested, setDeletionRequested] = useState(false);

  const motorista = data?.motorista;

  return (
    <Paper sx={{ p: 4 }}>
      {isLoading && <Typography>Carregando...</Typography>}
      {motorista && (
        <Stack spacing={3}>
          <Stack spacing={1} alignItems="center">
            {motorista.foto_perfil_url ? (
              <Avatar sx={{ width: 96, height: 96 }} src={`${motorista.foto_perfil_url.startsWith('http') ? '' : api.defaults.baseURL}${motorista.foto_perfil_url}?t=${Date.now()}`} />
            ) : (
              <Avatar sx={{ width: 96, height: 96 }}>{motorista.nome?.[0]}</Avatar>
            )}
            <Typography variant="h5" textAlign="center">{motorista.nome}</Typography>
            <Typography variant="body2" color="text.secondary" textAlign="center">{motorista.email}</Typography>
            <Stack direction="row" spacing={1} mt={1}>
              {motorista.status && <Chip size="small" label={motorista.status} />}
            </Stack>
          </Stack>
          <Box>
            <Typography variant="subtitle1" gutterBottom textAlign="center">Dados do Perfil</Typography>
            <Stack spacing={1} sx={{ maxWidth:600, mx:'auto' }}>
              {[ 
                { label:'Telefone', value: motorista.telefone },
                { label:'CPF', value: motorista.cpf },
                { label:'CNH', value: motorista.cnh },
                { label:'Categoria CNH', value: motorista.categoria_cnh },
                { label:'Validade CNH', value: motorista.validade_cnh ? new Date(motorista.validade_cnh).toLocaleDateString() : '' },
                { label:'Placa', value: motorista.placa_veiculo },
                { label:'Modelo', value: motorista.modelo_veiculo }
              ].map(f => (
                <FieldRow key={f.label} label={f.label} value={f.value} />
              ))}
            </Stack>
          </Box>
          <Divider flexItem />
          <Stack spacing={3} alignItems="center">
            {msg && <AppAlert severity="success" show onClose={() => setMsg('')}>{msg}</AppAlert>}
            {err && <AppAlert severity="error" show onClose={() => setErr('')}>{err}</AppAlert>}
            <Stack direction="row" spacing={2} justifyContent="center" flexWrap="wrap">
              {!editing && <AppButton onClick={() => { setTelefone(motorista.telefone); setEmail(motorista.email); setEditing(true);} } variant="outlined">Editar Dados</AppButton>}
              {!pwEditing && <AppButton onClick={() => setPwEditing(true)} variant="outlined">Alterar Senha</AppButton>}
            </Stack>
            {editing && (
              <Stack spacing={2} component="form" sx={{ width:'100%', maxWidth:400 }} onSubmit={async (e: React.FormEvent) => { e.preventDefault(); setErr(''); setMsg(''); try { const r = await api.put(`/api/profile/${id}`, { telefone, email }); setMsg(r.data.message); setEditing(false); queryClient.invalidateQueries({ queryKey:['motorista', id] }); } catch(e:any){ setErr(e.response?.data?.error || 'Erro ao atualizar'); } }}>
                <TextField label="Telefone" value={telefone} onChange={e => setTelefone(e.target.value)} size="small" />
                <TextField label="Email" value={email} onChange={e => setEmail(e.target.value)} size="small" />
                <Stack direction="row" spacing={1}>
                  <AppButton type="submit" variant="contained" size="small">Salvar</AppButton>
                  <AppButton variant="text" size="small" onClick={() => setEditing(false)}>Cancelar</AppButton>
                </Stack>
              </Stack>
            )}
            {pwEditing && (
              <Stack spacing={2} component="form" sx={{ width:'100%', maxWidth:400 }} onSubmit={async (e: React.FormEvent) => { e.preventDefault(); setErr(''); setMsg(''); try { const r = await api.put(`/api/profile/${id}/password`, { senha_atual: senhaAtual, nova_senha: novaSenha, confirmacao }); setMsg(r.data.message); setPwEditing(false); setSenhaAtual(''); setNovaSenha(''); setConfirmacao(''); } catch(e:any){ setErr(e.response?.data?.error || 'Erro ao alterar senha'); } }}>
                <TextField label="Senha Atual" type="password" value={senhaAtual} onChange={e => setSenhaAtual(e.target.value)} size="small" />
                <TextField label="Nova Senha" type="password" value={novaSenha} onChange={e => setNovaSenha(e.target.value)} size="small" />
                <TextField label="Confirmar" type="password" value={confirmacao} onChange={e => setConfirmacao(e.target.value)} size="small" />
                <Stack direction="row" spacing={1}>
                  <AppButton type="submit" variant="contained" size="small">Salvar</AppButton>
                  <AppButton variant="text" size="small" onClick={() => setPwEditing(false)}>Cancelar</AppButton>
                </Stack>
              </Stack>
            )}

            <Divider flexItem />
            <Typography variant="subtitle1" textAlign="center">Foto de Perfil</Typography>
            <Stack direction="row" spacing={2} justifyContent="center" alignItems="center" sx={{ width:'100%', maxWidth:400 }}>
              <AppButton component="label" variant="outlined" size="small" loading={photoUploading} loadingPosition="center">
                {photoUploading ? 'Enviando...' : 'Enviar Foto'}
                <input hidden type="file" accept="image/*" onChange={async e => { const file = e.target.files?.[0]; if(!file || !id) return; setErr(''); setMsg(''); setPhotoUploading(true); try { const fd = new FormData(); fd.append('foto', file); const r = await api.post(`/api/profile/${id}/photo`, fd, { headers:{'Content-Type':'multipart/form-data'} }); setMsg(r.data.message); queryClient.invalidateQueries({ queryKey:['motorista', id] }); } catch(er:any){ setErr(er.response?.data?.error || 'Erro no upload'); } finally { setPhotoUploading(false);} }} />
              </AppButton>
              {photoUploading && <LinearProgress sx={{ flex:1 }} />}
            </Stack>

            <Divider flexItem />
            <Typography variant="subtitle1" textAlign="center">Excluir Conta</Typography>
            {!deletionRequested && <AppButton color="error" variant="outlined" size="small" onClick={async ()=>{ if(!id) return; setErr(''); setMsg(''); try { await api.post(`/api/profile/${id}/request-deletion`); setDeletionRequested(true);} catch(e:any){ setErr(e.response?.data?.error || 'Erro ao solicitar exclusão'); } }} sx={{ mx:'auto' }}>Solicitar exclusão</AppButton>}
            {deletionRequested && (
              <AppAlert
                severity="warning"
                show
                sx={{
                  mx:'auto',
                  borderColor:'warning.main',
                  bgcolor: theme => theme.palette.mode === 'dark' ? 'rgba(255,183,77,0.15)' : 'rgba(255,183,77,0.12)',
                }}
              >
                Solicitação enviada. Verifique seu email.
              </AppAlert>
            )}
          </Stack>
        </Stack>
      )}
    </Paper>
  );
}
