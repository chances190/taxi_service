import { useParams, useNavigate } from 'react-router-dom';
import { useState } from 'react';
import { Alert, Button, Paper, Stack, Typography } from '@mui/material';
import api from '../../services/api';

export default function DocumentReviewPage() {
  const { id } = useParams();
  const navigate = useNavigate();
  const [status, setStatus] = useState<'pending' | 'approved' | 'rejected'>('pending');
  const [error, setError] = useState('');

  const validate = async () => {
    if (!id) return;
    try {
      await api.post(`/api/documents/${id}/validate`);
      setStatus('pending');
    } catch (err) {
      console.error(err);
      setError('Falha ao validar');
    }
  };

  const approve = async () => {
    if (!id) return;
    try {
      await api.put(`/api/documents/${id}/approve`);
      setStatus('approved');
    } catch (err) {
      console.error(err);
      setError('Falha ao aprovar');
    }
  };

  const reject = async () => {
    if (!id) return;
    try {
      await api.put(`/api/documents/${id}/reject`);
      setStatus('rejected');
    } catch (err) {
      console.error(err);
      setError('Falha ao rejeitar');
    }
  };

  return (
    <Paper sx={{ p: 4 }}>
      <Stack spacing={3}>
        <Typography variant="h5">Revisão de Documentos</Typography>
        {error && <Alert severity="error">{error}</Alert>}
        <Alert severity={status === 'approved' ? 'success' : status === 'rejected' ? 'error' : 'info'}>
          Status: {status === 'pending' ? 'Em validação' : status === 'approved' ? 'Aprovado' : 'Rejeitado'}
        </Alert>
        <Stack direction="row" spacing={2}>
          <Button variant="outlined" onClick={validate}>Validar</Button>
          <Button variant="contained" color="success" onClick={approve}>Aprovar</Button>
          <Button variant="contained" color="error" onClick={reject}>Rejeitar</Button>
          <Button variant="text" onClick={() => navigate(`/profile/${id}`)}>Voltar ao perfil</Button>
        </Stack>
      </Stack>
    </Paper>
  );
}
