import React from 'react';
import { Box, Typography } from '@mui/material';

export interface FieldRowProps {
  label: string;
  value?: React.ReactNode;
  minLabelWidth?: number;
}

const FieldRow: React.FC<FieldRowProps> = ({ label, value='-', minLabelWidth=120 }) => (
  <Box
    sx={{
      display:'flex',
      justifyContent:'space-between',
      gap:2,
      px:1.25,
      py:0.75,
      borderRadius:1,
      border:'1px solid',
      borderColor:'divider',
      bgcolor: theme => theme.palette.mode === 'dark' ? 'rgba(255,255,255,0.06)' : 'rgba(0,0,0,0.04)'
    }}
  >
    <Typography variant="caption" color="text.secondary" sx={{ minWidth:minLabelWidth, letterSpacing:0.25 }}>{label}</Typography>
    <Typography variant="body2" sx={{ textAlign:'right', flex:1, fontWeight:500 }}>{value || '-'}</Typography>
  </Box>
);

export default FieldRow;
