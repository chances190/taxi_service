import { useParams, useNavigate } from 'react-router-dom';
import { useState, useEffect, useMemo } from 'react';
import { Box, Paper, Stack, Typography, LinearProgress, Chip } from '@mui/material';
import AppButton from '../../components/ui/AppButton';
import AppAlert from '../../components/ui/AppAlert';
import api from '@services/api';

interface SingleDoc {
  file?: File;
  preview?: string;
  caminho_arquivo?: string;
  formato?: string;
  tamanho?: number;
}

interface UploadState {
  docs: Record<string, SingleDoc>; // CNH, CRLV, selfie_cnh
}

export default function DocumentUploadPage() {
  const { id } = useParams();
  const navigate = useNavigate();
  const [state, setState] = useState<UploadState>({ docs: { CNH: {}, CRLV: {}, selfie_cnh: {} } });
  const [loading, setLoading] = useState(false);
  const [error, setError] = useState('');
  const [success, setSuccess] = useState('');
  const [status, setStatus] = useState<string>('');
  const [documentos, setDocumentos] = useState<Array<{ tipo_documento: string; status: string }>>([]);

  const fetchStatus = async () => {
    if (!id) return;
    try {
      const r = await api.get(`/api/profile/${id}`);
      const m = r.data.motorista;
      setStatus(m?.status || '');
      setDocumentos((m?.documentos || []).map((d: any) => ({ tipo_documento: d.tipo_documento, status: d.status })));
    } catch {
      /* silent */
    }
  };
  const pendentes = useMemo(() => {
    const required = ['CNH', 'CRLV', 'selfie_cnh'];
    return required.filter(req => !documentos.find(d => d.tipo_documento === req));
  }, [documentos]);

  const todosEnviados = pendentes.length === 0 && documentos.length >= 3;

  useEffect(() => { fetchStatus(); }, [id]);

  const handleFile = (tipo: string) => (e: React.ChangeEvent<HTMLInputElement>) => {
    const file = e.target.files?.[0];
    if (!file) return;
    setState(s => ({
      ...s,
      docs: {
        ...s.docs,
        [tipo]: {
          file,
            preview: URL.createObjectURL(file),
            caminho_arquivo: file.name,
            formato: file.name.split('.').pop(),
            tamanho: file.size
        }
      }
    }));
  };

  const uploadBatch = async () => {
    if (!id) return;
    // montar FormData multipart
    const formData = new FormData();
    let index = 0;
    for (const tipo of Object.keys(state.docs)) {
      const d = (state.docs as any)[tipo] as SingleDoc;
      if (d.file) {
        formData.append('files', d.file);
        formData.append(`tipo_${index}`, tipo);
        index++;
      }
    }
    if (index === 0) return;
    setLoading(true);
    setError('');
    setSuccess('');
    try {
      const resp = await api.post(`/api/documents/${id}/upload/files`, formData, { headers: { 'Content-Type': 'multipart/form-data' } });
      setSuccess(resp.data?.message || 'Arquivos enviados');
      setState(s => ({ ...s, docs: { CNH: {}, CRLV: {}, selfie_cnh: {} } }));
      fetchStatus();
    } catch (e) {
      console.error(e);
      setError('Falha no upload de arquivos');
    } finally {
      setLoading(false);
    }
  };

  return (
    <Paper sx={{ p: 4 }}>
      <Stack spacing={3}>
        <Typography variant="h5">Enviar Documentos</Typography>
  {status && (
          <Stack direction="row" spacing={1} alignItems="center" flexWrap="wrap">
            <Typography variant="body2" color="text.secondary">Status atual:</Typography>
            <Chip size="small" label={status} color={status === 'documentos_em_analise' ? 'info' : status === 'aprovado' ? 'success' : 'default'} />
          </Stack>
        )}
        {!!pendentes.length && (
          <AppAlert severity="info" show sx={{ mb: 0 }}>Faltam: {pendentes.join(', ')}</AppAlert>
        )}
  <Stack spacing={2}>
          {(['CNH','CRLV','selfie_cnh'] as const).map(tipo => {
            const d = state.docs[tipo];
            return (
              <Paper key={tipo} variant="outlined" sx={{ p: 2, bgcolor: 'background.default' }}>
                <Stack spacing={1}>
                  <Stack direction="row" spacing={1} alignItems="center">
                    <Typography fontWeight={600}>{tipo === 'selfie_cnh' ? 'Selfie c/ CNH' : tipo}</Typography>
        {d.file && <Chip size="small" color="success" label="Selecionado" />}
                  </Stack>
                  <AppButton component="label" variant="outlined" size="small" sx={{ alignSelf: 'flex-start' }}>
                    {d.file ? 'Trocar arquivo' : 'Selecionar arquivo'}
                    <input hidden type="file" accept="image/*,.pdf" onChange={handleFile(tipo)} />
                  </AppButton>
                  {d.caminho_arquivo && <Typography variant="caption">{d.caminho_arquivo}</Typography>}
                  {d.preview && <Box component="img" src={d.preview} alt={tipo} sx={{ maxWidth: 240, borderRadius: 1, border: '1px solid', borderColor: 'divider' }} />}
                </Stack>
              </Paper>
            );
          })}
        </Stack>
        {loading && <LinearProgress />}
  {error && <AppAlert severity="error" show>{error}</AppAlert>}
  {success && <AppAlert severity="success" show>{success}</AppAlert>}
        <Stack direction="row" spacing={2} flexWrap="wrap">
          <AppButton onClick={uploadBatch} variant="contained" disabled={loading || !Object.values(state.docs).some(d => d.file)} loading={loading}>Enviar Documentos</AppButton>
          <AppButton variant="text" onClick={() => navigate(-1)}>Voltar</AppButton>
          {todosEnviados && status !== 'documentos_em_analise' && (
            <AppButton variant="outlined" color="info" onClick={fetchStatus}>Atualizar Status</AppButton>
          )}
        </Stack>
        {status === 'documentos_em_analise' && (
          <AppAlert severity="info" show>Todos os documentos enviados. Aguarde a análise.</AppAlert>
        )}
        {status === 'aprovado' && (
          <AppAlert severity="success" show>Cadastro aprovado!</AppAlert>
        )}
      </Stack>
    </Paper>
  );
}
