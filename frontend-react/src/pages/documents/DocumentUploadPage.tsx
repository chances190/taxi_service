import { useParams, useNavigate } from 'react-router-dom';
import { useState } from 'react';
import { Box, Button, Paper, Stack, Typography, Alert, LinearProgress } from '@mui/material';
import api from '../../services/api';

interface UploadState {
  file?: File;
  preview?: string;
}

export default function DocumentUploadPage() {
  const { id } = useParams();
  const navigate = useNavigate();
  const [state, setState] = useState<UploadState>({});
  const [loading, setLoading] = useState(false);
  const [error, setError] = useState('');

  const handleFile = (e: React.ChangeEvent<HTMLInputElement>) => {
    const file = e.target.files?.[0];
    if (!file) return;
    setState({ file, preview: URL.createObjectURL(file) });
  };

  const upload = async () => {
    if (!state.file || !id) return;
    setLoading(true);
    setError('');
    try {
      const formData = new FormData();
      formData.append('document', state.file);
      await api.post(`/api/documents/${id}/upload`, formData, { headers: { 'Content-Type': 'multipart/form-data' } });
      navigate(`/documents/${id}/review`);
    } catch (err) {
      console.error(err);
      setError('Falha no upload');
    } finally {
      setLoading(false);
    }
  };

  return (
    <Paper sx={{ p: 4 }}>
      <Stack spacing={3}>
        <Typography variant="h5">Enviar Documentos</Typography>
        <Button component="label" variant="outlined" size="large">
          Selecionar documento
          <input hidden type="file" accept="image/*,.pdf" onChange={handleFile} />
        </Button>
        {state.preview && (
          <Box component="img" src={state.preview} alt="preview" sx={{ maxWidth: '100%', borderRadius: 2, border: '1px solid', borderColor: 'divider' }} />
        )}
        {loading && <LinearProgress />}
        {error && <Alert severity="error">{error}</Alert>}
        <Stack direction="row" spacing={2}>
          <Button onClick={upload} variant="contained" disabled={!state.file || loading}>Enviar</Button>
          <Button variant="text" onClick={() => navigate(-1)}>Voltar</Button>
        </Stack>
      </Stack>
    </Paper>
  );
}
