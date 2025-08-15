import { ComponentType, LazyExoticComponent, Suspense, lazy } from 'react';
import { Box, CircularProgress } from '@mui/material';

export default function Loadable<T extends object>(factory: () => Promise<{ default: ComponentType<T> }>): LazyExoticComponent<ComponentType<T>> {
  const LazyComp = lazy(factory);
  // eslint-disable-next-line @typescript-eslint/no-explicit-any
  return ((props: any) => (
    <Suspense
      fallback={
        <Box sx={{ display: 'grid', placeItems: 'center', minHeight: '50vh' }}>
          <CircularProgress />
        </Box>
      }
    >
      <LazyComp {...props} />
    </Suspense>
  )) as unknown as LazyExoticComponent<ComponentType<T>>;
}
