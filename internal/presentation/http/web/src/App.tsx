import { Navigate, Route, Routes } from 'react-router-dom'
import { CreateHabitPage } from './pages/CreateHabitPage/CreateHabitPage'
import { HomePage } from './pages/HomePage/HomePage'
import { LoginPage } from './pages/LoginPage/LoginPage'
import { SignupPage } from './pages/SignupPage/SignupPage'

function App() {
    return (
        <Routes>
            <Route path="/" element={<HomePage />} />
            <Route path="/habits/new" element={<CreateHabitPage />} />
            <Route path="/login" element={<LoginPage />} />
            <Route path="/signup" element={<SignupPage />} />
            <Route path="*" element={<Navigate replace to="/" />} />
        </Routes>
    )
}

export default App
