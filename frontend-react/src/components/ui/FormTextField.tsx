import React from 'react';
import { TextField, TextFieldProps } from '@mui/material';

// Forward ref so react-hook-form register works (otherwise values stay undefined)
const FormTextField = React.forwardRef<HTMLInputElement, TextFieldProps>((props, ref) => {
  const { size, ...rest } = props;
  return (
    <TextField
      size={size || 'small'}
      fullWidth
      inputRef={ref}
      ref={ref as any}
      {...rest}
    />
  );
});

FormTextField.displayName = 'FormTextField';

export default FormTextField;
