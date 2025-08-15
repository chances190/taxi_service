import { useParams, useNavigate } from 'react-router-dom';
import { useState, useEffect } from 'react';
import { Paper, Stack, Typography, List, ListItem, ListItemText, CircularProgress } from '@mui/material';
import AppButton from '../../components/ui/AppButton';
import AppAlert from '../../components/ui/AppAlert';
import api from '@services/api';

export default function DocumentReviewPage() {
  const { id } = useParams();
  const navigate = useNavigate();
  const [status, setStatus] = useState<'approved' | 'rejected' | 'in_review'>('in_review');
  const [error, setError] = useState('');
  const [docs, setDocs] = useState<Array<{ tipo_documento: string; caminho_arquivo: string; status: string; formato: string; tamanho: number }>>([]);
  const [opening, setOpening] = useState<string>('');

  const load = async () => {
    if (!id) return;
    try {
      const r = await api.get(`/api/profile/${id}`);
      const m = r.data.motorista;
      setDocs(m?.documentos || []);
    } catch (e) {
      console.error(e);
      setError('Falha ao carregar documentos');
    }
  };

  useEffect(() => { load(); }, [id]);

  // Removida função de validação automática

  const approve = async () => {
    if (!id) return;
    try {
  await api.put(`/api/documents/${id}/approve`);
      setStatus('approved');
      load();
    } catch (err) {
      console.error(err);
      setError('Falha ao aprovar');
    }
  };

  const reject = async () => {
    if (!id) return;
    try {
  await api.put(`/api/documents/${id}/reject`, { motivo: 'Motivo de exemplo' });
      setStatus('rejected');
      load();
    } catch (err) {
      console.error(err);
      setError('Falha ao rejeitar');
    }
  };

  return (
    <Paper sx={{ p: 4 }}>
      <Stack spacing={3}>
  <Typography variant="h5">Revisão de Documentos (Admin)</Typography>
  <AppAlert severity="warning" show>Área restrita a administradores. Autenticação será aplicada quando middleware estiver disponível.</AppAlert>
        {error && <AppAlert severity="error" show>{error}</AppAlert>}
        <AppAlert severity={status === 'approved' ? 'success' : status === 'rejected' ? 'error' : 'info'} show>
          Status: {status === 'in_review' ? 'Em análise' : status === 'approved' ? 'Aprovado' : 'Rejeitado'}
        </AppAlert>
        <Stack direction="row" spacing={2}>
          <AppButton variant="contained" color="success" onClick={approve}>Aprovar</AppButton>
          <AppButton variant="contained" color="error" onClick={reject}>Rejeitar</AppButton>
          <AppButton variant="text" onClick={() => navigate(`/profile/${id}`)}>Voltar ao perfil</AppButton>
        </Stack>
        <Typography variant="h6">Documentos</Typography>
        <List dense>
          {docs.map(d => (
            <ListItem key={d.tipo_documento} divider secondaryAction={
              d.caminho_arquivo ? (
                <AppButton size="small" variant="outlined" onClick={async () => {
                  if (!id) return;
                  setError('');
                  setOpening(d.tipo_documento);
                  try {
                    const url = `${api.defaults.baseURL}/api/documents/${id}/file/${d.tipo_documento}`;
                    // Abre em nova aba com URL absoluto do backend (evita interceptação do router SPA)
                    window.open(url, '_blank', 'noopener');
                  } catch (e) {
                    console.error(e);
                    setError('Falha ao abrir documento');
                  } finally {
                    setOpening('');
                  }
                }} endIcon={opening===d.tipo_documento ? <CircularProgress size={14} /> : undefined}>Abrir</AppButton>
              ) : null
            }>
              <ListItemText
                primary={`${d.tipo_documento} (${d.status})`}
                secondary={`${(d.tamanho/1024).toFixed(1)} KB - ${d.formato}`}
              />
            </ListItem>
          ))}
          {!docs.length && <ListItem><ListItemText primary="Nenhum documento enviado" /></ListItem>}
        </List>
      </Stack>
    </Paper>
  );
}
