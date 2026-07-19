import { Navigate, Route, Routes } from 'react-router-dom'
import { DayPage } from './pages/DayPage/DayPage'
import { CreateHabitPage } from './pages/CreateHabitPage/CreateHabitPage'
import { HabitsPage } from './pages/HabitsPage/HabitsPage'
import { EditHabitPage } from './pages/EditHabitPage/EditHabitPage'
import { LoginPage } from './pages/LoginPage/LoginPage'
import { SignupPage } from './pages/SignupPage/SignupPage'

function App() {
    return (
        <Routes>
            <Route path="/" element={<DayPage date="today" />} />
            <Route path="/habits" element={<HabitsPage />} />
            <Route path="/habits/new" element={<CreateHabitPage />} />
            <Route path="/habits/:type/:id" element={<EditHabitPage />} />
            <Route path="/dates/:date" element={<DayPage />} />
            <Route path="/login" element={<LoginPage />} />
            <Route path="/signup" element={<SignupPage />} />
            <Route path="*" element={<Navigate replace to="/" />} />
        </Routes>
    )
}

export default App
