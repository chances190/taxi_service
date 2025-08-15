import React from 'react';
import { TextField, TextFieldProps } from '@mui/material';

const FormTextField: React.FC<TextFieldProps> = (props) => {
  return (
    <TextField size={props.size || 'small'} fullWidth {...props} />
  );
};

export default FormTextField;
