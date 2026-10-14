import { ConfirmDialogProvider } from '../shared/ui/ConfirmDialog'
import { AuthProvider } from '../shared/contexts/AuthContext'
import { LanguageProvider } from '../shared/contexts/LanguageContext'
import { AppRoutes } from './router/AppRoutes'

export default function App() {
  return (
    <LanguageProvider>
      <AuthProvider>
        <ConfirmDialogProvider>
          <AppRoutes />
        </ConfirmDialogProvider>
      </AuthProvider>
    </LanguageProvider>
  )
}
