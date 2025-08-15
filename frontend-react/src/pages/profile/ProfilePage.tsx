import { useParams } from 'react-router-dom';
import React from 'react';
import { useQuery } from '@tanstack/react-query';
import { Avatar, Box, Button, Chip, Divider, Grid2 as Grid, Paper, Stack, TextField, Typography } from '@mui/material';
import api from '../../services/api';
import { useState } from 'react';

interface Motorista {
  id: string;
  nome: string;
  email: string;
  aprovado?: boolean;
  documentosValidados?: boolean;
}

export default function ProfilePage() {
  const { id } = useParams();
  const { data, isLoading } = useQuery<Motorista>({
    queryKey: ['motorista', id],
  queryFn: () => api.get(`/api/profile/${id}`).then((r: { data: Motorista }): Motorista => r.data),
    enabled: !!id
  });

  const [editing, setEditing] = useState(false);
  const [nome, setNome] = useState('');

  const motorista = data;

  return (
    <Paper sx={{ p: 4 }}>
      {isLoading && <Typography>Carregando...</Typography>}
      {motorista && (
        <Stack spacing={3}>
          <Stack direction="row" spacing={2} alignItems="center">
            <Avatar sx={{ width: 72, height: 72 }}>{motorista.nome?.[0]}</Avatar>
            <Box>
              <Typography variant="h5">{motorista.nome}</Typography>
              <Typography variant="body2" color="text.secondary">{motorista.email}</Typography>
              <Stack direction="row" spacing={1} mt={1}>
                {motorista.aprovado && <Chip size="small" color="success" label="Aprovado" />}
                {motorista.documentosValidados && <Chip size="small" color="info" label="Docs validados" />}
              </Stack>
            </Box>
          </Stack>
          <Divider flexItem />
          <Stack spacing={2}>
            {!editing && <Button onClick={() => setEditing(true)} variant="outlined" sx={{ alignSelf: 'flex-start' }}>Editar Perfil</Button>}
            {editing && (
              <Stack spacing={2} component="form" onSubmit={(e: React.FormEvent) => { e.preventDefault(); setEditing(false); }}>
                <TextField label="Nome" value={nome || motorista.nome} onChange={(e: React.ChangeEvent<HTMLInputElement>) => setNome(e.target.value)} size="small" />
                <Button type="submit" variant="contained" size="small">Salvar</Button>
              </Stack>
            )}
          </Stack>
        </Stack>
      )}
    </Paper>
  );
}
