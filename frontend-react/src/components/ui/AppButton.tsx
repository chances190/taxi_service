import React from 'react';
import { Button, ButtonProps, CircularProgress } from '@mui/material';

export interface AppButtonProps extends ButtonProps {
  loading?: boolean;
  loadingPosition?: 'start' | 'end' | 'center';
}

/**
 * AppButton centraliza estilos de botões para manter consistência.
 * - padding/size padrão
 * - prop loading para exibir spinner
 */
const AppButton: React.FC<AppButtonProps> = ({ loading, loadingPosition='center', disabled, children, ...rest }) => {
  const spinner = <CircularProgress size={16} color={rest.color === 'inherit' ? 'primary' : rest.color} />;
  const showCenter = loading && loadingPosition === 'center';
  return (
    <Button
      size={rest.size || 'small'}
      disabled={disabled || loading}
      {...rest}
      sx={{ textTransform:'none', fontWeight:500, ...(rest.sx||{}) }}
    >
      {loading && loadingPosition==='start' && spinner}
      {showCenter ? spinner : children}
      {loading && loadingPosition==='end' && spinner}
    </Button>
  );
};

export default AppButton;
