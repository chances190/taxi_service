import { useParams, useNavigate } from 'react-router-dom';
import { useState, useEffect } from 'react';
import { Alert, Button, Paper, Stack, Typography, List, ListItem, ListItemText, CircularProgress } from '@mui/material';
import api from '../../services/api';

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
  <Alert severity="warning" variant="outlined">Área restrita a administradores. Autenticação será aplicada quando middleware estiver disponível.</Alert>
        {error && <Alert severity="error">{error}</Alert>}
        <Alert severity={status === 'approved' ? 'success' : status === 'rejected' ? 'error' : 'info'}>
          Status: {status === 'in_review' ? 'Em análise' : status === 'approved' ? 'Aprovado' : 'Rejeitado'}
        </Alert>
        <Stack direction="row" spacing={2}>
          <Button variant="contained" color="success" onClick={approve}>Aprovar</Button>
          <Button variant="contained" color="error" onClick={reject}>Rejeitar</Button>
          <Button variant="text" onClick={() => navigate(`/profile/${id}`)}>Voltar ao perfil</Button>
        </Stack>
        <Typography variant="h6">Documentos</Typography>
        <List dense>
          {docs.map(d => (
            <ListItem key={d.tipo_documento} divider secondaryAction={
              d.caminho_arquivo ? (
                <Button size="small" variant="outlined" onClick={async () => {
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
                }} endIcon={opening===d.tipo_documento ? <CircularProgress size={14} /> : undefined}>Abrir</Button>
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
