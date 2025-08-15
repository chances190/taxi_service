import { Suspense } from 'react';
import { Routes, Route, Navigate } from 'react-router-dom';
import Layout from './components/layout/Layout';
import Loadable from './components/common/Loadable';

const LoginPage = Loadable(() => import('./pages/auth/LoginPage'));
const RegisterPage = Loadable(() => import('./pages/auth/RegisterPage'));
const ProfilePage = Loadable(() => import('./pages/profile/ProfilePage'));
const DocumentUploadPage = Loadable(() => import('./pages/documents/DocumentUploadPage'));
const DocumentReviewPage = Loadable(() => import('./pages/documents/DocumentReviewPage'));

export default function App() {
  return (
    <Layout>
      <Suspense fallback={null}>
        <Routes>
          <Route path="/" element={<Navigate to="/login" replace />} />
          <Route path="/login" element={<LoginPage />} />
          <Route path="/register" element={<RegisterPage />} />
          <Route path="/profile/:id" element={<ProfilePage />} />
          <Route path="/documents/:id/upload" element={<DocumentUploadPage />} />
          <Route path="/documents/:id/review" element={<DocumentReviewPage />} />
          <Route path="*" element={<Navigate to="/login" replace />} />
        </Routes>
      </Suspense>
    </Layout>
  );
}
