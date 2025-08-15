import React from 'react';
import { Alert, AlertProps, Collapse } from '@mui/material';

export interface AppAlertProps extends AlertProps {
  show?: boolean;
  keepMounted?: boolean;
}

/**
 * AppAlert padroniza alertas (bordas suaves, densidade).
 */
const AppAlert: React.FC<AppAlertProps> = ({ show=true, keepMounted=false, variant='outlined', sx, ...rest }) => {
  if(!keepMounted && !show) return null;
  return (
    <Collapse in={show} appear unmountOnExit={!keepMounted}>
      <Alert
        variant={variant}
        sx={{ borderRadius:1, '& .MuiAlert-message':{ padding:'4px 0' }, ...(sx||{}) }}
        {...rest}
      />
    </Collapse>
  );
};

export default AppAlert;
