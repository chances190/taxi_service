import { createTheme } from '@mui/material/styles';

const theme = createTheme({
  cssVariables: true,
  palette: {
    mode: 'dark',
    primary: { main: '#00bcd4' },
    secondary: { main: '#ff9800' },
    background: {
      default: '#0d1117',
      paper: '#161b22'
    },
    success: { main: '#2e7d32' },
    error: { main: '#ef5350' },
    warning: { main: '#fbc02d' },
    info: { main: '#0288d1' }
  },
  shape: { borderRadius: 10 },
  typography: {
    fontFamily: 'Inter, Roboto, system-ui, Arial, sans-serif',
    h1: { fontSize: '2.2rem', fontWeight: 600 },
    h2: { fontSize: '1.8rem', fontWeight: 600 },
    h3: { fontSize: '1.4rem', fontWeight: 600 },
    button: { textTransform: 'none', fontWeight: 600 }
  },
  components: {
    MuiButton: { styleOverrides: { root: { borderRadius: 12 } } },
    MuiPaper: { styleOverrides: { root: { backgroundImage: 'none' } } },
    MuiCard: { styleOverrides: { root: { backgroundImage: 'none' } } }
  }
});

export default theme;
