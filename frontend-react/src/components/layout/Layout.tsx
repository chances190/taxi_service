import { PropsWithChildren } from 'react';
import { AppBar, Toolbar, Typography, Container, Box, IconButton } from '@mui/material';
import MenuIcon from '@mui/icons-material/Menu';

export default function Layout({ children }: PropsWithChildren) {
  return (
    <Box sx={{ display: 'flex', minHeight: '100dvh', flexDirection: 'column' }}>
      <AppBar position="static" color="transparent" elevation={0} enableColorOnDark>
        <Toolbar sx={{ gap: 2 }}>
          <IconButton edge="start" color="inherit" aria-label="menu" sx={{ display: { md: 'none' } }}>
            <MenuIcon />
          </IconButton>
          <Typography variant="h6" sx={{ fontWeight: 700 }}>
            Taxi Service
          </Typography>
        </Toolbar>
      </AppBar>
      <Container component="main" sx={{ flexGrow: 1, py: 4, width: '100%', maxWidth: 680 }}>
        {children}
      </Container>
      <Box component="footer" sx={{ py: 2, textAlign: 'center', opacity: 0.6, fontSize: 12 }}>
        &copy; {new Date().getFullYear()} Taxi Service
        <br/>
         Made with :'( and {'</>'} by Sanches
      </Box>
    </Box>
  );
}
