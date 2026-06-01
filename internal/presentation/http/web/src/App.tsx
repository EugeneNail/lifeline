import { Navigate, Route, Routes } from 'react-router-dom'
import { SignupPage } from './pages/SignupPage/SignupPage'

function App() {
    return (
        <Routes>
            <Route path="/signup" element={<SignupPage />} />
            <Route path="/" element={<Navigate replace to="/signup" />} />
            <Route path="*" element={<Navigate replace to="/signup" />} />
        </Routes>
    )
}

export default App
